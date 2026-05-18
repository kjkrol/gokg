package spatial

import (
	"fmt"

	"github.com/kjkrol/gokg/geom"
	"github.com/kjkrol/gokg/plane"
)

const defaultOpsBuffer = 4096

type GridIndexConfig struct {
	Resolution       Resolution
	BucketResolution Resolution
	BucketCapacity   int
	OpsBufferSize    int
}

type BucketDelta struct {
	Bucket  geom.AABB[uint32]
	Added   []EntryId
	Removed []EntryId
	Updated []EntryId
}

type GridIndexManager struct {
	bucketGrid   *bucketGrid
	space        plane.Space2D[uint32]
	opsCh        chan indexOp
	entries      map[uint64]entryCache
	bucketDeltas map[geom.AABB[uint32]]*bucketDelta
	maxGridCord  uint32
}

type entryCache struct {
	mask uint8
}

type bucketDelta struct {
	added   map[EntryId]struct{}
	removed map[EntryId]struct{}
	updated map[EntryId]struct{}
}

type opKind uint8

const (
	opInsert opKind = iota
	opRemove
	opUpdate
)

type indexOp struct {
	kind      opKind
	id        uint64
	aabb      plane.AABB[uint32] // Korzystamy z natywnego typu przestrzeni
	markDirty bool
}

func NewGridIndexManager(space plane.Space2D[uint32], cfg GridIndexConfig) (*GridIndexManager, error) {
	if space == nil {
		return nil, fmt.Errorf("space is required")
	}
	if cfg.Resolution == 0 {
		return nil, fmt.Errorf("world resolution is required")
	}
	if cfg.BucketResolution == 0 {
		return nil, fmt.Errorf("bucket resolution is required")
	}
	if cfg.Resolution < cfg.BucketResolution {
		return nil, fmt.Errorf("bucket resolution must be <= world resolution")
	}
	if cfg.BucketCapacity <= 0 {
		cfg.BucketCapacity = 2
	}
	index, err := NewBucketGrid(
		cfg.Resolution,
		cfg.BucketResolution,
		WithBucketCapacity(cfg.BucketCapacity),
	)
	if err != nil {
		return nil, err
	}
	grid, ok := index.(*bucketGrid)
	if !ok {
		return nil, fmt.Errorf("unexpected bucket grid type")
	}
	maxGridCord := grid.resolution.MaxCoord()
	opsBufferSize := cfg.OpsBufferSize
	if cfg.OpsBufferSize == 0 {
		opsBufferSize = defaultOpsBuffer
	}
	manager := &GridIndexManager{
		bucketGrid:   grid,
		space:        space,
		opsCh:        make(chan indexOp, opsBufferSize),
		entries:      make(map[uint64]entryCache),
		bucketDeltas: make(map[geom.AABB[uint32]]*bucketDelta),
		maxGridCord:  maxGridCord,
	}
	return manager, nil
}

func (m *GridIndexManager) QueueInsert(id uint64, aabb plane.AABB[uint32]) {
	m.opsCh <- indexOp{kind: opInsert, id: id, aabb: aabb, markDirty: true}
}

func (m *GridIndexManager) QueueRemove(id uint64) {
	m.opsCh <- indexOp{kind: opRemove, id: id}
}

func (m *GridIndexManager) QueueUpdate(id uint64, aabb plane.AABB[uint32], markDirty bool) {
	m.opsCh <- indexOp{kind: opUpdate, id: id, aabb: aabb, markDirty: markDirty}
}

func (m *GridIndexManager) Flush(onDirty func(geom.AABB[uint32])) {
	for {
		select {
		case op := <-m.opsCh:
			switch op.kind {
			case opInsert:
				m.applyInsert(op.id, op.aabb, op.markDirty, onDirty)
			case opRemove:
				m.applyRemove(op.id, onDirty)
			case opUpdate:
				m.applyUpdate(op.id, op.aabb, op.markDirty, onDirty)
			}
		default:
			return
		}
	}
}

func (m *GridIndexManager) EntryAABB(entryID EntryId) (geom.AABB[uint32], bool) {
	if m.bucketGrid == nil {
		return geom.AABB[uint32]{}, false
	}
	aabb, ok := m.bucketGrid.aabbById[entryID]
	if !ok {
		return geom.AABB[uint32]{}, false
	}
	return aabb, true
}

// QueryRange przyjmuje teraz czysty wycięty fragment z Broad Phase i sprawdza go bezpośrednio w gridzie
func (m *GridIndexManager) QueryRange(aabb geom.AABB[uint32], collector func(uint64, plane.FragPosition)) int {
	if m.bucketGrid == nil {
		return 0
	}
	if idxAABB, ok := m.indexAABB(aabb); ok {
		return m.bucketGrid.QueryRange(idxAABB, collector)
	}
	return 0
}

func (m *GridIndexManager) ConsumeBucketDeltas() []BucketDelta {
	if len(m.bucketDeltas) == 0 {
		return nil
	}
	out := make([]BucketDelta, 0, len(m.bucketDeltas))
	for rect, delta := range m.bucketDeltas {
		out = append(out, BucketDelta{
			Bucket:  rect,
			Added:   deltaKeys(delta.added),
			Removed: deltaKeys(delta.removed),
			Updated: deltaKeys(delta.updated),
		})
	}
	for rect := range m.bucketDeltas {
		delete(m.bucketDeltas, rect)
	}
	return out
}

func (m *GridIndexManager) recordBucketDelta(rect geom.AABB[uint32]) *bucketDelta {
	delta, ok := m.bucketDeltas[rect]
	if ok {
		return delta
	}
	delta = &bucketDelta{}
	m.bucketDeltas[rect] = delta
	return delta
}

func (d *bucketDelta) add(id EntryId) {
	if d.added == nil {
		d.added = make(map[EntryId]struct{})
	}
	if d.removed != nil {
		delete(d.removed, id)
	}
	if d.updated != nil {
		delete(d.updated, id)
	}
	d.added[id] = struct{}{}
}

func (d *bucketDelta) remove(id EntryId) {
	if d.added != nil {
		if _, ok := d.added[id]; ok {
			delete(d.added, id)
			return
		}
	}
	if d.removed == nil {
		d.removed = make(map[EntryId]struct{})
	}
	if d.updated != nil {
		delete(d.updated, id)
	}
	d.removed[id] = struct{}{}
}

func (d *bucketDelta) update(id EntryId) {
	if d.added != nil {
		if _, ok := d.added[id]; ok {
			return
		}
	}
	if d.removed != nil {
		if _, ok := d.removed[id]; ok {
			return
		}
	}
	if d.updated == nil {
		d.updated = make(map[EntryId]struct{})
	}
	d.updated[id] = struct{}{}
}

func deltaKeys(set map[EntryId]struct{}) []EntryId {
	if len(set) == 0 {
		return nil
	}
	out := make([]EntryId, 0, len(set))
	for id := range set {
		out = append(out, id)
	}
	return out
}

func (m *GridIndexManager) applyInsert(id uint64, shape plane.AABB[uint32], markDirty bool, onDirty func(geom.AABB[uint32])) {
	entries := make([]Entry, 0, 4)
	mask := uint8(0)

	// Ładujemy główne ciało obiektu na bit odpowiadający plane.FRAG_MAIN (0)
	if base, ok := m.indexAABB(shape.AABB); ok {
		entryID := NewEntryID(id, uint8(plane.FRAG_MAIN))
		entries = append(entries, Entry{
			AABB: base,
			Id:   entryID,
		})
		m.recordBucketAdds(entryID, base)
		mask |= 1 << plane.FRAG_MAIN
		if markDirty && onDirty != nil {
			onDirty(base)
		}
	}

	// Ładujemy pozostałe ucięte fragmenty (wykorzystując VisitFragments)
	shape.VisitFragments(func(pos plane.FragPosition, aabb geom.AABB[uint32]) bool {
		idx := uint8(pos)
		if frag, ok := m.indexAABB(aabb); ok {
			entryID := NewEntryID(id, idx)
			entries = append(entries, Entry{
				AABB: frag,
				Id:   entryID,
			})
			m.recordBucketAdds(entryID, frag)
			mask |= 1 << idx
			if markDirty && onDirty != nil {
				onDirty(frag)
			}
		}
		return true
	})

	if len(entries) > 0 {
		if m.bucketGrid != nil {
			m.bucketGrid.BulkInsert(entries)
		}
		m.entries[id] = entryCache{mask: mask}
	}
}

func (m *GridIndexManager) applyRemove(id uint64, onDirty func(geom.AABB[uint32])) {
	cache, ok := m.entries[id]
	if !ok {
		return
	}
	entries := make([]Entry, 0, 4)
	for idx := 0; idx < 4; idx++ {
		if cache.mask&(1<<idx) == 0 {
			continue
		}
		entryID := NewEntryID(id, uint8(idx))
		aabb, ok := m.bucketGrid.aabbById[entryID]
		if !ok {
			continue
		}
		m.recordBucketRemovals(entryID, aabb)
		entries = append(entries, Entry{
			AABB: aabb,
			Id:   entryID,
		})
		if onDirty != nil {
			onDirty(aabb)
		}
	}
	if len(entries) > 0 {
		if m.bucketGrid != nil {
			m.bucketGrid.BulkRemove(entries)
		}
	}
	delete(m.entries, id)
}

func (m *GridIndexManager) applyUpdate(id uint64, shape plane.AABB[uint32], markDirty bool, onDirty func(geom.AABB[uint32])) {
	oldCache, ok := m.entries[id]
	if !ok {
		m.applyInsert(id, shape, markDirty, onDirty)
		return
	}

	var newFrags [4]geom.AABB[uint32]
	newMask := uint8(0)

	if base, ok := m.indexAABB(shape.AABB); ok {
		newFrags[plane.FRAG_MAIN] = base
		newMask |= 1 << plane.FRAG_MAIN
	}

	shape.VisitFragments(func(pos plane.FragPosition, aabb geom.AABB[uint32]) bool {
		idx := uint8(pos)
		if frag, ok := m.indexAABB(aabb); ok {
			newFrags[idx] = frag
			newMask |= 1 << idx
		}
		return true
	})

	if newMask == 0 {
		m.applyRemove(id, onDirty)
		return
	}

	// Jeśli struktura fragmentacji i maska bitowa są identyczne, optymalnie przesuwamy istniejące wpisy
	if oldCache.mask == newMask {
		moves := NewEntriesMove(4)
		for idx := 0; idx < len(newFrags); idx++ {
			if newMask&(1<<idx) == 0 {
				continue
			}
			entryID := NewEntryID(id, uint8(idx))
			oldAABB := m.bucketGrid.aabbById[entryID]
			newAABB := newFrags[idx]
			m.recordBucketUpdates(entryID, oldAABB, newAABB)
			moves.Append(entryID, oldAABB, newAABB)
			if markDirty && onDirty != nil {
				onDirty(oldAABB)
				onDirty(newAABB)
			}
		}
		if len(moves.Old) > 0 && m.bucketGrid != nil {
			m.bucketGrid.BulkMove(moves)
		}
		return
	}

	// W przypadku zmiany liczby/układu fragmentów (np. obiekt wszedł na krawędź lub z niej zszedł) – resetujemy wpis
	m.applyRemove(id, onDirty)
	m.applyInsert(id, shape, markDirty, onDirty)
}

func (m *GridIndexManager) indexAABB(aabb geom.AABB[uint32]) (geom.AABB[uint32], bool) {
	minX := clampU32(aabb.TopLeft.X, m.maxGridCord)
	minY := clampU32(aabb.TopLeft.Y, m.maxGridCord)
	maxX := clampU32(aabb.BottomRight.X, m.maxGridCord)
	maxY := clampU32(aabb.BottomRight.Y, m.maxGridCord)
	if maxX < minX || maxY < minY {
		return geom.AABB[uint32]{}, false
	}
	return NewAABB(
		NewVec(minX, minY),
		NewVec(maxX, maxY),
	), true
}

func (m *GridIndexManager) recordBucketAdds(entryID EntryId, aabb geom.AABB[uint32]) {
	if m.bucketGrid == nil {
		return
	}
	m.bucketGrid.forEachBucketIndex(aabb, func(idx uint32) {
		m.recordBucketDelta(m.bucketRect(idx)).add(entryID)
	})
}

func (m *GridIndexManager) recordBucketRemovals(entryID EntryId, aabb geom.AABB[uint32]) {
	if m.bucketGrid == nil {
		return
	}
	m.bucketGrid.forEachBucketIndex(aabb, func(idx uint32) {
		m.recordBucketDelta(m.bucketRect(idx)).remove(entryID)
	})
}

func (m *GridIndexManager) recordBucketUpdates(entryID EntryId, oldAABB, newAABB geom.AABB[uint32]) {
	if oldAABB == newAABB {
		if m.bucketGrid == nil {
			return
		}
		m.bucketGrid.forEachBucketIndex(newAABB, func(idx uint32) {
			m.recordBucketDelta(m.bucketRect(idx)).update(entryID)
		})
		return
	}

	oldBuckets := m.bucketIndexSet(oldAABB)
	newBuckets := m.bucketIndexSet(newAABB)
	for idx := range newBuckets {
		if _, ok := oldBuckets[idx]; ok {
			m.recordBucketDelta(m.bucketRect(idx)).update(entryID)
		} else {
			m.recordBucketDelta(m.bucketRect(idx)).add(entryID)
		}
	}
	for idx := range oldBuckets {
		if _, ok := newBuckets[idx]; ok {
			continue
		}
		m.recordBucketDelta(m.bucketRect(idx)).remove(entryID)
	}
}

func (m *GridIndexManager) bucketRect(idx uint32) geom.AABB[uint32] {
	if m.bucketGrid == nil {
		return geom.AABB[uint32]{}
	}
	gridSide := m.bucketGrid.gridResolution.Side()
	if gridSide == 0 {
		return geom.AABB[uint32]{}
	}
	bucketSize := m.bucketGrid.bucketsResolution.Side()
	x := idx % gridSide
	y := idx / gridSide
	minX := x * bucketSize
	minY := y * bucketSize
	maxX := minX + bucketSize
	maxY := minY + bucketSize
	return NewAABB(
		NewVec(minX, minY),
		NewVec(maxX, maxY),
	)
}

func (m *GridIndexManager) bucketIndexSet(aabb geom.AABB[uint32]) map[uint32]struct{} {
	seen := make(map[uint32]struct{}, 16)
	if m.bucketGrid == nil {
		return seen
	}
	m.bucketGrid.forEachBucketIndex(aabb, func(idx uint32) {
		seen[idx] = struct{}{}
	})
	return seen
}

func entryID(id uint64, frag uint8) uint64 {
	return (id << 2) | uint64(frag&0x3)
}

func clampU32(val, max uint32) uint32 {
	if val > max {
		return max
	}
	return val
}

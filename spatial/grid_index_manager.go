package spatial

import (
	"fmt"
	"slices"

	"github.com/kjkrol/gokg/geom"
	"github.com/kjkrol/gokg/plane"
	"github.com/kjkrol/uid"
)

const defaultOpsBuffer = 4096

type GridIndexConfig struct {
	Resolution       Resolution
	BucketResolution Resolution
	BucketCapacity   int
	OpsBufferSize    int
	// CellCodec selects the grid-cell codec; zero value (LinearCellCodec)
	// matches today's default behavior.
	CellCodec CellCodecKind
}

type BucketDelta struct {
	Bucket  geom.AABB
	Added   []uid.UID64
	Removed []uid.UID64
	Updated []uid.UID64
}

// GridIndexManager buffers spatial updates and applies them in bulk.
// QueueInsert/QueueRemove/QueueUpdate are safe to call from any goroutine —
// they only send on opsCh. Flush, EntryAABB, and any other read of the
// manager's own state (entries, bucketGrid) are not synchronized and must
// only ever be called from the single goroutine that owns Flush.
type GridIndexManager struct {
	bucketGrid   *bucketGrid
	space        plane.Space2D
	opsCh        chan indexOp
	entries      map[uid.UID64]entryCache
	bucketDeltas map[geom.AABB]*bucketDelta
	maxGridCord  uint32

	// Scratch for the update path, reused rather than allocated per entity per
	// tick. Flush drains opsCh on one goroutine and nothing here is touched
	// from anywhere else, so a single buffer per manager is enough — the same
	// reasoning queryMapPool already applies to the query path. Making Flush
	// concurrent would break this.
	oldBuckets []uint32
	newBuckets []uint32
	moves      EntriesMove
}

type entryCache struct {
	mask uint8
}

type bucketDelta struct {
	added   map[uid.UID64]struct{}
	removed map[uid.UID64]struct{}
	updated map[uid.UID64]struct{}
}

type opKind uint8

const (
	opInsert opKind = iota
	opRemove
	opUpdate
)

type indexOp struct {
	kind      opKind
	id        uid.UID64
	aabb      plane.AABB // Korzystamy z natywnego typu przestrzeni
	markDirty bool
}

func NewGridIndexManager(space plane.Space2D, cfg GridIndexConfig) (*GridIndexManager, error) {
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
	opts := []Option{WithBucketCapacity(cfg.BucketCapacity)}
	if cfg.CellCodec == MortonCellCodec {
		opts = append(opts, WithMortonCodec())
	}
	index, err := NewBucketGrid(cfg.Resolution, cfg.BucketResolution, opts...)
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
		entries:      make(map[uid.UID64]entryCache),
		bucketDeltas: make(map[geom.AABB]*bucketDelta),
		maxGridCord:  maxGridCord,
	}
	return manager, nil
}

func (m *GridIndexManager) QueueInsert(id uid.UID64, aabb plane.AABB) {
	m.opsCh <- indexOp{kind: opInsert, id: id, aabb: aabb, markDirty: true}
}

func (m *GridIndexManager) QueueRemove(id uid.UID64) {
	m.opsCh <- indexOp{kind: opRemove, id: id}
}

func (m *GridIndexManager) QueueUpdate(id uid.UID64, aabb plane.AABB, markDirty bool) {
	m.opsCh <- indexOp{kind: opUpdate, id: id, aabb: aabb, markDirty: markDirty}
}

// Flush drains queued ops and applies them. Call it from a single, fixed
// goroutine only — see the GridIndexManager doc comment.
func (m *GridIndexManager) Flush(onDirty func(geom.AABB)) {
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

// EntryAABB reads m's own state directly (not via opsCh) — call it only
// from the same goroutine that calls Flush, see the GridIndexManager doc
// comment.
func (m *GridIndexManager) EntryAABB(entryID uid.UID64) (geom.AABB, bool) {
	if m.bucketGrid == nil {
		return geom.AABB{}, false
	}
	aabb, ok := m.bucketGrid.aabbById[entryID]
	if !ok {
		return geom.AABB{}, false
	}
	return aabb, true
}

// QueryRange przyjmuje teraz czysty wycięty fragment z Broad Phase i sprawdza go bezpośrednio w gridzie
func (m *GridIndexManager) QueryRange(aabb geom.AABB, collector func(uid.UID64, plane.FragPosition)) int {
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

func (m *GridIndexManager) recordBucketDelta(rect geom.AABB) *bucketDelta {
	delta, ok := m.bucketDeltas[rect]
	if ok {
		return delta
	}
	delta = &bucketDelta{}
	m.bucketDeltas[rect] = delta
	return delta
}

func (d *bucketDelta) add(id uid.UID64) {
	if d.added == nil {
		d.added = make(map[uid.UID64]struct{})
	}
	if d.removed != nil {
		delete(d.removed, id)
	}
	if d.updated != nil {
		delete(d.updated, id)
	}
	d.added[id] = struct{}{}
}

func (d *bucketDelta) remove(id uid.UID64) {
	if d.added != nil {
		if _, ok := d.added[id]; ok {
			delete(d.added, id)
			return
		}
	}
	if d.removed == nil {
		d.removed = make(map[uid.UID64]struct{})
	}
	if d.updated != nil {
		delete(d.updated, id)
	}
	d.removed[id] = struct{}{}
}

func (d *bucketDelta) update(id uid.UID64) {
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
		d.updated = make(map[uid.UID64]struct{})
	}
	d.updated[id] = struct{}{}
}

func deltaKeys(set map[uid.UID64]struct{}) []uid.UID64 {
	if len(set) == 0 {
		return nil
	}
	out := make([]uid.UID64, 0, len(set))
	for id := range set {
		out = append(out, id)
	}
	return out
}

func (m *GridIndexManager) applyInsert(id uid.UID64, shape plane.AABB, markDirty bool, onDirty func(geom.AABB)) {
	entries := make([]Entry, 0, 4)
	mask := uint8(0)

	// Ładujemy główne ciało obiektu na bit odpowiadający plane.FRAG_MAIN (0)
	if base, ok := m.indexAABB(shape.AABB); ok {
		entryID := withFrag(id, uint8(plane.FRAG_MAIN))
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
	shape.VisitFragments(func(pos plane.FragPosition, aabb geom.AABB) bool {
		idx := uint8(pos)
		if frag, ok := m.indexAABB(aabb); ok {
			entryID := withFrag(id, idx)
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

func (m *GridIndexManager) applyRemove(id uid.UID64, onDirty func(geom.AABB)) {
	cache, ok := m.entries[id]
	if !ok {
		return
	}
	entries := make([]Entry, 0, 4)
	for idx := 0; idx < 4; idx++ {
		if cache.mask&(1<<idx) == 0 {
			continue
		}
		entryID := withFrag(id, uint8(idx))
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

func (m *GridIndexManager) applyUpdate(id uid.UID64, shape plane.AABB, markDirty bool, onDirty func(geom.AABB)) {
	oldCache, ok := m.entries[id]
	if !ok {
		m.applyInsert(id, shape, markDirty, onDirty)
		return
	}

	var newFrags [4]geom.AABB
	newMask := uint8(0)

	if base, ok := m.indexAABB(shape.AABB); ok {
		newFrags[plane.FRAG_MAIN] = base
		newMask |= 1 << plane.FRAG_MAIN
	}

	shape.VisitFragments(func(pos plane.FragPosition, aabb geom.AABB) bool {
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
		moves := &m.moves
		moves.Old = moves.Old[:0]
		moves.New = moves.New[:0]
		for idx := 0; idx < len(newFrags); idx++ {
			if newMask&(1<<idx) == 0 {
				continue
			}
			entryID := withFrag(id, uint8(idx))
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
			// BulkMove only reads the batch, so handing it the reused buffers
			// is safe.
			m.bucketGrid.BulkMove(*moves)
		}
		return
	}

	// W przypadku zmiany liczby/układu fragmentów (np. obiekt wszedł na krawędź lub z niej zszedł) – resetujemy wpis
	m.applyRemove(id, onDirty)
	m.applyInsert(id, shape, markDirty, onDirty)
}

func (m *GridIndexManager) indexAABB(aabb geom.AABB) (geom.AABB, bool) {
	minX := clampToGrid(aabb.TopLeft.X, m.maxGridCord)
	minY := clampToGrid(aabb.TopLeft.Y, m.maxGridCord)
	maxX := clampToGrid(aabb.BottomRight.X, m.maxGridCord)
	maxY := clampToGrid(aabb.BottomRight.Y, m.maxGridCord)
	if maxX < minX || maxY < minY {
		return geom.AABB{}, false
	}
	return NewAABB(
		NewVec(minX, minY),
		NewVec(maxX, maxY),
	), true
}

func (m *GridIndexManager) recordBucketAdds(entryID uid.UID64, aabb geom.AABB) {
	if m.bucketGrid == nil {
		return
	}
	m.bucketGrid.forEachBucketIndex(aabb, func(idx uint32) {
		m.recordBucketDelta(m.bucketRect(idx)).add(entryID)
	})
}

func (m *GridIndexManager) recordBucketRemovals(entryID uid.UID64, aabb geom.AABB) {
	if m.bucketGrid == nil {
		return
	}
	m.bucketGrid.forEachBucketIndex(aabb, func(idx uint32) {
		m.recordBucketDelta(m.bucketRect(idx)).remove(entryID)
	})
}

func (m *GridIndexManager) recordBucketUpdates(entryID uid.UID64, oldAABB, newAABB geom.AABB) {
	if oldAABB == newAABB {
		if m.bucketGrid == nil {
			return
		}
		m.bucketGrid.forEachBucketIndex(newAABB, func(idx uint32) {
			m.recordBucketDelta(m.bucketRect(idx)).update(entryID)
		})
		return
	}

	m.splitBuckets(oldAABB, newAABB)

	for _, idx := range m.newBuckets {
		if slices.Contains(m.oldBuckets, idx) {
			m.recordBucketDelta(m.bucketRect(idx)).update(entryID)
		} else {
			m.recordBucketDelta(m.bucketRect(idx)).add(entryID)
		}
	}
	for _, idx := range m.oldBuckets {
		if !slices.Contains(m.newBuckets, idx) {
			m.recordBucketDelta(m.bucketRect(idx)).remove(entryID)
		}
	}
}

// splitBuckets collects into the manager's scratch which buckets each box
// covers, so the caller can tell what an entity joined, kept and left.
//
// The sets are tiny — a box smaller than a bucket lands in one, and even one
// straddling a corner reaches four — so duplicates go out by a linear scan.
// A map would cost more to allocate than the scan costs to run, and the slices
// also give the difference a stable order, which iterating a map did not.
func (m *GridIndexManager) splitBuckets(oldAABB, newAABB geom.AABB) {
	m.oldBuckets = m.oldBuckets[:0]
	m.newBuckets = m.newBuckets[:0]
	if m.bucketGrid == nil {
		return
	}
	m.bucketGrid.forEachBucketIndex(oldAABB, func(idx uint32) {
		if !slices.Contains(m.oldBuckets, idx) {
			m.oldBuckets = append(m.oldBuckets, idx)
		}
	})
	m.bucketGrid.forEachBucketIndex(newAABB, func(idx uint32) {
		if !slices.Contains(m.newBuckets, idx) {
			m.newBuckets = append(m.newBuckets, idx)
		}
	})
}

func (m *GridIndexManager) bucketRect(idx uint32) geom.AABB {
	if m.bucketGrid == nil {
		return geom.AABB{}
	}
	gridSide := m.bucketGrid.gridResolution.Side()
	if gridSide == 0 {
		return geom.AABB{}
	}
	bucketSize := m.bucketGrid.bucketsResolution.Side()
	x := idx % gridSide
	y := idx / gridSide
	minX := x * bucketSize
	minY := y * bucketSize
	maxX := minX + bucketSize
	maxY := minY + bucketSize
	return NewAABB(
		NewVec(float64(minX), float64(minY)),
		NewVec(float64(maxX), float64(maxY)),
	)
}

// clampToGrid holds a coordinate inside the indexable grid. Anything below zero
// or beyond the last cell — including NaN, which fails both comparisons — is
// pulled to the nearer edge rather than being converted out of range later.
func clampToGrid(val float64, max uint32) float64 {
	if !(val >= 0) {
		return 0
	}
	if val > float64(max) {
		return float64(max)
	}
	return val
}

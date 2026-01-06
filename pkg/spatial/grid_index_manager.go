package spatial

import (
	"fmt"

	"github.com/kjkrol/gokg/pkg/plane"
)

const defaultOpsBuffer = 4096

type GridIndexConfig struct {
	Resolution       Resolution
	BucketResolution Resolution
	BucketCapacity   int
	OpsBufferSize    int
}

type BucketDelta struct {
	Bucket  AABB
	Added   []uint64
	Removed []uint64
	Updated []uint64
}

type GridIndexManager struct {
	bucketGrid   *bucketGrid
	space        plane.Space2D[uint32]
	helper       gridHelper
	opsCh        chan indexOp
	entries      map[uint64]entryCache
	bucketDeltas map[AABB]*bucketDelta
	maxGridCord  uint32
}

type gridHelper interface {
	VisitWrappedAABB(aabb AABB, visit func(AABB))
	BuildFragments(shape AABB) (uint8, [4]AABB)
}

type toroidGridHelper struct {
	m *GridIndexManager
}

type standardGridHelper struct {
	m *GridIndexManager
}

type entryCache struct {
	mask uint8
}

type bucketDelta struct {
	added   map[uint64]struct{}
	removed map[uint64]struct{}
	updated map[uint64]struct{}
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
	aabb      AABB
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
		bucketDeltas: make(map[AABB]*bucketDelta),
		maxGridCord:  maxGridCord,
	}
	if space.Name() == "Toroidal2D" {
		manager.helper = &toroidGridHelper{m: manager}
	} else {
		manager.helper = &standardGridHelper{m: manager}
	}
	return manager, nil
}

func (m *GridIndexManager) QueueInsert(id uint64, aabb AABB) {
	m.opsCh <- indexOp{kind: opInsert, id: id, aabb: aabb, markDirty: true}
}

func (m *GridIndexManager) QueueRemove(id uint64) {
	m.opsCh <- indexOp{kind: opRemove, id: id}
}

func (m *GridIndexManager) QueueUpdate(id uint64, aabb AABB, markDirty bool) {
	m.opsCh <- indexOp{kind: opUpdate, id: id, aabb: aabb, markDirty: markDirty}
}

func (m *GridIndexManager) Flush(onDirty func(AABB)) {
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

func (m *GridIndexManager) EntryAABB(entryID uint64) (AABB, bool) {
	if m.bucketGrid == nil {
		return AABB{}, false
	}
	aabb, ok := m.bucketGrid.aabbById[entryID]
	if !ok {
		return AABB{}, false
	}
	return aabb, true
}

func (m *GridIndexManager) QueryRange(aabb AABB, collector func(uint64)) int {
	var count int
	m.VisitWrappedAABB(aabb, func(rect AABB) {
		if m.bucketGrid != nil {
			count += m.bucketGrid.QueryRange(rect, collector)
		}
	})
	return count
}

func (m *GridIndexManager) VisitWrappedAABB(aabb AABB, visit func(AABB)) {
	if m.helper == nil {
		return
	}
	m.helper.VisitWrappedAABB(aabb, visit)
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

func (m *GridIndexManager) recordBucketDelta(rect AABB) *bucketDelta {
	delta, ok := m.bucketDeltas[rect]
	if ok {
		return delta
	}
	delta = &bucketDelta{}
	m.bucketDeltas[rect] = delta
	return delta
}

func (d *bucketDelta) add(id uint64) {
	if d.added == nil {
		d.added = make(map[uint64]struct{})
	}
	if d.removed != nil {
		delete(d.removed, id)
	}
	if d.updated != nil {
		delete(d.updated, id)
	}
	d.added[id] = struct{}{}
}

func (d *bucketDelta) remove(id uint64) {
	if d.added != nil {
		if _, ok := d.added[id]; ok {
			delete(d.added, id)
			return
		}
	}
	if d.removed == nil {
		d.removed = make(map[uint64]struct{})
	}
	if d.updated != nil {
		delete(d.updated, id)
	}
	d.removed[id] = struct{}{}
}

func (d *bucketDelta) update(id uint64) {
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
		d.updated = make(map[uint64]struct{})
	}
	d.updated[id] = struct{}{}
}

func deltaKeys(set map[uint64]struct{}) []uint64 {
	if len(set) == 0 {
		return nil
	}
	out := make([]uint64, 0, len(set))
	for id := range set {
		out = append(out, id)
	}
	return out
}

func (m *GridIndexManager) applyInsert(id uint64, shape AABB, markDirty bool, onDirty func(AABB)) {
	mask, frags := m.buildFragments(shape)
	if mask == 0 {
		return
	}
	entries := make([]Entry, 0, 4)
	for idx := 0; idx < len(frags); idx++ {
		if mask&(1<<idx) == 0 {
			continue
		}
		entryID := entryID(id, uint8(idx))
		entries = append(entries, Entry{
			AABB: frags[idx],
			Id:   entryID,
		})
		m.recordBucketAdds(entryID, frags[idx])
	}
	if len(entries) > 0 {
		if m.bucketGrid != nil {
			m.bucketGrid.BulkInsert(entries)
		}
		m.entries[id] = entryCache{mask: mask}
	}
	if markDirty && onDirty != nil {
		for idx := 0; idx < len(frags); idx++ {
			if mask&(1<<idx) == 0 {
				continue
			}
			onDirty(frags[idx])
		}
	}
}

func (m *GridIndexManager) applyRemove(id uint64, onDirty func(AABB)) {
	cache, ok := m.entries[id]
	if !ok {
		return
	}
	entries := make([]Entry, 0, 4)
	for idx := 0; idx < 4; idx++ {
		if cache.mask&(1<<idx) == 0 {
			continue
		}
		entryID := entryID(id, uint8(idx))
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

func (m *GridIndexManager) applyUpdate(id uint64, shape AABB, markDirty bool, onDirty func(AABB)) {
	oldCache, ok := m.entries[id]
	if !ok {
		m.applyInsert(id, shape, markDirty, onDirty)
		return
	}
	newMask, newFrags := m.buildFragments(shape)
	if newMask == 0 {
		m.applyRemove(id, onDirty)
		return
	}
	if oldCache.mask == newMask {
		moves := NewEntriesMove(4)
		for idx := 0; idx < len(newFrags); idx++ {
			if newMask&(1<<idx) == 0 {
				continue
			}
			entryID := entryID(id, uint8(idx))
			oldAABB := m.bucketGrid.aabbById[entryID]
			newAABB := newFrags[idx]
			m.recordBucketUpdates(entryID, oldAABB, newAABB)
			moves.Append(entryID, oldAABB, newAABB)
			if markDirty && onDirty != nil {
				onDirty(oldAABB)
				onDirty(newAABB)
			}
		}
		if len(moves.Old) > 0 {
			if m.bucketGrid != nil {
				m.bucketGrid.BulkMove(moves)
			}
		}
		return
	}

	m.applyRemove(id, onDirty)
	m.applyInsert(id, shape, markDirty, onDirty)
}

func (m *GridIndexManager) buildFragments(shape AABB) (uint8, [4]AABB) {
	if m.helper == nil {
		return 0, [4]AABB{}
	}
	return m.helper.BuildFragments(shape)
}

func (m *GridIndexManager) indexAABB(aabb AABB) (AABB, bool) {
	minX := clampU32(aabb.TopLeft.X, m.maxGridCord)
	minY := clampU32(aabb.TopLeft.Y, m.maxGridCord)
	maxX := clampU32(aabb.BottomRight.X, m.maxGridCord)
	maxY := clampU32(aabb.BottomRight.Y, m.maxGridCord)
	if maxX < minX || maxY < minY {
		return AABB{}, false
	}
	return NewAABB(
		NewVec(minX, minY),
		NewVec(maxX, maxY),
	), true
}

func (m *GridIndexManager) recordBucketAdds(entryID uint64, aabb AABB) {
	if m.bucketGrid == nil {
		return
	}
	m.bucketGrid.forEachBucketIndex(aabb, func(idx uint32) {
		m.recordBucketDelta(m.bucketRect(idx)).add(entryID)
	})
}

func (m *GridIndexManager) recordBucketRemovals(entryID uint64, aabb AABB) {
	if m.bucketGrid == nil {
		return
	}
	m.bucketGrid.forEachBucketIndex(aabb, func(idx uint32) {
		m.recordBucketDelta(m.bucketRect(idx)).remove(entryID)
	})
}

func (m *GridIndexManager) recordBucketUpdates(entryID uint64, oldAABB, newAABB AABB) {
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

func (m *GridIndexManager) bucketRect(idx uint32) AABB {
	if m.bucketGrid == nil {
		return AABB{}
	}
	gridSide := m.bucketGrid.gridResolution.Side()
	if gridSide == 0 {
		return AABB{}
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

func (m *GridIndexManager) bucketIndexSet(aabb AABB) map[uint32]struct{} {
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

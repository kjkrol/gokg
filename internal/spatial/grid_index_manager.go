package spatial

import (
	"fmt"

	"github.com/kjkrol/aabbworld/geom"
	iplane "github.com/kjkrol/aabbworld/internal/plane"
	"github.com/kjkrol/aabbworld/plane"
	"github.com/kjkrol/uid"
)

type GridIndexConfig struct {
	Resolution       Resolution
	BucketResolution Resolution
	BucketCapacity   int
	// CellCodec selects the grid-cell codec; the zero value is LinearCellCodec.
	CellCodec CellCodecKind
}

// GridIndexManager buffers spatial updates and applies them in bulk, all from one goroutine.
type GridIndexManager struct {
	bucketGrid  *bucketGrid
	space       *iplane.Surface
	ops         []indexOp
	entries     map[uid.UID64]entryCache
	maxGridCord uint32

	// moves is scratch for the update path, owned by the goroutine that calls Flush.
	moves EntriesMove

	sweep pairSweep
}

type entryCache struct {
	mask uint8
}

type opKind uint8

const (
	opInsert opKind = iota
	opRemove
	opUpdate
	opSetCaps
)

type indexOp struct {
	kind      opKind
	id        uid.UID64
	aabb      plane.AABB // Korzystamy z natywnego typu przestrzeni
	caps      Capability
	markDirty bool
}

func NewGridIndexManager(space *iplane.Surface, cfg GridIndexConfig) (*GridIndexManager, error) {
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
	manager := &GridIndexManager{
		bucketGrid:  grid,
		space:       space,
		entries:     make(map[uid.UID64]entryCache),
		maxGridCord: maxGridCord,
	}
	return manager, nil
}

func (m *GridIndexManager) QueueInsert(id uid.UID64, aabb plane.AABB) {
	m.ops = append(m.ops, indexOp{kind: opInsert, id: id, aabb: aabb, markDirty: true})
}

func (m *GridIndexManager) QueueRemove(id uid.UID64) {
	m.ops = append(m.ops, indexOp{kind: opRemove, id: id})
}

func (m *GridIndexManager) QueueUpdate(id uid.UID64, aabb plane.AABB, markDirty bool) {
	m.ops = append(m.ops, indexOp{kind: opUpdate, id: id, aabb: aabb, markDirty: markDirty})
}

// QueueSetCapabilities queues id's capabilities, in order with every other change.
func (m *GridIndexManager) QueueSetCapabilities(id uid.UID64, c Capability) {
	m.ops = append(m.ops, indexOp{kind: opSetCaps, id: id, caps: c})
}

// Flush applies the queued ops; call it from a single, fixed goroutine.
func (m *GridIndexManager) Flush(onDirty func(geom.AABB)) {
	for i := range m.ops {
		op := &m.ops[i]
		switch op.kind {
		case opInsert:
			m.applyInsert(op.id, op.aabb, op.markDirty, onDirty)
		case opRemove:
			m.applyRemove(op.id, onDirty, true)
		case opUpdate:
			m.applyUpdate(op.id, op.aabb, op.markDirty, onDirty)
		case opSetCaps:
			m.bucketGrid.boxes.setCaps(op.id, op.caps)
		}
	}
	m.ops = m.ops[:0]
}

// EntryAABB returns the indexed box of entryID; call it from the goroutine that calls Flush.
func (m *GridIndexManager) EntryAABB(entryID uid.UID64) (geom.AABB, bool) {
	if m.bucketGrid == nil {
		return geom.AABB{}, false
	}
	return m.bucketGrid.boxes.get(entryID)
}

// QueryRange calls collector for every entry intersecting aabb and returns how many.
func (m *GridIndexManager) QueryRange(aabb geom.AABB, collector func(uid.UID64)) int {
	return m.QueryRangeWith(aabb, AnyCapability, collector)
}

// QueryRangeWith is QueryRange restricted to entries sharing a capability with want.
func (m *GridIndexManager) QueryRangeWith(aabb geom.AABB, want Capability, collector func(uid.UID64)) int {
	if m.bucketGrid == nil {
		return 0
	}
	if idxAABB, ok := m.indexAABB(aabb); ok {
		return m.bucketGrid.QueryRangeWith(idxAABB, want, collector)
	}
	return 0
}

func (m *GridIndexManager) applyInsert(id uid.UID64, shape plane.AABB, markDirty bool, onDirty func(geom.AABB)) {
	m.bucketGrid.boxes.setSize(id, shape.Size)
	entries := make([]Entry, 0, 4)
	mask := uint8(0)

	if base, ok := m.indexAABB(shape.AABB); ok {
		entryID := withFrag(id, uint8(plane.FRAG_MAIN))
		entries = append(entries, Entry{
			AABB: base,
			Id:   entryID,
		})
		mask |= 1 << plane.FRAG_MAIN
		if markDirty && onDirty != nil {
			onDirty(base)
		}
	}

	shape.VisitFragments(func(pos plane.FragPosition, aabb geom.AABB) bool {
		idx := uint8(pos)
		if frag, ok := m.indexAABB(aabb); ok {
			entryID := withFrag(id, idx)
			entries = append(entries, Entry{
				AABB: frag,
				Id:   entryID,
			})
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

// applyRemove drops every piece of id from the grid, and its capabilities too if clearCaps.
func (m *GridIndexManager) applyRemove(id uid.UID64, onDirty func(geom.AABB), clearCaps bool) {
	cache, ok := m.entries[id]
	if !ok {
		return
	}
	if clearCaps {
		defer m.bucketGrid.boxes.clearCaps(id)
	}
	entries := make([]Entry, 0, 4)
	for idx := 0; idx < 4; idx++ {
		if cache.mask&(1<<idx) == 0 {
			continue
		}
		entryID := withFrag(id, uint8(idx))
		aabb, ok := m.bucketGrid.boxes.get(entryID)
		if !ok {
			continue
		}
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
	m.bucketGrid.boxes.setSize(id, shape.Size)

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
		m.applyRemove(id, onDirty, false)
		return
	}

	if oldCache.mask == newMask {
		moves := &m.moves
		moves.Old = moves.Old[:0]
		moves.New = moves.New[:0]
		for idx := 0; idx < len(newFrags); idx++ {
			if newMask&(1<<idx) == 0 {
				continue
			}
			entryID := withFrag(id, uint8(idx))
			oldAABB, _ := m.bucketGrid.boxes.get(entryID)
			newAABB := newFrags[idx]
			moves.Append(entryID, oldAABB, newAABB)
			if markDirty && onDirty != nil {
				onDirty(oldAABB)
				onDirty(newAABB)
			}
		}
		if len(moves.Old) > 0 && m.bucketGrid != nil {
			m.bucketGrid.BulkMove(*moves)
		}
		return
	}

	m.applyRemove(id, onDirty, false)
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

// clampToGrid holds a coordinate inside the indexable grid, NaN included.
func clampToGrid(val float64, max uint32) float64 {
	if !(val >= 0) {
		return 0
	}
	if val > float64(max) {
		return float64(max)
	}
	return val
}

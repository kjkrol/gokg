package spatial

import (
	"github.com/kjkrol/gokg/pkg/geom"
)

// Index is a discrete spatial index over a 2D power-of-two grid.
// It stores objects at integer coordinates and supports point lookups,
// range queries (AABB) and bulk operations (insert, remove, move).
type (
	Vec  = geom.Vec[uint32]
	AABB = geom.AABB[uint32]

	Index interface {
		// BulkInsert – insert many objects at once.
		BulkInsert(entries []Entry)

		// BulkRemove – remove whatever is stored at the given positions.
		BulkRemove(entries []Entry)

		// BulkMove – update objects (typically same Value, different XY).
		BulkMove(moves EntriesMove)

		// QueryRange – all objects within the AABB.
		// Collector cannot modify Index.
		QueryRange(aabb AABB, collector func(uint64)) int

		// Count – number of objects in the structure.
		Count() int

		// Bounds – global bounds of the handled space.
		Bounds() AABB

		Optimize()

		Clear()
	}

	GridIndexer interface {
		Index
		CalculateGridIndex(vec Vec) int
	}

	Entry struct {
		AABB
		Id uint64
	}

	EntriesMove struct {
		Old []Entry
		New []Entry
	}
)

var (
	NewVec    = geom.NewVec[uint32]
	NewAABB   = geom.NewAABB[uint32]
	NewAABBAt = geom.NewAABBAt[uint32]
)

func NewEntriesMove(capHint int) EntriesMove {
	return EntriesMove{
		Old: make([]Entry, 0, capHint),
		New: make([]Entry, 0, capHint),
	}
}

func (u *EntriesMove) Append(id uint64, old, new AABB) {
	u.Old = append(u.Old, Entry{
		AABB: old,
		Id:   id,
	})

	u.New = append(u.New, Entry{
		AABB: new,
		Id:   id,
	})
}

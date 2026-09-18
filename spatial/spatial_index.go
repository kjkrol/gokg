package spatial

import (
	"github.com/kjkrol/gokg/geom"
	"github.com/kjkrol/gokg/plane"
	"github.com/kjkrol/uid"
)

// Index is a discrete spatial index over a 2D power-of-two grid.
// It stores objects at integer coordinates and supports point lookups,
// range queries (AABB) and bulk operations (insert, remove, move).
type (
	Vec  = geom.Vec
	AABB = geom.AABB

	Index interface {
		// BulkInsert – insert many objects at once.
		BulkInsert(entries []Entry)

		// BulkRemove – remove whatever is stored at the given positions.
		BulkRemove(entries []Entry)

		// BulkMove – update objects (typically same Value, different XY).
		BulkMove(moves EntriesMove)

		// QueryRange – all objects within the AABB.
		// Collector cannot modify Index.
		QueryRange(aabb AABB, collector func(uid.UID64, plane.FragPosition)) int

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
		Id uid.UID64
	}

	EntriesMove struct {
		Old []Entry
		New []Entry
	}
)

var (
	NewVec    = geom.NewVec
	NewAABB   = geom.NewAABB
	NewAABBAt = geom.NewAABBAt
)

func NewEntriesMove(capHint int) EntriesMove {
	return EntriesMove{
		Old: make([]Entry, 0, capHint),
		New: make([]Entry, 0, capHint),
	}
}

func (u *EntriesMove) Append(id uid.UID64, old, new AABB) {
	u.Old = append(u.Old, Entry{
		AABB: old,
		Id:   id,
	})

	u.New = append(u.New, Entry{
		AABB: new,
		Id:   id,
	})
}

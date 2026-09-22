package spatial

import (
	"github.com/kjkrol/aabbworld/plane"
	"testing"

	iplane "github.com/kjkrol/aabbworld/internal/plane"
	"github.com/kjkrol/uid"
)

// cellsOf lists the buckets that hold id.
func cellsOf(m *GridIndexManager, id uid.UID64) []int {
	var cells []int
	for i := range m.bucketGrid.buckets {
		for _, held := range m.bucketGrid.buckets[i].ids {
			if held == id {
				cells = append(cells, i)
			}
		}
	}
	return cells
}

func TestUpdate_ABoxMovingWithinItsCellsIsReBucketedOnlyWhenTheyChange(t *testing.T) {
	space := iplane.NewEuclidean2D(256, 256)
	m, err := NewGridIndexManager(space, GridIndexConfig{
		Resolution: Size256x256, BucketResolution: Size32x32, BucketCapacity: 8,
	})
	if err != nil {
		t.Fatalf("NewGridIndexManager: %v", err)
	}
	id := uid.UID64(7)
	at := func(x float64) plane.AABB { return space.WrapAABB(NewAABBAt(NewVec(x, 100), 8, 8)) }

	m.QueueInsert(id, at(10))
	m.Flush(nil)
	home := cellsOf(m, id)
	if len(home) != 1 {
		t.Fatalf("held in %v cells at rest, want one", home)
	}

	m.QueueUpdate(id, at(20), true)
	m.Flush(nil)
	if got := cellsOf(m, id); len(got) != 1 || got[0] != home[0] {
		t.Errorf("held in %v after a move within the cell, want still %v", got, home)
	}
	if box, ok := m.EntryAABB(id); !ok || box.TopLeft.X != 20 {
		t.Errorf("EntryAABB = %v (%v), want the box moved to x=20", box, ok)
	}

	m.QueueUpdate(id, at(100), true)
	m.Flush(nil)
	if got := cellsOf(m, id); len(got) != 1 || got[0] == home[0] {
		t.Errorf("held in %v after crossing into another cell, want exactly one new cell", got)
	}
	if n := m.QueryRange(NewAABBAt(NewVec(0, 90), 40, 20), func(uid.UID64) {}); n != 0 {
		t.Errorf("a query over the old place still finds %d entries", n)
	}
	if n := m.QueryRange(NewAABBAt(NewVec(95, 95), 20, 20), func(uid.UID64) {}); n != 1 {
		t.Errorf("a query over the new place finds %d entries, want 1", n)
	}

	m.QueueUpdate(id, space.WrapAABB(NewAABBAt(NewVec(28, 28), 8, 8)), true)
	m.Flush(nil)
	if got := cellsOf(m, id); len(got) != 4 {
		t.Errorf("held in %v straddling four cells, want four", got)
	}
}

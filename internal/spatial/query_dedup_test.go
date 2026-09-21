package spatial

import (
	"testing"

	"github.com/kjkrol/aabbworld/geom"
	"github.com/kjkrol/uid"
)

// countReports runs a query and returns how many times each entity reached the collector.
func countReports(t *testing.T, bg *bucketGrid, box AABB) map[uid.UID64]int {
	t.Helper()
	seen := map[uid.UID64]int{}
	bg.QueryRange(box, func(id uid.UID64) { seen[id]++ })
	return seen
}

func newGrid(t *testing.T, bucket Resolution) *bucketGrid {
	t.Helper()
	index, err := NewBucketGrid(Size256x256, bucket, WithBucketCapacity(8))
	if err != nil {
		t.Fatalf("NewBucketGrid: %v", err)
	}
	return index.(*bucketGrid)
}

func TestQueryRange_ReportsAnEntrySpanningFourCellsOnce(t *testing.T) {
	bg := newGrid(t, Size32x32)

	id := uid.UID64(1)
	bg.BulkInsert([]Entry{{Id: id, AABB: NewAABBAt(NewVec(56, 56), 16, 16)}})

	got := countReports(t, bg, geom.NewAABB(NewVec(40, 40), NewVec(88, 88)))

	if got[id] != 1 {
		t.Errorf("entity reported %d times, want exactly 1", got[id])
	}
}

func TestQueryRange_ReportsEveryEntryItShould(t *testing.T) {
	bg := newGrid(t, Size32x32)

	var entries []Entry
	for i := range 20 {
		entries = append(entries, Entry{
			Id:   uid.UID64(i + 1),
			AABB: NewAABBAt(NewVec(float64(i*12), float64(i*12)), 4, 4),
		})
	}
	bg.BulkInsert(entries)

	got := countReports(t, bg, geom.NewAABB(NewVec(0, 0), NewVec(255, 255)))

	for _, e := range entries {
		switch got[e.Id] {
		case 1:
		case 0:
			t.Errorf("entity %v was never reported, though it lies inside the query", e.Id)
		default:
			t.Errorf("entity %v reported %d times, want exactly 1", e.Id, got[e.Id])
		}
	}
}

func TestQueryRange_ReportsAnEntryStartingBeforeTheQuery(t *testing.T) {
	bg := newGrid(t, Size32x32)

	id := uid.UID64(1)
	bg.BulkInsert([]Entry{{Id: id, AABB: NewAABBAt(NewVec(24, 24), 16, 16)}})

	got := countReports(t, bg, geom.NewAABB(NewVec(33, 33), NewVec(90, 90)))

	if got[id] != 1 {
		t.Errorf("entity reported %d times, want exactly 1", got[id])
	}
}

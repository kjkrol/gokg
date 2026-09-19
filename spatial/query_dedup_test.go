package spatial

import (
	"testing"

	"github.com/kjkrol/gokg/geom"
	"github.com/kjkrol/gokg/plane"
	"github.com/kjkrol/uid"
)

// countReports runs a query and returns how many times each entity was handed
// to the collector — not how many distinct ones were found, which is what a
// map-based check would have measured and what would have hidden the bug.
func countReports(t *testing.T, bg *bucketGrid, box AABB) map[uid.UID64]int {
	t.Helper()
	seen := map[uid.UID64]int{}
	bg.QueryRange(box, func(id uid.UID64, _ plane.FragPosition) { seen[id]++ })
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

// An entry wide enough to sit in four cells at once used to be met four times
// by a query that also spans them, and a set of already-seen ids was what kept
// it from being reported four times over. The owner-cell rule replaces that
// set, so this is the test that says whether it works.
func TestQueryRange_ReportsAnEntrySpanningFourCellsOnce(t *testing.T) {
	bg := newGrid(t, Size32x32)

	// Straddling the corner where cells (1,1), (2,1), (1,2) and (2,2) meet.
	id := uid.UID64(1)
	bg.BulkInsert([]Entry{{Id: id, AABB: NewAABBAt(NewVec(56, 56), 16, 16)}})

	// A query box spanning the same four cells.
	got := countReports(t, bg, geom.NewAABB(NewVec(40, 40), NewVec(88, 88)))

	if got[id] != 1 {
		t.Errorf("entity reported %d times, want exactly 1", got[id])
	}
}

func TestQueryRange_ReportsEveryEntryItShould(t *testing.T) {
	bg := newGrid(t, Size32x32)

	var entries []Entry
	for i := range 20 {
		// A diagonal line of small boxes, each in its own cell.
		entries = append(entries, Entry{
			Id:   uid.UID64(i + 1),
			AABB: NewAABBAt(NewVec(float64(i*12), float64(i*12)), 4, 4),
		})
	}
	bg.BulkInsert(entries)

	// 255, not 256: the last addressable cell is 255>>5, and bucketGrid does
	// not clamp the way GridIndexManager does before calling it.
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

// An entry whose cells start before the query's do must still be reported, from
// the first cell the two have in common — the case a naive "report only from
// the entry's own first cell" rule would silently drop.
func TestQueryRange_ReportsAnEntryStartingBeforeTheQuery(t *testing.T) {
	bg := newGrid(t, Size32x32)

	id := uid.UID64(1)
	// Spans cells 0 and 1 on both axes.
	bg.BulkInsert([]Entry{{Id: id, AABB: NewAABBAt(NewVec(24, 24), 16, 16)}})

	// The query starts in cell 1, so the entry's own first cell (0) is outside it.
	got := countReports(t, bg, geom.NewAABB(NewVec(33, 33), NewVec(90, 90)))

	if got[id] != 1 {
		t.Errorf("entity reported %d times, want exactly 1", got[id])
	}
}

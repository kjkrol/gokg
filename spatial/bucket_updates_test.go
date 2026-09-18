package spatial

import (
	"testing"

	"github.com/kjkrol/gokg/geom"
	"github.com/kjkrol/gokg/plane"
	"github.com/kjkrol/uid"
)

// deltaTally is what one flush reported, summed over every bucket it touched.
// Which bucket is which does not matter here; how many an entity joined, kept
// and left is exactly what the set difference decides.
type deltaTally struct{ added, removed, updated int }

type mover struct {
	t     *testing.T
	m     *GridIndexManager
	space plane.Space2D
	id    uid.UID64
	size  float64
}

// newMover starts an entity of the given size at (x, y) in a 256x256 world cut
// into 32x32 buckets, and drains the insert so later moves report only
// themselves.
func newMover(t *testing.T, x, y, size float64) *mover {
	t.Helper()

	space := plane.NewToroidal2D(256, 256)
	m, err := NewGridIndexManager(space, GridIndexConfig{
		Resolution:       Size256x256,
		BucketResolution: Size32x32,
		BucketCapacity:   8,
	})
	if err != nil {
		t.Fatalf("NewGridIndexManager: %v", err)
	}

	e := &mover{t: t, m: m, space: space, id: uid.UID64(1), size: size}
	m.QueueInsert(e.id, space.WrapAABB(geom.NewAABBAt(geom.NewVec(x, y), size, size)))
	m.Flush(nil)
	m.ConsumeBucketDeltas()
	return e
}

// moveTo moves the entity and returns what that move reported.
func (e *mover) moveTo(x, y float64) deltaTally {
	e.t.Helper()

	e.m.QueueUpdate(e.id, e.space.WrapAABB(geom.NewAABBAt(geom.NewVec(x, y), e.size, e.size)), false)
	e.m.Flush(nil)

	var got deltaTally
	for _, d := range e.m.ConsumeBucketDeltas() {
		got.added += len(d.Added)
		got.removed += len(d.Removed)
		got.updated += len(d.Updated)
	}
	return got
}

// The set difference is the whole of what changed here, so this is the test
// that matters: an entity staying put in its bucket must not be reported as
// having joined or left one, and an entity crossing a boundary must be
// reported as leaving exactly what it left and joining exactly what it joined.
func TestBucketDeltas_ReportWhatTheEntityActuallyJoinedAndLeft(t *testing.T) {
	cases := map[string]struct {
		size     float64
		from, to [2]float64
		want     deltaTally
	}{
		"shifting inside one bucket touches nothing else": {
			size: 8, from: [2]float64{40, 40}, to: [2]float64{44, 44},
			want: deltaTally{updated: 1},
		},
		"crossing into the next bucket leaves one and joins one": {
			size: 8, from: [2]float64{40, 40}, to: [2]float64{72, 40},
			want: deltaTally{added: 1, removed: 1},
		},
		"crossing diagonally leaves one and joins one": {
			size: 8, from: [2]float64{40, 40}, to: [2]float64{72, 72},
			want: deltaTally{added: 1, removed: 1},
		},
		"sliding a four-bucket box sideways keeps the column it shares": {
			// 24 wide at (40,40) reaches 64 on both axes, so it covers the
			// four buckets around that corner. Moving to x=72 shifts it one
			// column: two buckets stay, two are joined and two are left.
			size: 24, from: [2]float64{40, 40}, to: [2]float64{72, 40},
			want: deltaTally{added: 2, removed: 2, updated: 2},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := newMover(t, tc.from[0], tc.from[1], tc.size)
			if got := e.moveTo(tc.to[0], tc.to[1]); got != tc.want {
				t.Errorf("move reported %+v, want %+v", got, tc.want)
			}
		})
	}
}

// Moving an entity back and forth must report the same thing every time — a
// scratch buffer carrying entries between moves would show up here first.
func TestBucketDeltas_RepeatedMovesReportTheSameThing(t *testing.T) {
	e := newMover(t, 40, 40, 8)

	for i := range 4 {
		out := e.moveTo(72, 40)
		back := e.moveTo(40, 40)
		want := deltaTally{added: 1, removed: 1}
		if out != want {
			t.Fatalf("round %d: moving out reported %+v, want %+v", i, out, want)
		}
		if back != want {
			t.Fatalf("round %d: moving back reported %+v, want %+v", i, back, want)
		}
	}
}

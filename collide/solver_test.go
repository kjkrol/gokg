package collide_test

import (
	"testing"

	"github.com/kjkrol/gokg/collide"
	"github.com/kjkrol/gokg/geom"
	"github.com/kjkrol/gokg/plane"
)

const iterations = 16

func at(x, y float64) plane.AABB { return plane.NewAABB(geom.NewVec(x, y), 10, 10) }

func flat() plane.Space2D { return plane.NewEuclidean2D(1000, 1000) }

// depth is how far two boxes interpenetrate on their shallower axis, negative
// when they stand apart. Resting exactly edge to edge is the right answer, so
// tests measure depth rather than asking whether the closed boxes touch.
func depth(a, b plane.AABB) float64 {
	x := min(a.BottomRight.X, b.BottomRight.X) - max(a.TopLeft.X, b.TopLeft.X)
	y := min(a.BottomRight.Y, b.BottomRight.Y) - max(a.TopLeft.Y, b.TopLeft.Y)
	return min(x, y)
}

func TestSolve_SeparatesAPairAndSplitsTheTravel(t *testing.T) {
	a, b := at(100, 100), at(104, 100)

	var s collide.Solver
	s.Add(collide.Pair{A: &a, B: &b})
	s.Solve(flat(), iterations, nil)

	if d := depth(a, b); d > 1e-9 {
		t.Errorf("still overlapping by %v after the solve", d)
	}
	// They overlapped by 6 and neither is pinned, so each should have given 3.
	if got := 100 - a.TopLeft.X; got != 3 {
		t.Errorf("A moved %v, want 3 — half the overlap", got)
	}
	if got := b.TopLeft.X - 104; got != 3 {
		t.Errorf("B moved %v, want 3 — the other half", got)
	}
}

// An odd penetration cannot lose anything to the split: half of 3 is 1.5 on
// each side, and both sides must end up clear.
func TestSolve_OddPenetrationLosesNothingInTheSplit(t *testing.T) {
	a, b := at(100, 100), at(107, 100) // overlapping by 3

	var s collide.Solver
	s.Add(collide.Pair{A: &a, B: &b})
	s.Solve(flat(), iterations, nil)

	travelled := (100 - a.TopLeft.X) + (b.TopLeft.X - 107)
	if travelled != 3 {
		t.Errorf("the pair travelled %v apart in total, want the whole 3 units of overlap", travelled)
	}
	if d := depth(a, b); d > 1e-9 {
		t.Errorf("still overlapping by %v", d)
	}
}

func TestSolve_PinnedSideNeverMoves(t *testing.T) {
	for name, tc := range map[string]struct{ staticA, staticB bool }{
		"B pinned": {false, true},
		"A pinned": {true, false},
	} {
		t.Run(name, func(t *testing.T) {
			a, b := at(100, 100), at(104, 100)
			startA, startB := a, b

			var s collide.Solver
			s.Add(collide.Pair{A: &a, B: &b, StaticA: tc.staticA, StaticB: tc.staticB})
			s.Solve(flat(), iterations, nil)

			if d := depth(a, b); d > 1e-9 {
				t.Errorf("still overlapping by %v", d)
			}
			if tc.staticA && a != startA {
				t.Errorf("pinned A moved from %v to %v", startA.AABB, a.AABB)
			}
			if tc.staticB && b != startB {
				t.Errorf("pinned B moved from %v to %v", startB.AABB, b.AABB)
			}
		})
	}
}

func TestSolve_ReportsButNeverSeparates(t *testing.T) {
	for name, p := range map[string]collide.Pair{
		"a sensor pair":     {Sensor: true},
		"two pinned bodies": {StaticA: true, StaticB: true},
	} {
		t.Run(name, func(t *testing.T) {
			a, b := at(100, 100), at(104, 100)
			p.A, p.B = &a, &b
			startA, startB := a, b

			reported := 0
			var s collide.Solver
			s.Add(p)
			s.Solve(flat(), iterations, func(int, geom.Vec) { reported++ })

			if reported != 1 {
				t.Errorf("contact reported %d times, want exactly 1", reported)
			}
			if a != startA || b != startB {
				t.Error("boxes moved, though this pair must only be reported")
			}
		})
	}
}

func TestSolve_ReportsEachContactOnceEvenAcrossPasses(t *testing.T) {
	// Deeply stacked, so the solver needs several passes and every pass finds
	// these pairs still overlapping.
	boxes := []plane.AABB{at(100, 100), at(102, 100), at(104, 100)}

	counts := map[int]int{}
	var s collide.Solver
	s.Add(collide.Pair{A: &boxes[0], B: &boxes[1]})
	s.Add(collide.Pair{A: &boxes[1], B: &boxes[2]})
	s.Add(collide.Pair{A: &boxes[0], B: &boxes[2]})
	s.Solve(flat(), iterations, func(i int, _ geom.Vec) { counts[i]++ })

	for i := range 3 {
		if counts[i] != 1 {
			t.Errorf("pair %d reported %d times, want exactly 1", i, counts[i])
		}
	}
}

// Separating one pair drives a box into the next, so the second pass has work
// the first could not have seen. This is what more than one iteration buys.
func TestSolve_KeepsGoingWhileSeparationCreatesNewOverlap(t *testing.T) {
	boxes := []plane.AABB{at(100, 100), at(102, 100), at(104, 100)}

	var s collide.Solver
	s.Add(collide.Pair{A: &boxes[0], B: &boxes[1]})
	s.Add(collide.Pair{A: &boxes[1], B: &boxes[2]})
	s.Add(collide.Pair{A: &boxes[0], B: &boxes[2]})
	s.Solve(flat(), iterations, nil)

	for _, pair := range [][2]int{{0, 1}, {1, 2}, {0, 2}} {
		if d := depth(boxes[pair[0]], boxes[pair[1]]); d > 1e-6 {
			t.Errorf("boxes %d and %d still overlap by %v", pair[0], pair[1], d)
		}
	}
}

// A pair meeting across a seam has to separate the short way, through the
// seam, not by travelling back across the whole world.
func TestSolve_SeparatesAcrossAToroidalSeam(t *testing.T) {
	space := plane.NewToroidal2D(200, 200)
	a := space.WrapAABB(geom.NewAABB(geom.NewVec(196, 100), geom.NewVec(206, 110)))
	b := space.WrapAABB(geom.NewAABB(geom.NewVec(2, 100), geom.NewVec(12, 110)))

	var s collide.Solver
	s.Add(collide.Pair{A: &a, B: &b})
	s.Solve(space, iterations, nil)

	if got, ok := a.DeepestOverlapWith(&b); ok {
		t.Errorf("still overlapping across the seam by %v", got.Penetration)
	}
	// Two units of overlap, split evenly, and nobody crossed the world.
	if moved := 196 - a.TopLeft.X; moved != 2 {
		t.Errorf("A moved %v, want 2 — the short way out", moved)
	}
}

func TestVisitMoved_NamesOnlyThePairsThatWerePushed(t *testing.T) {
	a, b := at(100, 100), at(104, 100) // overlapping
	c, d := at(500, 500), at(600, 500) // nowhere near each other
	e, f := at(700, 700), at(704, 700) // overlapping, but e is pinned

	var s collide.Solver
	hot := s.Add(collide.Pair{A: &a, B: &b})
	cold := s.Add(collide.Pair{A: &c, B: &d})
	pinned := s.Add(collide.Pair{A: &e, B: &f, StaticA: true})
	s.Solve(flat(), iterations, nil)

	seen := map[int][2]bool{}
	s.VisitMoved(func(i int, movedA, movedB bool) { seen[i] = [2]bool{movedA, movedB} })

	if got, ok := seen[hot]; !ok || got != [2]bool{true, true} {
		t.Errorf("overlapping pair reported as %v (present=%v), want both sides moved", got, ok)
	}
	if _, ok := seen[cold]; ok {
		t.Error("a pair that never touched was reported as moved")
	}
	if got, ok := seen[pinned]; !ok || got != [2]bool{false, true} {
		t.Errorf("pinned pair reported as %v (present=%v), want only B moved", got, ok)
	}
}

func TestReset_KeepsTheBuffersAndDropsTheBatch(t *testing.T) {
	a, b := at(100, 100), at(104, 100)

	var s collide.Solver
	s.Add(collide.Pair{A: &a, B: &b})
	s.Reset()

	if s.Len() != 0 {
		t.Errorf("Len = %d after Reset, want 0", s.Len())
	}
	start := a
	s.Solve(flat(), iterations, func(int, geom.Vec) { t.Error("a reset Solver reported a contact") })
	if a != start {
		t.Error("a reset Solver moved a box from the dropped batch")
	}
}

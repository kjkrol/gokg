package collide

import (
	"github.com/kjkrol/gokg/geom"
	"github.com/kjkrol/gokg/plane"
)

// Pair is two boxes that may be in contact.
//
// A and B point at boxes the caller owns; the solver moves them in place. A
// side marked static is never moved, and the whole push goes to the other one;
// a pair with both sides static is reported but never separated. A Sensor pair
// is reported and never separated either — the difference is intent, not
// mechanics: a sensor is asking to be told, a static pair has nowhere to go.
type Pair struct {
	A, B             *plane.AABB
	StaticA, StaticB bool
	Sensor           bool
}

// state is what the solver remembers about a pair across iterations, kept
// beside Pair rather than in it so the caller's input stays untouched.
type state struct {
	reported       bool
	movedA, movedB bool
}

// Solver separates a batch of pairs, reusing its buffers between batches.
//
// The zero value is ready to use. A Solver kept across ticks stops allocating
// once its buffers have grown. One Solver serves one goroutine.
type Solver struct {
	pairs  []Pair
	states []state
}

// Reset empties the Solver for a new batch, keeping the memory.
func (s *Solver) Reset() {
	s.pairs = s.pairs[:0]
	s.states = s.states[:0]
}

// Add enters a pair into the batch and returns its index, which is how Solve
// and VisitMoved name it afterwards. Indices follow Add order.
func (s *Solver) Add(p Pair) int {
	s.pairs = append(s.pairs, p)
	s.states = append(s.states, state{})
	return len(s.pairs) - 1
}

// Len reports how many pairs are in the current batch.
func (s *Solver) Len() int { return len(s.pairs) }

// Solve pushes every overlapping pair apart, up to iterations passes over the
// whole batch, moving boxes under surface's boundary rules.
//
// onContact fires at most once per pair, the first pass that finds it
// overlapping, with the penetration measured at that moment — later passes
// re-measure to separate but say nothing. It may be nil.
//
// More than one pass is needed because separating one pair can drive a box
// into another; a pass that separates nothing ends the walk, since every
// later pass would read the same unchanged geometry and reach the same
// verdict.
func (s *Solver) Solve(surface plane.Space2D, iterations int, onContact func(i int, pen geom.Vec)) {
	for range iterations {
		moved := false
		for i := range s.pairs {
			p := &s.pairs[i]
			st := &s.states[i]

			hit, ok := p.A.DeepestOverlapWith(p.B)
			if !ok {
				continue
			}

			if !st.reported {
				st.reported = true
				if onContact != nil {
					onContact(i, hit.Penetration)
				}
			}

			if p.Sensor {
				continue
			}
			pushA, pushB, ok := split(hit.Penetration, p.StaticA, p.StaticB)
			if !ok {
				continue
			}
			if pushA != (geom.Vec{}) {
				surface.Translate(p.A, pushA)
				st.movedA = true
				moved = true
			}
			if pushB != (geom.Vec{}) {
				surface.Translate(p.B, pushB)
				st.movedB = true
				moved = true
			}
		}
		if !moved {
			return
		}
	}
}

// VisitMoved calls fn for every pair the solver actually pushed, naming which
// sides moved. A caller keeping its boxes in a spatial index needs this: the
// index is worth telling once, about where a box came to rest, not once per
// pass about where it was passing through.
func (s *Solver) VisitMoved(fn func(i int, movedA, movedB bool)) {
	for i := range s.states {
		st := &s.states[i]
		if st.movedA || st.movedB {
			fn(i, st.movedA, st.movedB)
		}
	}
}

// split shares a penetration between the two sides: half each when both can
// move, all of it to whichever one can when the other is pinned, and nothing
// at all when neither can.
func split(pen geom.Vec, staticA, staticB bool) (pushA, pushB geom.Vec, ok bool) {
	switch {
	case staticA && staticB:
		return geom.Vec{}, geom.Vec{}, false
	case staticB:
		return pen, geom.Vec{}, true
	case staticA:
		return geom.Vec{}, negate(pen), true
	default:
		half := geom.NewVec(pen.X/2, pen.Y/2)
		return half, negate(half), true
	}
}

func negate(v geom.Vec) geom.Vec { return geom.NewVec(-v.X, -v.Y) }

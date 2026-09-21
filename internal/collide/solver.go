package collide

import (
	"github.com/kjkrol/aabbworld/geom"
	iplane "github.com/kjkrol/aabbworld/internal/plane"
	"github.com/kjkrol/aabbworld/plane"
)

// Pair is collide.Pair as the solver holds it — the two have to stay field for field the same.
type Pair struct {
	A, B             *plane.AABB
	KeyA, KeyB       uint32
	StaticA, StaticB bool
	Sensor           bool
}

// state is what the solver remembers about a pair across iterations, kept
// beside Pair rather than in it so the caller's input stays untouched.
type state struct {
	reported       bool
	movedA, movedB bool
	testedAt       uint32 // the solver's clock when this pair was last measured; 0 for never
}

// Solver is the state and the algorithm behind collide.NarrowPhase.
type Solver struct {
	pairs  []Pair
	states []state

	// clock counts measurements; movedAt is its reading when each box key was last pushed.
	clock   uint32
	movedAt []uint32
}

// Reset empties the Solver for a new batch, keeping the memory.
func (s *Solver) Reset() {
	s.pairs = s.pairs[:0]
	s.states = s.states[:0]
	s.clock = 0
	clear(s.movedAt)
}

// Add enters a pair into the batch and returns its index.
func (s *Solver) Add(p Pair) int {
	s.pairs = append(s.pairs, p)
	s.states = append(s.states, state{})
	if need := int(max(p.KeyA, p.KeyB)) + 1; need > len(s.movedAt) {
		s.movedAt = append(s.movedAt, make([]uint32, need-len(s.movedAt))...)
	}
	return len(s.pairs) - 1
}

// Solve reports each overlapping pair once and pushes it apart, in up to iterations passes.
func (s *Solver) Solve(surface *iplane.Surface, iterations int, onContact func(i int, pen geom.Vec)) {
	for range iterations {
		moved := false
		for i := range s.pairs {
			p := &s.pairs[i]
			st := &s.states[i]

			if st.testedAt > s.movedAt[p.KeyA] && st.testedAt > s.movedAt[p.KeyB] {
				continue
			}
			s.clock++
			st.testedAt = s.clock

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
				s.movedAt[p.KeyA] = s.clock
				st.movedA = true
				moved = true
			}
			if pushB != (geom.Vec{}) {
				surface.Translate(p.B, pushB)
				s.movedAt[p.KeyB] = s.clock
				st.movedB = true
				moved = true
			}
		}
		if !moved {
			return
		}
	}
}

// Pair is the i-th pair of the batch.
func (s *Solver) Pair(i int) *Pair { return &s.pairs[i] }

// VisitMoved calls fn for every pair the solver actually pushed, naming which sides moved.
func (s *Solver) VisitMoved(fn func(i int, movedA, movedB bool)) {
	for i := range s.states {
		st := &s.states[i]
		if st.movedA || st.movedB {
			fn(i, st.movedA, st.movedB)
		}
	}
}

// split shares a penetration between the sides that can move; ok is false when neither can.
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

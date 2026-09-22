package collide

import (
	"github.com/kjkrol/aabbworld/geom"
	iplane "github.com/kjkrol/aabbworld/internal/plane"
	"github.com/kjkrol/aabbworld/plane"
	"github.com/kjkrol/uid"
)

// Pair is two entities' boxes that may be in contact, with the flags the solver separates them by.
type Pair struct {
	A, B     *plane.AABB
	IDA, IDB uid.UID64
	Flags    uint8
}

const (
	StaticA uint8 = 1 << iota
	StaticB
	Sensor
	dropped
	reported
	movedA
	movedB
)

// state is what the solver remembers about a pair across passes.
type state struct {
	testedAt uint32 // the solver's clock when this pair was last measured; 0 for never
	flags    uint8
}

// Touch is asked once per pair, when its boxes first overlap; false drops the pair for the tick.
type Touch func(i int, pen geom.Vec) (geom.Vec, bool)

// Solver is the state and the algorithm behind collide.Engine.
type Solver struct {
	pairs  []Pair
	states []state

	// clock counts measurements; movedAt is its reading when each box was last pushed, by id index.
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
	s.states = append(s.states, state{flags: p.Flags})
	if need := int(max(p.IDA.Index(), p.IDB.Index())) + 1; need > len(s.movedAt) {
		s.movedAt = append(s.movedAt, make([]uint32, need-len(s.movedAt))...)
	}
	return len(s.pairs) - 1
}

// Len is how many pairs the batch holds.
func (s *Solver) Len() int { return len(s.pairs) }

// Pair is the i-th pair of the batch.
func (s *Solver) Pair(i int) *Pair { return &s.pairs[i] }

// Solve reports each overlapping pair once and pushes it apart, in up to iterations passes.
func (s *Solver) Solve(surface *iplane.Surface, iterations int, touch Touch, onContact func(i int, pen geom.Vec)) {
	for range iterations {
		moved := false
		for i := range s.pairs {
			p := &s.pairs[i]
			st := &s.states[i]
			if st.flags&dropped != 0 {
				continue
			}
			keyA, keyB := p.IDA.Index(), p.IDB.Index()

			if st.testedAt > s.movedAt[keyA] && st.testedAt > s.movedAt[keyB] {
				continue
			}
			s.clock++
			st.testedAt = s.clock

			hit, ok := p.A.DeepestOverlapWith(p.B)
			if !ok {
				continue
			}
			pen := hit.Penetration

			if st.flags&reported == 0 {
				if touch != nil {
					if pen, ok = touch(i, pen); !ok {
						st.flags |= dropped
						continue
					}
				}
				st.flags |= reported
				if onContact != nil {
					onContact(i, pen)
				}
			}

			if st.flags&Sensor != 0 {
				continue
			}
			pushA, pushB, ok := split(pen, st.flags&StaticA != 0, st.flags&StaticB != 0)
			if !ok {
				continue
			}
			if pushA != (geom.Vec{}) {
				surface.Translate(p.A, pushA)
				s.movedAt[keyA] = s.clock
				st.flags |= movedA
				moved = true
			}
			if pushB != (geom.Vec{}) {
				surface.Translate(p.B, pushB)
				s.movedAt[keyB] = s.clock
				st.flags |= movedB
				moved = true
			}
		}
		if !moved {
			return
		}
	}
}

// VisitMoved calls fn for every box the solver pushed, once per pair it was pushed in.
func (s *Solver) VisitMoved(fn func(id uid.UID64, box *plane.AABB)) {
	for i := range s.states {
		f := s.states[i].flags
		if f&movedA != 0 {
			fn(s.pairs[i].IDA, s.pairs[i].A)
		}
		if f&movedB != 0 {
			fn(s.pairs[i].IDB, s.pairs[i].B)
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

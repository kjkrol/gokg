package collide

import (
	"math/rand/v2"
	"testing"

	"github.com/kjkrol/gokg/geom"
	"github.com/kjkrol/gokg/plane"
)

// solveEverything is Solve as it was before pairs could be skipped: every pass
// measures every pair. It is what the skipping has to agree with, to the bit.
func solveEverything(pairs []Pair, surface plane.Space2D, iterations int, onContact func(i int)) (movedA, movedB []bool) {
	reported := make([]bool, len(pairs))
	movedA, movedB = make([]bool, len(pairs)), make([]bool, len(pairs))
	for range iterations {
		moved := false
		for i := range pairs {
			p := &pairs[i]
			hit, ok := p.A.DeepestOverlapWith(p.B)
			if !ok {
				continue
			}
			if !reported[i] {
				reported[i] = true
				onContact(i)
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
				movedA[i], moved = true, true
			}
			if pushB != (geom.Vec{}) {
				surface.Translate(p.B, pushB)
				movedB[i], moved = true, true
			}
		}
		if !moved {
			return
		}
	}
	return
}

// crowd is a random batch dense enough that pushes chain into one another:
// boxes scattered over a small field, every pair close enough to be a candidate
// entered, some sides pinned and some pairs only sensed.
type crowd struct {
	boxes []plane.AABB
	pairs []Pair
}

func newCrowd(rng *rand.Rand, surface plane.Space2D, count int, field float64) crowd {
	c := crowd{boxes: make([]plane.AABB, count)}
	pinned := make([]bool, count)
	for i := range c.boxes {
		c.boxes[i] = surface.WrapAABB(plane.NewAABB(geom.NewVec(rng.Float64()*field, rng.Float64()*field), 10, 10).AABB)
		pinned[i] = rng.IntN(10) == 0
	}
	for i := range c.boxes {
		for j := i + 1; j < count; j++ {
			dx := c.boxes[i].TopLeft.X - c.boxes[j].TopLeft.X
			dy := c.boxes[i].TopLeft.Y - c.boxes[j].TopLeft.Y
			if dx*dx+dy*dy > 30*30 {
				continue
			}
			c.pairs = append(c.pairs, Pair{
				KeyA: uint32(i), KeyB: uint32(j),
				StaticA: pinned[i], StaticB: pinned[j],
				Sensor: rng.IntN(12) == 0,
			})
		}
	}
	return c
}

// over is the crowd's pairs pointed at boxes — a copy of its own, so the same
// crowd can be solved more than once from the same start.
func (c crowd) over(boxes []plane.AABB, keyed bool) []Pair {
	copy(boxes, c.boxes)
	pairs := make([]Pair, len(c.pairs))
	for i, p := range c.pairs {
		p.A, p.B = &boxes[p.KeyA], &boxes[p.KeyB]
		if !keyed {
			p.KeyA, p.KeyB = 0, 0
		}
		pairs[i] = p
	}
	return pairs
}

func TestSolve_SkippingChangesNothing(t *testing.T) {
	surfaces := map[string]plane.Space2D{
		"bounded":  plane.NewEuclidean2D(400, 400),
		"toroidal": plane.NewToroidal2D(200, 200),
	}
	for name, surface := range surfaces {
		for _, keyed := range []bool{true, false} {
			rng := rand.New(rand.NewPCG(7, 11))
			for round := range 40 {
				c := newCrowd(rng, surface, 120, 200)

				want := make([]plane.AABB, len(c.boxes))
				var wantContacts []int
				wantA, wantB := solveEverything(c.over(want, false), surface, 16, func(i int) { wantContacts = append(wantContacts, i) })

				got := make([]plane.AABB, len(c.boxes))
				var gotContacts []int
				var s Solver
				for _, p := range c.over(got, keyed) {
					s.Add(p)
				}
				s.Solve(surface, 16, func(i int, _ geom.Vec) { gotContacts = append(gotContacts, i) })

				for i := range want {
					if got[i].AABB != want[i].AABB {
						t.Fatalf("%s keyed=%v round %d: box %d came to rest at %v, want %v", name, keyed, round, i, got[i].AABB, want[i].AABB)
					}
				}
				if len(gotContacts) != len(wantContacts) {
					t.Fatalf("%s keyed=%v round %d: %d contacts reported, want %d", name, keyed, round, len(gotContacts), len(wantContacts))
				}
				for k := range wantContacts {
					if gotContacts[k] != wantContacts[k] {
						t.Fatalf("%s keyed=%v round %d: contact %d was pair %d, want pair %d", name, keyed, round, k, gotContacts[k], wantContacts[k])
					}
				}
				s.VisitMoved(func(i int, a, b bool) { wantA[i], wantB[i] = wantA[i] != a, wantB[i] != b })
				for i := range wantA {
					if wantA[i] || wantB[i] {
						t.Fatalf("%s keyed=%v round %d: pair %d disagrees on which side moved", name, keyed, round, i)
					}
				}
			}
		}
	}
}

// The point of naming boxes: a candidate that stands clear of everything that
// moved is measured once, not once per pass.
func TestSolve_KeyedPairsAreMeasuredLessOften(t *testing.T) {
	surface := plane.NewEuclidean2D(400, 400)
	c := newCrowd(rand.New(rand.NewPCG(3, 5)), surface, 120, 200)
	boxes := make([]plane.AABB, len(c.boxes))

	measured := func(keyed bool) uint32 {
		var s Solver
		for _, p := range c.over(boxes, keyed) {
			s.Add(p)
		}
		s.Solve(surface, 16, nil)
		return s.clock
	}

	keyed, unkeyed := measured(true), measured(false)
	if keyed >= unkeyed {
		t.Errorf("measured %d pairs with keys and %d without — want naming the boxes to save measurements", keyed, unkeyed)
	}
}

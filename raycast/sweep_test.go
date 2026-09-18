package raycast_test

import (
	"math"
	"math/rand/v2"
	"testing"

	"github.com/kjkrol/gokg/geom"
	"github.com/kjkrol/gokg/raycast"
	"github.com/kjkrol/uid"
)

// bruteForce finds what a dense fan of rays meets, testing every box at every
// angle. It knows nothing of active sets or cone rejection, so it is an
// independent check on both.
//
// Only boxes the fan reaches across an arc wider than minArc are returned: the
// sweep resolves down to one unit at the far edge of the view, so a sliver
// thinner than that is below what it promises to find.
func bruteForce(boxes map[uid.UID64][4]float64, origin geom.Vec, cone raycast.Cone, steps int, minArc float64) map[uid.UID64]bool {
	dir := math.Atan2(cone.Direction.Y, cone.Direction.X)
	step := 2 * cone.HalfAngle / float64(steps)
	arcs := map[uid.UID64]float64{}
	for i := 0; i <= steps; i++ {
		a := dir - cone.HalfAngle + 2*cone.HalfAngle*float64(i)/float64(steps)
		dx, dy := math.Cos(a), math.Sin(a)

		best, hit := cone.Radius, uid.UID64(0)
		for id, b := range boxes {
			near, far := math.Inf(-1), math.Inf(1)
			for axis, ray := range [2][2]float64{{origin.X, dx}, {origin.Y, dy}} {
				o, d := ray[0], ray[1]
				lo, hi := b[axis], b[axis+2]
				if d == 0 {
					if o < lo || o > hi {
						near, far = 1, -1
					}
					continue
				}
				t1, t2 := (lo-o)/d, (hi-o)/d
				if t1 > t2 {
					t1, t2 = t2, t1
				}
				near, far = math.Max(near, t1), math.Min(far, t2)
			}
			if far >= near && far >= 0 && near <= best {
				best, hit = math.Max(near, 0), id
			}
		}
		if hit != 0 {
			arcs[hit] += step
		}
	}

	seen := map[uid.UID64]bool{}
	for id, a := range arcs {
		if a > minArc {
			seen[id] = true
		}
	}
	return seen
}

// The check runs one way: sampling can miss a sliver the exact sweep catches,
// but nothing the fan reaches across a resolvable arc may be missing from
// Visible.
func TestVisible_MissesNothingADenseFanOfRaysReaches(t *testing.T) {
	r := rand.New(rand.NewPCG(7, 11))
	for trial := range 20 {
		s := newFake(4000, 4000, false)
		s.put(eye, 2000, 2000, 10, 10)
		// One box per grid cell, so none overlaps another or swallows the
		// observer — both are preconditions of the angular arithmetic.
		boxes := map[uid.UID64][4]float64{}
		id := uid.UID64(100)
		for cx := range 8 {
			for cy := range 8 {
				if cx == 4 && cy == 4 {
					continue // the observer's own cell
				}
				x := float64(1600 + cx*100 + r.IntN(20))
				y := float64(1600 + cy*100 + r.IntN(20))
				w, h := float64(20+r.IntN(50)), float64(20+r.IntN(50))
				s.put(id, x, y, w, h)
				boxes[id] = [4]float64{float64(x), float64(y), float64(x + w), float64(y + h)}
				id++
			}
		}

		cone := raycast.Cone{
			Direction: geom.NewVec(math.Cos(float64(trial)), math.Sin(float64(trial))),
			HalfAngle: math.Pi / 3,
			Radius:    700,
		}
		origin := geom.NewVec(2005.0, 2005.0)

		got := map[uid.UID64]bool{}
		scan(t, s, eye, cone).Entities(func(id uid.UID64, _ float64) { got[id] = true })

		minArc := 2 * math.Atan2(1, cone.Radius)
		for id := range bruteForce(boxes, origin, cone, 4000, minArc) {
			if !got[id] {
				t.Errorf("trial %d: rays reach %d, but Visible does not report it", trial, id)
			}
		}
	}
}

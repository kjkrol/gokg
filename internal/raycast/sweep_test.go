package raycast_test

import (
	"math"
	"math/rand/v2"
	"testing"

	"github.com/kjkrol/aabbworld/geom"
	"github.com/kjkrol/aabbworld/internal/raycast"
	"github.com/kjkrol/uid"
)

// bruteForce is what a dense fan of rays meets across an arc wider than minArc.
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

func TestVisible_MissesNothingADenseFanOfRaysReaches(t *testing.T) {
	r := rand.New(rand.NewPCG(7, 11))
	for trial := range 20 {
		s := newFake(4000, 4000, false)
		s.put(eye, 2000, 2000, 10, 10)
		boxes := map[uid.UID64][4]float64{}
		id := uid.UID64(100)
		for cx := range 8 {
			for cy := range 8 {
				if cx == 4 && cy == 4 {
					continue
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

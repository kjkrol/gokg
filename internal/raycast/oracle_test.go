package raycast_test

import (
	"math"
	"sort"

	"github.com/kjkrol/aabbworld/geom"
	"github.com/kjkrol/uid"
)

// slabHit is the textbook ray-box test the oracles use: where a unit ray enters and leaves b.
func slabHit(b [4]float64, origin geom.Vec, dx, dy float64) (near, far float64, ok bool) {
	near, far = math.Inf(-1), math.Inf(1)
	for axis, ray := range [2][2]float64{{origin.X, dx}, {origin.Y, dy}} {
		o, d := ray[0], ray[1]
		lo, hi := b[axis], b[axis+2]
		if d == 0 {
			if o < lo || o > hi {
				return 0, 0, false
			}
			continue
		}
		t1, t2 := (lo-o)/d, (hi-o)/d
		if t1 > t2 {
			t1, t2 = t2, t1
		}
		near, far = math.Max(near, t1), math.Min(far, t2)
	}
	if far < near || far < 0 {
		return 0, 0, false
	}
	return math.Max(near, 0), far, true
}

// castRay is the sight model spelled out plainly over a map of boxes; taus nil makes every box
// block.
func castRay(boxes map[uid.UID64][4]float64, taus map[uid.UID64]float64, origin geom.Vec, angle, radius float64) (dist float64, id uid.UID64, hit bool) {
	dx, dy := math.Cos(angle), math.Sin(angle)
	type span struct{ near, far, tau float64 }
	var seeThrough []span
	wall, wallID, walled := radius, uid.UID64(0), false
	for bid, b := range boxes {
		near, far, ok := slabHit(b, origin, dx, dy)
		if !ok {
			continue
		}
		if tau := math.Min(taus[bid], 1); tau > 0 {
			seeThrough = append(seeThrough, span{near, far, tau})
			continue
		}
		if near <= wall {
			wall, wallID, walled = near, bid, true
		}
	}
	sort.Slice(seeThrough, func(i, j int) bool { return seeThrough[i].near < seeThrough[j].near })

	budget, pos := radius, 0.0
	for _, s := range seeThrough {
		if s.far <= pos {
			continue
		}
		start := math.Max(s.near, pos)
		if start-pos >= budget {
			break
		}
		budget -= start - pos
		pos = start
		if cost := (s.far - start) / s.tau; cost >= budget {
			pos += budget * s.tau
			budget = 0
			break
		} else {
			budget -= cost
			pos = s.far
		}
	}
	reach := pos + budget
	if walled && wall <= reach {
		return wall, wallID, true
	}
	return reach, 0, false
}

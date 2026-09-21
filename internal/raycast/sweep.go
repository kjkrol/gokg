package raycast

import (
	"math"
	"slices"

	"github.com/kjkrol/aabbworld/geom"
	"github.com/kjkrol/uid"
)

// sample is the reach of the view at one angle: how far it gets before meeting
// something, and what that something was.
type sample struct {
	angle float64 // relative to the cone's direction
	dist  float64
	id    uid.UID64
	hit   bool
}

// event marks a candidate entering or leaving the set of boxes whose angular
// span covers the sweep angle.
type event struct {
	angle float64
	idx   int
	enter bool
}

// sweep records what the view meets at every angle where the reach can change shape.
func sweep(origin geom.Vec, coneDir, halfAngle, radius float64, cands []candidate, sc *scratch) []sample {
	eps := math.Atan2(1, radius)

	angles := sc.angles
	add := func(a float64) {
		if a >= -halfAngle && a <= halfAngle {
			angles = append(angles, a)
		}
	}
	add(-halfAngle)
	add(halfAngle)

	events := sc.events
	for i, c := range cands {
		add(c.span.lo - eps)
		add(c.span.lo + eps)
		add(c.span.hi - eps)
		add(c.span.hi + eps)
		events = append(events,
			event{angle: c.span.lo - eps, idx: i, enter: true},
			event{angle: c.span.hi + eps, idx: i})
	}
	slices.Sort(angles)
	slices.SortFunc(events, func(a, b event) int {
		switch {
		case a.angle < b.angle:
			return -1
		case a.angle > b.angle:
			return 1
		default:
			return 0
		}
	})

	active, out := walk(origin, coneDir, radius, eps/8, angles, events, cands, sc.active, sc.samples)
	sc.angles, sc.events, sc.active, sc.samples = angles, events, active, out
	return out
}

// walk casts along ascending angles, one sample each, skipping any within minGap of the last.
func walk(
	origin geom.Vec,
	coneDir, radius, minGap float64,
	angles []float64,
	events []event,
	cands []candidate,
	active []int,
	out []sample,
) ([]int, []sample) {
	active = active[:0]
	next := 0

	for i, a := range angles {
		if i > 0 && a-angles[i-1] < minGap {
			continue
		}
		for next < len(events) && events[next].angle <= a {
			e := events[next]
			if e.enter {
				active = append(active, e.idx)
			} else if at := slices.Index(active, e.idx); at >= 0 {
				active = slices.Delete(active, at, at+1)
			}
			next++
		}
		out = append(out, castAt(origin, coneDir, a, radius, cands, active))
	}
	return active, out
}

// castAt finds the nearest active candidate along one angle, or the range limit.
func castAt(origin geom.Vec, coneDir, rel, radius float64, cands []candidate, active []int) sample {
	abs := coneDir + rel
	dir := geom.NewVec(math.Cos(abs), math.Sin(abs))

	best := sample{angle: rel, dist: radius}
	for _, i := range active {
		c := cands[i]
		d, ok := hitDistance(origin, dir, c.box)
		if !ok || d > best.dist {
			continue
		}
		best.dist, best.id, best.hit = d, c.id, true
	}
	return best
}

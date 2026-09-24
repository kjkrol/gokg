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

// arcStep is the angle between samples where the reach is curved: a free arc, a see-through box.
const arcStep = math.Pi / 90 // two degrees

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
		// The reach behind a see-through box is curved, so its span is sampled at the arc step too.
		if c.tau > 0 {
			for a := c.span.lo + arcStep; a < c.span.hi; a += arcStep {
				add(a)
			}
		}
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

	active, out := walk(origin, coneDir, radius, eps/8, angles, events, cands, sc.active, &sc.crossings, sc.samples)
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
	cross *[]crossing,
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
		out = append(out, castAt(origin, coneDir, a, radius, cands, active, cross))
	}
	return active, out
}

// crossing is one see-through box on a ray: entry, exit and the budget each unit inside costs.
type crossing struct {
	near, far float64
	rate      float64 // 1/tau
}

// castAt follows one angle: the nearest blocking candidate is the wall, see-through ones eat the
// budget. The sample is a hit only when the wall is reached before the budget runs out.
func castAt(origin geom.Vec, coneDir, rel, radius float64, cands []candidate, active []int, cross *[]crossing) sample {
	abs := coneDir + rel
	dir := geom.NewVec(math.Cos(abs), math.Sin(abs))

	best := sample{angle: rel, dist: radius}
	crossed := 0
	for _, i := range active {
		c := cands[i]
		near, far, ok := hitDistance(origin, dir, c.box)
		if !ok {
			continue
		}
		if c.tau > 0 {
			if crossed == 0 {
				*cross = (*cross)[:0]
			}
			*cross = insertCrossing(*cross, crossing{near: near, far: far, rate: 1 / c.tau})
			crossed++
			continue
		}
		if near > best.dist {
			continue
		}
		best.dist, best.id, best.hit = near, c.id, true
	}
	if crossed == 0 {
		return best
	}

	// An empty stretch costs its length, one inside a box length·rate; overlaps are charged once.
	budget, pos := radius, 0.0
	for _, x := range *cross {
		if x.far <= pos || x.near >= best.dist {
			continue
		}
		start := math.Max(x.near, pos)
		if gap := start - pos; gap >= budget {
			break
		}
		budget, pos = budget-(start-pos), start
		cost := (x.far - start) * x.rate
		if cost >= budget {
			budget, pos = 0, start+budget/x.rate
			break
		}
		budget, pos = budget-cost, x.far
	}
	reach := pos + budget
	if !best.hit || best.dist > reach {
		best.dist, best.id, best.hit = reach, 0, false
	}
	return best
}

// insertCrossing keeps cross sorted by where each box is entered; the list is a handful long.
func insertCrossing(cross []crossing, x crossing) []crossing {
	cross = append(cross, x)
	for i := len(cross) - 1; i > 0 && cross[i-1].near > x.near; i-- {
		cross[i] = cross[i-1]
		cross[i-1] = x
	}
	return cross
}

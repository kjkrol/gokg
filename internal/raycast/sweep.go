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

// arcStep is the angle between samples where the reach is curved: a free arc, a see-through box,
// any cone with heights.
const arcStep = math.Pi / 90 // two degrees

// sweep records what the view meets at every angle where the reach can change shape.
func sweep(origin geom.Vec, coneDir, halfAngle, radius float64, cands []candidate, elev *elevation, sc *scratch) []sample {
	eps := math.Atan2(1, radius)

	angles := sc.angles
	add := func(a float64) {
		if a >= -halfAngle && a <= halfAngle {
			angles = append(angles, a)
		}
	}
	add(-halfAngle)
	add(halfAngle)
	// Over uneven ground the reach changes everywhere, so the whole cone is sampled at the arc step.
	everywhere := elev != nil && elev.ground != nil
	if everywhere {
		for a := -halfAngle + arcStep; a < halfAngle; a += arcStep {
			add(a)
		}
	}

	events := sc.events
	for i, c := range cands {
		add(c.span.lo - eps)
		add(c.span.lo + eps)
		add(c.span.hi - eps)
		add(c.span.hi + eps)
		// The reach behind a see-through box, or any box with a height, is curved, so its span is
		// sampled at the arc step too.
		if !everywhere && (c.tau > 0 || elev != nil) {
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

	active, out := walk(origin, coneDir, radius, eps/8, angles, events, cands, sc.active, elev, sc, true, sc.samples)
	sc.angles, sc.events, sc.active, sc.samples = angles, events, active, out
	return out
}

// walk casts along ascending angles, one sample each, skipping any within minGap of the last;
// with mark the casts note on each candidate where they saw it.
func walk(
	origin geom.Vec,
	coneDir, radius, minGap float64,
	angles []float64,
	events []event,
	cands []candidate,
	active []int,
	elev *elevation,
	sc *scratch,
	mark bool,
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
		if elev != nil {
			out = append(out, castElevated(origin, coneDir, a, radius, cands, active, elev, sc, mark))
		} else {
			out = append(out, castFlat(origin, coneDir, a, radius, cands, active, &sc.crossings, mark))
		}
	}
	return active, out
}

// crossing is one box on a ray: entry, exit, the budget each unit inside costs (0 for a box that
// blocks) and which candidate it is.
type crossing struct {
	near, far float64
	rate      float64 // 1/tau, 0 when opaque
	idx       int
}

// castFlat follows one angle on a plane: the nearest blocking candidate is the wall, see-through
// ones eat the budget. The sample is a hit only when the wall is reached before the budget runs out.
func castFlat(origin geom.Vec, coneDir, rel, radius float64, cands []candidate, active []int, cross *[]crossing, mark bool) sample {
	abs := coneDir + rel
	dir := geom.NewVec(math.Cos(abs), math.Sin(abs))

	best, wall := sample{angle: rel, dist: radius}, -1
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
			*cross = insertCrossing(*cross, crossing{near: near, far: far, rate: 1 / c.tau, idx: i})
			crossed++
			continue
		}
		if near > best.dist {
			continue
		}
		best.dist, best.id, best.hit, wall = near, c.id, true, i
	}
	if crossed > 0 {
		// An empty stretch costs its length, one inside a box length·rate; overlaps are charged once.
		if reach := spend(*cross, best.dist, radius); !best.hit || best.dist > reach {
			best.dist, best.id, best.hit, wall = reach, 0, false, -1
		}
		for _, x := range *cross {
			if mark && x.near <= best.dist {
				cands[x.idx].see(x.near)
			}
		}
	}
	if mark && wall >= 0 {
		cands[wall].see(best.dist)
	}
	return best
}

// spend walks the see-through stretches sorted by entry, charging each once, and returns where the
// budget runs out; nothing at or past limit is charged.
func spend(cross []crossing, limit, budget float64) float64 {
	pos := 0.0
	for _, x := range cross {
		if x.far <= pos || x.near >= limit {
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
	return pos + budget
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

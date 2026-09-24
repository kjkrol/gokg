package raycast

import (
	"math"

	"github.com/kjkrol/aabbworld/geom"
)

// elevation is what a cast needs beyond the plane when the cone has heights: the eye, the ground
// and how far apart the ground is sampled along a ray.
type elevation struct {
	eye    float64
	ground func(geom.Vec) float64 // nil: flat at 0
	step   float64
	w, h   float64 // wrap sizes, 0 where the axis does not wrap
}

// at is the ground height under the point dist along the ray, folded back onto a wrapping world.
func (e *elevation) at(origin, dir geom.Vec, dist float64) float64 {
	if e.ground == nil {
		return 0
	}
	p := geom.NewVec(origin.X+dir.X*dist, origin.Y+dir.Y*dist)
	if e.w > 0 {
		if p.X = math.Mod(p.X, e.w); p.X < 0 {
			p.X += e.w
		}
	}
	if e.h > 0 {
		if p.Y = math.Mod(p.Y, e.h); p.Y < 0 {
			p.Y += e.h
		}
	}
	return e.ground(p)
}

// ray is one elevated cast in progress: the crossings entered so far are its blockers, horizon the
// steepest sightline to the ground yet, reach the farthest lit ground and cut the wall it ends at.
type ray struct {
	origin, dir geom.Vec
	radius      float64
	cands       []candidate
	cross       []crossing
	e           *elevation
	sc          *scratch
	mark        bool
	horizon     float64
	reach       float64
	cut         int
	walls       int // blocking crossings entered so far
	veils       int // see-through crossings entered so far
}

// castElevated follows one angle with heights: an entry is seen when the sightline to its top
// clears the ground and every nearer blocking box within the budget; the reach is the farthest
// lit ground, a standing box's foot included.
func castElevated(origin geom.Vec, coneDir, rel, radius float64, cands []candidate, active []int, e *elevation, sc *scratch, mark bool) sample {
	abs := coneDir + rel
	dir := geom.NewVec(math.Cos(abs), math.Sin(abs))

	cross := sc.crossings[:0]
	for _, i := range active {
		near, far, ok := hitDistance(origin, dir, cands[i].box)
		if !ok {
			continue
		}
		rate := 0.0
		if cands[i].tau > 0 {
			rate = 1 / cands[i].tau
		}
		cross = insertCrossing(cross, crossing{near: near, far: far, rate: rate, idx: i})
	}
	sc.crossings = cross

	r := ray{origin: origin, dir: dir, radius: radius, cands: cands, cross: cross, e: e, sc: sc, mark: mark, horizon: math.Inf(-1), cut: -1}
	j := 0
	for d := e.step; ; d += e.step {
		if d > radius {
			d = radius
		}
		for j < len(cross) && cross[j].near <= d {
			r.enter(j)
			if cross[j].rate == 0 {
				r.walls++
			} else {
				r.veils++
			}
			j++
		}
		r.ground(d, e.at(origin, dir, d), j, -1)
		if d >= radius {
			break
		}
	}
	out := sample{angle: rel, dist: r.reach}
	if r.cut >= 0 {
		out.id, out.hit = cands[r.cut].id, true
	}
	return out
}

// enter looks at crossing j as a target and, when it stands on the ground, lights its foot.
func (r *ray) enter(j int) {
	x := r.cross[j]
	c := &r.cands[x.idx]
	if x.near == 0 {
		r.see(c, 0) // the eye is inside the box: seen from where it stands
		return
	}
	if tan := (c.top - r.e.eye) / x.near; tan >= r.horizon && !r.blocked(tan, x.near, j) && r.spend(tan, x.near, j) >= x.near {
		r.see(c, x.near)
	}
	if !c.standing {
		return
	}
	cut := -1
	if x.rate == 0 {
		cut = x.idx
	}
	r.ground(x.near, c.foot, j, cut)
}

// ground takes the ground point at d as a target: lit when the sightline to it clears everything
// nearer, it becomes the reach, cut by candidate cut if any; hidden ground is passed over.
func (r *ray) ground(d, alt float64, j, cut int) {
	tan := (alt - r.e.eye) / d
	if tan < r.horizon {
		return
	}
	r.horizon = tan
	if r.blocked(tan, d, j) {
		return
	}
	if reach := r.spend(tan, d, j); reach < d {
		r.lit(reach, -1)
	} else {
		r.lit(d, cut)
	}
}

func (r *ray) see(c *candidate, dist float64) {
	if r.mark {
		c.see(dist)
	}
}

func (r *ray) lit(d float64, cut int) {
	if d >= r.reach {
		r.reach, r.cut = d, cut
	}
}

// blocked reports whether the sightline with tangent tan, followed to upto, passes through the
// band of a blocking box among the first j crossings.
func (r *ray) blocked(tan, upto float64, j int) bool {
	if r.walls == 0 {
		return false
	}
	for _, x := range r.cross[:j] {
		if x.rate != 0 {
			continue
		}
		c := &r.cands[x.idx]
		lo, hi := r.heights(tan, x.near, min(x.far, upto))
		if hi > c.bottom && lo < c.top {
			return true
		}
	}
	return false
}

// spend returns where the budget runs out along the sightline with tangent tan followed to upto,
// charging every see-through box among the first j crossings for the stretch the line spends in
// its band.
func (r *ray) spend(tan, upto float64, j int) float64 {
	if r.veils == 0 {
		return r.radius
	}
	st := r.sc.stretches[:0]
	for _, x := range r.cross[:j] {
		if x.rate == 0 {
			continue
		}
		c := &r.cands[x.idx]
		if lo, hi, ok := r.inBand(tan, x.near, min(x.far, upto), c.bottom, c.top); ok {
			st = insertCrossing(st, crossing{near: lo, far: hi, rate: x.rate})
		}
	}
	r.sc.stretches = st
	if len(st) == 0 {
		return r.radius
	}
	return spend(st, upto, r.radius)
}

// heights is the sightline's height at both ends of a stretch, lower first; a vertical line
// (an entry with no top) is as high as anything past the eye.
func (r *ray) heights(tan, from, to float64) (lo, hi float64) {
	lo, hi = r.at(tan, from), r.at(tan, to)
	if lo > hi {
		lo, hi = hi, lo
	}
	return lo, hi
}

func (r *ray) at(tan, d float64) float64 {
	if d == 0 {
		return r.e.eye
	}
	return r.e.eye + tan*d
}

// inBand is the part of [from, to] where the sightline runs between bottom and top.
func (r *ray) inBand(tan, from, to, bottom, top float64) (lo, hi float64, ok bool) {
	switch {
	case math.IsInf(tan, 0):
		return from, to, math.IsInf(top, 1) && from < to
	case tan == 0:
		return from, to, bottom <= r.e.eye && r.e.eye <= top && from < to
	}
	d1, d2 := (bottom-r.e.eye)/tan, (top-r.e.eye)/tan
	if d1 > d2 {
		d1, d2 = d2, d1
	}
	lo, hi = max(from, d1), min(to, d2)
	return lo, hi, lo < hi
}

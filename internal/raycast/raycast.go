// Package raycast is the state and the algorithm behind the public raycast package:
// the sweep, the shadows, and the View they fill.
package raycast

import (
	"math"
	"slices"

	"github.com/kjkrol/aabbworld/geom"
	"github.com/kjkrol/uid"
)

// QueryableSpace is the slice of a spatial index this package needs.
type QueryableSpace interface {
	Query(aabb geom.AABB, fn func(id uid.UID64)) int
	EntryAABB(id uid.UID64) (geom.AABB, bool)
	Bounds() (width, height uint32, wrapX, wrapY bool)
}

// Cone bounds a visibility query: a direction, a half-angle either side of it,
// and how far it reaches.
type Cone struct {
	Direction geom.Vec
	HalfAngle float64
	Radius    float64
}

// View is one observer's line of sight: scanned once, then read as many ways as needed.
// The zero value is ready; it must not be copied once scanned, and serves one goroutine.
type View struct {
	scratch

	col     collector
	colFn   func(uid.UID64)
	areas   []geom.AABB
	found   []seen
	uniform []float64
	depths  []sample

	origin    geom.Vec
	coneDir   float64
	halfAngle float64
	radius    float64
	valid     bool
}

// seen is one entity's nearest hit across a sweep.
type seen struct {
	id   uid.UID64
	dist float64
}

// Scan records what observer sees through cone and reports whether it could.
func (v *View) Scan(space QueryableSpace, observer uid.UID64, cone Cone) bool {
	v.scratch.reset()
	v.valid = false

	if cone.HalfAngle <= 0 || cone.HalfAngle >= math.Pi || cone.Radius <= 0 {
		return false
	}
	eye, ok := space.EntryAABB(observer)
	if !ok {
		return false
	}

	width, height, wrapX, wrapY := space.Bounds()
	w, h := float64(width), float64(height)
	eyeBox := eye

	v.origin = centerOf(eyeBox)
	v.coneDir = math.Atan2(cone.Direction.Y, cone.Direction.X)
	v.halfAngle = cone.HalfAngle
	v.radius = cone.Radius

	v.samples = sweep(v.origin, v.coneDir, cone.HalfAngle, cone.Radius, v.gather(space, observer, eyeBox, cone, w, h, wrapX, wrapY), &v.scratch)
	v.valid = true
	return true
}

// Entities calls fn for every entity in view, nearest first, and returns how many.
func (v *View) Entities(fn func(id uid.UID64, dist float64)) int {
	if !v.valid {
		return 0
	}

	found := v.found[:0]
	for _, s := range v.samples {
		if !s.hit {
			continue
		}
		at := -1
		for i := range found {
			if found[i].id == s.id {
				at = i
				break
			}
		}
		switch {
		case at < 0:
			found = append(found, seen{s.id, s.dist})
		case s.dist < found[at].dist:
			found[at].dist = s.dist
		}
	}
	v.found = found

	slices.SortFunc(found, func(a, b seen) int {
		switch {
		case a.dist < b.dist:
			return -1
		case a.dist > b.dist:
			return 1
		default:
			return 0
		}
	})
	for _, f := range found {
		fn(f.id, f.dist)
	}
	return len(found)
}

// Depths appends to dst the reach at k ≥ 2 angles evenly spaced across the cone.
func (v *View) Depths(k int, dst []float32) []float32 {
	if !v.valid || k < 2 {
		return dst
	}

	angles := v.uniform[:0]
	step := 2 * v.halfAngle / float64(k-1)
	for i := range k {
		angles = append(angles, -v.halfAngle+float64(i)*step)
	}
	v.uniform = angles

	v.active, v.depths = walk(v.origin, v.coneDir, v.radius, 0,
		angles, v.events, v.cands, v.active, v.depths[:0])

	for _, s := range v.depths {
		dst = append(dst, float32(s.dist))
	}
	return dst
}

// Outline appends the lit region to dst as a fan from the observer; maxArcStep 0 picks a default.
func (v *View) Outline(maxArcStep float64, dst []geom.Vec) []geom.Vec {
	if !v.valid {
		return dst
	}
	if maxArcStep <= 0 {
		maxArcStep = defaultArcStep
	}

	dst = append(dst, v.origin)
	at := func(rel, dist float64) geom.Vec {
		a := v.coneDir + rel
		return geom.NewVec(v.origin.X+dist*math.Cos(a), v.origin.Y+dist*math.Sin(a))
	}

	for i, s := range v.samples {
		if i > 0 {
			prev := v.samples[i-1]
			if !prev.hit && !s.hit {
				for a := prev.angle + maxArcStep; a < s.angle; a += maxArcStep {
					dst = append(dst, at(a, v.radius))
				}
			}
		}
		dst = append(dst, at(s.angle, s.dist))
	}
	return dst
}

// defaultArcStep keeps an unobstructed arc smooth enough to read as round.
const defaultArcStep = math.Pi / 90 // two degrees

// scratch holds the working buffers one scan needs — candidates, event list,
// samples — kept across scans so a repeated query stops allocating.
type scratch struct {
	cands   []candidate
	dedup   map[uid.UID64]struct{}
	angles  []float64
	events  []event
	active  []int
	samples []sample
}

func (s *scratch) reset() {
	if s.dedup == nil {
		s.dedup = map[uid.UID64]struct{}{}
	}
	s.cands = s.cands[:0]
	clear(s.dedup)
	s.angles = s.angles[:0]
	s.events = s.events[:0]
	s.active = s.active[:0]
	s.samples = s.samples[:0]
}

type candidate struct {
	id   uid.UID64
	dist float64
	span arc
	box  geom.AABB
}

// collector is what one gather needs while the index walks it.
type collector struct {
	space    QueryableSpace
	observer uid.UID64
	eyeBox   geom.AABB
	origin   geom.Vec
	cone     Cone
	coneDir  float64
	edges    wedge
	w, h     float64
	wrapX    bool
	wrapY    bool
	sc       *scratch
}

func (c *collector) take(id uid.UID64) {
	if id == c.observer {
		return
	}
	if _, dup := c.sc.dedup[id]; dup {
		return
	}
	raw, ok := c.space.EntryAABB(id)
	if !ok {
		return
	}
	box := raw
	if c.wrapX || c.wrapY {
		box = nearestImage(c.origin, box, wrapSize(c.w, c.wrapX), wrapSize(c.h, c.wrapY))
	}
	dist := boxDistance(c.eyeBox, box)
	if dist > c.cone.Radius || c.edges.excludes(c.origin, box) {
		return
	}
	span := subtendedArc(c.origin, box, c.coneDir, c.cone.HalfAngle)
	if span.empty() {
		return
	}
	c.sc.dedup[id] = struct{}{}
	c.sc.cands = append(c.sc.cands, candidate{id: id, dist: dist, span: span, box: box})
}

// gather collects the entities within range whose angular span overlaps the cone.
func (v *View) gather(
	space QueryableSpace,
	observer uid.UID64,
	eyeBox geom.AABB,
	cone Cone,
	w, h float64,
	wrapX, wrapY bool,
) []candidate {
	v.col = collector{
		space: space, observer: observer, eyeBox: eyeBox, origin: v.origin,
		cone: cone, coneDir: v.coneDir, edges: newWedge(v.coneDir, cone.HalfAngle),
		w: w, h: h, wrapX: wrapX, wrapY: wrapY, sc: &v.scratch,
	}
	if v.colFn == nil {
		v.colFn = v.col.take
	}

	v.areas = searchAreas(v.origin, cone.Radius, w, h, wrapX, wrapY, v.areas[:0])
	for _, r := range v.areas {
		space.Query(r, v.colFn)
	}
	return v.cands
}

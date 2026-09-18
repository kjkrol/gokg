// Package raycast answers "what can this entity see" over a spatial index:
// every entity within a cone and range, minus those hidden behind nearer ones.
//
// Every entity occludes. If C stands between A and B so that a ray from A stops
// at C, then C is visible and B is not. Entities are assumed not to overlap or
// contain one another.
//
// A query goes through [View]: scan once, then read the result as entities, as
// an outline, or as both.
package raycast

import (
	"math"
	"slices"

	"github.com/kjkrol/gokg/geom"
	"github.com/kjkrol/gokg/plane"
	"github.com/kjkrol/uid"
)

// QueryableSpace is the slice of a spatial index this package needs.
type QueryableSpace interface {
	Query(aabb geom.AABB, fn func(id uid.UID64, frag plane.FragPosition)) int
	EntryAABB(id uid.UID64) (geom.AABB, bool)
	Bounds() (width, height uint32, toroidal bool)
}

// Cone bounds a visibility query: a direction, a half-angle either side of it,
// and how far it reaches.
type Cone struct {
	Direction geom.Vec
	HalfAngle float64
	Radius    float64
}

// View is one observer's line of sight: scanned once, then read as many ways as
// needed. Its buffers belong to it, so a View kept across ticks stops
// allocating once they have grown — the outline included.
//
// The zero value is ready to use. A View must not be copied once scanned — its
// callback is bound to its own address — and one View serves one goroutine;
// scanning in parallel needs a View per worker.
type View struct {
	scratch

	col     collector
	colFn   func(uid.UID64, plane.FragPosition)
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

// Scan records what observer sees through cone, replacing whatever the View
// held before, and reports whether the query was answerable at all.
//
// HalfAngle must be below π; a wider cone wraps past ±π, which the angular
// arithmetic here does not model.
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

	width, height, toroidal := space.Bounds()
	w, h := float64(width), float64(height)
	eyeBox := eye

	v.origin = centerOf(eyeBox)
	v.coneDir = math.Atan2(cone.Direction.Y, cone.Direction.X)
	v.halfAngle = cone.HalfAngle
	v.radius = cone.Radius

	v.samples = sweep(v.origin, v.coneDir, cone.HalfAngle, cone.Radius, v.gather(space, observer, eyeBox, cone, w, h, toroidal), &v.scratch)
	v.valid = true
	return true
}

// Entities calls fn for every entity the View can see, nearest first, and
// returns how many there were. The observer itself is never reported.
//
// fn receives values, so nothing it keeps points back into the View and a later
// Scan cannot disturb it.
func (v *View) Entities(fn func(id uid.UID64, dist float64)) int {
	if !v.valid {
		return 0
	}

	// One report per entity. The set is small, so a linear scan beats a map.
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

// Depths appends k distances to dst, one per evenly spaced angle across the
// cone: angle i sits at HalfAngle*(2*i/(k-1) - 1) from the cone's direction, so
// the first and last land on its edges.
//
// Only the distances come back — the angles follow from the index, which is a
// quarter of the memory the same fan costs as points. k must be at least 2;
// below that, or without a successful Scan, dst comes back untouched.
//
// This is the approximate form: a shadow edge lands on the nearest sampled
// angle, and a sliver thinner than one step can be missed entirely. Widening
// the cone or reaching further needs a larger k to hold the same accuracy —
// the error at full range is Radius*2*HalfAngle/(k-1). Entities is unaffected;
// it reads the exact sweep.
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

	// The candidates and their span events are whatever the last Scan gathered;
	// only the angles they are tested at differ.
	v.active, v.depths = walk(v.origin, v.coneDir, v.radius, 0,
		angles, v.events, v.cands, v.active, v.depths[:0])

	for _, s := range v.depths {
		dst = append(dst, float32(s.dist))
	}
	return dst
}

// Outline appends the boundary of the lit region to dst as a fan around the
// observer: the cone clipped by whatever it meets, falling back to the range
// limit where nothing blocks. maxArcStep caps the angular gap between points
// along that limit, so the unobstructed part reads as an arc rather than a
// chord; pass zero for a sensible default.
//
// The first appended point is the observer itself, so the result can be drawn
// directly as a triangle fan. Points are copied into dst, so a later Scan
// leaves them alone.
//
// Passing buf[:0] reuses a buffer across frames. Assign the result back, as
// with any append: growing dst returns a different array, and dropping it would
// leave the caller holding the previous frame's points.
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
		// Where nothing was hit the true boundary is a circular arc, and a
		// straight chord between samples would cut the corner off.
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

// collector is what one gather needs while the index walks it. It lives in the
// View so the callback handed to Query can be bound once and reused, rather
// than built as a fresh closure on every scan.
type collector struct {
	space    QueryableSpace
	observer uid.UID64
	eyeBox   geom.AABB
	origin   geom.Vec
	cone     Cone
	coneDir  float64
	edges    wedge
	w, h     float64
	toroidal bool
	sc       *scratch
}

func (c *collector) take(id uid.UID64, _ plane.FragPosition) {
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
	if c.toroidal {
		box = nearestImage(c.origin, box, c.w, c.h)
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

// gather collects everything inside the cone's bounding square, keeping only
// entities within range whose angular span overlaps the cone.
func (v *View) gather(
	space QueryableSpace,
	observer uid.UID64,
	eyeBox geom.AABB,
	cone Cone,
	w, h float64,
	toroidal bool,
) []candidate {
	v.col = collector{
		space: space, observer: observer, eyeBox: eyeBox, origin: v.origin,
		cone: cone, coneDir: v.coneDir, edges: newWedge(v.coneDir, cone.HalfAngle),
		w: w, h: h, toroidal: toroidal, sc: &v.scratch,
	}
	if v.colFn == nil {
		// Bound to &v.col, which never moves, so later scans need no rebinding.
		v.colFn = v.col.take
	}

	v.areas = searchAreas(v.origin, cone.Radius, w, h, toroidal, v.areas[:0])
	for _, r := range v.areas {
		space.Query(r, v.colFn)
	}
	return v.cands
}

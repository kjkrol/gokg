package aabbworld

import (
	"github.com/kjkrol/aabbworld/geom"
	iraycast "github.com/kjkrol/aabbworld/internal/raycast"
	"github.com/kjkrol/aabbworld/internal/spatial"
	"github.com/kjkrol/uid"
)

// Cone bounds a visibility query: a direction, a half-angle either side of it, how far it reaches,
// how see-through each entity is and, when sight has heights, where the eye is and how the world
// stands. The budget model behind Transparency and the sightlines behind Eye, Elevation and
// Ground are in the package doc.
type Cone struct {
	Direction geom.Vec
	HalfAngle float64 // below π
	Radius    float64
	// Transparency is τ per entity: 1 as empty, 0.5 costs twice its depth, ≤ 0 blocks; nil blocks all.
	Transparency func(id uid.UID64) float64
	// Eye is the height the observer looks from; it matters once Elevation or Ground is set.
	Eye float64
	// Elevation is the bottom and top of each entity; nil spans every entity over all heights.
	Elevation func(id uid.UID64) (bottom, top float64)
	// Ground is the height of the ground at a point, sampled every GroundStep along a ray (0: a
	// sixteenth of Radius); nil is flat ground at 0.
	Ground     func(p geom.Vec) float64
	GroundStep float64
}

// View is one observer's line of sight: filled by Space.Scan, then read as many ways as needed.
// The zero value is ready; it must not be copied once scanned, and serves one goroutine.
type View struct {
	impl    iraycast.View
	shadows []iraycast.Shadow
}

// Entities calls fn for every entity in view, nearest first, and returns how many.
func (v *View) Entities(fn func(id uid.UID64, dist float64)) int { return v.impl.Entities(fn) }

// Depths appends to dst the reach at k ≥ 2 angles evenly spaced across the cone.
func (v *View) Depths(k int, dst []float32) []float32 { return v.impl.Depths(k, dst) }

// Shadow is a stretch of ground along one of the angles a View reads, From to To away from the
// observer, that the observer cannot see; Sample is the angle's index, as in Depths.
type Shadow struct {
	Sample   int
	From, To float32
}

// Shadows appends to dst the stretches of ground the observer cannot see within the radius, at the
// same k ≥ 2 angles as Depths: with heights every run of hidden ground — behind a crest, in a
// wall's shadow, past a cliff — up to where ground is lit again; on a plane, the reach to the radius.
func (v *View) Shadows(k int, dst []Shadow) []Shadow {
	v.shadows = v.impl.Shadows(k, v.shadows[:0])
	for _, s := range v.shadows {
		dst = append(dst, Shadow{Sample: s.Sample, From: s.From, To: s.To})
	}
	return dst
}

// Outline appends the lit region to dst as a fan from the observer; maxArcStep 0 picks a default.
func (v *View) Outline(maxArcStep float64, dst []geom.Vec) []geom.Vec {
	return v.impl.Outline(maxArcStep, dst)
}

// sight is a Space as a scan reads it. It is a type of its own so that what a
// scan needs of the index does not have to be part of Space's public face.
type sight struct{ s *Space }

func (e sight) Query(box geom.AABB, fn func(id uid.UID64)) int {
	return e.s.grid.Query(box, spatial.AnyCapability, fn)
}
func (e sight) EntryAABB(id uid.UID64) (geom.AABB, bool) { return e.s.grid.EntryAABB(id) }
func (e sight) Bounds() (width, height uint32, wrapX, wrapY bool) {
	return e.s.Width, e.s.Height, e.s.Edges.WrapsX(), e.s.Edges.WrapsY()
}

// Scan fills v with what observer sees through cone and reports whether it could.
func (s *Space) Scan(observer uid.UID64, cone Cone, v *View) bool {
	return v.impl.Scan(sight{s}, observer, iraycast.Cone(cone))
}

package aabbworld

import (
	"github.com/kjkrol/aabbworld/geom"
	iraycast "github.com/kjkrol/aabbworld/internal/raycast"
	"github.com/kjkrol/aabbworld/internal/spatial"
	"github.com/kjkrol/uid"
)

// Cone bounds a visibility query: a direction, a half-angle either side of it,
// and how far it reaches.
type Cone struct {
	Direction geom.Vec
	HalfAngle float64 // below π
	Radius    float64
}

// View is one observer's line of sight: filled by Space.Scan, then read as many ways as needed.
// The zero value is ready; it must not be copied once scanned, and serves one goroutine.
type View struct{ impl iraycast.View }

// Entities calls fn for every entity in view, nearest first, and returns how many.
func (v *View) Entities(fn func(id uid.UID64, dist float64)) int { return v.impl.Entities(fn) }

// Depths appends to dst the reach at k ≥ 2 angles evenly spaced across the cone.
func (v *View) Depths(k int, dst []float32) []float32 { return v.impl.Depths(k, dst) }

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

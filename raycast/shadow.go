package raycast

import (
	"math"

	"github.com/kjkrol/gokg/geom"
)

// arc is an angular interval, in radians, relative to a cone's direction.
type arc struct{ lo, hi float64 }

func (a arc) empty() bool { return a.lo > a.hi }

// wrapToPi folds an angle into (-π, π].
func wrapToPi(a float64) float64 {
	a = math.Mod(a+math.Pi, 2*math.Pi)
	if a <= 0 {
		a += 2 * math.Pi
	}
	return a - math.Pi
}

// subtendedArc returns the angular interval box occupies as seen from origin,
// expressed relative to coneDir and clipped to ±halfAngle. The returned arc is
// empty when the box falls outside the cone.
//
// origin must lie outside box, which keeps the box's angular span below π and
// the interval unambiguous. The two silhouette corners are picked with cross
// products — sign alone says which of two directions is the more clockwise —
// so only they need an atan2, rather than all four corners plus a reference.
func subtendedArc(origin geom.Vec[float64], box geom.AABB[float64], coneDir, halfAngle float64) arc {
	corners := [4][2]float64{
		{box.TopLeft.X - origin.X, box.TopLeft.Y - origin.Y},
		{box.BottomRight.X - origin.X, box.TopLeft.Y - origin.Y},
		{box.BottomRight.X - origin.X, box.BottomRight.Y - origin.Y},
		{box.TopLeft.X - origin.X, box.BottomRight.Y - origin.Y},
	}

	lo, hi := corners[0], corners[0]
	for _, c := range corners[1:] {
		if lo[0]*c[1]-lo[1]*c[0] < 0 {
			lo = c // c sits clockwise of the current start
		}
		if hi[0]*c[1]-hi[1]*c[0] > 0 {
			hi = c // c sits counter-clockwise of the current end
		}
	}

	start := math.Atan2(lo[1], lo[0])
	span := wrapToPi(math.Atan2(hi[1], hi[0]) - start)

	// Move the whole interval into the cone's frame, not each corner: the span
	// stays below π, so wrapping its midpoint keeps it contiguous.
	mid := wrapToPi(start + span/2 - coneDir)
	half := span / 2

	return arc{
		lo: math.Max(mid-half, -halfAngle),
		hi: math.Min(mid+half, halfAngle),
	}
}

// wedge holds the cone's two bounding edge directions so a candidate can be
// dismissed with cross products alone, before paying for any atan2.
//
// Only a cone narrower than a half-turn is bounded by the intersection of the
// two half-planes; a wider one keeps usable false and rejects nothing.
type wedge struct {
	lo, hi geom.Vec[float64]
	usable bool
}

func newWedge(coneDir, halfAngle float64) wedge {
	if halfAngle >= math.Pi/2 {
		return wedge{}
	}
	return wedge{
		lo:     geom.NewVec(math.Cos(coneDir-halfAngle), math.Sin(coneDir-halfAngle)),
		hi:     geom.NewVec(math.Cos(coneDir+halfAngle), math.Sin(coneDir+halfAngle)),
		usable: true,
	}
}

// excludes reports whether every corner of box lies beyond the same cone edge,
// which puts the whole box outside the cone.
func (w wedge) excludes(origin geom.Vec[float64], box geom.AABB[float64]) bool {
	if !w.usable {
		return false
	}
	pastLo, pastHi := true, true
	for _, dx := range [2]float64{box.TopLeft.X - origin.X, box.BottomRight.X - origin.X} {
		for _, dy := range [2]float64{box.TopLeft.Y - origin.Y, box.BottomRight.Y - origin.Y} {
			pastLo = pastLo && w.lo.X*dy-w.lo.Y*dx < 0
			pastHi = pastHi && w.hi.X*dy-w.hi.Y*dx > 0
			if !pastLo && !pastHi {
				return false
			}
		}
	}
	return true
}

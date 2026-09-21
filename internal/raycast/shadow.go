package raycast

import (
	"math"

	"github.com/kjkrol/aabbworld/geom"
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

// subtendedArc is the arc box spans as seen from origin outside it, clipped to ±halfAngle.
func subtendedArc(origin geom.Vec, box geom.AABB, coneDir, halfAngle float64) arc {
	corners := [4][2]float64{
		{box.TopLeft.X - origin.X, box.TopLeft.Y - origin.Y},
		{box.BottomRight.X - origin.X, box.TopLeft.Y - origin.Y},
		{box.BottomRight.X - origin.X, box.BottomRight.Y - origin.Y},
		{box.TopLeft.X - origin.X, box.BottomRight.Y - origin.Y},
	}

	lo, hi := corners[0], corners[0]
	for _, c := range corners[1:] {
		if lo[0]*c[1]-lo[1]*c[0] < 0 {
			lo = c
		}
		if hi[0]*c[1]-hi[1]*c[0] > 0 {
			hi = c
		}
	}

	start := math.Atan2(lo[1], lo[0])
	span := wrapToPi(math.Atan2(hi[1], hi[0]) - start)

	mid := wrapToPi(start + span/2 - coneDir)
	half := span / 2

	return arc{
		lo: math.Max(mid-half, -halfAngle),
		hi: math.Min(mid+half, halfAngle),
	}
}

// wedge holds the cone's two edge directions, dismissing a candidate by cross products alone.
// A cone of a half-turn or wider keeps usable false and rejects nothing.
type wedge struct {
	lo, hi geom.Vec
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

// excludes reports whether the whole box lies beyond one cone edge.
func (w wedge) excludes(origin geom.Vec, box geom.AABB) bool {
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

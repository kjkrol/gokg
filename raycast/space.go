package raycast

import (
	"math"

	"github.com/kjkrol/gokg/geom"
)

func centerOf(a geom.AABB) geom.Vec {
	return geom.NewVec(
		(a.TopLeft.X+a.BottomRight.X)/2,
		(a.TopLeft.Y+a.BottomRight.Y)/2,
	)
}

// boxDistance is the gap between two boxes, zero when they touch or overlap.
func boxDistance(a, b geom.AABB) float64 {
	return math.Hypot(a.AxisDistanceX(b), a.AxisDistanceY(b))
}

// nearestImage shifts box to whichever of its wrapped copies lies closest to
// origin, so angles and distances can be measured in one unwrapped frame.
func nearestImage(origin geom.Vec, box geom.AABB, w, h float64) geom.AABB {
	c := centerOf(box)
	dx := shortestDelta(c.X-origin.X, w)
	dy := shortestDelta(c.Y-origin.Y, h)
	d := geom.NewVec(origin.X+dx-c.X, origin.Y+dy-c.Y)
	return geom.NewAABB(box.TopLeft.Add(d), box.BottomRight.Add(d))
}

func shortestDelta(d, size float64) float64 {
	if size <= 0 {
		return d
	}
	half := size / 2
	switch {
	case d > half:
		return d - size
	case d < -half:
		return d + size
	default:
		return d
	}
}

// searchAreas returns the index rectangles to query for a cone of the given
// radius. A toroidal world needs up to four, because Space.Query does not wrap
// and would silently drop whatever lies across the seam.
func searchAreas(origin geom.Vec, radius, w, h float64, toroidal bool, dst []geom.AABB) []geom.AABB {
	lo := geom.NewVec(origin.X-radius, origin.Y-radius)
	hi := geom.NewVec(origin.X+radius, origin.Y+radius)

	if !toroidal {
		return append(dst, clampRect(lo, hi, w, h))
	}

	for _, xs := range splitAxis(lo.X, hi.X, w) {
		for _, ys := range splitAxis(lo.Y, hi.Y, h) {
			dst = append(dst, clampRect(
				geom.NewVec(xs[0], ys[0]),
				geom.NewVec(xs[1], ys[1]), w, h))
		}
	}
	return dst
}

// splitAxis wraps [lo, hi] into [0, size), cutting it in two when it crosses
// the seam.
func splitAxis(lo, hi, size float64) [][2]float64 {
	if size <= 0 || hi-lo >= size {
		return [][2]float64{{0, size}}
	}
	span := hi - lo
	lo = math.Mod(lo, size)
	if lo < 0 {
		lo += size
	}
	hi = lo + span
	if hi <= size {
		return [][2]float64{{lo, hi}}
	}
	return [][2]float64{{lo, size}, {0, hi - size}}
}

func clampRect(lo, hi geom.Vec, w, h float64) geom.AABB {
	cl := func(v, max float64) float64 {
		switch {
		case v < 0:
			return 0
		case v > max:
			return max
		default:
			return v
		}
	}
	return geom.NewAABB(
		geom.NewVec(cl(lo.X, w), cl(lo.Y, h)),
		geom.NewVec(cl(hi.X, w), cl(hi.Y, h)),
	)
}

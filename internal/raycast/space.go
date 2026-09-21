package raycast

import (
	"math"

	"github.com/kjkrol/aabbworld/geom"
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

// nearestImage shifts box to whichever of its wrapped copies lies closest to origin.
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

// searchAreas appends the index rectangles covering a cone of radius, split at each wrapping seam.
func searchAreas(origin geom.Vec, radius, w, h float64, wrapX, wrapY bool, dst []geom.AABB) []geom.AABB {
	lo := geom.NewVec(origin.X-radius, origin.Y-radius)
	hi := geom.NewVec(origin.X+radius, origin.Y+radius)

	for _, xs := range axisSpans(lo.X, hi.X, w, wrapX) {
		for _, ys := range axisSpans(lo.Y, hi.Y, h, wrapY) {
			dst = append(dst, clampRect(
				geom.NewVec(xs[0], ys[0]),
				geom.NewVec(xs[1], ys[1]), w, h))
		}
	}
	return dst
}

func axisSpans(lo, hi, size float64, wraps bool) [][2]float64 {
	if !wraps {
		return [][2]float64{{lo, hi}}
	}
	return splitAxis(lo, hi, size)
}

// wrapSize is size on a wrapping axis and zero on any other, which shortestDelta leaves alone.
func wrapSize(size float64, wraps bool) float64 {
	if wraps {
		return size
	}
	return 0
}

// splitAxis wraps [lo, hi] into [0, size), cutting it in two where it crosses the seam.
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

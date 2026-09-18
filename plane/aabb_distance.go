package plane

import "github.com/kjkrol/gokg/geom"

// AABBDistance measures gaps between
// axis-aligned bounding boxes using the metric defined by the provided plane.
type AABBDistance func(aabb1, aabb2 geom.AABB) float64

func newAABBDistance(metric Metric) AABBDistance {
	return func(aabb1, aabb2 geom.AABB) float64 {
		if aabb1.Intersects(aabb2) {
			return 0
		}
		dx := aabb1.AxisDistanceX(aabb2)
		dy := aabb1.AxisDistanceY(aabb2)
		return metric(geom.NewVec(dx, dy), geom.NewVec(0, 0))
	}
}

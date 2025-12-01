package plane

import "github.com/kjkrol/gokg/pkg/geom"

// AABBDistance measures gaps between
// axis-aligned bounding boxes using the metric defined by the provided plane.
type AABBDistance[T geom.Numeric] func(aabb1, aabb2 geom.AABB[T]) T

func newAABBDistance[T geom.Numeric](metric Metric[T]) AABBDistance[T] {
	return func(aabb1, aabb2 geom.AABB[T]) T {
		if aabb1.Intersects(aabb2) {
			return 0
		}
		dx := aabb1.AxisDistanceX(aabb2)
		dy := aabb1.AxisDistanceY(aabb2)
		return metric(geom.NewVec(dx, dy), geom.NewVec[T](0, 0))
	}
}

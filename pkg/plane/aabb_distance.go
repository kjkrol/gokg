package plane

import "github.com/kjkrol/gokg/pkg/geom"

// AABBDistance measures gaps between
// axis-aligned bounding boxes using the metric defined by the provided plane.
type AABBDistance[T geom.Numeric] func(aa, bb geom.AABB[T]) T

func newAABBDistance[T geom.Numeric](metric Metric[T]) AABBDistance[T] {
	return func(aa, bb geom.AABB[T]) T {
		if aa.Intersects(bb) {
			return 0
		}
		dx := aa.AxisDistanceTo(bb, func(v geom.Vec[T]) T { return v.X })
		dy := aa.AxisDistanceTo(bb, func(v geom.Vec[T]) T { return v.Y })
		return metric(geom.NewVec(dx, dy), geom.NewVec[T](0, 0))
	}
}

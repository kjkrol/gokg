package geometry

// BoundingBoxDistance measures gaps between
// axis-aligned bounding boxes using the metric defined by the provided plane.
type BoundingBoxDistance[T SupportedNumeric] func(aa, bb BoundingBox[T]) T

func newBoundingBoxDistance[T SupportedNumeric](metric Metric[T]) BoundingBoxDistance[T] {
	return func(aa, bb BoundingBox[T]) T {
		if aa.Intersects(bb) {
			return 0
		}
		dx := aa.axisDistanceTo(bb, func(v Vec[T]) T { return v.X })
		dy := aa.axisDistanceTo(bb, func(v Vec[T]) T { return v.Y })
		return metric(NewVec(dx, dy), NewVec[T](0, 0))
	}
}

package geometry

// Distance represents a strategy that computes the gap between two spatial objects.
type Distance[T SupportedNumeric] func(a, b BoundingBox[T]) T

// BoundingBoxDistanceForPlane builds a Distance strategy that measures gaps between
// axis-aligned bounding boxes using the metric defined by the provided plane.
func BoundingBoxDistance[T SupportedNumeric](plane Plane[T]) Distance[T] {
	return func(a, b BoundingBox[T]) T {
		return boundingBoxDistance(a, b, plane.metric)
	}
}

func boundingBoxDistance[T SupportedNumeric](
	aa BoundingBox[T],
	bb BoundingBox[T],
	metric func(Vec[T], Vec[T]) T,
) T {
	if aa.Intersects(bb) {
		return 0
	}

	dx := aa.axisDistanceTo(bb, func(v Vec[T]) T { return v.X })
	dy := aa.axisDistanceTo(bb, func(v Vec[T]) T { return v.Y })

	return metric(Vec[T]{X: dx, Y: dy}, Vec[T]{X: 0, Y: 0})
}

//-------------------------------------------------------------------------

package geometry

// Distance represents a strategy that computes the gap between two spatial objects.
type Distance[T SupportedNumeric] func(a, b Shape[T]) T

// BoundingBoxDistanceForPlane builds a Distance strategy that measures gaps between
// axis-aligned bounding boxes using the metric defined by the provided plane.
func BoundingBoxDistanceForPlane[T SupportedNumeric](plane Plane[T]) Distance[T] {
	return BoundingBoxDistance(plane.Metric)
}

// BoundingBoxDistance returns a Distance strategy evaluating the separation between
// objects via their axis-aligned bounding boxes and the supplied metric.
func BoundingBoxDistance[T SupportedNumeric](metric func(Vec[T], Vec[T]) T) Distance[T] {
	return func(a, b Shape[T]) T {
		return boundingBoxDistance(a, b, metric)
	}
}

func boundingBoxDistance[T SupportedNumeric](
	a Shape[T],
	b Shape[T],
	metric func(Vec[T], Vec[T]) T,
) T {
	boundsA := a.Bounds()
	boundsB := b.Bounds()

	if boundsA.Intersects(boundsB) {
		return 0
	}

	dx := boundsA.AxisDistanceTo(boundsB, func(v Vec[T]) T { return v.X })
	dy := boundsA.AxisDistanceTo(boundsB, func(v Vec[T]) T { return v.Y })

	return metric(Vec[T]{X: dx, Y: dy}, Vec[T]{X: 0, Y: 0})
}

package geometry

import (
	s "github.com/kjkrol/gokg/pkg/geometry/spatial"
)

// Distance represents a strategy that computes the gap between two spatial objects.
type Distance[T supportedNumeric] func(a, b s.Spatial[T]) T

// BoundingBoxDistanceForPlane builds a Distance strategy that measures gaps between
// axis-aligned bounding boxes using the metric defined by the provided plane.
func BoundingBoxDistanceForPlane[T supportedNumeric](plane Plane[T]) Distance[T] {
	return BoundingBoxDistance(plane.Metric)
}

// BoundingBoxDistance returns a Distance strategy evaluating the separation between
// objects via their axis-aligned bounding boxes and the supplied metric.
func BoundingBoxDistance[T supportedNumeric](metric func(s.Vec[T], s.Vec[T]) T) Distance[T] {
	return func(a, b s.Spatial[T]) T {
		return boundingBoxDistance(a, b, metric)
	}
}

func boundingBoxDistance[T supportedNumeric](
	a s.Spatial[T],
	b s.Spatial[T],
	metric func(s.Vec[T], s.Vec[T]) T,
) T {
	boundsA := a.Bounds()
	boundsB := b.Bounds()

	if boundsA.Intersects(boundsB) {
		return 0
	}

	dx := boundsA.AxisDistanceTo(boundsB, func(v vec[T]) T { return v.X })
	dy := boundsA.AxisDistanceTo(boundsB, func(v vec[T]) T { return v.Y })

	return metric(vec[T]{X: dx, Y: dy}, vec[T]{X: 0, Y: 0})
}

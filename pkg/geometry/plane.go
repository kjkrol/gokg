package geometry

import (
	s "github.com/kjkrol/gokg/pkg/geometry/spatial"
)

type Plane[T supportedNumeric] struct {
	size       vec[T]
	vectorMath VectorMath[T]
	normalize  func(*vec[T])
	metric     func(v1, v2 vec[T]) T
	name       string
}

// -----------------------------------------------------------------------------

func (p Plane[T]) Size() s.Vec[T] { return p.size }

func (p Plane[T]) Translate(vec *s.Vec[T], delta s.Vec[T]) {
	vec.AddMutable(delta)
	p.normalize(vec)
}

func (p Plane[T]) TranslateSpatial(spatialItem s.Spatial[T], delta s.Vec[T]) {
	if spatialItem == nil {
		return
	}

	translateInPlace(spatialItem, delta)

	if p.name != "cyclic" {
		for _, v := range spatialItem.Vertices() {
			if v != nil {
				p.normalize(v)
			}
		}
		spatialItem.SetFragments(nil)
		return
	}

	spatialItem.SetFragments(wrapSpatialFragments(spatialItem, p.size, p.vectorMath))
}

func (p Plane[T]) Metric(v1, v2 s.Vec[T]) T { return p.metric(v1, v2) }

func (p Plane[T]) Contains(vec s.Vec[T]) bool {
	return vec.X >= 0 && vec.X < p.size.X && vec.Y >= 0 && vec.Y < p.size.Y
}

func (p Plane[T]) Normalize(vec *s.Vec[T]) { p.normalize(vec) }

func (p Plane[T]) relativeMetric(v1, v2 s.Vec[T]) T {
	delta := v1.Sub(v2)
	p.normalize(&delta)
	return p.vectorMath.Length(delta)
}

func (p Plane[T]) Name() string { return p.name }

// -----------------------------------------------------------------------------

func NewBoundedPlane[T supportedNumeric](sizeX, sizeY T) Plane[T] {
	plane := Plane[T]{
		name:       "bounded",
		size:       s.NewVec(sizeX, sizeY),
		vectorMath: VectorMathByType[T](),
	}
	plane.normalize = func(v *vec[T]) { plane.vectorMath.Clamp(v, plane.size) }
	plane.metric = func(v1, v2 vec[T]) T { return max(plane.relativeMetric(v1, v2), plane.relativeMetric(v2, v1)) }
	return plane
}

// -----------------------------------------------------------------------------

func NewCyclicBoundedPlane[T supportedNumeric](sizeX, sizeY T) Plane[T] {
	plane := Plane[T]{
		name:       "cyclic",
		size:       s.NewVec(sizeX, sizeY),
		vectorMath: VectorMathByType[T](),
	}
	plane.normalize = func(v *vec[T]) { plane.vectorMath.Wrap(v, plane.size) }
	plane.metric = func(v1, v2 vec[T]) T { return min(plane.relativeMetric(v1, v2), plane.relativeMetric(v2, v1)) }
	return plane
}

// -----------------------------------------------------------------------------

package geometry

type Plane[T SupportedNumeric] struct {
	size       Vec[T]
	vectorMath VectorMath[T]
	normalize  func(*Vec[T])
	metric     func(v1, v2 Vec[T]) T
	name       string
}

// -----------------------------------------------------------------------------

func (b Plane[T]) Size() Vec[T] { return b.size }

func (b Plane[T]) Translate(vec *Vec[T], delta Vec[T]) {
	vec.AddMutable(delta)
	b.normalize(vec)
}

func (b Plane[T]) Metric(v1, v2 Vec[T]) T { return b.metric(v1, v2) }

func (b Plane[T]) Contains(vec Vec[T]) bool {
	return vec.X >= 0 && vec.X < b.size.X && vec.Y >= 0 && vec.Y < b.size.Y
}

func (b Plane[T]) Normalize(vec *Vec[T]) { b.normalize(vec) }

func (b Plane[T]) relativeMetric(v1, v2 Vec[T]) T {
	delta := v1.Sub(v2)
	b.normalize(&delta)
	return b.vectorMath.Length(delta)
}

func (b Plane[T]) Name() string { return b.name }

// -----------------------------------------------------------------------------

func NewBoundedPlane[T SupportedNumeric](sizeX, sizeY T) Plane[T] {
	plane := Plane[T]{
		name:       "bounded",
		size:       Vec[T]{sizeX, sizeY},
		vectorMath: VectorMathByType[T](),
	}
	plane.normalize = func(v *Vec[T]) { plane.vectorMath.Clamp(v, plane.size) }
	plane.metric = func(v1, v2 Vec[T]) T { return max(plane.relativeMetric(v1, v2), plane.relativeMetric(v2, v1)) }
	return plane
}

// -----------------------------------------------------------------------------

func NewCyclicBoundedPlane[T SupportedNumeric](sizeX, sizeY T) Plane[T] {
	plane := Plane[T]{
		name:       "cyclic",
		size:       Vec[T]{sizeX, sizeY},
		vectorMath: VectorMathByType[T](),
	}
	plane.normalize = func(v *Vec[T]) { plane.vectorMath.Wrap(v, plane.size) }
	plane.metric = func(v1, v2 Vec[T]) T { return min(plane.relativeMetric(v1, v2), plane.relativeMetric(v2, v1)) }
	return plane
}

// -----------------------------------------------------------------------------

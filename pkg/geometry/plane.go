package geometry

const (
	BOUNDED = "bounded"
	CYCLIC  = "cyclic"
)

type Plane[T SupportedNumeric] struct {
	size       Vec[T]
	vectorMath VectorMath[T]
	normalize  func(*Vec[T])
	metric     func(v1, v2 Vec[T]) T
	name       string
}

// -----------------------------------------------------------------------------

func (p Plane[T]) Size() Vec[T] { return p.size }

func (p Plane[T]) Translate(shape Shape[T], delta Vec[T]) {
	if shape == nil {
		return
	}
	translate(shape, delta, p)
}

func (p Plane[T]) Metric(v1, v2 Vec[T]) T { return p.metric(v1, v2) }

func (p Plane[T]) Contains(vec Vec[T]) bool {
	return vec.X >= 0 && vec.X < p.size.X && vec.Y >= 0 && vec.Y < p.size.Y
}

func (p Plane[T]) Normalize(vec *Vec[T]) { p.normalize(vec) }

func (p Plane[T]) relativeMetric(v1, v2 Vec[T]) T {
	delta := v1.Sub(v2)
	p.normalize(&delta)
	return p.vectorMath.Length(delta)
}

func (p Plane[T]) Name() string { return p.name }

// -----------------------------------------------------------------------------

func NewBoundedPlane[T SupportedNumeric](sizeX, sizeY T) Plane[T] {
	plane := Plane[T]{
		name:       BOUNDED,
		size:       NewVec(sizeX, sizeY),
		vectorMath: VectorMathByType[T](),
	}
	plane.normalize = func(v *Vec[T]) { plane.vectorMath.Clamp(v, plane.size) }
	plane.metric = func(v1, v2 Vec[T]) T { return max(plane.relativeMetric(v1, v2), plane.relativeMetric(v2, v1)) }
	return plane
}

// -----------------------------------------------------------------------------

func NewCyclicBoundedPlane[T SupportedNumeric](sizeX, sizeY T) Plane[T] {
	plane := Plane[T]{
		name:       CYCLIC,
		size:       NewVec(sizeX, sizeY),
		vectorMath: VectorMathByType[T](),
	}
	plane.normalize = func(v *Vec[T]) { plane.vectorMath.Wrap(v, plane.size) }
	plane.metric = func(v1, v2 Vec[T]) T { return min(plane.relativeMetric(v1, v2), plane.relativeMetric(v2, v1)) }
	return plane
}

// -----------------------------------------------------------------------------

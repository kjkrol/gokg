package geometry

type Plane[T Number] interface {
	Translate(vec *Vec[T], delta Vec[T])
	Metric(v1, v2 Vec[T]) T
	Contains(vec Vec[T]) bool
}

// ---- Bounded Plane ----

type BoundedPlane[T Number] struct {
	Min, Max Vec[T]
	Geometry Geometry[T]
}

func (b BoundedPlane[T]) clamp(v *Vec[T]) {
	if v.X > b.Max.X-1 {
		v.X = b.Max.X - 1
	} else if v.X < b.Min.X {
		v.X = b.Min.X
	}
	if v.Y > b.Max.Y-1 {
		v.Y = b.Max.Y - 1
	} else if v.Y < b.Min.Y {
		v.Y = b.Min.Y
	}
}

func (b BoundedPlane[T]) Translate(vec *Vec[T], delta Vec[T]) {
	vec.AddMutable(delta)
	b.clamp(vec)
}

func (b BoundedPlane[T]) Metric(v1, v2 Vec[T]) T {
	delta1 := v1.Sub(v2)
	b.clamp(&delta1)
	delta2 := v2.Sub(v1)
	b.clamp(&delta2)
	return max(b.Geometry.Length(delta1), b.Geometry.Length(delta2))
}

func (b BoundedPlane[T]) Contains(vec Vec[T]) bool {
	return vec.X >= b.Min.X && vec.X <= b.Max.X &&
		vec.Y >= b.Min.Y && vec.Y <= b.Max.Y
}

func newBoundedPlane[T Number](sizeX, sizeY T, geometry Geometry[T]) BoundedPlane[T] {
	return BoundedPlane[T]{
		Min:      Vec[T]{0, 0},
		Max:      Vec[T]{sizeX, sizeY},
		Geometry: geometry,
	}
}

// ---- Cyclic bounded plane ----

type CyclicBoundedPlane[T Number] struct {
	Min, Max Vec[T]
	geometry Geometry[T]
}

func (c CyclicBoundedPlane[T]) wrap(v *Vec[T]) {
	c.geometry.ModMutable(v, c.Max)
	v.AddMutable(c.Max)
	c.geometry.ModMutable(v, c.Max)
}

func (c CyclicBoundedPlane[T]) Translate(vec *Vec[T], delta Vec[T]) {
	vec.AddMutable(delta)
	c.wrap(vec)
}

func (c CyclicBoundedPlane[T]) Metric(v1, v2 Vec[T]) T {
	delta1 := v1.Sub(v2)
	c.wrap(&delta1)
	delta2 := v2.Sub(v1)
	c.wrap(&delta2)
	return min(c.geometry.Length(delta1), c.geometry.Length(delta2))
}

func (c CyclicBoundedPlane[T]) Contains(vec Vec[T]) bool {
	return vec.X >= c.Min.X && vec.X < c.Max.X &&
		vec.Y >= c.Min.Y && vec.Y < c.Max.Y
}

func newCyclicBoundedPlane[T Number](sizeX, sizeY T, geometry Geometry[T]) CyclicBoundedPlane[T] {
	return CyclicBoundedPlane[T]{
		Min:      Vec[T]{0, 0},
		Max:      Vec[T]{sizeX, sizeY},
		geometry: geometry,
	}
}

// ---- discrete numbers planes

func NewDiscreteCyclicBoundedPlane(sizeX int, sizeY int) CyclicBoundedPlane[int] {
	return newCyclicBoundedPlane(sizeX, sizeY, IntGeometry{})
}

func NewDiscreteBoundedPlane(sizeX int, sizeY int) BoundedPlane[int] {
	return newBoundedPlane(sizeX, sizeY, IntGeometry{})
}

// ---- continuous numbers planes

func NewContinuousCyclicBoundedPlane(sizeX float64, sizeY float64) CyclicBoundedPlane[float64] {
	return newCyclicBoundedPlane(sizeX, sizeY, Float64Geometry{})
}

func NewContinuousBoundedPlane(sizeX float64, sizeY float64) BoundedPlane[float64] {
	return newBoundedPlane(sizeX, sizeY, Float64Geometry{})
}

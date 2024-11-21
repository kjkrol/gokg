package geometry

// Plane represents a geometric plane that supports translation, metric calculation, and containment checks.
// T is a numeric type that is supported by the plane.
// Plane represents a geometric plane with operations defined for a numeric type T.
// T must satisfy the SupportedNumeric constraint.
//
// Methods:
//   - Translate(vec *Vec[T], delta Vec[T]): Translates the given vector by the specified delta.
//   - Metric(v1, v2 Vec[T]) T: Computes the metric (distance or other measure) between two vectors.
//   - Contains(vec Vec[T]) bool: Checks if the plane contains the given vector.
//   - Max() Vec[T]: Returns the vector with the maximum coordinates in the plane.
//   - Min() Vec[T]: Returns the vector with the minimum coordinates in the plane.

type Plane[T SupportedNumeric] interface {
	Translate(vec *Vec[T], delta Vec[T])
	Metric(v1, v2 Vec[T]) T
	Contains(vec Vec[T]) bool
	Max() Vec[T]
	Min() Vec[T]
}

// boundedPlane represents a plane in a geometric space that is bounded by minimum and maximum vectors.
// The type parameter T must satisfy the SupportedNumeric constraint, which ensures that the plane's
// coordinates are numeric types.
//
// Fields:
// - min: The minimum vector defining one corner of the bounded plane.
// - max: The maximum vector defining the opposite corner of the bounded plane.
// - gometry: An embedded field that provides geometric operations and properties for the plane.
type boundedPlane[T SupportedNumeric] struct {
	min, max Vec[T]
	geometry Geometry[T]
}

func (c boundedPlane[T]) Max() Vec[T] { return c.max }

func (c boundedPlane[T]) Min() Vec[T] { return c.min }

func (b boundedPlane[T]) clamp(v *Vec[T]) {
	if v.X > b.max.X {
		v.X = b.max.X
	} else if v.X < b.min.X {
		v.X = b.min.X
	}
	if v.Y > b.max.Y {
		v.Y = b.max.Y
	} else if v.Y < b.min.Y {
		v.Y = b.min.Y
	}
}

func (b boundedPlane[T]) Translate(vec *Vec[T], delta Vec[T]) {
	vec.AddMutable(delta)
	b.clamp(vec)
}

func (b boundedPlane[T]) Metric(v1, v2 Vec[T]) T {
	delta1 := v1.Sub(v2)
	b.clamp(&delta1)
	delta2 := v2.Sub(v1)
	b.clamp(&delta2)
	return max(b.geometry.Length(delta1), b.geometry.Length(delta2))
}

func (b boundedPlane[T]) Contains(vec Vec[T]) bool {
	return vec.X >= b.min.X && vec.X <= b.max.X &&
		vec.Y >= b.min.Y && vec.Y <= b.max.Y
}

func newBoundedPlane[T SupportedNumeric](sizeX, sizeY T, geometry Geometry[T]) boundedPlane[T] {
	return boundedPlane[T]{
		min:      Vec[T]{0, 0},
		max:      Vec[T]{sizeX - 1, sizeY - 1},
		geometry: geometry,
	}
}

// cyclicBoundedPlane represents a plane in a cyclic bounded space.
// It is defined by minimum and maximum vectors, and a geometry object.
//
// T is a type parameter that must satisfy the SupportedNumeric constraint.
//
// Fields:
// - min: The minimum vector defining one corner of the plane.
// - max: The maximum vector defining the opposite corner of the plane.
// - geometry: The geometry object associated with the plane.
type cyclicBoundedPlane[T SupportedNumeric] struct {
	min, max, size Vec[T]
	Geometry       Geometry[T]
}

func (c cyclicBoundedPlane[T]) Max() Vec[T] { return c.max }

func (c cyclicBoundedPlane[T]) Min() Vec[T] { return c.min }

func (c cyclicBoundedPlane[T]) wrap(v *Vec[T]) {
	c.Geometry.ModMutable(v, c.size)
	v.AddMutable(c.size)
	c.Geometry.ModMutable(v, c.size)
}

func (c cyclicBoundedPlane[T]) Translate(vec *Vec[T], delta Vec[T]) {
	vec.AddMutable(delta)
	c.wrap(vec)
}

func (c cyclicBoundedPlane[T]) Metric(v1, v2 Vec[T]) T {
	delta1 := v1.Sub(v2)
	c.wrap(&delta1)
	delta2 := v2.Sub(v1)
	c.wrap(&delta2)
	return min(c.Geometry.Length(delta1), c.Geometry.Length(delta2))
}

func (c cyclicBoundedPlane[T]) Contains(vec Vec[T]) bool {
	return vec.X >= c.min.X && vec.X <= c.max.X &&
		vec.Y >= c.min.Y && vec.Y <= c.max.Y
}

func newCyclicBoundedPlane[T SupportedNumeric](sizeX, sizeY T, geometry Geometry[T]) Plane[T] {
	return cyclicBoundedPlane[T]{
		min:      Vec[T]{0, 0},
		max:      Vec[T]{sizeX - 1, sizeY - 1},
		size:     Vec[T]{sizeX, sizeY},
		Geometry: geometry,
	}
}

// NewDiscreteCyclicBoundedPlane creates a new CyclicBoundedPlane with discrete integer coordinates.
// The plane is bounded by the specified size in the X and Y dimensions.
// The plane wraps around cyclically, meaning that coordinates exceeding the bounds will wrap around to the other side.
//
// Parameters:
//   - sizeX: The size of the plane in the X dimension.
//   - sizeY: The size of the plane in the Y dimension.
//
// Returns:
//
//	A CyclicBoundedPlane with integer coordinates.
func NewDiscreteCyclicBoundedPlane(sizeX int, sizeY int) Plane[int] {
	return newCyclicBoundedPlane(sizeX, sizeY, INT_GEOMETRY)
}

// NewDiscreteBoundedPlane creates a new BoundedPlane with discrete integer coordinates.
// The plane is bounded by the specified size in the X and Y dimensions.
//
// Parameters:
//   - sizeX: The size of the plane in the X dimension.
//   - sizeY: The size of the plane in the Y dimension.
//
// Returns:
//
//	A BoundedPlane with integer coordinates.
func NewDiscreteBoundedPlane(sizeX int, sizeY int) Plane[int] {
	return newBoundedPlane(sizeX, sizeY, INT_GEOMETRY)
}

// NewContinuousCyclicBoundedPlane creates a new instance of CyclicBoundedPlane with the specified
// dimensions sizeX and sizeY. The plane is continuous and cyclic, meaning that it wraps around
// at the boundaries. The function uses Float64Geometry for the underlying geometry calculations.
//
// Parameters:
//   - sizeX: The width of the plane.
//   - sizeY: The height of the plane.
//
// Returns:
//
//	A new CyclicBoundedPlane with the specified dimensions.
func NewContinuousCyclicBoundedPlane(sizeX float64, sizeY float64) Plane[float64] {
	return newCyclicBoundedPlane(sizeX, sizeY, FLOAT_64_GEOMETRY)
}

// NewContinuousBoundedPlane creates a new BoundedPlane with the specified dimensions.
// The plane is continuous and bounded by the given sizeX and sizeY parameters.
//
// Parameters:
//   - sizeX: The size of the plane along the X-axis.
//   - sizeY: The size of the plane along the Y-axis.
//
// Returns:
//
//	A BoundedPlane instance with float64 precision.
func NewContinuousBoundedPlane(sizeX float64, sizeY float64) Plane[float64] {
	return newBoundedPlane(sizeX, sizeY, FLOAT_64_GEOMETRY)
}

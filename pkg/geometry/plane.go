package geometry

// Plane represents a geometric plane that supports translation, metric calculation, and containment checks.
// T is a numeric type that is supported by the plane.
type Plane[T SupportedNumeric] interface {
	// Translate moves the given vector by the specified delta.
	//
	// Parameters:
	//   - vec: A pointer to the vector that will be translated.
	//   - delta: The vector representing the translation delta.
	//
	// Returns:
	//   - None.
	Translate(vec *Vec[T], delta Vec[T])

	// Metric calculates a metric (such as distance) between two vectors.
	//
	// Parameters:
	//   - v1: The first vector.
	//   - v2: The second vector.
	//
	// Returns:
	//   - The calculated metric value between the two vectors.
	Metric(v1, v2 Vec[T]) T

	// Contains checks if the given vector is contained within a certain context (e.g., within a plane or a shape).
	//
	// Parameters:
	//   - vec: The vector to check for containment.
	//
	// Returns:
	//   - True if the vector is contained within the context, otherwise false.
	Contains(vec Vec[T]) bool
}

// BoundedPlane represents a plane in a geometric space that is bounded by minimum and maximum vectors.
// The type parameter T must satisfy the SupportedNumeric constraint, allowing for various numeric types.
//
// Fields:
// - Min: The minimum vector defining one corner of the bounding box.
// - Max: The maximum vector defining the opposite corner of the bounding box.
// - Geometry: The geometric properties associated with the bounded plane.
type BoundedPlane[T SupportedNumeric] struct {
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

func newBoundedPlane[T SupportedNumeric](sizeX, sizeY T, geometry Geometry[T]) BoundedPlane[T] {
	return BoundedPlane[T]{
		Min:      Vec[T]{0, 0},
		Max:      Vec[T]{sizeX, sizeY},
		Geometry: geometry,
	}
}

// CyclicBoundedPlane represents a plane in a cyclic bounded space.
// It is defined by minimum and maximum vectors, and a geometry object.
//
// T is a type parameter that must satisfy the SupportedNumeric constraint.
//
// Fields:
// - Min: The minimum vector defining one corner of the plane.
// - Max: The maximum vector defining the opposite corner of the plane.
// - geometry: The geometry object associated with the plane.
type CyclicBoundedPlane[T SupportedNumeric] struct {
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

func newCyclicBoundedPlane[T SupportedNumeric](sizeX, sizeY T, geometry Geometry[T]) CyclicBoundedPlane[T] {
	return CyclicBoundedPlane[T]{
		Min:      Vec[T]{0, 0},
		Max:      Vec[T]{sizeX, sizeY},
		geometry: geometry,
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
func NewDiscreteCyclicBoundedPlane(sizeX int, sizeY int) CyclicBoundedPlane[int] {
	return newCyclicBoundedPlane(sizeX, sizeY, IntGeometry{})
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
func NewDiscreteBoundedPlane(sizeX int, sizeY int) BoundedPlane[int] {
	return newBoundedPlane(sizeX, sizeY, IntGeometry{})
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
func NewContinuousCyclicBoundedPlane(sizeX float64, sizeY float64) CyclicBoundedPlane[float64] {
	return newCyclicBoundedPlane(sizeX, sizeY, Float64Geometry{})
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
func NewContinuousBoundedPlane(sizeX float64, sizeY float64) BoundedPlane[float64] {
	return newBoundedPlane(sizeX, sizeY, Float64Geometry{})
}

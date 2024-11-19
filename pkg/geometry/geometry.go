package geometry

import "math"

// Geometry is a generic interface for geometric operations on vectors.
//
// Type Parameters:
//   - T: A numeric type that can be either int or float64.
//
// Methods:
//   - Length(v Vec[T]) T: Calculates the length (magnitude) of a vector.
//   - Distance(v1, v2 Vec[T]) T: Calculates the distance between two vectors.
//   - Mod(v1, v2 Vec[T]) Vec[T]: Returns a new vector that is the result of applying the modulus operation between two vectors.
//   - ModMutable(v1 *Vec[T], v2 Vec[T]): Modifies the first vector by applying the modulus operation with the second vector.
type Geometry[T int | float64] interface {
	// Length calculates the length (magnitude) of a vector.
	//
	// Parameters:
	//   - v: The vector for which to calculate the length.
	//
	// Returns:
	//   - The length of the vector.
	Length(v Vec[T]) T

	// Distance calculates the distance between two vectors.
	//
	// Parameters:
	//   - v1: The first vector.
	//   - v2: The second vector.
	//
	// Returns:
	//   - The distance between the two vectors.
	Distance(v1, v2 Vec[T]) T

	// Mod returns a new vector that is the result of applying the modulus operation between two vectors.
	//
	// Parameters:
	//   - v1: The first vector.
	//   - v2: The second vector.
	//
	// Returns:
	//   - A new vector that is the result of the modulus operation.
	Mod(v1, v2 Vec[T]) Vec[T]

	// ModMutable modifies the first vector by applying the modulus operation with the second vector.
	//
	// Parameters:
	//   - v1: A pointer to the first vector, which will be modified.
	//   - v2: The second vector.
	//
	// Returns:
	//   - None.
	ModMutable(v1 *Vec[T], v2 Vec[T])
}

// Float64Geometry represents a geometric structure with floating-point precision.
// This type can be used to perform various geometric calculations and operations
// involving floating-point numbers.
type Float64Geometry struct{}

func (g Float64Geometry) Length(v Vec[float64]) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (g Float64Geometry) Distance(v1, v2 Vec[float64]) float64 {
	delta := v1.Sub(v2)
	return g.Length(delta)
}

func (g Float64Geometry) ModMutable(v1 *Vec[float64], v2 Vec[float64]) {
	v1.X = math.Mod(v1.X, v2.X)
	v1.Y = math.Mod(v1.Y, v2.Y)
}

func (g Float64Geometry) Mod(v1, v2 Vec[float64]) Vec[float64] {
	return Vec[float64]{math.Mod(v1.X, v2.X), math.Mod(v1.Y, v2.Y)}
}

// IntGeometry represents a geometric structure with integer coordinates.
// It provides methods to perform various geometric operations and calculations.
type IntGeometry struct{}

func (g IntGeometry) Length(v Vec[int]) int {
	return int(math.Ceil(math.Sqrt(float64(v.X*v.X + v.Y*v.Y))))
}

func (g IntGeometry) Distance(v1, v2 Vec[int]) int {
	delta := v1.Sub(v2)
	return g.Length(delta)
}

func (g IntGeometry) ModMutable(v1 *Vec[int], v2 Vec[int]) { v1.X %= v2.X; v1.Y %= v2.Y }

func (g IntGeometry) Mod(v1, v2 Vec[int]) Vec[int] { return Vec[int]{v1.X % v2.X, v1.Y % v2.Y} }

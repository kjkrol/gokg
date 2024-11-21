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
	Length(v Vec[T]) T
	Distance(v1, v2 Vec[T]) T
	Mod(v1, v2 Vec[T]) Vec[T]
	ModMutable(v1 *Vec[T], v2 Vec[T])
}

var FLOAT_64_GEOMETRY = float64Geometry{}

var INT_GEOMETRY = intGeometry{}

// float64Geometry represents a geometric structure with floating-point precision.
// This type can be used to perform various geometric calculations and operations
// involving floating-point numbers.
type float64Geometry struct{}

func (g float64Geometry) Length(v Vec[float64]) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (g float64Geometry) Distance(v1, v2 Vec[float64]) float64 {
	delta := v1.Sub(v2)
	return g.Length(delta)
}

func (g float64Geometry) ModMutable(v1 *Vec[float64], v2 Vec[float64]) {
	v1.X = math.Mod(v1.X, v2.X)
	v1.Y = math.Mod(v1.Y, v2.Y)
}

func (g float64Geometry) Mod(v1, v2 Vec[float64]) Vec[float64] {
	return Vec[float64]{math.Mod(v1.X, v2.X), math.Mod(v1.Y, v2.Y)}
}

// intGeometry represents a geometric structure with integer coordinates.
// It provides methods to perform various geometric operations and calculations.
type intGeometry struct{}

func (g intGeometry) Length(v Vec[int]) int {
	return int(math.Ceil(math.Sqrt(float64(v.X*v.X + v.Y*v.Y))))
}
func (g intGeometry) Distance(v1, v2 Vec[int]) int {
	delta := v1.Sub(v2)
	return g.Length(delta)
}

func (g intGeometry) ModMutable(v1 *Vec[int], v2 Vec[int]) { v1.X %= v2.X; v1.Y %= v2.Y }

func (g intGeometry) Mod(v1, v2 Vec[int]) Vec[int] { return Vec[int]{v1.X % v2.X, v1.Y % v2.Y} }

// Package geometry provides a set of operations for 2D graphics
package geometry

import (
	"fmt"
)

type (
	SupportedNumeric interface{ int | float64 }

	// Vec is a generic 2D vector with components X and Y.
	//
	// Type Parameters:
	//   - T: A numeric type that satisfies the Number constraint.
	Vec[T SupportedNumeric] struct{ X, Y T }
)

var (
	ZERO_INT_VEC   = Vec[int]{0, 0}
	ZERO_FLOAT_VEC = Vec[float64]{0, 0}
)

func NewVec[T SupportedNumeric](X, Y T) Vec[T] {
	return Vec[T]{X, Y}
}

// Add computes the sum of the current vector and the given vector.
// It returns a new vector containing the result without modifying the current vector.
//
// Parameters:
//   - v2: The vector to add to the current vector.
//
// Returns:
//   - A new vector representing the sum of the current vector and v2.
func (v Vec[T]) Add(v2 Vec[T]) Vec[T] { return Vec[T]{v.X + v2.X, v.Y + v2.Y} }

// Sub computes the difference between the current vector and the given vector.
// It returns a new vector containing the result without modifying the current vector.
//
// Parameters:
//   - v2: The vector to subtract from the current vector.
//
// Returns:
//   - A new vector representing the difference between the current vector and v2.
func (v Vec[T]) Sub(v2 Vec[T]) Vec[T] { return Vec[T]{v.X - v2.X, v.Y - v2.Y} }

// AddMutable adds the given vector to the current vector, modifying the current vector in place.
//
// Parameters:
//   - v2: The vector to add to the current vector.
//
// Effect:
//   - Updates the X and Y components of the current vector by adding the corresponding components of v2.
func (v *Vec[T]) AddMutable(v2 Vec[T]) { v.X += v2.X; v.Y += v2.Y }

// SubMutable subtracts the given vector from the current vector, modifying the current vector in place.
//
// Parameters:
//   - v2: The vector to subtract from the current vector.
//
// Effect:
//   - Updates the X and Y components of the current vector by subtracting the corresponding components of v2.
func (v *Vec[T]) SubMutable(v2 Vec[T]) { v.X -= v2.X; v.Y -= v2.Y }

// Equals checks if the current vector is equal to another vector.
// It compares the X and Y components of both vectors for equality.
//
// Parameters:
//   - v2: The vector to compare with the current vector.
//
// Returns:
//   - true if both vectors have the same X and Y values, false otherwise.
func (v Vec[T]) Equals(v2 Vec[T]) bool { return v.X == v2.X && v.Y == v2.Y }

// String returns a string representation of the vector.
// The format of the string is "(X, Y)", where X and Y are the vector's components.
//
// Returns:
//   - A string in the format "(X, Y)" representing the vector.
func (v Vec[T]) String() string { return fmt.Sprintf("(%v,%v)", v.X, v.Y) }

// Bounds returns the zero-area rectangle representing the point.
func (v Vec[T]) Bounds() AABB[T] {
	return BuildAABB(v, 0)
}

// Vertices returns the address of the vector so callers can mutate it in place.
func (v *Vec[T]) Vertices() []*Vec[T] {
	if v == nil {
		return nil
	}
	return []*Vec[T]{v}
}

func (v Vec[T]) Fragments() []Shape[T] { return nil }

func (v Vec[T]) SetFragments(_ []Shape[T]) {}

func (v Vec[T]) Clone() Shape[T] {
	copy := v
	return &copy
}

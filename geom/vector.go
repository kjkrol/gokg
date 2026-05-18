// Package geometry provides a set of operations for 2D graphics
package geom

import (
	"fmt"
)

func NewVec[T Numeric](X, Y T) Vec[T] {
	return Vec[T]{X, Y}
}

// Add returns a new vector that is the sum of v and v2.
func (v Vec[T]) Add(v2 Vec[T]) Vec[T] { return Vec[T]{v.X + v2.X, v.Y + v2.Y} }

// Sub returns a new vector that subtracts v2 from v.
func (v Vec[T]) Sub(v2 Vec[T]) Vec[T] { return Vec[T]{v.X - v2.X, v.Y - v2.Y} }

// AddMutable adds v2 to v in place.
func (v *Vec[T]) AddMutable(v2 Vec[T]) { v.X += v2.X; v.Y += v2.Y }

// SubMutable subtracts v2 from v in place.
func (v *Vec[T]) SubMutable(v2 Vec[T]) { v.X -= v2.X; v.Y -= v2.Y }

// Invert directions
func (v *Vec[T]) Invert() { v.X = -v.X; v.Y = -v.Y }

// Multiply by factor
func (v *Vec[T]) Multiply(factor T) { v.X = v.X * factor; v.Y = v.Y * factor }

// Equals reports whether v and v2 have the same components.
func (v Vec[T]) Equals(v2 Vec[T]) bool { return v.X == v2.X && v.Y == v2.Y }

// String formats v as "(X,Y)".
func (v Vec[T]) String() string { return fmt.Sprintf("(%v,%v)", v.X, v.Y) }

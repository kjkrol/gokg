package geometry

import (
	"fmt"
)

type Number interface{ int | float64 }

// ---- Generic 2d vector ----
type Vec[T Number] struct{ X, Y T }

func (v Vec[T]) Add(v2 Vec[T]) Vec[T] { return Vec[T]{v.X + v2.X, v.Y + v2.Y} }
func (v Vec[T]) Sub(v2 Vec[T]) Vec[T] { return Vec[T]{v.X - v2.X, v.Y - v2.Y} }

func (v *Vec[T]) AddMutable(v2 Vec[T]) { v.X += v2.X; v.Y += v2.Y }
func (v *Vec[T]) SubMutable(v2 Vec[T]) { v.X -= v2.X; v.Y -= v2.Y }

func (v Vec[T]) Equals(v2 Vec[T]) bool { return v.X == v2.X && v.Y == v2.Y }
func (v Vec[T]) String() string        { return fmt.Sprintf("(%v,%v)", v.X, v.Y) }

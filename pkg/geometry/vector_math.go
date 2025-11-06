// Package geometry provides helpers for vector math and spatial bounds.
package geometry

import "math"

// VectorMath exposes vector operations for numeric components.
type VectorMath[T int | float64] interface {
	// Length returns the Euclidean magnitude of v.
	Length(v Vec[T]) T
	// Clamp bounds v1 to the inclusive box [0,size].
	Clamp(v1 *Vec[T], size Vec[T])
	// Wrap folds v1 back into [0,size) using modulo semantics appropriate for T.
	Wrap(v1 *Vec[T], size Vec[T])
	OverlapEpsilon() T
}

var (
	// FLOAT_64_VEC_MATH provides vector operations for float64 components.
	FLOAT_64_VEC_MATH VectorMath[float64] = float64VectorMath{}
	// INT_VEC_MATH provides vector operations for int components.
	INT_VEC_MATH VectorMath[int] = intVectorMath{}
)

// VectorMathByType returns the VectorMath implementation matching the generic type T.
func VectorMathByType[T SupportedNumeric]() VectorMath[T] {
	var zero T
	if _, ok := any(zero).(float64); ok {
		return any(FLOAT_64_VEC_MATH).(VectorMath[T])
	}
	return any(INT_VEC_MATH).(VectorMath[T])
}

// -----------------------------------------------------------------------------

type float64VectorMath struct{}

func (m float64VectorMath) Length(v Vec[float64]) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (m float64VectorMath) Clamp(v *Vec[float64], size Vec[float64]) { clampClosed(v, size) }

func (m float64VectorMath) Wrap(v *Vec[float64], size Vec[float64]) {
	modMutable := func(v1 *Vec[float64], v2 Vec[float64]) {
		if v2.X != 0 {
			v1.X = math.Mod(v1.X, v2.X)
		}
		if v2.Y != 0 {
			v1.Y = math.Mod(v1.Y, v2.Y)
		}

	}
	wrap(v, size, modMutable)
}
func (m float64VectorMath) OverlapEpsilon() float64 { return 1e-9 }

// -----------------------------------------------------------------------------

type intVectorMath struct{}

func (m intVectorMath) Length(v Vec[int]) int {
	return int(math.Ceil(math.Sqrt(float64(v.X*v.X + v.Y*v.Y))))
}

func (m intVectorMath) Clamp(v *Vec[int], size Vec[int]) { clampClosed(v, size) }

func (m intVectorMath) Wrap(v *Vec[int], size Vec[int]) {
	modMutable := func(v1 *Vec[int], v2 Vec[int]) {
		if v2.X != 0 {
			v1.X %= v2.X
		}
		if v2.Y != 0 {
			v1.Y %= v2.Y
		}
	}
	wrap(v, size, modMutable)
}

func (m intVectorMath) OverlapEpsilon() int { return 0 }

//-----------------------------------------------------------------------------

func clampClosed[T SupportedNumeric](v *Vec[T], bounds Vec[T]) {
	if v.X > bounds.X {
		v.X = bounds.X
	} else if v.X < 0 {
		v.X = 0
	}
	if v.Y > bounds.Y {
		v.Y = bounds.Y
	} else if v.Y < 0 {
		v.Y = 0
	}
}

func wrap[T SupportedNumeric](v *Vec[T], bounds Vec[T], modMutable func(*Vec[T], Vec[T])) {
	modMutable(v, bounds)
	v.AddMutable(bounds)
	modMutable(v, bounds)
}

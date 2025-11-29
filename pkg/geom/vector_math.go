// Package geometry provides helpers for vector math and spatial bounds.
package geom

import (
	"fmt"
	"math"
)

// VectorMath exposes vector operations for numeric components.
type (
	VectorMath[T Numeric] interface {
		// Length returns the Euclidean magnitude of v.
		Length(v Vec[T]) T
		// Clamp bounds v1 to the inclusive box [0,size].
		Clamp(v1 *Vec[T], size Vec[T])
		// Wrap folds v1 back into [0,size) using modulo semantics appropriate for T.
		Wrap(v1 *Vec[T], size Vec[T])
		OverlapEpsilon() T
		Sub(v1, v2 Vec[T]) Vec[T]
	}

	float64VectorMath = floatVectorMath[float64]
	intVectorMath     = signedIntVectorMath[int]
	int64VectorMath   = signedIntVectorMath[int64]
	uint32VectorMath  = unsignedIntVectorMath[uint32]
)

var (
	FLOAT_64_VEC_MATH VectorMath[float64] = float64VectorMath{}
	INT_VEC_MATH      VectorMath[int]     = intVectorMath{}
	INT64_VEC_MATH    VectorMath[int64]   = int64VectorMath{}
	UINT32_VEC_MATH   VectorMath[uint32]  = uint32VectorMath{}
)

// VectorMathByType returns the VectorMath implementation matching the generic type T.
func VectorMathByType[T Numeric]() VectorMath[T] {
	var zero T
	switch any(zero).(type) {
	case float64:
		return any(FLOAT_64_VEC_MATH).(VectorMath[T])
	case int:
		return any(INT_VEC_MATH).(VectorMath[T])
	case int64:
		return any(INT64_VEC_MATH).(VectorMath[T])
	case uint32:
		return any(UINT32_VEC_MATH).(VectorMath[T])
	default:
		panic(fmt.Sprintf("no VectorMath implementation for %T", zero))
	}
}

// -----------------------------------------------------------------------------

type floatVectorMath[T Floating] struct{}

func (m floatVectorMath[T]) Length(v Vec[T]) T {
	return T(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

func (m floatVectorMath[T]) Clamp(v *Vec[T], size Vec[T]) { clampClosed(v, size) }

func (m floatVectorMath[T]) Wrap(v *Vec[T], size Vec[T]) {
	modMutable := func(v1 *Vec[T], v2 Vec[T]) {
		if v2.X != 0 {
			v1.X = T(math.Mod(float64(v1.X), float64(v2.X)))
		}
		if v2.Y != 0 {
			v1.Y = T(math.Mod(float64(v1.Y), float64(v2.Y)))
		}

	}
	wrap(v, size, modMutable)
}
func (m floatVectorMath[T]) OverlapEpsilon() T { return 1e-9 }

func (m floatVectorMath[T]) Sub(v1, v2 Vec[T]) Vec[T] {
	return v1.Sub(v2)
}

// -----------------------------------------------------------------------------

type signedIntVectorMath[T SignedInt] struct{}

func (m signedIntVectorMath[T]) Length(v Vec[T]) T {
	return T(math.Ceil(math.Sqrt(float64(v.X*v.X + v.Y*v.Y))))
}

func (m signedIntVectorMath[T]) Clamp(v *Vec[T], size Vec[T]) { clampClosed(v, size) }

func (m signedIntVectorMath[T]) Wrap(v *Vec[T], size Vec[T]) {
	modMutable := func(v1 *Vec[T], v2 Vec[T]) {
		if v2.X != 0 {
			v1.X %= v2.X
		}
		if v2.Y != 0 {
			v1.Y %= v2.Y
		}
	}
	wrap(v, size, modMutable)
}

func (m signedIntVectorMath[T]) OverlapEpsilon() T { return 0 }

func (m signedIntVectorMath[T]) Sub(v1, v2 Vec[T]) Vec[T] { return v1.Sub(v2) }

//-----------------------------------------------------------------------------

type unsignedIntVectorMath[T UnsignedInt] struct{}

func (m unsignedIntVectorMath[T]) Length(v Vec[T]) T {
	return T(math.Ceil(math.Sqrt(float64(v.X*v.X + v.Y*v.Y))))
}

func (m unsignedIntVectorMath[T]) Clamp(v *Vec[T], size Vec[T]) {
	if v.X > size.X {
		v.X = size.X
	}
	if v.Y > size.Y {
		v.Y = size.Y
	}
}

func (m unsignedIntVectorMath[T]) Wrap(v *Vec[T], size Vec[T]) {
	signed := Vec[int64]{int64(v.X), int64(v.Y)}
	bounds := Vec[int64]{int64(size.X), int64(size.Y)}

	modMutable := func(v1 *Vec[int64], v2 Vec[int64]) {
		if v2.X != 0 {
			v1.X %= v2.X
		}
		if v2.Y != 0 {
			v1.Y %= v2.Y
		}
	}

	signed.AddMutable(bounds)
	modMutable(&signed, bounds)

	v.X = T(signed.X)
	v.Y = T(signed.Y)
}

func (m unsignedIntVectorMath[T]) OverlapEpsilon() T { return 0 }

func (m unsignedIntVectorMath[T]) Sub(v1, v2 Vec[T]) Vec[T] {
	sx := int64(v1.X) - int64(v2.X)
	sy := int64(v1.Y) - int64(v2.Y)
	return NewVec(T(sx), T(sy))
}

//-----------------------------------------------------------------------------

func clampClosed[T Numeric](v *Vec[T], bounds Vec[T]) {
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

func wrap[T Numeric](v *Vec[T], bounds Vec[T], modMutable func(*Vec[T], Vec[T])) {
	modMutable(v, bounds)
	v.AddMutable(bounds)
	modMutable(v, bounds)
}

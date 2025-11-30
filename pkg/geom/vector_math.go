// Package geometry provides helpers for vector math and spatial bounds.
package geom

import (
	"fmt"
	"math"
)

// VectorMath exposes vector operations for numeric components.
type (
	SignedInt   interface{ ~int | ~int64 }
	UnsignedInt interface{ ~uint32 }
	Floating    interface{ ~float32 | ~float64 }
	Numeric     interface {
		SignedInt | UnsignedInt | Floating
	}

	// Vec is a generic 2D vector with components X and Y.
	//
	// Type Parameters:
	//   - T: A numeric type that satisfies the Number constraint.
	Vec[T Numeric] struct{ X, Y T }

	VectorMath[T Numeric] interface {
		// Length returns the Euclidean magnitude of v.
		Length(v Vec[T]) T
		// Clamp bounds v1 to the inclusive box [0,size].
		Clamp(v1 *Vec[T], size Vec[T])
		// Wrap folds v1 back into [0,size) using modulo semantics appropriate for T.
		Wrap(v1 *Vec[T], size Vec[T])
		Sub(v1, v2 Vec[T]) Vec[T]
	}
)

const eps = 1e-9

var (
	float32VecMath = FloatVectorMath[float32]{}
	float64VecMath = FloatVectorMath[float64]{}
	intVecMath     = SignedIntVectorMath[int]{}
	int64VecMath   = SignedIntVectorMath[int64]{}
	uint32VecMath  = UnsignedIntVectorMath[uint32]{}
)

func VectorMathByType[T Numeric]() VectorMath[T] {
	var zero T

	switch any(zero).(type) {
	case float32:
		return any(float32VecMath).(VectorMath[T])
	case float64:
		return any(float64VecMath).(VectorMath[T])
	case int:
		return any(intVecMath).(VectorMath[T])
	case int64:
		return any(int64VecMath).(VectorMath[T])
	case uint32:
		return any(uint32VecMath).(VectorMath[T])
	default:
		panic(fmt.Sprintf("no VectorMath implementation for %T", zero))
	}
}

// -----------------------------------------------------------------------------

type FloatVectorMath[T Floating] struct{}

func (m FloatVectorMath[T]) Length(v Vec[T]) T {
	return T(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

func (m FloatVectorMath[T]) Clamp(v *Vec[T], size Vec[T]) { clampClosed(v, size) }

func (m FloatVectorMath[T]) Wrap(v *Vec[T], size Vec[T]) {
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

func (m FloatVectorMath[T]) Sub(v1, v2 Vec[T]) Vec[T] {
	return v1.Sub(v2)
}

// -----------------------------------------------------------------------------

type SignedIntVectorMath[T SignedInt] struct{}

func (m SignedIntVectorMath[T]) Length(v Vec[T]) T {
	return T(math.Ceil(math.Sqrt(float64(v.X*v.X + v.Y*v.Y))))
}

func (m SignedIntVectorMath[T]) Clamp(v *Vec[T], size Vec[T]) { clampClosed(v, size) }

func (m SignedIntVectorMath[T]) Wrap(v *Vec[T], size Vec[T]) {
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

func (m SignedIntVectorMath[T]) Sub(v1, v2 Vec[T]) Vec[T] { return v1.Sub(v2) }

//-----------------------------------------------------------------------------

type UnsignedIntVectorMath[T UnsignedInt] struct{}

func (m UnsignedIntVectorMath[T]) Length(v Vec[T]) T {
	return T(math.Ceil(math.Sqrt(float64(v.X*v.X + v.Y*v.Y))))
}

func (m UnsignedIntVectorMath[T]) Clamp(v *Vec[T], size Vec[T]) {
	// odczytaj komponenty jako signed, żeby podwijane ujemne wartości nie lądowały na „9”
	sx := int64(int32(v.X))
	sy := int64(int32(v.Y))

	clamp := func(val, max int64) int64 {
		if val < 0 {
			return 0
		}
		if val > max {
			return max
		}
		return val
	}

	v.X = T(clamp(sx, int64(size.X)))
	v.Y = T(clamp(sy, int64(size.Y)))
}

func (m UnsignedIntVectorMath[T]) Wrap(v *Vec[T], size Vec[T]) {
	// reinterpretuj komponenty jako signed (int32 -> int64), żeby -8 nie stało się 4294967288
	signed := Vec[int64]{int64(int32(v.X)), int64(int32(v.Y))}
	bounds := Vec[int64]{int64(size.X), int64(size.Y)}

	// modulo dodatnie pozwala znormalizować także duże wartości ujemne (np. -17 przy size=9 -> 1)
	if bounds.X != 0 {
		signed.X = ((signed.X % bounds.X) + bounds.X) % bounds.X
	}
	if bounds.Y != 0 {
		signed.Y = ((signed.Y % bounds.Y) + bounds.Y) % bounds.Y
	}

	v.X = T(signed.X)
	v.Y = T(signed.Y)
}

func (m UnsignedIntVectorMath[T]) Sub(v1, v2 Vec[T]) Vec[T] {
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

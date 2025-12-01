// Package geometry provides helpers for vector math and spatial bounds.
package geom

import (
	"fmt"
	"math"
)

// VectorMath exposes vector operations for numeric components.
type (
	SignedInt   interface{ ~int | ~int64 | ~int32 }
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
		// Clamp bounds v1 to the inclusive box [0,size] and returns the bounded vector.
		Clamp(v1 Vec[T], size Vec[T]) Vec[T]
		// Wrap folds v1 back into [0,size) using modulo semantics appropriate for T and returns the wrapped vector.
		Wrap(v1 Vec[T], size Vec[T]) Vec[T]
		Sub(v1, v2 Vec[T]) Vec[T]
	}
)

const eps = 1e-9

var (
	float32VecMath = FloatVectorMath[float32]{}
	float64VecMath = FloatVectorMath[float64]{}
	intVecMath     = SignedIntVectorMath[int]{}
	int32VecMath   = SignedIntVectorMath[int32]{}
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
	case int32:
		return any(int32VecMath).(VectorMath[T])
	case uint32:
		return any(uint32VecMath).(VectorMath[T])
	default:
		panic(fmt.Sprintf("no VectorMath implementation for %T", zero))
	}
}

// -----------------------------------------------------------------------------

type FloatVectorMath[T Floating] struct{}

func (m FloatVectorMath[T]) Length(v Vec[T]) T {
	l := lengthFloat64XY(v.X, v.Y)
	return T(l)
}

func (m FloatVectorMath[T]) Clamp(v Vec[T], size Vec[T]) Vec[T] {
	return Vec[T]{
		X: clampSigned(v.X, size.X),
		Y: clampSigned(v.Y, size.Y),
	}
}

func (m FloatVectorMath[T]) Wrap(v Vec[T], size Vec[T]) Vec[T] {
	return Vec[T]{
		X: wrapFloat(v.X, size.X),
		Y: wrapFloat(v.Y, size.Y),
	}
}

func (m FloatVectorMath[T]) Sub(v1, v2 Vec[T]) Vec[T] {
	return v1.Sub(v2)
}

// -----------------------------------------------------------------------------

type SignedIntVectorMath[T SignedInt] struct{}

func (m SignedIntVectorMath[T]) Length(v Vec[T]) T {
	l := lengthFloat64XY(v.X, v.Y)
	return T(math.Ceil(l))
}

func (m SignedIntVectorMath[T]) Clamp(v Vec[T], size Vec[T]) Vec[T] {
	return Vec[T]{
		X: clampSigned(v.X, size.X),
		Y: clampSigned(v.Y, size.Y),
	}
}

func (m SignedIntVectorMath[T]) Wrap(v Vec[T], size Vec[T]) Vec[T] {
	return Vec[T]{
		X: wrapSigned(v.X, size.X),
		Y: wrapSigned(v.Y, size.Y),
	}
}

func (m SignedIntVectorMath[T]) Sub(v1, v2 Vec[T]) Vec[T] { return v1.Sub(v2) }

//-----------------------------------------------------------------------------

type UnsignedIntVectorMath[T UnsignedInt] struct{}

func (m UnsignedIntVectorMath[T]) Length(v Vec[T]) T {
	l := lengthFloat64XY(v.X, v.Y)
	return T(math.Ceil(l))
}

func (m UnsignedIntVectorMath[T]) Clamp(v Vec[T], size Vec[T]) Vec[T] {
	// odczytaj komponenty jako signed (int32 -> int64),
	// żeby np. 0xFFFF_FFF8 traktować jako -8, a nie 4_294_967_288
	sx := int64(int32(v.X))
	sy := int64(int32(v.Y))

	maxX := int64(size.X)
	maxY := int64(size.Y)

	sx = clampSigned(sx, maxX)
	sy = clampSigned(sy, maxY)

	return Vec[T]{T(sx), T(sy)}
}

func (m UnsignedIntVectorMath[T]) Wrap(v Vec[T], size Vec[T]) Vec[T] {
	// reinterpretacja do signed (int32 -> int64)
	signed := Vec[int64]{
		int64(int32(v.X)),
		int64(int32(v.Y)),
	}
	bounds := Vec[int64]{
		int64(size.X),
		int64(size.Y),
	}

	signed.X = wrapSigned(signed.X, bounds.X)
	signed.Y = wrapSigned(signed.Y, bounds.Y)

	return Vec[T]{T(signed.X), T(signed.Y)}
}

func (m UnsignedIntVectorMath[T]) Sub(v1, v2 Vec[T]) Vec[T] {
	sx := int64(v1.X) - int64(v2.X)
	sy := int64(v1.Y) - int64(v2.Y)
	return NewVec(T(sx), T(sy))
}

//-----------------------------------------------------------------------------

func clampSigned[T SignedInt | Floating](val, max T) T {
	if val > max {
		return max
	}
	if val < 0 {
		return 0
	}
	return val
}

func wrapSigned[T SignedInt](val, max T) T {
	if max == 0 {
		return val
	}

	if max < 0 {
		max = -max
	}

	// szybka ścieżka: już w [0, max)
	if val >= 0 && val < max {
		return val
	}

	// szybka ścieżka dla potęgi dwójki
	if max&(max-1) == 0 {
		mask := max - 1
		return val & mask
	}

	// ogólny przypadek: jedno modulo + korekta ujemnego wyniku
	r := val % max
	if r < 0 {
		r += max
	}
	return r
}

func wrapFloat[T Floating](val, max T) T {
	if max == 0 {
		return val
	}
	if max < 0 {
		max = -max
	}

	// szybka ścieżka: już w [0, max)
	if val >= 0 && val < max {
		return val
	}

	// dodatkowa szybka ścieżka: niewielkie odchylenie
	twice := max + max
	if val >= 0 && val < twice {
		if val >= max {
			return val - max
		}
		return val
	}
	if val < 0 {
		if val >= -max {
			return val + max
		}
		// dla (-2*max, -max) trzeba dodać 2*max albo spaść do math.Mod
	}

	// ogólny przypadek
	r := T(math.Mod(float64(val), float64(max)))
	if r < 0 {
		r += max
	}
	return r
}

func lengthFloat64XY[T Numeric](x, y T) float64 {
	dx := float64(x)
	dy := float64(y)
	return math.Sqrt(dx*dx + dy*dy)
}

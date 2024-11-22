// VectorMath is a generic interface that defines mathematical operations for vectors.
// The type parameter T can be either int or float64.
//
// Methods:
//   - Length(v Vec[T]) T: Calculates the length (magnitude) of the vector v.
//   - Clamp(v1 *Vec[T], size Vec[T]): Clamps the vector v1 within the bounds defined by size.
//   - Wrap(v1 *Vec[T], size Vec[T]): Wraps the vector v1 within the bounds defined by size.
package geometry

import "math"

type VectorMath[T int | float64] interface {
	Length(v Vec[T]) T
	Clamp(v1 *Vec[T], size Vec[T])
	Wrap(v1 *Vec[T], size Vec[T])
}

// FLOAT_64_VEC_MATH is an instance of float64VectorMath that provides
// various mathematical operations for vectors with float64 components.
var (
	FLOAT_64_VEC_MATH VectorMath[float64] = float64VectorMath{}
	INT_VEC_MATH      VectorMath[int]     = intVectorMath{}
)

func VectorMathByType[T SupportedNumeric]() VectorMath[T] {
	switch any(*new(T)).(type) {
	case float64:
		return any(FLOAT_64_VEC_MATH).(VectorMath[T])
	case int:
		return any(INT_VEC_MATH).(VectorMath[T])
	default:
		panic("unsupported type")
	}
}

// -----------------------------------------------------------------------------

type float64VectorMath struct{}

func (m float64VectorMath) Length(v Vec[float64]) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (m float64VectorMath) Clamp(v *Vec[float64], size Vec[float64]) { clamp(v, size, 0.0001) }

func (m float64VectorMath) Wrap(v *Vec[float64], size Vec[float64]) {
	modMutable := func(v1 *Vec[float64], v2 Vec[float64]) { v1.X, v1.Y = math.Mod(v1.X, v2.X), math.Mod(v1.Y, v2.Y) }
	wrap(v, size, modMutable)
}

// -----------------------------------------------------------------------------

type intVectorMath struct{}

func (m intVectorMath) Length(v Vec[int]) int {
	return int(math.Ceil(math.Sqrt(float64(v.X*v.X + v.Y*v.Y))))
}

func (m intVectorMath) Clamp(v *Vec[int], size Vec[int]) { clamp(v, size, 1) }

func (m intVectorMath) Wrap(v *Vec[int], size Vec[int]) {
	modMutable := func(v1 *Vec[int], v2 Vec[int]) { v1.X %= v2.X; v1.Y %= v2.Y }
	wrap(v, size, modMutable)
}

//-----------------------------------------------------------------------------

func clamp[T SupportedNumeric](v *Vec[T], size Vec[T], delta T) {
	if v.X > size.X-delta {
		v.X = size.X - delta
	} else if v.X < 0 {
		v.X = 0
	}
	if v.Y > size.Y-delta {
		v.Y = size.Y - delta
	} else if v.Y < 0 {
		v.Y = 0
	}
}

func wrap[T SupportedNumeric](v *Vec[T], size Vec[T], modMutable func(*Vec[T], Vec[T])) {
	modMutable(v, size)
	v.AddMutable(size)
	modMutable(v, size)
}

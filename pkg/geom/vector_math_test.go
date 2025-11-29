package geom

import (
	"testing"
)

func TestVectorMath_Length(t *testing.T) {
	runLengthTest(t, "int", INT_VEC_MATH)
	runLengthTest(t, "uint32", UINT32_VEC_MATH)
	runLengthTest(t, "float64", FLOAT_64_VEC_MATH)
}

func runLengthTest[T Numeric](t *testing.T, name string, math VectorMath[T]) {
	t.Run(name, func(t *testing.T) {
		v := NewVec(T(3), T(4))
		expected := T(5)

		if result := math.Length(v); result != expected {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})
}

func TestVectorMath_Clamp(t *testing.T) {
	runClampTest(t, "int", INT_VEC_MATH)
	runClampTest(t, "uint32", UINT32_VEC_MATH)
	runClampTest(t, "float64", FLOAT_64_VEC_MATH)
}

func runClampTest[T Numeric](t *testing.T, name string, math VectorMath[T]) {
	t.Run(name, func(t *testing.T) {
		v := NewVec(T(5), T(7))
		bounds := NewVec(T(4), T(6))
		expected := NewVec(T(4), T(6))

		math.Clamp(&v, bounds)
		if v != expected {
			t.Errorf("expected %v, got %v", expected, v)
		}
	})
}

func TestVectorMath_Wrap(t *testing.T) {
	runWrapTest(t, "int", INT_VEC_MATH)
	runWrapTest(t, "uint32", UINT32_VEC_MATH)
	runWrapTest(t, "float64", FLOAT_64_VEC_MATH)
}

func runWrapTest[T Numeric](t *testing.T, name string, math VectorMath[T]) {
	t.Run(name, func(t *testing.T) {
		v := NewVec(T(5), T(7))
		bounds := NewVec(T(4), T(4))
		cases := []struct {
			offset   Vec[T]
			expected Vec[T]
		}{
			{NewVec(T(4), T(6)), NewVec(T(1), T(1))},
			{NewVec(bounds.X, 0), NewVec(T(1), T(7))},
			{NewVec(0, bounds.Y), NewVec(T(5), T(3))},
			{NewVec(bounds.X, bounds.Y), NewVec(T(1), T(3))},
		}

		for _, c := range cases {
			vec := NewVec(v.X, v.Y)
			math.Wrap(&vec, c.offset)
			if vec != c.expected {
				t.Errorf("expected %v, got %v", c.expected, vec)
			}
		}
	})
}

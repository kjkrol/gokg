package geom

import (
	"testing"
)

func TestVectorMath_Length(t *testing.T) {
	runLengthTest(t, "int", SignedIntVectorMath[int]{})
	runLengthTest(t, "uint32", UnsignedIntVectorMath[uint32]{})
	runLengthTest(t, "float64", FloatVectorMath[float64]{})
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
	runClampTest(t, "int", SignedIntVectorMath[int]{})
	runClampTest(t, "uint32", UnsignedIntVectorMath[uint32]{})
	runClampTest(t, "float64", FloatVectorMath[float64]{})
}

func runClampTest[T Numeric](t *testing.T, name string, math VectorMath[T]) {
	t.Run(name, func(t *testing.T) {
		bounds := NewVec(T(4), T(6))
		negX := int32(-3)
		negY := int32(-1)
		negStart := NewVec(T(negX), T(negY))
		cases := []struct {
			start    Vec[T]
			expected Vec[T]
		}{
			{NewVec(T(5), T(7)), NewVec(T(4), T(6))},
			{negStart, NewVec(T(0), T(0))}, // ujemne wartości muszą zostać przycięte do 0 także dla uint32
		}

		for _, c := range cases {
			got := math.Clamp(c.start, bounds)
			if got != c.expected {
				t.Errorf("expected %v, got %v", c.expected, got)
			}
		}
	})
}

func TestVectorMath_Wrap(t *testing.T) {
	runWrapTest(t, "int", SignedIntVectorMath[int]{})
	runWrapTest(t, "uint32", UnsignedIntVectorMath[uint32]{})
	runWrapTest(t, "float64", FloatVectorMath[float64]{})
}

func runWrapTest[T Numeric](t *testing.T, name string, math VectorMath[T]) {
	t.Run(name, func(t *testing.T) {
		v := NewVec(T(5), T(7))
		bounds := NewVec(T(4), T(4))
		negX := int32(-3)
		negY := int32(-1)
		negStart := NewVec(T(negX), T(negY))
		cases := []struct {
			offset   Vec[T]
			expected Vec[T]
			start    Vec[T]
		}{
			{NewVec(T(4), T(6)), NewVec(T(1), T(1)), Vec[T]{}},
			{NewVec(bounds.X, 0), NewVec(T(1), T(7)), Vec[T]{}},
			{NewVec(0, bounds.Y), NewVec(T(5), T(3)), Vec[T]{}},
			{NewVec(bounds.X, bounds.Y), NewVec(T(1), T(3)), Vec[T]{}},
			// ujemne wejście powinno się prawidłowo zawinąć dla wszystkich typów
			{bounds, NewVec(T(1), T(3)), negStart},
		}

		for _, c := range cases {
			vec := NewVec(v.X, v.Y)
			if c.start != (Vec[T]{}) {
				vec = c.start
			}
			got := math.Wrap(vec, c.offset)
			if got != c.expected {
				t.Errorf("expected %v, got %v", c.expected, got)
			}
		}
	})
}

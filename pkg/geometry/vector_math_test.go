package geometry

import (
	"testing"
)

func TestFloat64VectorMath_Length(t *testing.T) {
	v := Vec[float64]{X: 3, Y: 4}
	expected := 5.0
	result := FLOAT_64_VEC_MATH.Length(v)
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestIntVectorMath_Length(t *testing.T) {
	v := Vec[int]{X: 3, Y: 4}
	expected := 5
	result := INT_VEC_MATH.Length(v)
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestFloat64VectorMath_Clamp(t *testing.T) {
	v := Vec[float64]{X: 5, Y: 7}
	bounds := Vec[float64]{X: 4, Y: 6}
	expected := Vec[float64]{X: 3.9999, Y: 5.9999}
	FLOAT_64_VEC_MATH.Clamp(&v, bounds)
	if v != expected {
		t.Errorf("expected %v, got %v", expected, v)
	}
}

func TestIntVectorMath_Clamp(t *testing.T) {
	v := Vec[int]{X: 5, Y: 7}
	bounds := Vec[int]{X: 4, Y: 6}
	expected := Vec[int]{X: 3, Y: 5}
	INT_VEC_MATH.Clamp(&v, bounds)
	if v != expected {
		t.Errorf("expected %v, got %v", expected, v)
	}
}

func TestFloat64VectorMath_Wrap(t *testing.T) {
	v := Vec[float64]{X: 5, Y: 7}
	bounds := Vec[float64]{X: 4, Y: 4}
	cases := []struct {
		offset   Vec[float64]
		expected Vec[float64]
	}{
		{Vec[float64]{X: 4, Y: 6}, Vec[float64]{X: 1, Y: 1}},
		{Vec[float64]{X: bounds.X, Y: 0}, Vec[float64]{X: 1, Y: 7}},
		{Vec[float64]{X: 0, Y: bounds.Y}, Vec[float64]{X: 5, Y: 3}},
		{Vec[float64]{X: bounds.X, Y: bounds.Y}, Vec[float64]{X: 1, Y: 3}},
	}

	for _, c := range cases {
		vec := Vec[float64]{X: v.X, Y: v.Y}
		FLOAT_64_VEC_MATH.Wrap(&vec, c.offset)
		if vec != c.expected {
			t.Errorf("expected %v, got %v", c.expected, vec)
		}
	}
}

func TestIntVectorMath_Wrap(t *testing.T) {
	v := Vec[int]{X: 5, Y: 7}
	bounds := Vec[int]{X: 4, Y: 4}
	cases := []struct {
		offset   Vec[int]
		expected Vec[int]
	}{
		{Vec[int]{X: 4, Y: 6}, Vec[int]{X: 1, Y: 1}},
		{Vec[int]{X: bounds.X, Y: 0}, Vec[int]{X: 1, Y: 7}},
		{Vec[int]{X: 0, Y: bounds.Y}, Vec[int]{X: 5, Y: 3}},
		{Vec[int]{X: bounds.X, Y: bounds.Y}, Vec[int]{X: 1, Y: 3}},
	}

	for _, c := range cases {
		vec := Vec[int]{X: v.X, Y: v.Y}
		INT_VEC_MATH.Wrap(&vec, c.offset)
		if vec != c.expected {
			t.Errorf("expected %v, got %v", c.expected, vec)
		}
	}
}

package geometry

import (
	"testing"
)

func TestFloat64VectorMath_Length(t *testing.T) {
	v := vec[float64]{X: 3, Y: 4}
	expected := 5.0
	result := FLOAT_64_VEC_MATH.Length(v)
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestIntVectorMath_Length(t *testing.T) {
	v := vec[int]{X: 3, Y: 4}
	expected := 5
	result := INT_VEC_MATH.Length(v)
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestFloat64VectorMath_Clamp(t *testing.T) {
	v := vec[float64]{X: 5, Y: 7}
	bounds := vec[float64]{X: 4, Y: 6}
	expected := vec[float64]{X: 3.9999, Y: 5.9999}
	FLOAT_64_VEC_MATH.Clamp(&v, bounds)
	if v != expected {
		t.Errorf("expected %v, got %v", expected, v)
	}
}

func TestIntVectorMath_Clamp(t *testing.T) {
	v := vec[int]{X: 5, Y: 7}
	bounds := vec[int]{X: 4, Y: 6}
	expected := vec[int]{X: 3, Y: 5}
	INT_VEC_MATH.Clamp(&v, bounds)
	if v != expected {
		t.Errorf("expected %v, got %v", expected, v)
	}
}

func TestFloat64VectorMath_Wrap(t *testing.T) {
	v := vec[float64]{X: 5, Y: 7}
	bounds := vec[float64]{X: 4, Y: 4}
	cases := []struct {
		offset   vec[float64]
		expected vec[float64]
	}{
		{vec[float64]{X: 4, Y: 6}, vec[float64]{X: 1, Y: 1}},
		{vec[float64]{X: bounds.X, Y: 0}, vec[float64]{X: 1, Y: 7}},
		{vec[float64]{X: 0, Y: bounds.Y}, vec[float64]{X: 5, Y: 3}},
		{vec[float64]{X: bounds.X, Y: bounds.Y}, vec[float64]{X: 1, Y: 3}},
	}

	for _, c := range cases {
		vec := vec[float64]{X: v.X, Y: v.Y}
		FLOAT_64_VEC_MATH.Wrap(&vec, c.offset)
		if vec != c.expected {
			t.Errorf("expected %v, got %v", c.expected, vec)
		}
	}
}

func TestIntVectorMath_Wrap(t *testing.T) {
	v := vec[int]{X: 5, Y: 7}
	bounds := vec[int]{X: 4, Y: 4}
	cases := []struct {
		offset   vec[int]
		expected vec[int]
	}{
		{vec[int]{X: 4, Y: 6}, vec[int]{X: 1, Y: 1}},
		{vec[int]{X: bounds.X, Y: 0}, vec[int]{X: 1, Y: 7}},
		{vec[int]{X: 0, Y: bounds.Y}, vec[int]{X: 5, Y: 3}},
		{vec[int]{X: bounds.X, Y: bounds.Y}, vec[int]{X: 1, Y: 3}},
	}

	for _, c := range cases {
		vec := vec[int]{X: v.X, Y: v.Y}
		INT_VEC_MATH.Wrap(&vec, c.offset)
		if vec != c.expected {
			t.Errorf("expected %v, got %v", c.expected, vec)
		}
	}
}

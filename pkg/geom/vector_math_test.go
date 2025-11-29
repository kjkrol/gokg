package geom

import (
	"testing"
)

func TestFloat64VectorMath_Length(t *testing.T) {
	v := NewVec(3.0, 4.0)
	expected := 5.0
	result := FLOAT_64_VEC_MATH.Length(v)
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestIntVectorMath_Length(t *testing.T) {
	v := NewVec(3, 4)
	expected := 5
	result := INT_VEC_MATH.Length(v)
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestFloat64VectorMath_Clamp(t *testing.T) {
	v := NewVec(5.0, 7.0)
	bounds := NewVec(4.0, 6.0)
	expected := NewVec(4.0, 6.0)
	FLOAT_64_VEC_MATH.Clamp(&v, bounds)
	if v != expected {
		t.Errorf("expected %v, got %v", expected, v)
	}
}

func TestIntVectorMath_Clamp(t *testing.T) {
	v := NewVec(5, 7)
	bounds := NewVec(4, 6)
	expected := NewVec(4, 6)
	INT_VEC_MATH.Clamp(&v, bounds)
	if v != expected {
		t.Errorf("expected %v, got %v", expected, v)
	}
}

func TestFloat64VectorMath_Wrap(t *testing.T) {
	v := NewVec(5.0, 7.0)
	bounds := NewVec(4.0, 4.0)
	cases := []struct {
		offset   Vec[float64]
		expected Vec[float64]
	}{
		{NewVec(4.0, 6.0), NewVec(1.0, 1.0)},
		{NewVec(bounds.X, 0.0), NewVec(1.0, 7.0)},
		{NewVec(0.0, bounds.Y), NewVec(5.0, 3.0)},
		{NewVec(bounds.X, bounds.Y), NewVec(1.0, 3.0)},
	}

	for _, c := range cases {
		vec := NewVec(v.X, v.Y)
		FLOAT_64_VEC_MATH.Wrap(&vec, c.offset)
		if vec != c.expected {
			t.Errorf("expected %v, got %v", c.expected, vec)
		}
	}
}

func TestIntVectorMath_Wrap(t *testing.T) {
	v := NewVec(5, 7)
	bounds := NewVec(4, 4)
	cases := []struct {
		offset   Vec[int]
		expected Vec[int]
	}{
		{NewVec(4, 6), NewVec(1, 1)},
		{NewVec(bounds.X, 0), NewVec(1, 7)},
		{NewVec(0, bounds.Y), NewVec(5, 3)},
		{NewVec(bounds.X, bounds.Y), NewVec(1, 3)},
	}

	for _, c := range cases {
		vec := NewVec(v.X, v.Y)
		INT_VEC_MATH.Wrap(&vec, c.offset)
		if vec != c.expected {
			t.Errorf("expected %v, got %v", c.expected, vec)
		}
	}
}

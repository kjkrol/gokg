package geometry

import (
	"testing"

	s "github.com/kjkrol/gokg/pkg/geometry/spatial"
)

func TestDiscreteCyclicPlaneTranslate(t *testing.T) {
	plane := NewCyclicBoundedPlane(5, 5)
	for _, test := range []struct {
		arg1     vec[int]
		arg2     vec[int]
		expected vec[int]
	}{
		{s.NewVec(2, 3), s.NewVec(-1, -2), s.NewVec(1, 1)},
		{s.NewVec(1, 2), s.NewVec(-1, -2), s.NewVec(0, 0)},
		{s.NewVec(0, 0), s.NewVec(-4, -4), s.NewVec(1, 1)},
		{s.NewVec(4, 0), s.NewVec(-1, -0), s.NewVec(3, 0)},
		{s.NewVec(1, 0), s.NewVec(-4, -0), s.NewVec(2, 0)},
	} {
		if plane.Translate(&test.arg1, test.arg2); !test.arg1.Equals(test.expected) {
			t.Errorf("result %v not equal to expected %v", test.arg1, test.expected)
		}
	}
}

func TestDiscreteCyclicPlaneMetric(t *testing.T) {
	plane := NewCyclicBoundedPlane(9, 9)
	for _, test := range []struct {
		arg1     vec[int]
		arg2     vec[int]
		expected int
	}{
		{s.NewVec(1, 2), s.NewVec(2, 3), 2},
		{s.NewVec(1, 2), s.NewVec(1, 2), 0},
		{s.NewVec(0, 0), s.NewVec(1, 1), 2},
		{s.NewVec(0, 0), s.NewVec(2, 2), 3},
		{s.NewVec(0, 0), s.NewVec(8, 8), 2},
		{s.NewVec(0, 0), s.NewVec(9, 9), 0}, // vec(9,9) has been wrapped to vec(0,0)
	} {
		if output := plane.Metric(test.arg1, test.arg2); output != test.expected {
			t.Errorf("vectors: %v, %v, metric %v not equal to expected %v", test.arg1, test.arg2, output, test.expected)
		}
	}
}

func TestDiscreteBoundedPlaneTranslate(t *testing.T) {
	plane := NewBoundedPlane(9, 9)
	for _, test := range []struct {
		arg1     vec[int]
		arg2     vec[int]
		expected vec[int]
	}{
		{s.NewVec(2, 3), s.NewVec(-1, -2), s.NewVec(1, 1)},
		{s.NewVec(1, 2), s.NewVec(-1, -2), s.NewVec(0, 0)},
		{s.NewVec(0, 0), s.NewVec(-4, -4), s.NewVec(0, 0)},
		{s.NewVec(4, 0), s.NewVec(-1, -0), s.NewVec(3, 0)},
		{s.NewVec(6, 0), s.NewVec(-4, -0), s.NewVec(2, 0)},
	} {
		if plane.Translate(&test.arg1, test.arg2); !test.arg1.Equals(test.expected) {
			t.Errorf("result %v not equal to expected %v", test.arg1, test.expected)
		}
	}
}

func TestDiscreteBoundedPlaneMetric(t *testing.T) {
	plane := NewBoundedPlane(9, 9)
	for _, test := range []struct {
		arg1     vec[int]
		arg2     vec[int]
		expected int
	}{
		{s.NewVec(1, 2), s.NewVec(2, 3), 2},
		{s.NewVec(1, 2), s.NewVec(1, 2), 0},
		{s.NewVec(0, 0), s.NewVec(1, 1), 2},
		{s.NewVec(0, 0), s.NewVec(2, 2), 3},
		{s.NewVec(0, 0), s.NewVec(8, 8), 12},
		{s.NewVec(0, 0), s.NewVec(9, 9), 12}, // vec(9,9) has been clamped to vec(8,8)
	} {
		if output := plane.Metric(test.arg1, test.arg2); output != test.expected {
			t.Errorf("vectors: %v, %v, metric %v not equal to expected %v", test.arg1, test.arg2, output, test.expected)
		}
	}
}

// -----------------------------------------------------------------------------

func TestDiscreteCyclicPlaneTranslateFloat64(t *testing.T) {
	plane := NewCyclicBoundedPlane(5.0, 5.0)
	for _, test := range []struct {
		arg1     vec[float64]
		arg2     vec[float64]
		expected vec[float64]
	}{
		{s.NewVec(2.0, 3.0), s.NewVec(-1.0, -2.0), s.NewVec(1.0, 1.0)},
		{s.NewVec(1.0, 2.0), s.NewVec(-1.0, -2.0), s.NewVec(0.0, 0.0)},
		{s.NewVec(0.0, 0.0), s.NewVec(-4.0, -4.0), s.NewVec(1.0, 1.0)},
		{s.NewVec(4.0, 0.0), s.NewVec(-1.0, 0.0), s.NewVec(3.0, 0.0)},
		{s.NewVec(1.0, 0.0), s.NewVec(-4.0, 0.0), s.NewVec(2.0, 0.0)},
	} {
		if plane.Translate(&test.arg1, test.arg2); !test.arg1.Equals(test.expected) {
			t.Errorf("result %v not equal to expected %v", test.arg1, test.expected)
		}
	}
}

func TestDiscreteCyclicPlaneMetricFloat64(t *testing.T) {
	plane := NewCyclicBoundedPlane(9.0, 9.0)
	for _, test := range []struct {
		arg1     vec[float64]
		arg2     vec[float64]
		expected float64
	}{
		{s.NewVec(1.0, 2.0), s.NewVec(2.0, 3.0), 1.4142135623730951},
		{s.NewVec(1.0, 2.0), s.NewVec(1.0, 2.0), 0.0},
		{s.NewVec(0.0, 0.0), s.NewVec(1.0, 1.0), 1.4142135623730951},
		{s.NewVec(0.0, 0.0), s.NewVec(2.0, 2.0), 2.8284271247461903},
		{s.NewVec(0.0, 0.0), s.NewVec(8.0, 8.0), 1.4142135623730951},
		{s.NewVec(0.0, 0.0), s.NewVec(9.0, 9.0), 0.0}, // vec(9,9) has been wrapped to vec(0,0)
	} {
		if output := plane.Metric(test.arg1, test.arg2); output != test.expected {
			t.Errorf("vectors: %v, %v, metric %v not equal to expected %v", test.arg1, test.arg2, output, test.expected)
		}
	}
}

func TestDiscreteBoundedPlaneTranslateFloat64(t *testing.T) {
	plane := NewBoundedPlane(9.0, 9.0)
	for _, test := range []struct {
		arg1     vec[float64]
		arg2     vec[float64]
		expected vec[float64]
	}{
		{s.NewVec(2.0, 3.0), s.NewVec(-1.0, -2.0), s.NewVec(1.0, 1.0)},
		{s.NewVec(1.0, 2.0), s.NewVec(-1.0, -2.0), s.NewVec(0.0, 0.0)},
		{s.NewVec(0.0, 0.0), s.NewVec(-4.0, -4.0), s.NewVec(0.0, 0.0)},
		{s.NewVec(4.0, 0.0), s.NewVec(-1.0, 0.0), s.NewVec(3.0, 0.0)},
		{s.NewVec(6.0, 0.0), s.NewVec(-4.0, 0.0), s.NewVec(2.0, 0.0)},
	} {
		if plane.Translate(&test.arg1, test.arg2); !test.arg1.Equals(test.expected) {
			t.Errorf("result %v not equal to expected %v", test.arg1, test.expected)
		}
	}
}

func TestDiscreteBoundedPlaneMetricFloat64(t *testing.T) {
	plane := NewBoundedPlane(9.0, 9.0)
	for _, test := range []struct {
		arg1     vec[float64]
		arg2     vec[float64]
		expected float64
	}{
		{s.NewVec(1.0, 2.0), s.NewVec(2.0, 3.0), 1.4142135623730951},
		{s.NewVec(1.0, 2.0), s.NewVec(1.0, 2.0), 0.0},
		{s.NewVec(0.0, 0.0), s.NewVec(1.0, 1.0), 1.4142135623730951},
		{s.NewVec(0.0, 0.0), s.NewVec(2.0, 2.0), 2.8284271247461903},
		{s.NewVec(0.0, 0.0), s.NewVec(8.0, 8.0), 11.313708498984761},
		{s.NewVec(0.0, 0.0), s.NewVec(9.0, 9.0), 12.727780640001619}, // Vec(9,9) has been clamped to vec(8,8)
		{s.NewVec(0.0, 0.0), s.NewVec(8.5, 0.0), 8.5},
	} {
		if output := plane.Metric(test.arg1, test.arg2); output != test.expected {
			t.Errorf("vectors: %v, %v, metric %v not equal to expected %v", test.arg1, test.arg2, output, test.expected)
		}
	}
}

func TestPlaneNormalize(t *testing.T) {
	plane := NewCyclicBoundedPlane(5, 5)
	vec := s.NewVec(7, -2)
	plane.Normalize(&vec)
	expected := s.NewVec(2, 3)
	if vec != expected {
		t.Errorf("expected normalized vector %v, got %v", expected, vec)
	}
}

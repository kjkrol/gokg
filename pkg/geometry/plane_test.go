package geometry

import (
	"testing"
)

func TestDiscreteCyclicPlaneTranslate(t *testing.T) {
	plane := NewCyclicBoundedPlane(5, 5)
	for _, test := range []struct {
		arg1     Vec[int]
		arg2     Vec[int]
		expected Vec[int]
	}{
		{Vec[int]{2, 3}, Vec[int]{-1, -2}, Vec[int]{1, 1}},
		{Vec[int]{1, 2}, Vec[int]{-1, -2}, Vec[int]{0, 0}},
		{Vec[int]{0, 0}, Vec[int]{-4, -4}, Vec[int]{1, 1}},
		{Vec[int]{4, 0}, Vec[int]{-1, -0}, Vec[int]{3, 0}},
		{Vec[int]{1, 0}, Vec[int]{-4, -0}, Vec[int]{2, 0}},
	} {
		if plane.Translate(&test.arg1, test.arg2); !test.arg1.Equals(test.expected) {
			t.Errorf("result %v not equal to expected %v", test.arg1, test.expected)
		}
	}
}

func TestDiscreteCyclicPlaneMetric(t *testing.T) {
	plane := NewCyclicBoundedPlane(9, 9)
	for _, test := range []struct {
		arg1     Vec[int]
		arg2     Vec[int]
		expected int
	}{
		{Vec[int]{1, 2}, Vec[int]{2, 3}, 2},
		{Vec[int]{1, 2}, Vec[int]{1, 2}, 0},
		{Vec[int]{0, 0}, Vec[int]{1, 1}, 2},
		{Vec[int]{0, 0}, Vec[int]{2, 2}, 3},
		{Vec[int]{0, 0}, Vec[int]{8, 8}, 2},
		{Vec[int]{0, 0}, Vec[int]{9, 9}, 0}, // vec(9,9) has been wrapped to vec(0,0)
	} {
		if output := plane.Metric(test.arg1, test.arg2); output != test.expected {
			t.Errorf("vectors: %v, %v, metric %v not equal to expected %v", test.arg1, test.arg2, output, test.expected)
		}
	}
}

func TestDiscreteBoundedPlaneTranslate(t *testing.T) {
	plane := NewBoundedPlane(9, 9)
	for _, test := range []struct {
		arg1     Vec[int]
		arg2     Vec[int]
		expected Vec[int]
	}{
		{Vec[int]{2, 3}, Vec[int]{-1, -2}, Vec[int]{1, 1}},
		{Vec[int]{1, 2}, Vec[int]{-1, -2}, Vec[int]{0, 0}},
		{Vec[int]{0, 0}, Vec[int]{-4, -4}, Vec[int]{0, 0}},
		{Vec[int]{4, 0}, Vec[int]{-1, -0}, Vec[int]{3, 0}},
		{Vec[int]{6, 0}, Vec[int]{-4, -0}, Vec[int]{2, 0}},
	} {
		if plane.Translate(&test.arg1, test.arg2); !test.arg1.Equals(test.expected) {
			t.Errorf("result %v not equal to expected %v", test.arg1, test.expected)
		}
	}
}

func TestDiscreteBoundedPlaneMetric(t *testing.T) {
	plane := NewBoundedPlane(9, 9)
	for _, test := range []struct {
		arg1     Vec[int]
		arg2     Vec[int]
		expected int
	}{
		{Vec[int]{1, 2}, Vec[int]{2, 3}, 2},
		{Vec[int]{1, 2}, Vec[int]{1, 2}, 0},
		{Vec[int]{0, 0}, Vec[int]{1, 1}, 2},
		{Vec[int]{0, 0}, Vec[int]{2, 2}, 3},
		{Vec[int]{0, 0}, Vec[int]{8, 8}, 12},
		{Vec[int]{0, 0}, Vec[int]{9, 9}, 12}, // vec(9,9) has been clamped to vec(8,8)
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
		arg1     Vec[float64]
		arg2     Vec[float64]
		expected Vec[float64]
	}{
		{Vec[float64]{2.0, 3.0}, Vec[float64]{-1.0, -2.0}, Vec[float64]{1.0, 1.0}},
		{Vec[float64]{1.0, 2.0}, Vec[float64]{-1.0, -2.0}, Vec[float64]{0.0, 0.0}},
		{Vec[float64]{0.0, 0.0}, Vec[float64]{-4.0, -4.0}, Vec[float64]{1.0, 1.0}},
		{Vec[float64]{4.0, 0.0}, Vec[float64]{-1.0, 0.0}, Vec[float64]{3.0, 0.0}},
		{Vec[float64]{1.0, 0.0}, Vec[float64]{-4.0, 0.0}, Vec[float64]{2.0, 0.0}},
	} {
		if plane.Translate(&test.arg1, test.arg2); !test.arg1.Equals(test.expected) {
			t.Errorf("result %v not equal to expected %v", test.arg1, test.expected)
		}
	}
}

func TestDiscreteCyclicPlaneMetricFloat64(t *testing.T) {
	plane := NewCyclicBoundedPlane(9.0, 9.0)
	for _, test := range []struct {
		arg1     Vec[float64]
		arg2     Vec[float64]
		expected float64
	}{
		{Vec[float64]{1.0, 2.0}, Vec[float64]{2.0, 3.0}, 1.4142135623730951},
		{Vec[float64]{1.0, 2.0}, Vec[float64]{1.0, 2.0}, 0.0},
		{Vec[float64]{0.0, 0.0}, Vec[float64]{1.0, 1.0}, 1.4142135623730951},
		{Vec[float64]{0.0, 0.0}, Vec[float64]{2.0, 2.0}, 2.8284271247461903},
		{Vec[float64]{0.0, 0.0}, Vec[float64]{8.0, 8.0}, 1.4142135623730951},
		{Vec[float64]{0.0, 0.0}, Vec[float64]{9.0, 9.0}, 0.0}, // vec(9,9) has been wrapped to vec(0,0)
	} {
		if output := plane.Metric(test.arg1, test.arg2); output != test.expected {
			t.Errorf("vectors: %v, %v, metric %v not equal to expected %v", test.arg1, test.arg2, output, test.expected)
		}
	}
}

func TestDiscreteBoundedPlaneTranslateFloat64(t *testing.T) {
	plane := NewBoundedPlane(9.0, 9.0)
	for _, test := range []struct {
		arg1     Vec[float64]
		arg2     Vec[float64]
		expected Vec[float64]
	}{
		{Vec[float64]{2.0, 3.0}, Vec[float64]{-1.0, -2.0}, Vec[float64]{1.0, 1.0}},
		{Vec[float64]{1.0, 2.0}, Vec[float64]{-1.0, -2.0}, Vec[float64]{0.0, 0.0}},
		{Vec[float64]{0.0, 0.0}, Vec[float64]{-4.0, -4.0}, Vec[float64]{0.0, 0.0}},
		{Vec[float64]{4.0, 0.0}, Vec[float64]{-1.0, 0.0}, Vec[float64]{3.0, 0.0}},
		{Vec[float64]{6.0, 0.0}, Vec[float64]{-4.0, 0.0}, Vec[float64]{2.0, 0.0}},
	} {
		if plane.Translate(&test.arg1, test.arg2); !test.arg1.Equals(test.expected) {
			t.Errorf("result %v not equal to expected %v", test.arg1, test.expected)
		}
	}
}

func TestDiscreteBoundedPlaneMetricFloat64(t *testing.T) {
	plane := NewBoundedPlane(9.0, 9.0)
	for _, test := range []struct {
		arg1     Vec[float64]
		arg2     Vec[float64]
		expected float64
	}{
		{Vec[float64]{1.0, 2.0}, Vec[float64]{2.0, 3.0}, 1.4142135623730951},
		{Vec[float64]{1.0, 2.0}, Vec[float64]{1.0, 2.0}, 0.0},
		{Vec[float64]{0.0, 0.0}, Vec[float64]{1.0, 1.0}, 1.4142135623730951},
		{Vec[float64]{0.0, 0.0}, Vec[float64]{2.0, 2.0}, 2.8284271247461903},
		{Vec[float64]{0.0, 0.0}, Vec[float64]{8.0, 8.0}, 11.313708498984761},
		{Vec[float64]{0.0, 0.0}, Vec[float64]{9.0, 9.0}, 12.727780640001619}, // vec(9,9) has been clamped to vec(8,8)
		{Vec[float64]{0.0, 0.0}, Vec[float64]{8.5, 0.0}, 8.5},
	} {
		if output := plane.Metric(test.arg1, test.arg2); output != test.expected {
			t.Errorf("vectors: %v, %v, metric %v not equal to expected %v", test.arg1, test.arg2, output, test.expected)
		}
	}
}

func TestPlaneNormalize(t *testing.T) {
	plane := NewCyclicBoundedPlane(5, 5)
	vec := Vec[int]{X: 7, Y: -2}
	plane.Normalize(&vec)
	expected := Vec[int]{X: 2, Y: 3}
	if vec != expected {
		t.Errorf("expected normalized vector %v, got %v", expected, vec)
	}
}

package geometry

import (
	"testing"
)

func TestDiscreteCyclicPlaneTranslate(t *testing.T) {
	plane := NewDiscreteCyclicBoundedPlane(5, 5)
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
	plane := NewDiscreteCyclicBoundedPlane(9, 9)
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
	plane := NewDiscreteBoundedPlane(9, 9)
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
	plane := NewDiscreteBoundedPlane(9, 9)
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

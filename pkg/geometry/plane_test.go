package geometry

import (
	"testing"
)

func TestDiscreteCyclicPlaneNormalizeAfterAdd(t *testing.T) {
	plane := NewCyclicBoundedPlane(5, 5)
	for _, test := range []struct {
		arg1     Vec[int]
		arg2     Vec[int]
		expected Vec[int]
	}{
		{NewVec(2, 3), NewVec(-1, -2), NewVec(1, 1)},
		{NewVec(1, 2), NewVec(-1, -2), NewVec(0, 0)},
		{NewVec(0, 0), NewVec(-4, -4), NewVec(1, 1)},
		{NewVec(4, 0), NewVec(-1, -0), NewVec(3, 0)},
		{NewVec(1, 0), NewVec(-4, -0), NewVec(2, 0)},
	} {
		result := test.arg1
		result.AddMutable(test.arg2)
		plane.normalize(&result)
		if !result.Equals(test.expected) {
			t.Errorf("result %v not equal to expected %v", result, test.expected)
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
		{NewVec(1, 2), NewVec(2, 3), 2},
		{NewVec(1, 2), NewVec(1, 2), 0},
		{NewVec(0, 0), NewVec(1, 1), 2},
		{NewVec(0, 0), NewVec(2, 2), 3},
		{NewVec(0, 0), NewVec(8, 8), 2},
		{NewVec(0, 0), NewVec(9, 9), 0}, // vec(9,9) has been wrapped to vec(0,0)
	} {
		if output := plane.Metric(test.arg1, test.arg2); output != test.expected {
			t.Errorf("vectors: %v, %v, metric %v not equal to expected %v", test.arg1, test.arg2, output, test.expected)
		}
	}
}

func TestDiscreteBoundedPlaneNormalizeAfterAdd(t *testing.T) {
	plane := NewBoundedPlane(9, 9)
	for _, test := range []struct {
		arg1     Vec[int]
		arg2     Vec[int]
		expected Vec[int]
	}{
		{NewVec(2, 3), NewVec(-1, -2), NewVec(1, 1)},
		{NewVec(1, 2), NewVec(-1, -2), NewVec(0, 0)},
		{NewVec(0, 0), NewVec(-4, -4), NewVec(0, 0)},
		{NewVec(4, 0), NewVec(-1, -0), NewVec(3, 0)},
		{NewVec(6, 0), NewVec(-4, -0), NewVec(2, 0)},
	} {
		result := test.arg1
		result.AddMutable(test.arg2)
		plane.normalize(&result)
		if !result.Equals(test.expected) {
			t.Errorf("result %v not equal to expected %v", result, test.expected)
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
		{NewVec(1, 2), NewVec(2, 3), 2},
		{NewVec(1, 2), NewVec(1, 2), 0},
		{NewVec(0, 0), NewVec(1, 1), 2},
		{NewVec(0, 0), NewVec(2, 2), 3},
		{NewVec(0, 0), NewVec(8, 8), 12},
		{NewVec(0, 0), NewVec(9, 9), 13}, // vec(9,9) stays on the boundary
	} {
		if output := plane.Metric(test.arg1, test.arg2); output != test.expected {
			t.Errorf("vectors: %v, %v, metric %v not equal to expected %v", test.arg1, test.arg2, output, test.expected)
		}
	}
}

// -----------------------------------------------------------------------------

func TestDiscreteCyclicPlaneNormalizeAfterAddFloat64(t *testing.T) {
	plane := NewCyclicBoundedPlane(5.0, 5.0)
	for _, test := range []struct {
		arg1     Vec[float64]
		arg2     Vec[float64]
		expected Vec[float64]
	}{
		{NewVec(2.0, 3.0), NewVec(-1.0, -2.0), NewVec(1.0, 1.0)},
		{NewVec(1.0, 2.0), NewVec(-1.0, -2.0), NewVec(0.0, 0.0)},
		{NewVec(0.0, 0.0), NewVec(-4.0, -4.0), NewVec(1.0, 1.0)},
		{NewVec(4.0, 0.0), NewVec(-1.0, 0.0), NewVec(3.0, 0.0)},
		{NewVec(1.0, 0.0), NewVec(-4.0, 0.0), NewVec(2.0, 0.0)},
	} {
		result := test.arg1
		result.AddMutable(test.arg2)
		plane.normalize(&result)
		if !result.Equals(test.expected) {
			t.Errorf("result %v not equal to expected %v", result, test.expected)
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
		{NewVec(1.0, 2.0), NewVec(2.0, 3.0), 1.4142135623730951},
		{NewVec(1.0, 2.0), NewVec(1.0, 2.0), 0.0},
		{NewVec(0.0, 0.0), NewVec(1.0, 1.0), 1.4142135623730951},
		{NewVec(0.0, 0.0), NewVec(2.0, 2.0), 2.8284271247461903},
		{NewVec(0.0, 0.0), NewVec(8.0, 8.0), 1.4142135623730951},
		{NewVec(0.0, 0.0), NewVec(9.0, 9.0), 0.0}, // vec(9,9) has been wrapped to vec(0,0)
	} {
		if output := plane.Metric(test.arg1, test.arg2); output != test.expected {
			t.Errorf("vectors: %v, %v, metric %v not equal to expected %v", test.arg1, test.arg2, output, test.expected)
		}
	}
}

func TestDiscreteBoundedPlaneNormalizeAfterAddFloat64(t *testing.T) {
	plane := NewBoundedPlane(9.0, 9.0)
	for _, test := range []struct {
		arg1     Vec[float64]
		arg2     Vec[float64]
		expected Vec[float64]
	}{
		{NewVec(2.0, 3.0), NewVec(-1.0, -2.0), NewVec(1.0, 1.0)},
		{NewVec(1.0, 2.0), NewVec(-1.0, -2.0), NewVec(0.0, 0.0)},
		{NewVec(0.0, 0.0), NewVec(-4.0, -4.0), NewVec(0.0, 0.0)},
		{NewVec(4.0, 0.0), NewVec(-1.0, 0.0), NewVec(3.0, 0.0)},
		{NewVec(6.0, 0.0), NewVec(-4.0, 0.0), NewVec(2.0, 0.0)},
	} {
		result := test.arg1
		result.AddMutable(test.arg2)
		plane.normalize(&result)
		if !result.Equals(test.expected) {
			t.Errorf("result %v not equal to expected %v", result, test.expected)
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
		{NewVec(1.0, 2.0), NewVec(2.0, 3.0), 1.4142135623730951},
		{NewVec(1.0, 2.0), NewVec(1.0, 2.0), 0.0},
		{NewVec(0.0, 0.0), NewVec(1.0, 1.0), 1.4142135623730951},
		{NewVec(0.0, 0.0), NewVec(2.0, 2.0), 2.8284271247461903},
		{NewVec(0.0, 0.0), NewVec(8.0, 8.0), 11.313708498984761},
		{NewVec(0.0, 0.0), NewVec(9.0, 9.0), 12.727922061357855}, // Vec(9,9) stays on the boundary
		{NewVec(0.0, 0.0), NewVec(8.5, 0.0), 8.5},
	} {
		if output := plane.Metric(test.arg1, test.arg2); output != test.expected {
			t.Errorf("vectors: %v, %v, metric %v not equal to expected %v", test.arg1, test.arg2, output, test.expected)
		}
	}
}

func TestPlaneNormalize(t *testing.T) {
	plane := NewCyclicBoundedPlane(5, 5)
	vec := NewVec(7, -2)
	plane.normalize(&vec)
	expected := NewVec(2, 3)
	if vec != expected {
		t.Errorf("expected normalized vector %v, got %v", expected, vec)
	}
}

func TestBoundedPlane_TransformBackAndForth(t *testing.T) {
	plane := NewBoundedPlane(10, 10)

	planeBox := newPlaneBox(NewVec(0, 0), 2, 2)

	shift := NewVec(2, 2)
	plane.Translate(&planeBox, shift)
	expectPlaneBoxState(t, planeBox, NewVec(2, 2), NewVec(4, 4), map[FragPosition][2]Vec[int]{})

	plane.Expand(&planeBox, 2)
	expectPlaneBoxState(t, planeBox, NewVec(0, 0), NewVec(6, 6), map[FragPosition][2]Vec[int]{})

	plane.Expand(&planeBox, -2)
	expectPlaneBoxState(t, planeBox, NewVec(2, 2), NewVec(4, 4), map[FragPosition][2]Vec[int]{})

	shift.Invert()
	plane.Translate(&planeBox, shift)
	expectPlaneBoxState(t, planeBox, NewVec(0, 0), NewVec(2, 2), map[FragPosition][2]Vec[int]{})
}

func TestCyclicPlane_TransformBackAndForth(t *testing.T) {
	plane := NewCyclicBoundedPlane(10, 10)

	planeBox := newPlaneBox(NewVec(0, 0), 2, 2)

	shift := NewVec(-1, -1)
	plane.Translate(&planeBox, shift)
	expectPlaneBoxState(t, planeBox, NewVec(9, 9), NewVec(10, 10), map[FragPosition][2]Vec[int]{
		FRAG_RIGHT:        {NewVec(0, 9), NewVec(1, 10)},
		FRAG_BOTTOM:       {NewVec(9, 0), NewVec(10, 1)},
		FRAG_BOTTOM_RIGHT: {NewVec(0, 0), NewVec(1, 1)},
	})

	plane.Expand(&planeBox, 2)
	expectPlaneBoxState(t, planeBox, NewVec(7, 7), NewVec(10, 10), map[FragPosition][2]Vec[int]{
		FRAG_RIGHT:        {NewVec(0, 7), NewVec(3, 10)},
		FRAG_BOTTOM:       {NewVec(7, 0), NewVec(10, 3)},
		FRAG_BOTTOM_RIGHT: {NewVec(0, 0), NewVec(3, 3)},
	})

	plane.Expand(&planeBox, -2)
	expectPlaneBoxState(t, planeBox, NewVec(9, 9), NewVec(10, 10), map[FragPosition][2]Vec[int]{
		FRAG_RIGHT:        {NewVec(0, 9), NewVec(1, 10)},
		FRAG_BOTTOM:       {NewVec(9, 0), NewVec(10, 1)},
		FRAG_BOTTOM_RIGHT: {NewVec(0, 0), NewVec(1, 1)},
	})

	shift.Invert()
	plane.Translate(&planeBox, shift)
	expectPlaneBoxState(t, planeBox, NewVec(0, 0), NewVec(2, 2), map[FragPosition][2]Vec[int]{})
}

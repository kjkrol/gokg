package plane

import (
	"testing"

	"github.com/kjkrol/gokg/pkg/geom"
)

func TestTorusNormalizeVec_int(t *testing.T) {
	torus := NewTorus(5, 5)
	for _, test := range []struct {
		arg1     geom.Vec[int]
		arg2     geom.Vec[int]
		expected geom.Vec[int]
	}{
		{geom.NewVec(2, 3), geom.NewVec(-1, -2), geom.NewVec(1, 1)},
		{geom.NewVec(1, 2), geom.NewVec(-1, -2), geom.NewVec(0, 0)},
		{geom.NewVec(0, 0), geom.NewVec(-4, -4), geom.NewVec(1, 1)},
		{geom.NewVec(4, 0), geom.NewVec(-1, -0), geom.NewVec(3, 0)},
		{geom.NewVec(1, 0), geom.NewVec(-4, -0), geom.NewVec(2, 0)},
	} {
		result := test.arg1
		result.AddMutable(test.arg2)
		torus.(space2d[int]).normalizeVec(&result)
		if !result.Equals(test.expected) {
			t.Errorf("result %v not equal to expected %v", result, test.expected)
		}
	}
}

func TestTorusMetric_int(t *testing.T) {
	torus := NewTorus(9, 9)
	for _, test := range []struct {
		arg1     geom.Vec[int]
		arg2     geom.Vec[int]
		expected int
	}{
		{geom.NewVec(1, 2), geom.NewVec(2, 3), 2},
		{geom.NewVec(1, 2), geom.NewVec(1, 2), 0},
		{geom.NewVec(0, 0), geom.NewVec(1, 1), 2},
		{geom.NewVec(0, 0), geom.NewVec(2, 2), 3},
		{geom.NewVec(0, 0), geom.NewVec(8, 8), 2},
		{geom.NewVec(0, 0), geom.NewVec(9, 9), 0}, // vec(9,9) has been wrapped to vec(0,0)
	} {
		if output := torus.(space2d[int]).metric(test.arg1, test.arg2); output != test.expected {
			t.Errorf("vectors: %v, %v, metric %v not equal to expected %v", test.arg1, test.arg2, output, test.expected)
		}
	}
}

func TestCartesianNormalizedVec_int(t *testing.T) {
	cartesian := NewCartesian(9, 9)
	for _, test := range []struct {
		arg1     geom.Vec[int]
		arg2     geom.Vec[int]
		expected geom.Vec[int]
	}{
		{geom.NewVec(2, 3), geom.NewVec(-1, -2), geom.NewVec(1, 1)},
		{geom.NewVec(1, 2), geom.NewVec(-1, -2), geom.NewVec(0, 0)},
		{geom.NewVec(0, 0), geom.NewVec(-4, -4), geom.NewVec(0, 0)},
		{geom.NewVec(4, 0), geom.NewVec(-1, -0), geom.NewVec(3, 0)},
		{geom.NewVec(6, 0), geom.NewVec(-4, -0), geom.NewVec(2, 0)},
	} {
		result := test.arg1
		result.AddMutable(test.arg2)
		cartesian.(space2d[int]).normalizeVec(&result)
		if !result.Equals(test.expected) {
			t.Errorf("result %v not equal to expected %v", result, test.expected)
		}
	}
}

func TestCartesianMetric_int(t *testing.T) {
	cartesian := NewCartesian(9, 9)
	for _, test := range []struct {
		arg1     geom.Vec[int]
		arg2     geom.Vec[int]
		expected int
	}{
		{geom.NewVec(1, 2), geom.NewVec(2, 3), 2},
		{geom.NewVec(1, 2), geom.NewVec(1, 2), 0},
		{geom.NewVec(0, 0), geom.NewVec(1, 1), 2},
		{geom.NewVec(0, 0), geom.NewVec(2, 2), 3},
		{geom.NewVec(0, 0), geom.NewVec(8, 8), 12},
		{geom.NewVec(0, 0), geom.NewVec(9, 9), 13}, // vec(9,9) stays on the boundary
	} {
		if output := cartesian.(space2d[int]).metric(test.arg1, test.arg2); output != test.expected {
			t.Errorf("vectors: %v, %v, metric %v not equal to expected %v", test.arg1, test.arg2, output, test.expected)
		}
	}
}

// -----------------------------------------------------------------------------

func TestTorusNormalizedVec_float(t *testing.T) {
	torus := NewTorus(5.0, 5.0)
	for _, test := range []struct {
		arg1     geom.Vec[float64]
		arg2     geom.Vec[float64]
		expected geom.Vec[float64]
	}{
		{geom.NewVec(2.0, 3.0), geom.NewVec(-1.0, -2.0), geom.NewVec(1.0, 1.0)},
		{geom.NewVec(1.0, 2.0), geom.NewVec(-1.0, -2.0), geom.NewVec(0.0, 0.0)},
		{geom.NewVec(0.0, 0.0), geom.NewVec(-4.0, -4.0), geom.NewVec(1.0, 1.0)},
		{geom.NewVec(4.0, 0.0), geom.NewVec(-1.0, 0.0), geom.NewVec(3.0, 0.0)},
		{geom.NewVec(1.0, 0.0), geom.NewVec(-4.0, 0.0), geom.NewVec(2.0, 0.0)},
	} {
		result := test.arg1
		result.AddMutable(test.arg2)
		torus.(space2d[float64]).normalizeVec(&result)
		if !result.Equals(test.expected) {
			t.Errorf("result %v not equal to expected %v", result, test.expected)
		}
	}
}

func TestTorusMetric_float(t *testing.T) {
	torus := NewTorus(9.0, 9.0)
	for _, test := range []struct {
		arg1     geom.Vec[float64]
		arg2     geom.Vec[float64]
		expected float64
	}{
		{geom.NewVec(1.0, 2.0), geom.NewVec(2.0, 3.0), 1.4142135623730951},
		{geom.NewVec(1.0, 2.0), geom.NewVec(1.0, 2.0), 0.0},
		{geom.NewVec(0.0, 0.0), geom.NewVec(1.0, 1.0), 1.4142135623730951},
		{geom.NewVec(0.0, 0.0), geom.NewVec(2.0, 2.0), 2.8284271247461903},
		{geom.NewVec(0.0, 0.0), geom.NewVec(8.0, 8.0), 1.4142135623730951},
		{geom.NewVec(0.0, 0.0), geom.NewVec(9.0, 9.0), 0.0}, // vec(9,9) has been wrapped to vec(0,0)
	} {
		if output := torus.(space2d[float64]).metric(test.arg1, test.arg2); output != test.expected {
			t.Errorf("vectors: %v, %v, metric %v not equal to expected %v", test.arg1, test.arg2, output, test.expected)
		}
	}
}

func TestCartesianNormalizeVec_float(t *testing.T) {
	cartesian := NewCartesian(9.0, 9.0)
	for _, test := range []struct {
		arg1     geom.Vec[float64]
		arg2     geom.Vec[float64]
		expected geom.Vec[float64]
	}{
		{geom.NewVec(2.0, 3.0), geom.NewVec(-1.0, -2.0), geom.NewVec(1.0, 1.0)},
		{geom.NewVec(1.0, 2.0), geom.NewVec(-1.0, -2.0), geom.NewVec(0.0, 0.0)},
		{geom.NewVec(0.0, 0.0), geom.NewVec(-4.0, -4.0), geom.NewVec(0.0, 0.0)},
		{geom.NewVec(4.0, 0.0), geom.NewVec(-1.0, 0.0), geom.NewVec(3.0, 0.0)},
		{geom.NewVec(6.0, 0.0), geom.NewVec(-4.0, 0.0), geom.NewVec(2.0, 0.0)},
	} {
		result := test.arg1
		result.AddMutable(test.arg2)
		cartesian.(space2d[float64]).normalizeVec(&result)
		if !result.Equals(test.expected) {
			t.Errorf("result %v not equal to expected %v", result, test.expected)
		}
	}
}

func TestCartesianMetric_float(t *testing.T) {
	cartesian := NewCartesian(9.0, 9.0)
	for _, test := range []struct {
		arg1     geom.Vec[float64]
		arg2     geom.Vec[float64]
		expected float64
	}{
		{geom.NewVec(1.0, 2.0), geom.NewVec(2.0, 3.0), 1.4142135623730951},
		{geom.NewVec(1.0, 2.0), geom.NewVec(1.0, 2.0), 0.0},
		{geom.NewVec(0.0, 0.0), geom.NewVec(1.0, 1.0), 1.4142135623730951},
		{geom.NewVec(0.0, 0.0), geom.NewVec(2.0, 2.0), 2.8284271247461903},
		{geom.NewVec(0.0, 0.0), geom.NewVec(8.0, 8.0), 11.313708498984761},
		{geom.NewVec(0.0, 0.0), geom.NewVec(9.0, 9.0), 12.727922061357855}, // Vec(9,9) stays on the boundary
		{geom.NewVec(0.0, 0.0), geom.NewVec(8.5, 0.0), 8.5},
	} {
		if output := cartesian.(space2d[float64]).metric(test.arg1, test.arg2); output != test.expected {
			t.Errorf("vectors: %v, %v, metric %v not equal to expected %v", test.arg1, test.arg2, output, test.expected)
		}
	}
}

func TestSpace2dNormalizeVec(t *testing.T) {
	plane := NewTorus(5, 5)
	vec := geom.NewVec(7, -2)
	plane.(space2d[int]).normalizeVec(&vec)
	expected := geom.NewVec(2, 3)
	if vec != expected {
		t.Errorf("expected normalized vector %v, got %v", expected, vec)
	}
}

func TestBoundedPlane_TransformBackAndForth(t *testing.T) {
	plane := NewCartesian(10, 10)

	planeBox := newAABB(geom.NewVec(0, 0), 2, 2)

	shift := geom.NewVec(2, 2)
	plane.Translate(&planeBox, shift)
	expectPlaneBoxState(t, planeBox, geom.NewVec(2, 2), geom.NewVec(4, 4), map[FragPosition][2]geom.Vec[int]{})

	plane.Expand(&planeBox, 2)
	expectPlaneBoxState(t, planeBox, geom.NewVec(0, 0), geom.NewVec(6, 6), map[FragPosition][2]geom.Vec[int]{})

	plane.Expand(&planeBox, -2)
	expectPlaneBoxState(t, planeBox, geom.NewVec(2, 2), geom.NewVec(4, 4), map[FragPosition][2]geom.Vec[int]{})

	shift.Invert()
	plane.Translate(&planeBox, shift)
	expectPlaneBoxState(t, planeBox, geom.NewVec(0, 0), geom.NewVec(2, 2), map[FragPosition][2]geom.Vec[int]{})
}

func TestCyclicPlane_TransformBackAndForth(t *testing.T) {
	plane := NewTorus(10, 10)

	planeBox := newAABB(geom.NewVec(0, 0), 2, 2)

	shift := geom.NewVec(-1, -1)
	plane.Translate(&planeBox, shift)
	expectPlaneBoxState(t, planeBox, geom.NewVec(9, 9), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec[int]{
		FRAG_RIGHT:        {geom.NewVec(0, 9), geom.NewVec(1, 10)},
		FRAG_BOTTOM:       {geom.NewVec(9, 0), geom.NewVec(10, 1)},
		FRAG_BOTTOM_RIGHT: {geom.NewVec(0, 0), geom.NewVec(1, 1)},
	})

	plane.Expand(&planeBox, 2)
	expectPlaneBoxState(t, planeBox, geom.NewVec(7, 7), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec[int]{
		FRAG_RIGHT:        {geom.NewVec(0, 7), geom.NewVec(3, 10)},
		FRAG_BOTTOM:       {geom.NewVec(7, 0), geom.NewVec(10, 3)},
		FRAG_BOTTOM_RIGHT: {geom.NewVec(0, 0), geom.NewVec(3, 3)},
	})

	plane.Expand(&planeBox, -2)
	expectPlaneBoxState(t, planeBox, geom.NewVec(9, 9), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec[int]{
		FRAG_RIGHT:        {geom.NewVec(0, 9), geom.NewVec(1, 10)},
		FRAG_BOTTOM:       {geom.NewVec(9, 0), geom.NewVec(10, 1)},
		FRAG_BOTTOM_RIGHT: {geom.NewVec(0, 0), geom.NewVec(1, 1)},
	})

	shift.Invert()
	plane.Translate(&planeBox, shift)
	expectPlaneBoxState(t, planeBox, geom.NewVec(0, 0), geom.NewVec(2, 2), map[FragPosition][2]geom.Vec[int]{})
}

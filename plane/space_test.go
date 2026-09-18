package plane

import (
	"math"
	"testing"

	"github.com/kjkrol/gokg/geom"
)

func TestToroidal2DNormalizeVec(t *testing.T) {
	toroidal := NewToroidal2D(5, 5)
	for _, test := range []struct {
		arg1     geom.Vec
		arg2     geom.Vec
		expected geom.Vec
	}{
		{vec(2, 3), vec(4, 4), vec(1, 2)},
		{vec(1, 2), vec(0, 0), vec(1, 2)},
		{vec(0, 0), vec(6, 6), vec(1, 1)},
		{vec(4, 0), vec(3, 0), vec(2, 0)},
		{vec(3, 4), vec(7, 1), vec(0, 0)},
	} {
		result := test.arg1
		result.AddMutable(test.arg2)
		result = toroidal.(*toroidal2d).normalizeVec(result)
		if !result.Equals(test.expected) {
			t.Errorf("result %v not equal to expected %v", result, test.expected)
		}
	}
}

func TestToroidal2DMetric(t *testing.T) {
	toroidal := NewToroidal2D(9, 9)
	// The distance is now the real one. It used to be carried three times over,
	// once per coordinate type, because an integer metric truncated sqrt(2) to
	// 1 — that rounding was the only reason the cases differed.
	for _, test := range []struct {
		arg1, arg2 [2]float64
		want       float64
	}{
		{arg1: [2]float64{1, 2}, arg2: [2]float64{2, 3}, want: math.Sqrt2},
		{arg1: [2]float64{1, 2}, arg2: [2]float64{1, 2}, want: 0},
		{arg1: [2]float64{0, 0}, arg2: [2]float64{8, 8}, want: math.Sqrt2},
		{arg1: [2]float64{0, 0}, arg2: [2]float64{9, 9}, want: 0}, // vec(9,9) has been wrapped to vec(0,0)
	} {
		expected := test.want
		arg1 := vec(test.arg1[0], test.arg1[1])
		arg2 := vec(test.arg2[0], test.arg2[1])
		if output := toroidal.(*toroidal2d).metric(arg1, arg2); math.Abs(output-expected) > 1e-12 {
			t.Errorf("vectors: %v, %v, metric %v not equal to expected %v", arg1, arg2, output, expected)
		}
	}
}

func TestEuclidean2DNormalizeVec(t *testing.T) {
	euclidean := NewEuclidean2D(9, 9)
	for _, test := range []struct {
		arg1     geom.Vec
		arg2     geom.Vec
		expected geom.Vec
	}{
		{vec(2, 3), vec(4, 4), vec(6, 7)},
		{vec(1, 2), vec(0, 0), vec(1, 2)},
		{vec(0, 0), vec(15, 15), vec(9, 9)},
		{vec(4, 0), vec(9, 0), vec(9, 0)},
		{vec(6, 1), vec(3, 10), vec(9, 9)},
	} {
		result := test.arg1
		result.AddMutable(test.arg2)
		result = euclidean.(*euclidean2d).normalizeVec(result)
		if !result.Equals(test.expected) {
			t.Errorf("result %v not equal to expected %v", result, test.expected)
		}
	}
}

func TestEuclidean2DMetric(t *testing.T) {
	euclidean := NewEuclidean2D(9, 9)
	for _, test := range []struct {
		arg1, arg2 [2]float64
		want       float64
	}{
		{arg1: [2]float64{1, 2}, arg2: [2]float64{2, 3}, want: math.Sqrt2},
		{arg1: [2]float64{1, 2}, arg2: [2]float64{1, 2}, want: 0},
		{arg1: [2]float64{0, 0}, arg2: [2]float64{8, 8}, want: 11.313708498984761},
		{arg1: [2]float64{0, 0}, arg2: [2]float64{9, 9}, want: 12.727922061357855}, // Vec(9,9) stays on the boundary
		{arg1: [2]float64{0, 0}, arg2: [2]float64{5, 0}, want: 5},
	} {
		expected := test.want
		arg1 := vec(test.arg1[0], test.arg1[1])
		arg2 := vec(test.arg2[0], test.arg2[1])
		if output := euclidean.(*euclidean2d).metric(arg1, arg2); math.Abs(output-expected) > 1e-12 {
			t.Errorf("vectors: %v, %v, metric %v not equal to expected %v", arg1, arg2, output, expected)
		}
	}
}

func TestSpace2dNormalizeVec(t *testing.T) {
	toroidal := NewToroidal2D(5, 5)
	v := toroidal.(*toroidal2d).normalizeVec(vec(7, 13))
	expected := vec(2, 3)
	if v != expected {
		t.Errorf("expected normalized vector %v, got %v", expected, v)
	}
}

func vec(x, y float64) geom.Vec {
	return geom.NewVec(x, y)
}

func TestEuclidean2DSpace_TransformBackAndForth(t *testing.T) {
	euclidean := NewEuclidean2D(10, 10)

	box := NewAABB(geom.NewVec(0, 0), 2, 2)

	shift := geom.NewVec(2, 2)
	euclidean.Translate(&box, shift)
	expectAABBState(t, box, geom.NewVec(2, 2), geom.NewVec(4, 4), map[FragPosition][2]geom.Vec{})

	euclidean.Expand(&box, 2)
	expectAABBState(t, box, geom.NewVec(0, 0), geom.NewVec(6, 6), map[FragPosition][2]geom.Vec{})

	euclidean.Expand(&box, -2)
	expectAABBState(t, box, geom.NewVec(2, 2), geom.NewVec(4, 4), map[FragPosition][2]geom.Vec{})

	shift.Invert()
	euclidean.Translate(&box, shift)
	expectAABBState(t, box, geom.NewVec(0, 0), geom.NewVec(2, 2), map[FragPosition][2]geom.Vec{})
}

func TestToroidal2DSpace_TransformBackAndForth(t *testing.T) {
	toroidal := NewToroidal2D(10, 10)

	box := NewAABB(geom.NewVec(0, 0), 2, 2)

	shift := geom.NewVec(-1, -1)
	toroidal.Translate(&box, shift)
	expectAABBState(t, box, geom.NewVec(9, 9), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec{
		FRAG_RIGHT:        {geom.NewVec(0, 9), geom.NewVec(1, 10)},
		FRAG_BOTTOM:       {geom.NewVec(9, 0), geom.NewVec(10, 1)},
		FRAG_BOTTOM_RIGHT: {geom.NewVec(0, 0), geom.NewVec(1, 1)},
	})

	toroidal.Expand(&box, 2)
	expectAABBState(t, box, geom.NewVec(7, 7), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec{
		FRAG_RIGHT:        {geom.NewVec(0, 7), geom.NewVec(3, 10)},
		FRAG_BOTTOM:       {geom.NewVec(7, 0), geom.NewVec(10, 3)},
		FRAG_BOTTOM_RIGHT: {geom.NewVec(0, 0), geom.NewVec(3, 3)},
	})

	toroidal.Expand(&box, -2)
	expectAABBState(t, box, geom.NewVec(9, 9), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec{
		FRAG_RIGHT:        {geom.NewVec(0, 9), geom.NewVec(1, 10)},
		FRAG_BOTTOM:       {geom.NewVec(9, 0), geom.NewVec(10, 1)},
		FRAG_BOTTOM_RIGHT: {geom.NewVec(0, 0), geom.NewVec(1, 1)},
	})

	shift.Invert()
	toroidal.Translate(&box, shift)
	expectAABBState(t, box, geom.NewVec(0, 0), geom.NewVec(2, 2), map[FragPosition][2]geom.Vec{})
}

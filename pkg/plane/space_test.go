package plane

import (
	"math"
	"testing"

	"github.com/kjkrol/gokg/pkg/geom"
)

func TestToroidal2DNormalizeVec(t *testing.T) {
	runToroidal2DNormalizeVecTest[int](t, "int")
	runToroidal2DNormalizeVecTest[uint32](t, "uint32")
	runToroidal2DNormalizeVecTest[float64](t, "float64")
}

func runToroidal2DNormalizeVecTest[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		toroidal := NewToroidal2D(T(5), T(5))
		for _, test := range []struct {
			arg1     geom.Vec[T]
			arg2     geom.Vec[T]
			expected geom.Vec[T]
		}{
			{vec[T](2, 3), vec[T](4, 4), vec[T](1, 2)},
			{vec[T](1, 2), vec[T](0, 0), vec[T](1, 2)},
			{vec[T](0, 0), vec[T](6, 6), vec[T](1, 1)},
			{vec[T](4, 0), vec[T](3, 0), vec[T](2, 0)},
			{vec[T](3, 4), vec[T](7, 1), vec[T](0, 0)},
		} {
			result := test.arg1
			result.AddMutable(test.arg2)
			result = toroidal.(*toroidal2d[T]).normalizeVec(result)
			if !result.Equals(test.expected) {
				t.Errorf("result %v not equal to expected %v", result, test.expected)
			}
		}
	})
}

func TestToroidal2DMetric(t *testing.T) {
	runToroidal2DMetricTest[int](t, "int")
	runToroidal2DMetricTest[uint32](t, "uint32")
	runToroidal2DMetricTest[float64](t, "float64")
}

func runToroidal2DMetricTest[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		toroidal := NewToroidal2D(T(9), T(9))
		for _, test := range []struct {
			arg1, arg2   [2]int
			wantInt      int
			wantUnsigned int
			wantFloat    float64 // math.Sqrt2 is the constant value of sqrt(2)
		}{
			{arg1: [2]int{1, 2}, arg2: [2]int{2, 3}, wantInt: 2, wantUnsigned: 2, wantFloat: math.Sqrt2},
			{arg1: [2]int{1, 2}, arg2: [2]int{1, 2}, wantInt: 0, wantUnsigned: 0, wantFloat: 0},
			{arg1: [2]int{0, 0}, arg2: [2]int{8, 8}, wantInt: 2, wantUnsigned: 2, wantFloat: math.Sqrt2},
			{arg1: [2]int{0, 0}, arg2: [2]int{9, 9}, wantInt: 0, wantUnsigned: 0, wantFloat: 0}, // vec(9,9) has been wrapped to vec(0,0)
		} {
			expected := chooseExpected[T](test.wantInt, test.wantUnsigned, test.wantFloat)
			arg1 := vec[T](test.arg1[0], test.arg1[1])
			arg2 := vec[T](test.arg2[0], test.arg2[1])
			if output := toroidal.(*toroidal2d[T]).metric(arg1, arg2); output != expected {
				t.Errorf("vectors: %v, %v, metric %v not equal to expected %v", arg1, arg2, output, expected)
			}
		}
	})
}

func TestEuclidean2DNormalizeVec(t *testing.T) {
	runEuclidean2DNormalizeVecTest[int](t, "int")
	runEuclidean2DNormalizeVecTest[uint32](t, "uint32")
	runEuclidean2DNormalizeVecTest[float64](t, "float64")
}

func runEuclidean2DNormalizeVecTest[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		euclidean := NewEuclidean2D(T(9), T(9))
		for _, test := range []struct {
			arg1     geom.Vec[T]
			arg2     geom.Vec[T]
			expected geom.Vec[T]
		}{
			{vec[T](2, 3), vec[T](4, 4), vec[T](6, 7)},
			{vec[T](1, 2), vec[T](0, 0), vec[T](1, 2)},
			{vec[T](0, 0), vec[T](15, 15), vec[T](9, 9)},
			{vec[T](4, 0), vec[T](9, 0), vec[T](9, 0)},
			{vec[T](6, 1), vec[T](3, 10), vec[T](9, 9)},
		} {
			result := test.arg1
			result.AddMutable(test.arg2)
			result = euclidean.(*euclidean2d[T]).normalizeVec(result)
			if !result.Equals(test.expected) {
				t.Errorf("result %v not equal to expected %v", result, test.expected)
			}
		}
	})
}

func TestEuclidean2DMetric(t *testing.T) {
	runEuclidean2DMetricTest[int](t, "int")
	runEuclidean2DMetricTest[uint32](t, "uint32")
	runEuclidean2DMetricTest[float64](t, "float64")
}

func runEuclidean2DMetricTest[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		euclidean := NewEuclidean2D(T(9), T(9))
		for _, test := range []struct {
			arg1, arg2   [2]int
			wantInt      int
			wantUnsigned int
			wantFloat    float64
		}{
			{arg1: [2]int{1, 2}, arg2: [2]int{2, 3}, wantInt: 2, wantUnsigned: 2, wantFloat: math.Sqrt2},
			{arg1: [2]int{1, 2}, arg2: [2]int{1, 2}, wantInt: 0, wantUnsigned: 0, wantFloat: 0},
			{arg1: [2]int{0, 0}, arg2: [2]int{8, 8}, wantInt: 12, wantUnsigned: 12, wantFloat: 11.313708498984761},
			{arg1: [2]int{0, 0}, arg2: [2]int{9, 9}, wantInt: 13, wantUnsigned: 13, wantFloat: 12.727922061357855}, // Vec(9,9) stays on the boundary
			{arg1: [2]int{0, 0}, arg2: [2]int{5, 0}, wantInt: 5, wantUnsigned: 5, wantFloat: 5},
		} {
			expected := chooseExpected[T](test.wantInt, test.wantUnsigned, test.wantFloat)
			arg1 := vec[T](test.arg1[0], test.arg1[1])
			arg2 := vec[T](test.arg2[0], test.arg2[1])
			if output := euclidean.(*euclidean2d[T]).metric(arg1, arg2); output != expected {
				t.Errorf("vectors: %v, %v, metric %v not equal to expected %v", arg1, arg2, output, expected)
			}
		}
	})
}

func TestSpace2dNormalizeVec(t *testing.T) {
	runSpace2dNormalizeVecTest[int](t, "int")
	runSpace2dNormalizeVecTest[uint32](t, "uint32")
	runSpace2dNormalizeVecTest[float64](t, "float64")
}

func runSpace2dNormalizeVecTest[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		toroidal := NewToroidal2D(T(5), T(5))
		v := toroidal.(*toroidal2d[T]).normalizeVec(vec[T](7, 13))
		expected := vec[T](2, 3)
		if v != expected {
			t.Errorf("expected normalized vector %v, got %v", expected, v)
		}
	})
}

func chooseExpected[T geom.Numeric](wantInt, wantUnsigned int, wantFloat float64) T {
	var zero T
	switch any(zero).(type) {
	case float64:
		return T(wantFloat)
	case uint32:
		return T(wantUnsigned)
	default:
		return T(wantInt)
	}
}

func vec[T geom.Numeric](x, y int) geom.Vec[T] {
	return geom.NewVec(T(x), T(y))
}

func TestEuclidean2DSpace_TransformBackAndForth(t *testing.T) {
	euclidean := NewEuclidean2D(10, 10)

	box := newAABB(geom.NewVec(0, 0), 2, 2)

	shift := geom.NewVec(2, 2)
	euclidean.Translate(&box, shift)
	expectAABBState(t, box, geom.NewVec(2, 2), geom.NewVec(4, 4), map[FragPosition][2]geom.Vec[int]{})

	euclidean.Expand(&box, 2)
	expectAABBState(t, box, geom.NewVec(0, 0), geom.NewVec(6, 6), map[FragPosition][2]geom.Vec[int]{})

	euclidean.Expand(&box, -2)
	expectAABBState(t, box, geom.NewVec(2, 2), geom.NewVec(4, 4), map[FragPosition][2]geom.Vec[int]{})

	shift.Invert()
	euclidean.Translate(&box, shift)
	expectAABBState(t, box, geom.NewVec(0, 0), geom.NewVec(2, 2), map[FragPosition][2]geom.Vec[int]{})
}

func TestToroidal2DSpace_TransformBackAndForth(t *testing.T) {
	toroidal := NewToroidal2D(10, 10)

	box := newAABB(geom.NewVec(0, 0), 2, 2)

	shift := geom.NewVec(-1, -1)
	toroidal.Translate(&box, shift)
	expectAABBState(t, box, geom.NewVec(9, 9), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec[int]{
		FRAG_RIGHT:        {geom.NewVec(0, 9), geom.NewVec(1, 10)},
		FRAG_BOTTOM:       {geom.NewVec(9, 0), geom.NewVec(10, 1)},
		FRAG_BOTTOM_RIGHT: {geom.NewVec(0, 0), geom.NewVec(1, 1)},
	})

	toroidal.Expand(&box, 2)
	expectAABBState(t, box, geom.NewVec(7, 7), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec[int]{
		FRAG_RIGHT:        {geom.NewVec(0, 7), geom.NewVec(3, 10)},
		FRAG_BOTTOM:       {geom.NewVec(7, 0), geom.NewVec(10, 3)},
		FRAG_BOTTOM_RIGHT: {geom.NewVec(0, 0), geom.NewVec(3, 3)},
	})

	toroidal.Expand(&box, -2)
	expectAABBState(t, box, geom.NewVec(9, 9), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec[int]{
		FRAG_RIGHT:        {geom.NewVec(0, 9), geom.NewVec(1, 10)},
		FRAG_BOTTOM:       {geom.NewVec(9, 0), geom.NewVec(10, 1)},
		FRAG_BOTTOM_RIGHT: {geom.NewVec(0, 0), geom.NewVec(1, 1)},
	})

	shift.Invert()
	toroidal.Translate(&box, shift)
	expectAABBState(t, box, geom.NewVec(0, 0), geom.NewVec(2, 2), map[FragPosition][2]geom.Vec[int]{})
}

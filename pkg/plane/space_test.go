package plane

import (
	"math"
	"testing"

	"github.com/kjkrol/gokg/pkg/geom"
)

func TestTorusNormalizeVec(t *testing.T) {
	runTorusNormalizeVecTest[int](t, "int")
	runTorusNormalizeVecTest[uint32](t, "uint32")
	runTorusNormalizeVecTest[float64](t, "float64")
}

func runTorusNormalizeVecTest[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		torus := NewTorus(T(5), T(5))
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
			torus.(space2d[T]).normalizeVec(&result)
			if !result.Equals(test.expected) {
				t.Errorf("result %v not equal to expected %v", result, test.expected)
			}
		}
	})
}

func TestTorusMetric(t *testing.T) {
	runTorusMetricTest[int](t, "int")
	runTorusMetricTest[uint32](t, "uint32")
	runTorusMetricTest[float64](t, "float64")
}

func runTorusMetricTest[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		torus := NewTorus(T(9), T(9))
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
			if output := torus.(space2d[T]).metric(arg1, arg2); output != expected {
				t.Errorf("vectors: %v, %v, metric %v not equal to expected %v", arg1, arg2, output, expected)
			}
		}
	})
}

func TestCartesianNormalizeVec(t *testing.T) {
	runCartesianNormalizeVecTest[int](t, "int")
	runCartesianNormalizeVecTest[uint32](t, "uint32")
	runCartesianNormalizeVecTest[float64](t, "float64")
}

func runCartesianNormalizeVecTest[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		cartesian := NewCartesian(T(9), T(9))
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
			cartesian.(space2d[T]).normalizeVec(&result)
			if !result.Equals(test.expected) {
				t.Errorf("result %v not equal to expected %v", result, test.expected)
			}
		}
	})
}

func TestCartesianMetric(t *testing.T) {
	runCartesianMetricTest[int](t, "int")
	runCartesianMetricTest[uint32](t, "uint32")
	runCartesianMetricTest[float64](t, "float64")
}

func runCartesianMetricTest[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		cartesian := NewCartesian(T(9), T(9))
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
			if output := cartesian.(space2d[T]).metric(arg1, arg2); output != expected {
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
		plane := NewTorus(T(5), T(5))
		v := vec[T](7, 13)
		plane.(space2d[T]).normalizeVec(&v)
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

func TestBoundedPlane_TransformBackAndForth(t *testing.T) {
	plane := NewCartesian(10, 10)

	planeBox := newAABB(geom.NewVec(0, 0), 2, 2)

	shift := geom.NewVec(2, 2)
	plane.Translate(&planeBox, shift)
	expectAABBState(t, planeBox, geom.NewVec(2, 2), geom.NewVec(4, 4), map[FragPosition][2]geom.Vec[int]{})

	plane.Expand(&planeBox, 2)
	expectAABBState(t, planeBox, geom.NewVec(0, 0), geom.NewVec(6, 6), map[FragPosition][2]geom.Vec[int]{})

	plane.Expand(&planeBox, -2)
	expectAABBState(t, planeBox, geom.NewVec(2, 2), geom.NewVec(4, 4), map[FragPosition][2]geom.Vec[int]{})

	shift.Invert()
	plane.Translate(&planeBox, shift)
	expectAABBState(t, planeBox, geom.NewVec(0, 0), geom.NewVec(2, 2), map[FragPosition][2]geom.Vec[int]{})
}

func TestCyclicPlane_TransformBackAndForth(t *testing.T) {
	plane := NewTorus(10, 10)

	planeBox := newAABB(geom.NewVec(0, 0), 2, 2)

	shift := geom.NewVec(-1, -1)
	plane.Translate(&planeBox, shift)
	expectAABBState(t, planeBox, geom.NewVec(9, 9), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec[int]{
		FRAG_RIGHT:        {geom.NewVec(0, 9), geom.NewVec(1, 10)},
		FRAG_BOTTOM:       {geom.NewVec(9, 0), geom.NewVec(10, 1)},
		FRAG_BOTTOM_RIGHT: {geom.NewVec(0, 0), geom.NewVec(1, 1)},
	})

	plane.Expand(&planeBox, 2)
	expectAABBState(t, planeBox, geom.NewVec(7, 7), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec[int]{
		FRAG_RIGHT:        {geom.NewVec(0, 7), geom.NewVec(3, 10)},
		FRAG_BOTTOM:       {geom.NewVec(7, 0), geom.NewVec(10, 3)},
		FRAG_BOTTOM_RIGHT: {geom.NewVec(0, 0), geom.NewVec(3, 3)},
	})

	plane.Expand(&planeBox, -2)
	expectAABBState(t, planeBox, geom.NewVec(9, 9), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec[int]{
		FRAG_RIGHT:        {geom.NewVec(0, 9), geom.NewVec(1, 10)},
		FRAG_BOTTOM:       {geom.NewVec(9, 0), geom.NewVec(10, 1)},
		FRAG_BOTTOM_RIGHT: {geom.NewVec(0, 0), geom.NewVec(1, 1)},
	})

	shift.Invert()
	plane.Translate(&planeBox, shift)
	expectAABBState(t, planeBox, geom.NewVec(0, 0), geom.NewVec(2, 2), map[FragPosition][2]geom.Vec[int]{})
}

package plane

import (
	"testing"

	"github.com/kjkrol/gokg/pkg/geom"
)

func TestAxisDistance(t *testing.T) {
	runAxisDistanceTest[int](t, "int")
	runAxisDistanceTest[uint32](t, "uint32")
	runAxisDistanceTest[float64](t, "float64")
}

func runAxisDistanceTest[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		aa := newAABB(geom.NewVec(T(0), T(0)), T(2), T(2))
		bb := newAABB(geom.NewVec(T(5), T(0)), T(2), T(2))

		dx := aa.AxisDistanceX(bb.AABB)
		dy := aa.AxisDistanceY(bb.AABB)

		if dx != T(3) {
			t.Errorf("expected dx=3, got %v", dx)
		}
		if dy != T(0) {
			t.Errorf("expected dy=0, got %v", dy)
		}
	})
}

func TestAABBDistance_CartesianSpace(t *testing.T) {
	runAABBDistanceCartesianTest[int](t, "int")
	runAABBDistanceCartesianTest[uint32](t, "uint32")
	runAABBDistanceCartesianTest[float64](t, "float64")
}

func runAABBDistanceCartesianTest[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		rectA := newAABB(vec[T](0, 0), T(2), T(2))
		rectB := newAABB(vec[T](4, 5), T(2), T(2))

		cartesian := NewCartesian(T(20), T(20))
		distance := cartesian.AABBDistance()(rectA.AABB, rectB.AABB)

		expected := cartesian.(space2d[T]).metric(vec[T](2, 3), geom.NewVec[T](0, 0))
		if distance != expected {
			t.Errorf("expected distance %v, got %v", expected, distance)
		}
	})
}

func TestAABBDistance_ReturnsZeroOnIntersection(t *testing.T) {
	runAABBDistanceZeroOnIntersectionTest[int](t, "int")
	runAABBDistanceZeroOnIntersectionTest[uint32](t, "uint32")
	runAABBDistanceZeroOnIntersectionTest[float64](t, "float64")
}

func runAABBDistanceZeroOnIntersectionTest[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		rectA := newAABB(vec[T](0, 0), T(4), T(4))
		rectB := newAABB(vec[T](2, 2), T(4), T(4))

		cartesian := NewCartesian(T(20), T(20))
		distance := cartesian.AABBDistance()(rectA.AABB, rectB.AABB)
		if distance != T(0) {
			t.Errorf("expected distance 0 for intersecting rectangles, got %v", distance)
		}
	})
}

func TestAABBDistance_To_Vector_Cartesian(t *testing.T) {
	runAABBDistanceToVectorCartesianTest[int](t, "int")
	runAABBDistanceToVectorCartesianTest[uint32](t, "uint32")
	runAABBDistanceToVectorCartesianTest[float64](t, "float64")
}

func runAABBDistanceToVectorCartesianTest[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		testCases := []struct {
			name   string
			first  geom.AABB[T]
			second geom.AABB[T]
			delta  geom.Vec[T]
		}{
			{
				name:   "vectorAsPointLeft",
				first:  geom.NewAABBAt(vec[T](0, 0), T(0), T(0)),
				second: newAABB(vec[T](4, 0), T(2), T(2)).AABB,
				delta:  vec[T](4, 0),
			},
			{
				name:   "vectorAsPointRight",
				first:  newAABB(vec[T](0, 0), T(2), T(2)).AABB,
				second: geom.NewAABBAt(vec[T](5, 6), T(0), T(0)),
				delta:  vec[T](3, 4),
			},
		}

		cartesian := NewCartesian(T(100), T(100))
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				distance := cartesian.AABBDistance()(tc.first, tc.second)
				expected := cartesian.(space2d[T]).metric(tc.delta, geom.NewVec[T](0, 0))
				if distance != expected {
					t.Errorf("expected distance %v, got %v", expected, distance)
				}
			})
		}
	})
}

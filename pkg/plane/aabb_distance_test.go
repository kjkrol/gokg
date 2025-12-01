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

func TestAABBDistance_Euclidean2DSpace(t *testing.T) {
	runAABBDistanceEuclidean2DTest[int](t, "int")
	runAABBDistanceEuclidean2DTest[uint32](t, "uint32")
	runAABBDistanceEuclidean2DTest[float64](t, "float64")
}

func runAABBDistanceEuclidean2DTest[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		rectA := newAABB(vec[T](0, 0), T(2), T(2))
		rectB := newAABB(vec[T](4, 5), T(2), T(2))

		euclidean := NewEuclidean2D(T(20), T(20))
		distance := euclidean.AABBDistance()(rectA.AABB, rectB.AABB)

		expected := euclidean.(*euclidean2d[T]).metric(vec[T](2, 3), geom.NewVec[T](0, 0))
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

		euclidean := NewEuclidean2D(T(20), T(20))
		distance := euclidean.AABBDistance()(rectA.AABB, rectB.AABB)
		if distance != T(0) {
			t.Errorf("expected distance 0 for intersecting rectangles, got %v", distance)
		}
	})
}

func TestAABBDistance_To_Vector_Euclidean2D(t *testing.T) {
	runAABBDistanceToVectorEuclidean2DTest[int](t, "int")
	runAABBDistanceToVectorEuclidean2DTest[uint32](t, "uint32")
	runAABBDistanceToVectorEuclidean2DTest[float64](t, "float64")
}

func runAABBDistanceToVectorEuclidean2DTest[T geom.Numeric](t *testing.T, name string) {
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

		euclidean := NewEuclidean2D(T(100), T(100))
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				distance := euclidean.AABBDistance()(tc.first, tc.second)
				expected := euclidean.(*euclidean2d[T]).metric(tc.delta, geom.NewVec[T](0, 0))
				if distance != expected {
					t.Errorf("expected distance %v, got %v", expected, distance)
				}
			})
		}
	})
}

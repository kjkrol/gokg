package plane

import (
	"testing"

	"github.com/kjkrol/gokg/pkg/geom"
)

func TestEuclidean2DExpand(t *testing.T) {
	runEuclidean2DExpandTest[int](t, "int")
	runEuclidean2DExpandTest[uint32](t, "uint32")
	runEuclidean2DExpandTest[float64](t, "float64")
}

func runEuclidean2DExpandTest[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		euclidean := NewEuclidean2D(T(10), T(10))
		aabb := newAABB(vec[T](2, 3), T(3), T(4))
		euclidean.Expand(&aabb, T(2))
		expectAABBState(t, aabb, vec[T](0, 1), vec[T](7, 9), map[FragPosition][2]geom.Vec[T]{})
	})
}

func TestEuclidean2DExpandCornerCase(t *testing.T) {
	runEuclidean2DExpandCornerCase[int](t, "int")
	runEuclidean2DExpandCornerCase[uint32](t, "uint32")
	runEuclidean2DExpandCornerCase[float64](t, "float64")
}

func runEuclidean2DExpandCornerCase[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		euclidean := NewEuclidean2D(T(10), T(10))
		aabb := newAABB(vec[T](0, 0), T(2), T(2))
		euclidean.Expand(&aabb, T(2))
		expectAABBState(t, aabb, vec[T](0, 0), vec[T](4, 4), map[FragPosition][2]geom.Vec[T]{})
	})
}

func TestToroidal2DExpandCornerCase(t *testing.T) {
	runToroidal2DExpandCornerCase[int](t, "int")
	runToroidal2DExpandCornerCase[uint32](t, "uint32")
	runToroidal2DExpandCornerCase[float64](t, "float64")
}

func runToroidal2DExpandCornerCase[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		toroidal := NewToroidal2D(T(10), T(10))
		aabb := newAABB(vec[T](0, 0), T(2), T(2))
		toroidal.Expand(&aabb, T(2))
		expectAABBState(t, aabb, vec[T](8, 8), vec[T](10, 10), convertFragments[T](map[FragPosition][2]geom.Vec[int]{
			FRAG_RIGHT:        {geom.NewVec(0, 8), geom.NewVec(4, 10)},
			FRAG_BOTTOM:       {geom.NewVec(8, 0), geom.NewVec(10, 4)},
			FRAG_BOTTOM_RIGHT: {geom.NewVec(0, 0), geom.NewVec(4, 4)},
		}))
	})
}

func TestToroidal2DExpandThenIntersects(t *testing.T) {
	runToroidal2DExpandThenIntersects[int](t, "int")
	runToroidal2DExpandThenIntersects[uint32](t, "uint32")
	runToroidal2DExpandThenIntersects[float64](t, "float64")
}

func runToroidal2DExpandThenIntersects[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		toroidal := NewToroidal2D(T(100), T(100))

		aabb1 := newAABB(vec[T](5, 5), T(10), T(10))
		aabb2 := newAABB(vec[T](96, 96), T(10), T(10))

		toroidal.Expand(&aabb2, T(0))

		if intersects := aabb1.IntersectsWithFrags(aabb2); intersects != true {
			t.Errorf("unexpected intersection result. got %t, want %t for boxes %v and %v", intersects, true, aabb1, aabb2)
		}
	})
}

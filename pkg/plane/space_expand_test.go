package plane

import (
	"testing"

	"github.com/kjkrol/gokg/pkg/geom"
)

func TestCartesianExpand(t *testing.T) {
	runCartesianExpandTest[int](t, "int")
	runCartesianExpandTest[uint32](t, "uint32")
	runCartesianExpandTest[float64](t, "float64")
}

func runCartesianExpandTest[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		cartesian := NewCartesian(T(10), T(10))
		aabb := newAABB(vec[T](2, 3), T(3), T(4))
		cartesian.Expand(&aabb, T(2))
		expectAABBState(t, aabb, vec[T](0, 1), vec[T](7, 9), map[FragPosition][2]geom.Vec[T]{})
	})
}

func TestCartesianExpandCornerCase(t *testing.T) {
	runCartesianExpandCornerCase[int](t, "int")
	runCartesianExpandCornerCase[uint32](t, "uint32")
	runCartesianExpandCornerCase[float64](t, "float64")
}

func runCartesianExpandCornerCase[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		cartesian := NewCartesian(T(10), T(10))
		aabb := newAABB(vec[T](0, 0), T(2), T(2))
		cartesian.Expand(&aabb, T(2))
		expectAABBState(t, aabb, vec[T](0, 0), vec[T](4, 4), map[FragPosition][2]geom.Vec[T]{})
	})
}

func TestTorusExpandCornerCase(t *testing.T) {
	runTorusExpandCornerCase[int](t, "int")
	runTorusExpandCornerCase[uint32](t, "uint32")
	runTorusExpandCornerCase[float64](t, "float64")
}

func runTorusExpandCornerCase[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		torus := NewTorus(T(10), T(10))
		aabb := newAABB(vec[T](0, 0), T(2), T(2))
		torus.Expand(&aabb, T(2))
		expectAABBState(t, aabb, vec[T](8, 8), vec[T](10, 10), convertFragments[T](map[FragPosition][2]geom.Vec[int]{
			FRAG_RIGHT:        {geom.NewVec(0, 8), geom.NewVec(4, 10)},
			FRAG_BOTTOM:       {geom.NewVec(8, 0), geom.NewVec(10, 4)},
			FRAG_BOTTOM_RIGHT: {geom.NewVec(0, 0), geom.NewVec(4, 4)},
		}))
	})
}

func TestTorusExpandThenIntersects(t *testing.T) {
	runTorusExpandThenIntersects[int](t, "int")
	runTorusExpandThenIntersects[uint32](t, "uint32")
	runTorusExpandThenIntersects[float64](t, "float64")
}

func runTorusExpandThenIntersects[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		torus := NewTorus(T(100), T(100))

		aabb1 := newAABB(vec[T](5, 5), T(10), T(10))
		aabb2 := newAABB(vec[T](96, 96), T(10), T(10))

		torus.Expand(&aabb2, T(0))

		if intersects := aabb1.Intersects(aabb2); intersects != true {
			t.Errorf("unexpected intersection result. got %t, want %t for boxes %v and %v", intersects, true, aabb1, aabb2)
		}
	})
}

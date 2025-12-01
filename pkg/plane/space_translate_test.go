package plane

import (
	"testing"

	"github.com/kjkrol/gokg/pkg/geom"
)

func TestToroidal2DTranslate(t *testing.T) {
	runToroidal2DTranslateTest[int](t, "int")
	runToroidal2DTranslateTest[uint32](t, "uint32")
	runToroidal2DTranslateTest[float64](t, "float64")
}

func runToroidal2DTranslateTest[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		t.Run("CrossesRightEdge", func(t *testing.T) {
			toroidal := NewToroidal2D(T(10), T(10))
			aabb := newAABB(vec[T](8, 4), T(4), T(2))

			toroidal.Translate(&aabb, vec[T](0, 0))

			expectAABBState(t, aabb, vec[T](8, 4), vec[T](10, 6), map[FragPosition][2]geom.Vec[T]{
				FRAG_RIGHT: {vec[T](0, 4), vec[T](2, 6)},
			})
		})

		t.Run("HugeShift", func(t *testing.T) {
			toroidal := NewToroidal2D(T(10), T(10))
			aabb := newAABB(vec[T](8, 4), T(4), T(2))

			toroidal.Translate(&aabb, vec[T](100, 100))
			expectAABBState(t, aabb, vec[T](8, 4), vec[T](10, 6), map[FragPosition][2]geom.Vec[T]{
				FRAG_RIGHT: {vec[T](0, 4), vec[T](2, 6)},
			})

			toroidal.Translate(&aabb, vec[T](100, 100))
			expectAABBState(t, aabb, vec[T](8, 4), vec[T](10, 6), map[FragPosition][2]geom.Vec[T]{
				FRAG_RIGHT: {vec[T](0, 4), vec[T](2, 6)},
			})
		})

		t.Run("BackAndForth", func(t *testing.T) {
			toroidal := NewToroidal2D(T(10), T(10))
			aabb := newAABB(vec[T](8, 8), T(2), T(2))

			toroidal.Translate(&aabb, vec[T](5, 5))
			expectAABBState(t, aabb, vec[T](3, 3), vec[T](5, 5), map[FragPosition][2]geom.Vec[T]{})

			toroidal.Translate(&aabb, vec[T](5, 5))
			expectAABBState(t, aabb, vec[T](8, 8), vec[T](10, 10), map[FragPosition][2]geom.Vec[T]{})
		})

		t.Run("CrossesBottomEdge", func(t *testing.T) {
			toroidal := NewToroidal2D(T(10), T(10))
			aabb := newAABB(vec[T](4, 8), T(2), T(4))

			toroidal.Translate(&aabb, vec[T](0, 0))

			expectAABBState(t, aabb, vec[T](4, 8), vec[T](6, 10), map[FragPosition][2]geom.Vec[T]{
				FRAG_BOTTOM: {vec[T](4, 0), vec[T](6, 2)},
			})
		})

		t.Run("CrossesCorner", func(t *testing.T) {
			toroidal := NewToroidal2D(T(10), T(10))
			aabb := newAABB(vec[T](9, 9), T(2), T(2))

			toroidal.Translate(&aabb, vec[T](0, 0))

			expectAABBState(t, aabb, vec[T](9, 9), vec[T](10, 10), map[FragPosition][2]geom.Vec[T]{
				FRAG_RIGHT:        {vec[T](0, 9), vec[T](1, 10)},
				FRAG_BOTTOM:       {vec[T](9, 0), vec[T](10, 1)},
				FRAG_BOTTOM_RIGHT: {vec[T](0, 0), vec[T](1, 1)},
			})
		})

		t.Run("ClearsFragmentsWhenNotWrapping", func(t *testing.T) {
			toroidal := NewToroidal2D(T(10), T(10))
			aabb := newAABB(vec[T](8, 4), T(4), T(2))

			toroidal.Translate(&aabb, vec[T](0, 0))
			expectAABBState(t, aabb, vec[T](8, 4), vec[T](10, 6), map[FragPosition][2]geom.Vec[T]{
				FRAG_RIGHT: {vec[T](0, 4), vec[T](2, 6)},
			})

			toroidal.Translate(&aabb, vec[T](8, 0))
			expectAABBState(t, aabb, vec[T](6, 4), vec[T](10, 6), map[FragPosition][2]geom.Vec[T]{})
		})

		t.Run("ThroughEdge", func(t *testing.T) {
			toroidal := NewToroidal2D(T(10), T(10))
			aabb := newAABB(vec[T](2, 2), T(2), T(2))

			toroidal.Translate(&aabb, vec[T](8, 0))

			expectAABBState(t, aabb, vec[T](0, 2), vec[T](2, 4), map[FragPosition][2]geom.Vec[T]{})
		})

		t.Run("FragmentsMergeSequence", func(t *testing.T) {
			toroidal := NewToroidal2D(T(10), T(10))
			aabb := newAABB(vec[T](2, 2), T(2), T(2))

			toroidal.Translate(&aabb, vec[T](7, 0))
			expectAABBState(t, aabb, vec[T](9, 2), vec[T](10, 4), map[FragPosition][2]geom.Vec[T]{
				FRAG_RIGHT: {vec[T](0, 2), vec[T](1, 4)},
			})

			toroidal.Translate(&aabb, vec[T](0, 7))
			expectAABBState(t, aabb, vec[T](9, 9), vec[T](10, 10), map[FragPosition][2]geom.Vec[T]{
				FRAG_RIGHT:        {vec[T](0, 9), vec[T](1, 10)},
				FRAG_BOTTOM:       {vec[T](9, 0), vec[T](10, 1)},
				FRAG_BOTTOM_RIGHT: {vec[T](0, 0), vec[T](1, 1)},
			})

			toroidal.Translate(&aabb, vec[T](3, 3))
			expectAABBState(t, aabb, vec[T](2, 2), vec[T](4, 4), map[FragPosition][2]geom.Vec[T]{})
		})
	})
}

func TestEuclidean2DTranslate(t *testing.T) {
	runEuclidean2DTranslateTest[int](t, "int")
	runEuclidean2DTranslateTest[uint32](t, "uint32")
	runEuclidean2DTranslateTest[float64](t, "float64")
}

func runEuclidean2DTranslateTest[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		t.Run("HugeShift", func(t *testing.T) {
			euclidean := NewEuclidean2D(T(10), T(10))
			aabb := newAABB(vec[T](8, 4), T(4), T(2))

			euclidean.Translate(&aabb, vec[T](100, 100))
			expectAABBState(t, aabb, vec[T](10, 10), vec[T](10, 10), map[FragPosition][2]geom.Vec[T]{})
		})

		t.Run("TranslateWithinBounds", func(t *testing.T) {
			euclidean := NewEuclidean2D(T(10), T(10))
			aabb := newAABB(vec[T](1, 1), T(2), T(2))

			euclidean.Translate(&aabb, vec[T](2, 2))
			expectAABBState(t, aabb, vec[T](3, 3), vec[T](5, 5), map[FragPosition][2]geom.Vec[T]{})
		})

		t.Run("ClampAtBoundary", func(t *testing.T) {
			euclidean := NewEuclidean2D(T(10), T(10))
			aabb := newAABB(vec[T](9, 9), T(3), T(3))

			euclidean.Translate(&aabb, vec[T](2, 2))
			expectAABBState(t, aabb, vec[T](10, 10), vec[T](10, 10), map[FragPosition][2]geom.Vec[T]{})
		})
	})
}

func expectAABBState[T geom.Numeric](
	t *testing.T,
	b AABB[T],
	expectedPos geom.Vec[T],
	expectedBottomRight geom.Vec[T],
	expectedFragments map[FragPosition][2]geom.Vec[T],
) {
	t.Helper()

	if !b.TopLeft.Equals(expectedPos) {
		t.Fatalf("expected bound position %v, got %v", expectedPos, b.TopLeft)
	}

	if !b.BottomRight.Equals(expectedBottomRight) {
		t.Fatalf("expected bound bottom-right %v, got %v", expectedBottomRight, b.BottomRight)
	}

	expectAABBFragments(t, b, expectedFragments)
}

func expectAABBFragments[T geom.Numeric](t *testing.T, b AABB[T], expected map[FragPosition][2]geom.Vec[T]) {
	t.Helper()

	actual := map[FragPosition][2]geom.Vec[T]{}
	(&b).VisitFragments(func(pos FragPosition, box geom.AABB[T]) bool {
		actual[pos] = [2]geom.Vec[T]{box.TopLeft, box.BottomRight}
		return true
	})

	if len(actual) != len(expected) {
		t.Fatalf("expected %d fragments, got %d", len(expected), len(actual))
	}

	for pos, want := range expected {
		frag, ok := actual[pos]
		if !ok {
			t.Fatalf("missing fragment at %d", pos)
		}
		if !frag[0].Equals(want[0]) || !frag[1].Equals(want[1]) {
			t.Fatalf("fragment at %d has bounds %v..%v, expected %v..%v",
				pos, frag[0], frag[1], want[0], want[1])
		}
	}
}

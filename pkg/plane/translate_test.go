package plane

import (
	"testing"

	"github.com/kjkrol/gokg/pkg/geom"
)

func TestTorusTranslate(t *testing.T) {
	runTorusTranslateTest[int](t, "int")
	runTorusTranslateTest[uint32](t, "uint32")
	runTorusTranslateTest[float64](t, "float64")
}

func runTorusTranslateTest[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		t.Run("CrossesRightEdge", func(t *testing.T) {
			torus := NewTorus(T(10), T(10))
			aabb := newAABB(vec[T](8, 4), T(4), T(2))

			torus.Translate(&aabb, vec[T](0, 0))

			expectAABBState(t, aabb, vec[T](8, 4), vec[T](10, 6), fragmentsForType(map[FragPosition][2]geom.Vec[T]{
				FRAG_RIGHT: {vec[T](0, 4), vec[T](2, 6)},
			}))
		})

		t.Run("HugeShift", func(t *testing.T) {
			torus := NewTorus(T(10), T(10))
			aabb := newAABB(vec[T](8, 4), T(4), T(2))

			torus.Translate(&aabb, vec[T](100, 100))
			expectAABBState(t, aabb, vec[T](8, 4), vec[T](10, 6), fragmentsForType(map[FragPosition][2]geom.Vec[T]{
				FRAG_RIGHT: {vec[T](0, 4), vec[T](2, 6)},
			}))

			torus.Translate(&aabb, vec[T](100, 100))
			expectAABBState(t, aabb, vec[T](8, 4), vec[T](10, 6), fragmentsForType(map[FragPosition][2]geom.Vec[T]{
				FRAG_RIGHT: {vec[T](0, 4), vec[T](2, 6)},
			}))
		})

		t.Run("BackAndForth", func(t *testing.T) {
			torus := NewTorus(T(10), T(10))
			aabb := newAABB(vec[T](8, 8), T(2), T(2))

			torus.Translate(&aabb, vec[T](5, 5))
			expectAABBState(t, aabb, vec[T](3, 3), vec[T](5, 5), map[FragPosition][2]geom.Vec[T]{})

			torus.Translate(&aabb, vec[T](5, 5))
			expectAABBState(t, aabb, vec[T](8, 8), vec[T](10, 10), map[FragPosition][2]geom.Vec[T]{})
		})

		t.Run("CrossesBottomEdge", func(t *testing.T) {
			torus := NewTorus(T(10), T(10))
			aabb := newAABB(vec[T](4, 8), T(2), T(4))

			torus.Translate(&aabb, vec[T](0, 0))

			expectAABBState(t, aabb, vec[T](4, 8), vec[T](6, 10), fragmentsForType(map[FragPosition][2]geom.Vec[T]{
				FRAG_BOTTOM: {vec[T](4, 0), vec[T](6, 2)},
			}))
		})

		t.Run("CrossesCorner", func(t *testing.T) {
			torus := NewTorus(T(10), T(10))
			aabb := newAABB(vec[T](9, 9), T(2), T(2))

			torus.Translate(&aabb, vec[T](0, 0))

			expectAABBState(t, aabb, vec[T](9, 9), vec[T](10, 10), fragmentsForType(map[FragPosition][2]geom.Vec[T]{
				FRAG_RIGHT:        {vec[T](0, 9), vec[T](1, 10)},
				FRAG_BOTTOM:       {vec[T](9, 0), vec[T](10, 1)},
				FRAG_BOTTOM_RIGHT: {vec[T](0, 0), vec[T](1, 1)},
			}))
		})

		t.Run("ClearsFragmentsWhenNotWrapping", func(t *testing.T) {
			torus := NewTorus(T(10), T(10))
			aabb := newAABB(vec[T](8, 4), T(4), T(2))

			torus.Translate(&aabb, vec[T](0, 0))
			expectAABBState(t, aabb, vec[T](8, 4), vec[T](10, 6), fragmentsForType(map[FragPosition][2]geom.Vec[T]{
				FRAG_RIGHT: {vec[T](0, 4), vec[T](2, 6)},
			}))

			torus.Translate(&aabb, vec[T](8, 0))
			expectAABBState(t, aabb, vec[T](6, 4), vec[T](10, 6), nil)
		})

		t.Run("ThroughEdge", func(t *testing.T) {
			torus := NewTorus(T(10), T(10))
			aabb := newAABB(vec[T](2, 2), T(2), T(2))

			torus.Translate(&aabb, vec[T](8, 0))

			expectAABBState(t, aabb, vec[T](0, 2), vec[T](2, 4), nil)
		})

		t.Run("FragmentsMergeSequence", func(t *testing.T) {
			torus := NewTorus(T(10), T(10))
			aabb := newAABB(vec[T](2, 2), T(2), T(2))

			torus.Translate(&aabb, vec[T](7, 0))
			expectAABBState(t, aabb, vec[T](9, 2), vec[T](10, 4), fragmentsForType(map[FragPosition][2]geom.Vec[T]{
				FRAG_RIGHT: {vec[T](0, 2), vec[T](1, 4)},
			}))

			torus.Translate(&aabb, vec[T](0, 7))
			expectAABBState(t, aabb, vec[T](9, 9), vec[T](10, 10), fragmentsForType(map[FragPosition][2]geom.Vec[T]{
				FRAG_RIGHT:        {vec[T](0, 9), vec[T](1, 10)},
				FRAG_BOTTOM:       {vec[T](9, 0), vec[T](10, 1)},
				FRAG_BOTTOM_RIGHT: {vec[T](0, 0), vec[T](1, 1)},
			}))

			torus.Translate(&aabb, vec[T](3, 3))
			expectAABBState(t, aabb, vec[T](2, 2), vec[T](4, 4), nil)
		})
	})
}

func TestCartesianTranslate(t *testing.T) {
	runCartesianTranslateTest[int](t, "int")
	runCartesianTranslateTest[uint32](t, "uint32")
	runCartesianTranslateTest[float64](t, "float64")
}

func runCartesianTranslateTest[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		t.Run("HugeShift", func(t *testing.T) {
			cartesian := NewCartesian(T(10), T(10))
			aabb := newAABB(vec[T](8, 4), T(4), T(2))

			cartesian.Translate(&aabb, vec[T](100, 100))
			expectAABBState(t, aabb, vec[T](10, 10), vec[T](10, 10), map[FragPosition][2]geom.Vec[T]{})
		})

		t.Run("TranslateWithinBounds", func(t *testing.T) {
			cartesian := NewCartesian(T(10), T(10))
			aabb := newAABB(vec[T](1, 1), T(2), T(2))

			cartesian.Translate(&aabb, vec[T](2, 2))
			expectAABBState(t, aabb, vec[T](3, 3), vec[T](5, 5), map[FragPosition][2]geom.Vec[T]{})
		})

		t.Run("ClampAtBoundary", func(t *testing.T) {
			cartesian := NewCartesian(T(10), T(10))
			aabb := newAABB(vec[T](9, 9), T(3), T(3))

			cartesian.Translate(&aabb, vec[T](2, 2))
			expectAABBState(t, aabb, vec[T](10, 10), vec[T](10, 10), map[FragPosition][2]geom.Vec[T]{})
		})
	})
}

func fragmentsForType[T geom.Numeric](frags map[FragPosition][2]geom.Vec[T]) map[FragPosition][2]geom.Vec[T] {
	if frags == nil {
		return nil
	}

	var zero T
	if _, ok := any(zero).(uint32); ok {
		return map[FragPosition][2]geom.Vec[T]{}
	}

	return frags
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

	if len(b.frags) != len(expected) {
		t.Fatalf("expected %d fragments, got %d", len(expected), len(b.frags))
	}

	for pos, want := range expected {
		frag, ok := b.frags[pos]
		if !ok {
			t.Fatalf("missing fragment at %d", pos)
		}
		if !frag.TopLeft.Equals(want[0]) || !frag.BottomRight.Equals(want[1]) {
			t.Fatalf("fragment at %d has bounds %v..%v, expected %v..%v",
				pos, frag.TopLeft, frag.BottomRight, want[0], want[1])
		}
	}
}

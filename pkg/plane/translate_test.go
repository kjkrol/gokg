package plane

import (
	"testing"

	"github.com/kjkrol/gokg/pkg/geom"
)

func TestCyclicPlane_Translate_CrossesRightEdge(t *testing.T) {
	plane := NewTorus(10, 10)
	planeBox := newAABB(geom.NewVec(8, 4), 4, 2)

	plane.Translate(&planeBox, geom.NewVec(0, 0))

	expectPlaneBoxState(t, planeBox, geom.NewVec(8, 4), geom.NewVec(10, 6), map[FragPosition][2]geom.Vec[int]{
		FRAG_RIGHT: {geom.NewVec(0, 4), geom.NewVec(2, 6)},
	})
}

func TestCyclicPlane_Translate_HugeShift(t *testing.T) {
	plane := NewTorus(10, 10)
	planeBox := newAABB(geom.NewVec(8, 4), 4, 2)

	plane.Translate(&planeBox, geom.NewVec(100, 100))
	expectPlaneBoxState(t, planeBox, geom.NewVec(8, 4), geom.NewVec(10, 6), map[FragPosition][2]geom.Vec[int]{
		FRAG_RIGHT: {geom.NewVec(0, 4), geom.NewVec(2, 6)},
	})

	plane.Translate(&planeBox, geom.NewVec(-100, -100))
	expectPlaneBoxState(t, planeBox, geom.NewVec(8, 4), geom.NewVec(10, 6), map[FragPosition][2]geom.Vec[int]{
		FRAG_RIGHT: {geom.NewVec(0, 4), geom.NewVec(2, 6)},
	})
}

func TestBoundedPlane_Translate_HugeShift(t *testing.T) {
	plane := NewCartesian(10, 10)
	planeBox := newAABB(geom.NewVec(8, 4), 4, 2)

	plane.Translate(&planeBox, geom.NewVec(100, 100))
	expectPlaneBoxState(t, planeBox, geom.NewVec(10, 10), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec[int]{})

	plane.Translate(&planeBox, geom.NewVec(-100, -100))
	expectPlaneBoxState(t, planeBox, geom.NewVec(0, 0), geom.NewVec(0, 0), map[FragPosition][2]geom.Vec[int]{})
}
func TestBoundedPlane_Translate_BackAndForth(t *testing.T) {
	plane := NewCartesian(10, 10)
	planeBox := newAABB(geom.NewVec(8, 8), 2, 2)

	plane.Translate(&planeBox, geom.NewVec(5, 5))
	expectPlaneBoxState(t, planeBox, geom.NewVec(10, 10), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec[int]{})

	plane.Translate(&planeBox, geom.NewVec(-5, -5))
	expectPlaneBoxState(t, planeBox, geom.NewVec(5, 5), geom.NewVec(7, 7), map[FragPosition][2]geom.Vec[int]{})
}

func TestCyclicPlane_Translate_BackAndForth(t *testing.T) {
	plane := NewTorus(10, 10)
	planeBox := newAABB(geom.NewVec(8, 8), 2, 2)

	plane.Translate(&planeBox, geom.NewVec(5, 5))
	expectPlaneBoxState(t, planeBox, geom.NewVec(3, 3), geom.NewVec(5, 5), map[FragPosition][2]geom.Vec[int]{})

	plane.Translate(&planeBox, geom.NewVec(-5, -5))
	expectPlaneBoxState(t, planeBox, geom.NewVec(8, 8), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec[int]{})
}

func TestCyclicPlane_Translate_CrossesBottomEdge(t *testing.T) {
	plane := NewTorus(10, 10)
	planeBox := newAABB(geom.NewVec(4, 8), 2, 4)

	plane.Translate(&planeBox, geom.NewVec(0, 0))

	expectPlaneBoxState(t, planeBox, geom.NewVec(4, 8), geom.NewVec(6, 10), map[FragPosition][2]geom.Vec[int]{
		FRAG_BOTTOM: {geom.NewVec(4, 0), geom.NewVec(6, 2)},
	})
}

func TestCyclicPlane_Translate_CrossesCorner(t *testing.T) {
	plane := NewTorus(10, 10)
	planeBox := newAABB(geom.NewVec(9, 9), 2, 2)

	plane.Translate(&planeBox, geom.NewVec(0, 0))

	expectPlaneBoxState(t, planeBox, geom.NewVec(9, 9), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec[int]{
		FRAG_RIGHT:        {geom.NewVec(0, 9), geom.NewVec(1, 10)},
		FRAG_BOTTOM:       {geom.NewVec(9, 0), geom.NewVec(10, 1)},
		FRAG_BOTTOM_RIGHT: {geom.NewVec(0, 0), geom.NewVec(1, 1)},
	})
}

func TestCyclicPlane_Translate_ClearsFragmentsWhenNotWrapping(t *testing.T) {
	plane := NewTorus(10, 10)
	planeBox := newAABB(geom.NewVec(8, 4), 4, 2)

	plane.Translate(&planeBox, geom.NewVec(0, 0))
	expectPlaneBoxState(t, planeBox, geom.NewVec(8, 4), geom.NewVec(10, 6), map[FragPosition][2]geom.Vec[int]{
		FRAG_RIGHT: {geom.NewVec(0, 4), geom.NewVec(2, 6)},
	})

	plane.Translate(&planeBox, geom.NewVec(-2, 0))
	expectPlaneBoxState(t, planeBox, geom.NewVec(6, 4), geom.NewVec(10, 6), nil)
}

func TestCyclicPlane_Translate_ThroughEdge(t *testing.T) {
	plane := NewTorus(10, 10)
	planeBox := newAABB(geom.NewVec(2, 2), 2, 2)

	plane.Translate(&planeBox, geom.NewVec(8, 0))

	expectPlaneBoxState(t, planeBox, geom.NewVec(0, 2), geom.NewVec(2, 4), nil)
}

func TestCyclicPlane_Translate_FragmentsMergeSequence(t *testing.T) {
	plane := NewTorus(10, 10)
	planeBox := newAABB(geom.NewVec(2, 2), 2, 2)

	plane.Translate(&planeBox, geom.NewVec(-3, 0))
	expectPlaneBoxState(t, planeBox, geom.NewVec(9, 2), geom.NewVec(10, 4), map[FragPosition][2]geom.Vec[int]{
		FRAG_RIGHT: {geom.NewVec(0, 2), geom.NewVec(1, 4)},
	})

	plane.Translate(&planeBox, geom.NewVec(0, -3))
	expectPlaneBoxState(t, planeBox, geom.NewVec(9, 9), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec[int]{
		FRAG_RIGHT:        {geom.NewVec(0, 9), geom.NewVec(1, 10)},
		FRAG_BOTTOM:       {geom.NewVec(9, 0), geom.NewVec(10, 1)},
		FRAG_BOTTOM_RIGHT: {geom.NewVec(0, 0), geom.NewVec(1, 1)},
	})

	plane.Translate(&planeBox, geom.NewVec(3, 3))
	expectPlaneBoxState(t, planeBox, geom.NewVec(2, 2), geom.NewVec(4, 4), nil)
}

func expectPlaneBoxState(
	t *testing.T,
	b AABB[int],
	expectedPos geom.Vec[int],
	expectedBottomRight geom.Vec[int],
	expectedFragments map[FragPosition][2]geom.Vec[int],
) {
	t.Helper()

	if !b.TopLeft.Equals(expectedPos) {
		t.Fatalf("expected bound position %v, got %v", expectedPos, b.TopLeft)
	}

	if !b.BottomRight.Equals(expectedBottomRight) {
		t.Fatalf("expected bound bottom-right %v, got %v", expectedBottomRight, b.BottomRight)
	}

	expectPlaneBoxFragments(t, b, expectedFragments)
}

func expectPlaneBoxFragments(t *testing.T, b AABB[int], expected map[FragPosition][2]geom.Vec[int]) {
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

package plane

import (
	"testing"

	"github.com/kjkrol/gokg/pkg/geom"
)

func TestTorus_Translate_CrossesRightEdge(t *testing.T) {
	torus := NewTorus(10, 10)
	aabb := newAABB(geom.NewVec(8, 4), 4, 2)

	torus.Translate(&aabb, geom.NewVec(0, 0))

	expectAABBState(t, aabb, geom.NewVec(8, 4), geom.NewVec(10, 6), map[FragPosition][2]geom.Vec[int]{
		FRAG_RIGHT: {geom.NewVec(0, 4), geom.NewVec(2, 6)},
	})
}

func TestTorus_Translate_HugeShift(t *testing.T) {
	torus := NewTorus(10, 10)
	aabb := newAABB(geom.NewVec(8, 4), 4, 2)

	torus.Translate(&aabb, geom.NewVec(100, 100))
	expectAABBState(t, aabb, geom.NewVec(8, 4), geom.NewVec(10, 6), map[FragPosition][2]geom.Vec[int]{
		FRAG_RIGHT: {geom.NewVec(0, 4), geom.NewVec(2, 6)},
	})

	torus.Translate(&aabb, geom.NewVec(-100, -100))
	expectAABBState(t, aabb, geom.NewVec(8, 4), geom.NewVec(10, 6), map[FragPosition][2]geom.Vec[int]{
		FRAG_RIGHT: {geom.NewVec(0, 4), geom.NewVec(2, 6)},
	})
}

func TestCartesian_Translate_HugeShift(t *testing.T) {
	cartesian := NewCartesian(10, 10)
	aabb := newAABB(geom.NewVec(8, 4), 4, 2)

	cartesian.Translate(&aabb, geom.NewVec(100, 100))
	expectAABBState(t, aabb, geom.NewVec(10, 10), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec[int]{})

	cartesian.Translate(&aabb, geom.NewVec(-100, -100))
	expectAABBState(t, aabb, geom.NewVec(0, 0), geom.NewVec(0, 0), map[FragPosition][2]geom.Vec[int]{})
}
func TestCartesian_Translate_BackAndForth(t *testing.T) {
	cartesian := NewCartesian(10, 10)
	aabb := newAABB(geom.NewVec(8, 8), 2, 2)

	cartesian.Translate(&aabb, geom.NewVec(5, 5))
	expectAABBState(t, aabb, geom.NewVec(10, 10), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec[int]{})

	cartesian.Translate(&aabb, geom.NewVec(-5, -5))
	expectAABBState(t, aabb, geom.NewVec(5, 5), geom.NewVec(7, 7), map[FragPosition][2]geom.Vec[int]{})
}

func TestTorus_Translate_BackAndForth(t *testing.T) {
	torus := NewTorus(10, 10)
	aabb := newAABB(geom.NewVec(8, 8), 2, 2)

	torus.Translate(&aabb, geom.NewVec(5, 5))
	expectAABBState(t, aabb, geom.NewVec(3, 3), geom.NewVec(5, 5), map[FragPosition][2]geom.Vec[int]{})

	torus.Translate(&aabb, geom.NewVec(-5, -5))
	expectAABBState(t, aabb, geom.NewVec(8, 8), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec[int]{})
}

func TestTorus_Translate_CrossesBottomEdge(t *testing.T) {
	torus := NewTorus(10, 10)
	aabb := newAABB(geom.NewVec(4, 8), 2, 4)

	torus.Translate(&aabb, geom.NewVec(0, 0))

	expectAABBState(t, aabb, geom.NewVec(4, 8), geom.NewVec(6, 10), map[FragPosition][2]geom.Vec[int]{
		FRAG_BOTTOM: {geom.NewVec(4, 0), geom.NewVec(6, 2)},
	})
}

func TestTorus_Translate_CrossesCorner(t *testing.T) {
	torus := NewTorus(10, 10)
	aabb := newAABB(geom.NewVec(9, 9), 2, 2)

	torus.Translate(&aabb, geom.NewVec(0, 0))

	expectAABBState(t, aabb, geom.NewVec(9, 9), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec[int]{
		FRAG_RIGHT:        {geom.NewVec(0, 9), geom.NewVec(1, 10)},
		FRAG_BOTTOM:       {geom.NewVec(9, 0), geom.NewVec(10, 1)},
		FRAG_BOTTOM_RIGHT: {geom.NewVec(0, 0), geom.NewVec(1, 1)},
	})
}

func TestTorus_Translate_ClearsFragmentsWhenNotWrapping(t *testing.T) {
	torus := NewTorus(10, 10)
	aabb := newAABB(geom.NewVec(8, 4), 4, 2)

	torus.Translate(&aabb, geom.NewVec(0, 0))
	expectAABBState(t, aabb, geom.NewVec(8, 4), geom.NewVec(10, 6), map[FragPosition][2]geom.Vec[int]{
		FRAG_RIGHT: {geom.NewVec(0, 4), geom.NewVec(2, 6)},
	})

	torus.Translate(&aabb, geom.NewVec(-2, 0))
	expectAABBState(t, aabb, geom.NewVec(6, 4), geom.NewVec(10, 6), nil)
}

func TestTorus_Translate_ThroughEdge(t *testing.T) {
	torus := NewTorus(10, 10)
	aabb := newAABB(geom.NewVec(2, 2), 2, 2)

	torus.Translate(&aabb, geom.NewVec(8, 0))

	expectAABBState(t, aabb, geom.NewVec(0, 2), geom.NewVec(2, 4), nil)
}

func TestTorus_Translate_FragmentsMergeSequence(t *testing.T) {
	torus := NewTorus(10, 10)
	aabb := newAABB(geom.NewVec(2, 2), 2, 2)

	torus.Translate(&aabb, geom.NewVec(-3, 0))
	expectAABBState(t, aabb, geom.NewVec(9, 2), geom.NewVec(10, 4), map[FragPosition][2]geom.Vec[int]{
		FRAG_RIGHT: {geom.NewVec(0, 2), geom.NewVec(1, 4)},
	})

	torus.Translate(&aabb, geom.NewVec(0, -3))
	expectAABBState(t, aabb, geom.NewVec(9, 9), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec[int]{
		FRAG_RIGHT:        {geom.NewVec(0, 9), geom.NewVec(1, 10)},
		FRAG_BOTTOM:       {geom.NewVec(9, 0), geom.NewVec(10, 1)},
		FRAG_BOTTOM_RIGHT: {geom.NewVec(0, 0), geom.NewVec(1, 1)},
	})

	torus.Translate(&aabb, geom.NewVec(3, 3))
	expectAABBState(t, aabb, geom.NewVec(2, 2), geom.NewVec(4, 4), nil)
}

func expectAABBState(
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

	expectAABBFragments(t, b, expectedFragments)
}

func expectAABBFragments(t *testing.T, b AABB[int], expected map[FragPosition][2]geom.Vec[int]) {
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

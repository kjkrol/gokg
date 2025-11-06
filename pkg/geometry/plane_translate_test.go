package geometry

import "testing"

func TestAABBTranslate_CrossesRightEdge(t *testing.T) {
	plane := NewCyclicBoundedPlane(10, 10)
	planeBox := NewPlaneBox(NewVec(8, 4), 4, 2)

	plane.Translate(&planeBox, NewVec(0, 0))

	expectAABBState(t, planeBox, NewVec(8, 4), NewVec(10, 6), map[FragPosition][2]Vec[int]{
		FRAG_RIGHT: {NewVec(0, 4), NewVec(2, 6)},
	})
}

func TestAABBTranslate_CrossesBottomEdge(t *testing.T) {
	plane := NewCyclicBoundedPlane(10, 10)
	planeBox := NewPlaneBox(NewVec(4, 8), 2, 4)

	plane.Translate(&planeBox, NewVec(0, 0))

	expectAABBState(t, planeBox, NewVec(4, 8), NewVec(6, 10), map[FragPosition][2]Vec[int]{
		FRAG_BOTTOM: {NewVec(4, 0), NewVec(6, 2)},
	})
}

func TestAABBTranslate_CrossesCorner(t *testing.T) {
	plane := NewCyclicBoundedPlane(10, 10)
	planeBox := NewPlaneBox(NewVec(9, 9), 2, 2)

	plane.Translate(&planeBox, NewVec(0, 0))

	expectAABBState(t, planeBox, NewVec(9, 9), NewVec(10, 10), map[FragPosition][2]Vec[int]{
		FRAG_RIGHT:        {NewVec(0, 9), NewVec(1, 10)},
		FRAG_BOTTOM:       {NewVec(9, 0), NewVec(10, 1)},
		FRAG_BOTTOM_RIGHT: {NewVec(0, 0), NewVec(1, 1)},
	})
}

func TestAABBTranslate_ClearsFragmentsWhenNotWrapping(t *testing.T) {
	plane := NewCyclicBoundedPlane(10, 10)
	planeBox := NewPlaneBox(NewVec(8, 4), 4, 2)

	plane.Translate(&planeBox, NewVec(0, 0))
	expectAABBState(t, planeBox, NewVec(8, 4), NewVec(10, 6), map[FragPosition][2]Vec[int]{
		FRAG_RIGHT: {NewVec(0, 4), NewVec(2, 6)},
	})

	plane.Translate(&planeBox, NewVec(-2, 0))
	expectAABBState(t, planeBox, NewVec(6, 4), NewVec(10, 6), nil)
}

func TestAABBTranslate_ThroughEdge(t *testing.T) {
	plane := NewCyclicBoundedPlane(10, 10)
	planeBox := NewPlaneBox(NewVec(2, 2), 2, 2)

	plane.Translate(&planeBox, NewVec(8, 0))

	expectAABBState(t, planeBox, NewVec(0, 2), NewVec(2, 4), nil)
}

func TestAABBTranslate_FragmentsMergeSequence(t *testing.T) {
	plane := NewCyclicBoundedPlane(10, 10)
	planeBox := NewPlaneBox(NewVec(2, 2), 2, 2)

	plane.Translate(&planeBox, NewVec(-3, 0))
	expectAABBState(t, planeBox, NewVec(9, 2), NewVec(10, 4), map[FragPosition][2]Vec[int]{
		FRAG_RIGHT: {NewVec(0, 2), NewVec(1, 4)},
	})

	plane.Translate(&planeBox, NewVec(0, -3))
	expectAABBState(t, planeBox, NewVec(9, 9), NewVec(10, 10), map[FragPosition][2]Vec[int]{
		FRAG_RIGHT:        {NewVec(0, 9), NewVec(1, 10)},
		FRAG_BOTTOM:       {NewVec(9, 0), NewVec(10, 1)},
		FRAG_BOTTOM_RIGHT: {NewVec(0, 0), NewVec(1, 1)},
	})

	plane.Translate(&planeBox, NewVec(3, 3))
	expectAABBState(t, planeBox, NewVec(2, 2), NewVec(4, 4), nil)
}

func expectAABBState(
	t *testing.T,
	b PlaneBox[int],
	expectedPos Vec[int],
	expectedBottomRight Vec[int],
	expectedFragments map[FragPosition][2]Vec[int],
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

func expectAABBFragments(t *testing.T, b PlaneBox[int], expected map[FragPosition][2]Vec[int]) {
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

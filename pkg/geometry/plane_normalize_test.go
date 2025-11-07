package geometry

import "testing"

func TestPlaneBounded_normalizeBox(t *testing.T) {
	plane := NewBoundedPlane(10, 10)

	planeBox := newPlaneBox(NewVec(-2, -2), 2, 2)
	plane.normalizeBox(&planeBox)
	expectPlaneBoxState(t, planeBox, NewVec(0, 0), NewVec(0, 0), map[FragPosition][2]Vec[int]{})

	planeBox = newPlaneBox(NewVec(-1, -1), 2, 2)
	plane.normalizeBox(&planeBox)
	expectPlaneBoxState(t, planeBox, NewVec(0, 0), NewVec(1, 1), map[FragPosition][2]Vec[int]{})

	planeBox = newPlaneBox(NewVec(9, 9), 2, 2)
	plane.normalizeBox(&planeBox)
	expectPlaneBoxState(t, planeBox, NewVec(9, 9), NewVec(10, 10), map[FragPosition][2]Vec[int]{})

	planeBox = newPlaneBox(NewVec(-11, -11), 2, 2)
	plane.normalizeBox(&planeBox)
	expectPlaneBoxState(t, planeBox, NewVec(0, 0), NewVec(0, 0), map[FragPosition][2]Vec[int]{})

	planeBox = newPlaneBox(NewVec(19, 19), 2, 2)
	plane.normalizeBox(&planeBox)
	expectPlaneBoxState(t, planeBox, NewVec(10, 10), NewVec(10, 10), map[FragPosition][2]Vec[int]{})
}

func TestPlaneCyclic_normalizeBox(t *testing.T) {
	plane := NewCyclicBoundedPlane(10, 10)

	planeBox := newPlaneBox(NewVec(-2, -2), 2, 2)
	plane.normalizeBox(&planeBox)
	expectPlaneBoxState(t, planeBox, NewVec(8, 8), NewVec(10, 10), map[FragPosition][2]Vec[int]{})

	planeBox = newPlaneBox(NewVec(-1, -1), 2, 2)
	plane.normalizeBox(&planeBox)
	expectPlaneBoxState(t, planeBox, NewVec(9, 9), NewVec(10, 10), map[FragPosition][2]Vec[int]{
		FRAG_RIGHT:        {NewVec(0, 9), NewVec(1, 10)},
		FRAG_BOTTOM:       {NewVec(9, 0), NewVec(10, 1)},
		FRAG_BOTTOM_RIGHT: {NewVec(0, 0), NewVec(1, 1)},
	})

	planeBox = newPlaneBox(NewVec(9, 9), 2, 2)
	plane.normalizeBox(&planeBox)
	expectPlaneBoxState(t, planeBox, NewVec(9, 9), NewVec(10, 10), map[FragPosition][2]Vec[int]{
		FRAG_RIGHT:        {NewVec(0, 9), NewVec(1, 10)},
		FRAG_BOTTOM:       {NewVec(9, 0), NewVec(10, 1)},
		FRAG_BOTTOM_RIGHT: {NewVec(0, 0), NewVec(1, 1)},
	})

	planeBox = newPlaneBox(NewVec(-11, -11), 2, 2)
	plane.normalizeBox(&planeBox)
	expectPlaneBoxState(t, planeBox, NewVec(9, 9), NewVec(10, 10), map[FragPosition][2]Vec[int]{
		FRAG_RIGHT:        {NewVec(0, 9), NewVec(1, 10)},
		FRAG_BOTTOM:       {NewVec(9, 0), NewVec(10, 1)},
		FRAG_BOTTOM_RIGHT: {NewVec(0, 0), NewVec(1, 1)},
	})

	planeBox = newPlaneBox(NewVec(19, 19), 2, 2)
	plane.normalizeBox(&planeBox)
	expectPlaneBoxState(t, planeBox, NewVec(9, 9), NewVec(10, 10), map[FragPosition][2]Vec[int]{
		FRAG_RIGHT:        {NewVec(0, 9), NewVec(1, 10)},
		FRAG_BOTTOM:       {NewVec(9, 0), NewVec(10, 1)},
		FRAG_BOTTOM_RIGHT: {NewVec(0, 0), NewVec(1, 1)},
	})
}

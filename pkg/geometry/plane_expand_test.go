package geometry

import (
	"testing"
)

func TestPlane_Expand_On_BoundedPlane(t *testing.T) {
	plane := NewBoundedPlane(10, 10)
	planeBox := newPlaneBox(NewVec(2, 3), 3, 4)
	plane.Expand(&planeBox, 2)
	expectPlaneBoxState(t, planeBox, NewVec(0, 1), NewVec(7, 9), map[FragPosition][2]Vec[int]{})
}

func TestPlane_Expand_On_BoundedPlane_CornerCase(t *testing.T) {
	plane := NewBoundedPlane(10, 10)
	planeBox := newPlaneBox(NewVec(0, 0), 2, 2)
	plane.Expand(&planeBox, 2)
	expectPlaneBoxState(t, planeBox, NewVec(0, 0), NewVec(4, 4), map[FragPosition][2]Vec[int]{})
}

func TestPlane_Expand_On_CyclicPlane_CornerCase(t *testing.T) {
	plane := NewCyclicBoundedPlane(10, 10)
	planeBox := newPlaneBox(NewVec(0, 0), 2, 2)
	plane.Expand(&planeBox, 2)
	expectPlaneBoxState(t, planeBox, NewVec(8, 8), NewVec(10, 10), map[FragPosition][2]Vec[int]{
		FRAG_RIGHT:        {NewVec(0, 8), NewVec(4, 10)},
		FRAG_BOTTOM:       {NewVec(8, 0), NewVec(10, 4)},
		FRAG_BOTTOM_RIGHT: {NewVec(0, 0), NewVec(4, 4)},
	})
}

func TestPlane_Exapnad_ThenIntersects(t *testing.T) {
	plane := NewCyclicBoundedPlane(100, 100)

	rect1 := newPlaneBox(NewVec(5, 5), 10, 10)
	rect2 := newPlaneBox(NewVec(96, 96), 10, 10)

	plane.Expand(&rect2, 0)

	if !rect1.Intersects(rect2) {
		t.Errorf("rect1 %v should intersect with rect2 %v", rect1, rect2)
	}

}

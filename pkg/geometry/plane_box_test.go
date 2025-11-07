package geometry

import "testing"

func TestPlaneBox_NewPlaneBox(t *testing.T) {
	planeBox := newPlaneBox(ZERO_INT_VEC, 10, 10)
	expected := NewVec(10, 10)
	if planeBox.BottomRight != expected {
		t.Errorf("center %v not equal to expected %v", planeBox.BottomRight, expected)
	}
}

func TestPlaneBox_IntersectsIncludingFrags_ReturnsTrue(t *testing.T) {
	base := newPlaneBox(NewVec(0, 0), 2, 2)
	other := newPlaneBox(NewVec(4, 4), 1, 1)
	other.frags[FRAG_RIGHT] = NewBoundingBox(NewVec(0, 4), NewVec(1, 5))
	other.frags[FRAG_BOTTOM] = NewBoundingBox(NewVec(4, 0), NewVec(5, 1))
	other.frags[FRAG_BOTTOM_RIGHT] = NewBoundingBox(NewVec(0, 0), NewVec(1, 1))

	if !base.Intersects(other) {
		t.Errorf("expected IntersectsAny to return true, but got false")
	}
}

func TestPlaneBox_IntersectsIncludingFrags_ReturnsFalse(t *testing.T) {
	base := newPlaneBox(NewVec(0, 0), 2, 2)
	other := newPlaneBox(NewVec(4, 4), 2, 2)
	other.frags[FRAG_RIGHT] = NewBoundingBox(NewVec(0, 4), NewVec(1, 6))

	if base.Intersects(other) {
		t.Errorf("expected IntersectsAny to return false, but got true")
	}
}

func TestPlaneBox_Contains(t *testing.T) {
	outer := newPlaneBox(NewVec(0, 0), 10, 10)
	inner := newPlaneBox(NewVec(2, 2), 6, 6)
	onlyTopLeftInside := newPlaneBox(NewVec(-1, -1), 4, 4)
	onlyBottomRightInside := newPlaneBox(NewVec(5, 5), 7, 7)

	if !outer.Contains(inner) {
		t.Errorf("expected outer to contain inner")
	}
	if outer.Contains(onlyTopLeftInside) {
		t.Errorf("expected outer not to contain rectangle with outside top-left")
	}
	if outer.Contains(onlyBottomRightInside) {
		t.Errorf("expected outer not to contain rectangle with outside bottom-right")
	}
}

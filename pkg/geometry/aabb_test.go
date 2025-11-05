package geometry

import "testing"

func TestAABB_NewAABB(t *testing.T) {
	aabb := NewAABB(ZERO_INT_VEC, 10, 10)
	expected := NewVec(10, 10)
	if aabb.BottomRight != expected {
		t.Errorf("center %v not equal to expected %v", aabb.BottomRight, expected)
	}
}

func TestAABB_BuildAABB(t *testing.T) {
	center := NewVec(5, 5)
	aabb := BuildAABB(center, 2)
	expectedTopLeft := NewVec(3, 3)
	expectedBottomRight := NewVec(7, 7)
	if aabb.TopLeft != expectedTopLeft {
		t.Errorf("topLeft %v not equal to expected %v", aabb.TopLeft, expectedTopLeft)
	}
	if aabb.BottomRight != expectedBottomRight {
		t.Errorf("bottomRight %v not equal to expected %v", aabb.BottomRight, expectedBottomRight)
	}
}

func TestAABB_Split(t *testing.T) {
	parent := NewAABB(ZERO_INT_VEC, 10, 10)
	splitted := parent.Split()

	expected := [4]AABB[int]{
		NewAABB(ZERO_INT_VEC, 5, 5),
		NewAABB(NewVec(5, 0), 5, 5),
		NewAABB(NewVec(0, 5), 5, 5),
		NewAABB(NewVec(5, 5), 5, 5),
	}
	for i := 0; i < 4; i++ {
		if !splitted[i].Equals(expected[i]) {
			t.Errorf("split %v not equal to expected %v", splitted[i], expected[i])
		}
	}
}

func TestAABB_Intersects(t *testing.T) {
	intersects := []struct{ rect1, rect2 AABB[int] }{
		{
			rect1: NewAABB(Vec[int]{X: 3, Y: 3}, 2, 2),
			rect2: NewAABB(Vec[int]{X: 2, Y: 2}, 2, 2),
		},
		{
			rect1: NewAABB(Vec[int]{X: 3, Y: 3}, 2, 2),
			rect2: NewAABB(Vec[int]{X: 4, Y: 2}, 2, 2),
		},
		{
			rect1: NewAABB(Vec[int]{X: 3, Y: 3}, 2, 2),
			rect2: NewAABB(Vec[int]{X: 2, Y: 4}, 2, 2),
		},
		{
			rect1: NewAABB(Vec[int]{X: 3, Y: 3}, 2, 2),
			rect2: NewAABB(Vec[int]{X: 4, Y: 4}, 2, 2),
		},
		{
			rect1: NewAABB(Vec[int]{X: 3, Y: 3}, 2, 2),
			rect2: NewAABB(Vec[int]{X: 2, Y: 2}, 2, 2),
		},
		{
			rect1: NewAABB(Vec[int]{X: 3, Y: 3}, 2, 2),
			rect2: NewAABB(Vec[int]{X: 2, Y: 2}, 2, 2),
		},
		{
			rect1: NewAABB(Vec[int]{X: 3, Y: 3}, 2, 2),
			rect2: NewAABB(Vec[int]{X: 2, Y: 4}, 2, 2),
		},
		{
			rect1: NewAABB(Vec[int]{X: 3, Y: 3}, 2, 2),
			rect2: NewAABB(Vec[int]{X: 2, Y: 2}, 2, 2),
		},
		{
			rect1: NewAABB(Vec[int]{X: 3, Y: 3}, 2, 2),
			rect2: NewAABB(Vec[int]{X: 2, Y: 5}, 2, 2),
		},
		{
			rect1: NewAABB(Vec[int]{X: 3, Y: 3}, 2, 2),
			rect2: NewAABB(Vec[int]{X: 5, Y: 2}, 2, 2),
		},
		{
			rect1: NewAABB(Vec[int]{X: 3, Y: 3}, 2, 2),
			rect2: NewAABB(Vec[int]{X: 5, Y: 2}, 2, 2),
		},
	}

	for _, intersection := range intersects {
		if !intersection.rect1.Intersects(intersection.rect2) {
			t.Errorf("rect1 %v should intersect with rect2 %v", intersection.rect1, intersection.rect2)
		}
		if !intersection.rect2.Intersects(intersection.rect1) {
			t.Errorf("rect2 %v should intersect with rect1 %v", intersection.rect2, intersection.rect1)
		}
	}

	notIntersects := []struct{ rect1, rect2 AABB[int] }{
		{
			rect1: NewAABB(Vec[int]{X: 3, Y: 3}, 2, 2),
			rect2: NewAABB(Vec[int]{X: 1, Y: 1}, 1, 1),
		},
		{
			rect1: NewAABB(Vec[int]{X: 3, Y: 3}, 2, 2),
			rect2: NewAABB(Vec[int]{X: 6, Y: 0}, 3, 9),
		},
		{
			rect1: NewAABB(Vec[int]{X: 3, Y: 3}, 2, 2),
			rect2: NewAABB(Vec[int]{X: 0, Y: 6}, 9, 3),
		},
	}
	for _, intersection := range notIntersects {
		if intersection.rect1.Intersects(intersection.rect2) {
			t.Errorf("rect1 %v should not intersect with rect2 %v", intersection.rect1, intersection.rect2)
		}
		if intersection.rect2.Intersects(intersection.rect1) {
			t.Errorf("rect2 %v should not intersect with rect1 %v", intersection.rect2, intersection.rect1)
		}
	}
}

func TestAABB_IntersectsIncludingFrags_ReturnsTrue(t *testing.T) {
	base := NewAABB(NewVec(0, 0), 2, 2)
	other := NewAABB(NewVec(4, 4), 1, 1)
	other.frags[BOUND_RIGHT] = newAABBFrag(NewVec(0, 4), NewVec(1, 5))
	other.frags[BOUND_BOTTOM] = newAABBFrag(NewVec(4, 0), NewVec(5, 1))
	other.frags[BOUND_BOTTOM_RIGHT] = newAABBFrag(NewVec(0, 0), NewVec(1, 1))

	if !base.IntersectsIncludingFrags(other) {
		t.Errorf("expected IntersectsAny to return true, but got false")
	}
}

func TestAABB_IntersectsIncludingFrags_ReturnsFalse(t *testing.T) {
	base := NewAABB(NewVec(0, 0), 2, 2)
	other := NewAABB(NewVec(4, 4), 2, 2)
	other.frags[BOUND_RIGHT] = newAABBFrag(NewVec(0, 4), NewVec(1, 6))

	if base.IntersectsIncludingFrags(other) {
		t.Errorf("expected IntersectsAny to return false, but got true")
	}
}

func TestAABB_Contains(t *testing.T) {
	outer := NewAABB(NewVec(0, 0), 10, 10)
	inner := NewAABB(NewVec(2, 2), 6, 6)
	onlyTopLeftInside := NewAABB(NewVec(-1, -1), 4, 4)
	onlyBottomRightInside := NewAABB(NewVec(5, 5), 7, 7)

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

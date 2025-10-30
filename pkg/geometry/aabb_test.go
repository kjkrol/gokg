package geometry

import "testing"

func TestRectangle_NewRectangle(t *testing.T) {
	rect := NewAABB(ZERO_INT_VEC, Vec[int]{X: 10, Y: 10})
	expected := Vec[int]{X: 5, Y: 5}
	if rect.Center != expected {
		t.Errorf("center %v not equal to expected %v", rect.Center, expected)
	}
}

func TestRectangle_BuildRectangle(t *testing.T) {
	center := Vec[int]{X: 5, Y: 5}
	rect := BuildAABB(center, 2)
	expectedTopLeft := Vec[int]{X: 3, Y: 3}
	expectedBottomRight := Vec[int]{X: 7, Y: 7}
	if rect.TopLeft != expectedTopLeft {
		t.Errorf("topLeft %v not equal to expected %v", rect.TopLeft, expectedTopLeft)
	}
	if rect.BottomRight != expectedBottomRight {
		t.Errorf("bottomRight %v not equal to expected %v", rect.BottomRight, expectedBottomRight)
	}
}

func TestRectangle_Split(t *testing.T) {
	parent := NewAABB(ZERO_INT_VEC, Vec[int]{X: 10, Y: 10})
	splitted := parent.Split()

	expected := [4]AABB[int]{
		NewAABB(ZERO_INT_VEC, Vec[int]{X: 5, Y: 5}),
		NewAABB(Vec[int]{X: 5, Y: 0}, Vec[int]{X: 10, Y: 5}),
		NewAABB(Vec[int]{X: 0, Y: 5}, Vec[int]{X: 5, Y: 10}),
		NewAABB(Vec[int]{X: 5, Y: 5}, Vec[int]{X: 10, Y: 10}),
	}
	for i := 0; i < 4; i++ {
		if !splitted[i].Equals(expected[i]) {
			t.Errorf("split %v not equal to expected %v", splitted[i], expected[i])
		}
	}
}

func TestRectangle_Intersects(t *testing.T) {
	intersects := []struct{ rect1, rect2 AABB[int] }{
		{
			rect1: NewAABB(Vec[int]{X: 3, Y: 3}, Vec[int]{X: 5, Y: 5}),
			rect2: NewAABB(Vec[int]{X: 2, Y: 2}, Vec[int]{X: 4, Y: 4}),
		},
		{
			rect1: NewAABB(Vec[int]{X: 3, Y: 3}, Vec[int]{X: 5, Y: 5}),
			rect2: NewAABB(Vec[int]{X: 4, Y: 2}, Vec[int]{X: 6, Y: 4}),
		},
		{
			rect1: NewAABB(Vec[int]{X: 3, Y: 3}, Vec[int]{X: 5, Y: 5}),
			rect2: NewAABB(Vec[int]{X: 2, Y: 4}, Vec[int]{X: 4, Y: 6}),
		},
		{
			rect1: NewAABB(Vec[int]{X: 3, Y: 3}, Vec[int]{X: 5, Y: 5}),
			rect2: NewAABB(Vec[int]{X: 4, Y: 4}, Vec[int]{X: 6, Y: 6}),
		},
		{
			rect1: NewAABB(Vec[int]{X: 3, Y: 3}, Vec[int]{X: 5, Y: 5}),
			rect2: NewAABB(Vec[int]{X: 2, Y: 2}, Vec[int]{X: 6, Y: 6}),
		},
		{
			rect1: NewAABB(Vec[int]{X: 3, Y: 3}, Vec[int]{X: 5, Y: 5}),
			rect2: NewAABB(Vec[int]{X: 2, Y: 2}, Vec[int]{X: 3, Y: 4}),
		},
		{
			rect1: NewAABB(Vec[int]{X: 3, Y: 3}, Vec[int]{X: 5, Y: 5}),
			rect2: NewAABB(Vec[int]{X: 2, Y: 4}, Vec[int]{X: 3, Y: 6}),
		},
		{
			rect1: NewAABB(Vec[int]{X: 3, Y: 3}, Vec[int]{X: 5, Y: 5}),
			rect2: NewAABB(Vec[int]{X: 2, Y: 2}, Vec[int]{X: 6, Y: 3}),
		},
		{
			rect1: NewAABB(Vec[int]{X: 3, Y: 3}, Vec[int]{X: 5, Y: 5}),
			rect2: NewAABB(Vec[int]{X: 2, Y: 5}, Vec[int]{X: 6, Y: 6}),
		},
		{
			rect1: NewAABB(Vec[int]{X: 3, Y: 3}, Vec[int]{X: 5, Y: 5}),
			rect2: NewAABB(Vec[int]{X: 5, Y: 2}, Vec[int]{X: 6, Y: 6}),
		},
		{
			rect1: NewAABB(Vec[int]{X: 3, Y: 3}, Vec[int]{X: 5, Y: 5}),
			rect2: NewAABB(Vec[int]{X: 5, Y: 2}, Vec[int]{X: 6, Y: 6}),
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
			rect1: NewAABB(Vec[int]{X: 3, Y: 3}, Vec[int]{X: 5, Y: 5}),
			rect2: NewAABB(Vec[int]{X: 1, Y: 1}, Vec[int]{X: 2, Y: 2}),
		},
		{
			rect1: NewAABB(Vec[int]{X: 3, Y: 3}, Vec[int]{X: 5, Y: 5}),
			rect2: NewAABB(Vec[int]{X: 6, Y: 0}, Vec[int]{X: 9, Y: 9}),
		},
		{
			rect1: NewAABB(Vec[int]{X: 3, Y: 3}, Vec[int]{X: 5, Y: 5}),
			rect2: NewAABB(Vec[int]{X: 0, Y: 6}, Vec[int]{X: 9, Y: 9}),
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

func TestRectangle_IntersectsAny_ReturnsFalse(t *testing.T) {
	base := AABB[int]{
		TopLeft:     Vec[int]{X: 0, Y: 0},
		BottomRight: Vec[int]{X: 10, Y: 10},
		Center:      Vec[int]{X: 5, Y: 5},
	}

	others := []AABB[int]{
		{TopLeft: Vec[int]{X: 20, Y: 20}, BottomRight: Vec[int]{X: 30, Y: 30}},
		{TopLeft: Vec[int]{X: 40, Y: 0}, BottomRight: Vec[int]{X: 50, Y: 10}},
		{TopLeft: Vec[int]{X: 0, Y: 40}, BottomRight: Vec[int]{X: 10, Y: 50}},
	}

	if base.IntersectsAny(others) {
		t.Errorf("expected IntersectsAny to return false, but got true")
	}
}

func TestRectangle_Contains(t *testing.T) {
	outer := NewAABB(Vec[int]{X: 0, Y: 0}, Vec[int]{X: 10, Y: 10})
	inner := NewAABB(Vec[int]{X: 2, Y: 2}, Vec[int]{X: 8, Y: 8})
	onlyTopLeftInside := NewAABB(Vec[int]{X: -1, Y: -1}, Vec[int]{X: 5, Y: 5})
	onlyBottomRightInside := NewAABB(Vec[int]{X: 5, Y: 5}, Vec[int]{X: 12, Y: 12})

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

func TestRectangle_Expand(t *testing.T) {
	rect := NewAABB(Vec[int]{X: 2, Y: 3}, Vec[int]{X: 5, Y: 7})
	expanded := rect.Expand(2)
	expectedTopLeft := Vec[int]{X: 0, Y: 1}
	expectedBottomRight := Vec[int]{X: 7, Y: 9}
	if expanded.TopLeft != expectedTopLeft || expanded.BottomRight != expectedBottomRight {
		t.Errorf("unexpected expanded rectangle: %+v", expanded)
	}
}

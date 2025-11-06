package geometry

import "testing"

func TestBox_Split(t *testing.T) {
	parent := NewBoundingBoxAt(ZERO_INT_VEC, 10, 10)
	splitted := parent.Split()

	expected := [4]BoundingBox[int]{
		NewBoundingBoxAt(ZERO_INT_VEC, 5, 5),
		NewBoundingBoxAt(NewVec(5, 0), 5, 5),
		NewBoundingBoxAt(NewVec(0, 5), 5, 5),
		NewBoundingBoxAt(NewVec(5, 5), 5, 5),
	}
	for i := 0; i < 4; i++ {
		if !splitted[i].Equals(expected[i]) {
			t.Errorf("split %v not equal to expected %v", splitted[i], expected[i])
		}
	}
}

func TestBox_BoxAround(t *testing.T) {
	center := NewVec(5, 5)
	box := NewBoundingBoxAround(center, 2)
	expectedTopLeft := NewVec(3, 3)
	expectedBottomRight := NewVec(7, 7)
	if box.TopLeft != expectedTopLeft {
		t.Errorf("topLeft %v not equal to expected %v", box.TopLeft, expectedTopLeft)
	}
	if box.BottomRight != expectedBottomRight {
		t.Errorf("bottomRight %v not equal to expected %v", box.BottomRight, expectedBottomRight)
	}
}

func TestBox_Intersects(t *testing.T) {
	intersects := []struct{ box1, box2 BoundingBox[int] }{
		{
			box1: NewBoundingBoxAt(Vec[int]{X: 3, Y: 3}, 2, 2),
			box2: NewBoundingBoxAt(Vec[int]{X: 2, Y: 2}, 2, 2),
		},
		{
			box1: NewBoundingBoxAt(Vec[int]{X: 3, Y: 3}, 2, 2),
			box2: NewBoundingBoxAt(Vec[int]{X: 4, Y: 2}, 2, 2),
		},
		{
			box1: NewBoundingBoxAt(Vec[int]{X: 3, Y: 3}, 2, 2),
			box2: NewBoundingBoxAt(Vec[int]{X: 2, Y: 4}, 2, 2),
		},
		{
			box1: NewBoundingBoxAt(Vec[int]{X: 3, Y: 3}, 2, 2),
			box2: NewBoundingBoxAt(Vec[int]{X: 4, Y: 4}, 2, 2),
		},
		{
			box1: NewBoundingBoxAt(Vec[int]{X: 3, Y: 3}, 2, 2),
			box2: NewBoundingBoxAt(Vec[int]{X: 2, Y: 2}, 2, 2),
		},
		{
			box1: NewBoundingBoxAt(Vec[int]{X: 3, Y: 3}, 2, 2),
			box2: NewBoundingBoxAt(Vec[int]{X: 2, Y: 2}, 2, 2),
		},
		{
			box1: NewBoundingBoxAt(Vec[int]{X: 3, Y: 3}, 2, 2),
			box2: NewBoundingBoxAt(Vec[int]{X: 2, Y: 4}, 2, 2),
		},
		{
			box1: NewBoundingBoxAt(Vec[int]{X: 3, Y: 3}, 2, 2),
			box2: NewBoundingBoxAt(Vec[int]{X: 2, Y: 2}, 2, 2),
		},
		{
			box1: NewBoundingBoxAt(Vec[int]{X: 3, Y: 3}, 2, 2),
			box2: NewBoundingBoxAt(Vec[int]{X: 2, Y: 5}, 2, 2),
		},
		{
			box1: NewBoundingBoxAt(Vec[int]{X: 3, Y: 3}, 2, 2),
			box2: NewBoundingBoxAt(Vec[int]{X: 5, Y: 2}, 2, 2),
		},
		{
			box1: NewBoundingBoxAt(Vec[int]{X: 3, Y: 3}, 2, 2),
			box2: NewBoundingBoxAt(Vec[int]{X: 5, Y: 2}, 2, 2),
		},
	}

	for _, intersection := range intersects {
		if !intersection.box1.Intersects(intersection.box2) {
			t.Errorf("box1 %v should intersect with box2 %v", intersection.box1, intersection.box2)
		}
		if !intersection.box2.Intersects(intersection.box1) {
			t.Errorf("box2 %v should intersect with box1 %v", intersection.box2, intersection.box1)
		}
	}

	notIntersects := []struct{ box1, box2 BoundingBox[int] }{
		{
			box1: NewBoundingBoxAt(Vec[int]{X: 3, Y: 3}, 2, 2),
			box2: NewBoundingBoxAt(Vec[int]{X: 1, Y: 1}, 1, 1),
		},
		{
			box1: NewBoundingBoxAt(Vec[int]{X: 3, Y: 3}, 2, 2),
			box2: NewBoundingBoxAt(Vec[int]{X: 6, Y: 0}, 3, 9),
		},
		{
			box1: NewBoundingBoxAt(Vec[int]{X: 3, Y: 3}, 2, 2),
			box2: NewBoundingBoxAt(Vec[int]{X: 0, Y: 6}, 9, 3),
		},
	}
	for _, intersection := range notIntersects {
		if intersection.box1.Intersects(intersection.box2) {
			t.Errorf("box1 %v should not intersect with box2 %v", intersection.box1, intersection.box2)
		}
		if intersection.box2.Intersects(intersection.box1) {
			t.Errorf("box2 %v should not intersect with box1 %v", intersection.box2, intersection.box1)
		}
	}
}

package geom

import "testing"

func TestBox_Split(t *testing.T) {
	parent := NewAABBAt(NewVec(0, 0), 10, 10)
	splitted := parent.Split()

	expected := [4]AABB[int]{
		NewAABBAt(NewVec(0, 0), 5, 5),
		NewAABBAt(NewVec(5, 0), 5, 5),
		NewAABBAt(NewVec(0, 5), 5, 5),
		NewAABBAt(NewVec(5, 5), 5, 5),
	}
	for i := 0; i < 4; i++ {
		if !splitted[i].Equals(expected[i]) {
			t.Errorf("split %v not equal to expected %v", splitted[i], expected[i])
		}
	}
}

func TestBox_BoxAround(t *testing.T) {
	center := NewVec(5, 5)
	box := NewAABBAround(center, 2)
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
	intersects := []struct{ box1, box2 AABB[int] }{
		{
			box1: NewAABBAt(Vec[int]{X: 3, Y: 3}, 2, 2),
			box2: NewAABBAt(Vec[int]{X: 2, Y: 2}, 2, 2),
		},
		{
			box1: NewAABBAt(Vec[int]{X: 3, Y: 3}, 2, 2),
			box2: NewAABBAt(Vec[int]{X: 4, Y: 2}, 2, 2),
		},
		{
			box1: NewAABBAt(Vec[int]{X: 3, Y: 3}, 2, 2),
			box2: NewAABBAt(Vec[int]{X: 2, Y: 4}, 2, 2),
		},
		{
			box1: NewAABBAt(Vec[int]{X: 3, Y: 3}, 2, 2),
			box2: NewAABBAt(Vec[int]{X: 4, Y: 4}, 2, 2),
		},
		{
			box1: NewAABBAt(Vec[int]{X: 3, Y: 3}, 2, 2),
			box2: NewAABBAt(Vec[int]{X: 2, Y: 2}, 2, 2),
		},
		{
			box1: NewAABBAt(Vec[int]{X: 3, Y: 3}, 2, 2),
			box2: NewAABBAt(Vec[int]{X: 2, Y: 2}, 2, 2),
		},
		{
			box1: NewAABBAt(Vec[int]{X: 3, Y: 3}, 2, 2),
			box2: NewAABBAt(Vec[int]{X: 2, Y: 4}, 2, 2),
		},
		{
			box1: NewAABBAt(Vec[int]{X: 3, Y: 3}, 2, 2),
			box2: NewAABBAt(Vec[int]{X: 2, Y: 2}, 2, 2),
		},
		{
			box1: NewAABBAt(Vec[int]{X: 3, Y: 3}, 2, 2),
			box2: NewAABBAt(Vec[int]{X: 2, Y: 5}, 2, 2),
		},
		{
			box1: NewAABBAt(Vec[int]{X: 3, Y: 3}, 2, 2),
			box2: NewAABBAt(Vec[int]{X: 5, Y: 2}, 2, 2),
		},
		{
			box1: NewAABBAt(Vec[int]{X: 3, Y: 3}, 2, 2),
			box2: NewAABBAt(Vec[int]{X: 5, Y: 2}, 2, 2),
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

	notIntersects := []struct{ box1, box2 AABB[int] }{
		{
			box1: NewAABBAt(Vec[int]{X: 3, Y: 3}, 2, 2),
			box2: NewAABBAt(Vec[int]{X: 1, Y: 1}, 1, 1),
		},
		{
			box1: NewAABBAt(Vec[int]{X: 3, Y: 3}, 2, 2),
			box2: NewAABBAt(Vec[int]{X: 6, Y: 0}, 3, 9),
		},
		{
			box1: NewAABBAt(Vec[int]{X: 3, Y: 3}, 2, 2),
			box2: NewAABBAt(Vec[int]{X: 0, Y: 6}, 9, 3),
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

func TestSortBoxesBy(t *testing.T) {
	mk := func(x1, y1, x2, y2 int) AABB[int] {
		return NewAABB(NewVec(x1, y1), NewVec(x2, y2))
	}

	byLeftX := func(box AABB[int]) int { return box.TopLeft.X }
	byLeftY := func(box AABB[int]) int { return box.TopLeft.Y }
	byWidth := func(box AABB[int]) int { return box.BottomRight.X - box.TopLeft.X }

	testCases := []struct {
		name       string
		a, b       AABB[int]
		keyFns     []func(AABB[int]) int
		wantFirst  AABB[int]
		wantSecond AABB[int]
	}{
		{
			name:       "ordersByFirstKey",
			a:          mk(0, 0, 2, 2),
			b:          mk(5, 0, 7, 2),
			keyFns:     []func(AABB[int]) int{byLeftX},
			wantFirst:  mk(0, 0, 2, 2),
			wantSecond: mk(5, 0, 7, 2),
		},
		{
			name:       "reversesWhenFirstGreater",
			a:          mk(10, 0, 12, 2),
			b:          mk(3, 0, 5, 2),
			keyFns:     []func(AABB[int]) int{byLeftX},
			wantFirst:  mk(3, 0, 5, 2),
			wantSecond: mk(10, 0, 12, 2),
		},
		{
			name: "fallsBackToNextKey",
			a:    mk(1, 5, 3, 7),
			b:    mk(1, 2, 3, 4),
			keyFns: []func(AABB[int]) int{
				byLeftX,
				byLeftY,
			},
			wantFirst:  mk(1, 2, 3, 4),
			wantSecond: mk(1, 5, 3, 7),
		},
		{
			name: "keepsOriginalWhenAllKeysEqual",
			a:    mk(0, 0, 2, 2),
			b:    mk(0, 0, 2, 2),
			keyFns: []func(AABB[int]) int{
				byLeftX,
				byWidth,
			},
			wantFirst:  mk(0, 0, 2, 2),
			wantSecond: mk(0, 0, 2, 2),
		},
		{
			name:       "noKeysKeepsOriginal",
			a:          mk(2, 2, 4, 4),
			b:          mk(1, 1, 3, 3),
			wantFirst:  mk(2, 2, 4, 4),
			wantSecond: mk(1, 1, 3, 3),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			first, second := SortBoxesBy(tc.a, tc.b, tc.keyFns...)
			if first != tc.wantFirst {
				t.Fatalf("first box mismatch: got %v want %v", first, tc.wantFirst)
			}
			if second != tc.wantSecond {
				t.Fatalf("second box mismatch: got %v want %v", second, tc.wantSecond)
			}
		})
	}
}

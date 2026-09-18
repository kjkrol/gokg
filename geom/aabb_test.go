package geom

import "testing"

func TestAABB_Split(t *testing.T) {
	parent := NewAABBAt(NewVec(0, 0), 10, 10)
	splitted := parent.Split()

	expected := [4]AABB{
		NewAABBAt(NewVec(0, 0), 5, 5),
		NewAABBAt(NewVec(5, 0), 5, 5),
		NewAABBAt(NewVec(0, 5), 5, 5),
		NewAABBAt(NewVec(5, 5), 5, 5),
	}

	for i := range expected {
		if !splitted[i].Equals(expected[i]) {
			t.Errorf("split %v not equal to expected %v", splitted[i], expected[i])
		}
	}
}

func TestAABB_NewAABBAround(t *testing.T) {
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

func TestAABB_Intersects(t *testing.T) {
	mk := func(x, y, w, h float64) AABB {
		return NewAABBAt(NewVec(x, y), w, h)
	}

	base := mk(3, 3, 2, 2)
	testCases := []struct {
		name       string
		box1, box2 AABB
		want       bool
	}{
		{name: "identical", box1: base, box2: base, want: true},
		{name: "overlapLeft", box1: base, box2: mk(2, 3, 2, 2), want: true},
		{name: "overlapRight", box1: base, box2: mk(4, 3, 2, 2), want: true},
		{name: "overlapTop", box1: base, box2: mk(3, 2, 2, 2), want: true},
		{name: "overlapBottom", box1: base, box2: mk(3, 4, 2, 2), want: true},
		{name: "containsOther", box1: base, box2: mk(3, 3, 1, 1), want: true},
		{name: "containedInOther", box1: base, box2: mk(2, 2, 4, 4), want: true},
		{name: "touchLeftEdge", box1: base, box2: mk(1, 3, 2, 2), want: true},
		{name: "touchRightEdge", box1: base, box2: mk(5, 3, 2, 2), want: true},
		{name: "touchTopEdge", box1: base, box2: mk(3, 1, 2, 2), want: true},
		{name: "touchBottomEdge", box1: base, box2: mk(3, 5, 2, 2), want: true},
		{name: "touchTopLeftCorner", box1: base, box2: mk(1, 1, 2, 2), want: true},
		{name: "touchTopRightCorner", box1: base, box2: mk(5, 1, 2, 2), want: true},
		{name: "touchBottomLeftCorner", box1: base, box2: mk(1, 5, 2, 2), want: true},
		{name: "touchBottomRightCorner", box1: base, box2: mk(5, 5, 2, 2), want: true},
		{name: "separateLeft", box1: base, box2: mk(0, 3, 2, 2), want: false},
		{name: "separateRight", box1: base, box2: mk(6, 3, 2, 2), want: false},
		{name: "separateAbove", box1: base, box2: mk(3, 0, 2, 2), want: false},
		{name: "separateBelow", box1: base, box2: mk(3, 6, 2, 2), want: false},
		{name: "separateDiagonal", box1: base, box2: mk(0, 0, 2, 2), want: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.box1.Intersects(tc.box2); got != tc.want {
				t.Fatalf("box1 %v intersects box2 %v = %v, want %v", tc.box1, tc.box2, got, tc.want)
			}
			if got := tc.box2.Intersects(tc.box1); got != tc.want {
				t.Fatalf("box2 %v intersects box1 %v = %v, want %v", tc.box2, tc.box1, got, tc.want)
			}
		})
	}
}

func TestAABB_ContainsVec(t *testing.T) {
	box := NewAABBAt(NewVec(3, 3), 2, 2) // [3,3] .. [5,5]
	testCases := []struct {
		name string
		vec  Vec
		want bool
	}{
		{name: "inside", vec: NewVec(4, 4), want: true},
		{name: "topLeftCorner", vec: NewVec(3, 3), want: true},
		{name: "bottomRightCorner", vec: NewVec(5, 5), want: true},
		{name: "onLeftEdge", vec: NewVec(3, 4), want: true},
		{name: "onTopEdge", vec: NewVec(4, 3), want: true},
		{name: "outsideLeft", vec: NewVec(2, 4), want: false},
		{name: "outsideBelow", vec: NewVec(4, 6), want: false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := box.ContainsVec(tc.vec); got != tc.want {
				t.Fatalf("box %v ContainsVec(%v) = %v, want %v", box, tc.vec, got, tc.want)
			}
		})
	}
}

func TestSortAABBsBy(t *testing.T) {
	mk := func(x1, y1, x2, y2 float64) AABB {
		return NewAABB(NewVec(x1, y1), NewVec(x2, y2))
	}

	byLeftX := func(box AABB) float64 { return box.TopLeft.X }
	byLeftY := func(box AABB) float64 { return box.TopLeft.Y }
	byWidth := func(box AABB) float64 { return box.BottomRight.X - box.TopLeft.X }

	testCases := []struct {
		name       string
		a, b       AABB
		keyFns     []func(AABB) float64
		wantFirst  AABB
		wantSecond AABB
	}{
		{
			name:       "ordersByFirstKey",
			a:          mk(0, 0, 2, 2),
			b:          mk(5, 0, 7, 2),
			keyFns:     []func(AABB) float64{byLeftX},
			wantFirst:  mk(0, 0, 2, 2),
			wantSecond: mk(5, 0, 7, 2),
		},
		{
			name:       "reversesWhenFirstGreater",
			a:          mk(10, 0, 12, 2),
			b:          mk(3, 0, 5, 2),
			keyFns:     []func(AABB) float64{byLeftX},
			wantFirst:  mk(3, 0, 5, 2),
			wantSecond: mk(10, 0, 12, 2),
		},
		{
			name: "fallsBackToNextKey",
			a:    mk(1, 5, 3, 7),
			b:    mk(1, 2, 3, 4),
			keyFns: []func(AABB) float64{
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
			keyFns: []func(AABB) float64{
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
			first, second := SortAABBsBy(tc.a, tc.b, tc.keyFns...)
			if first != tc.wantFirst {
				t.Fatalf("first box mismatch: got %v want %v", first, tc.wantFirst)
			}
			if second != tc.wantSecond {
				t.Fatalf("second box mismatch: got %v want %v", second, tc.wantSecond)
			}
		})
	}
}

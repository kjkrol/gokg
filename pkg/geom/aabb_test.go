package geom

import "testing"

func TestAABB_Split(t *testing.T) {
	runAABBSplitTest[int](t, "int")
	runAABBSplitTest[uint32](t, "uint32")
	runAABBSplitTest[float64](t, "float64")
}

func runAABBSplitTest[T Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		parent := NewAABBAt(NewVec(T(0), T(0)), T(10), T(10))
		splitted := parent.Split()

		expected := [4]AABB[T]{
			NewAABBAt(NewVec(T(0), T(0)), T(5), T(5)),
			NewAABBAt(NewVec(T(5), T(0)), T(5), T(5)),
			NewAABBAt(NewVec(T(0), T(5)), T(5), T(5)),
			NewAABBAt(NewVec(T(5), T(5)), T(5), T(5)),
		}

		for i := range expected {
			if !splitted[i].Equals(expected[i]) {
				t.Errorf("split %v not equal to expected %v", splitted[i], expected[i])
			}
		}
	})
}

func TestAABB_NewAABBAround(t *testing.T) {
	runAABBAroundTest[int](t, "int")
	runAABBAroundTest[uint32](t, "uint32")
	runAABBAroundTest[float64](t, "float64")
}

func runAABBAroundTest[T Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		center := NewVec(T(5), T(5))
		box := NewAABBAround(center, T(2))
		expectedTopLeft := NewVec(T(3), T(3))
		expectedBottomRight := NewVec(T(7), T(7))

		if box.TopLeft != expectedTopLeft {
			t.Errorf("topLeft %v not equal to expected %v", box.TopLeft, expectedTopLeft)
		}
		if box.BottomRight != expectedBottomRight {
			t.Errorf("bottomRight %v not equal to expected %v", box.BottomRight, expectedBottomRight)
		}
	})
}

func TestAABB_Intersects(t *testing.T) {
	runAABBIntersectsTest[int](t, "int")
	runAABBIntersectsTest[uint32](t, "uint32")
	runAABBIntersectsTest[float64](t, "float64")
}

func runAABBIntersectsTest[T Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		mk := func(x, y, w, h T) AABB[T] {
			return NewAABBAt(NewVec(x, y), w, h)
		}

		base := mk(T(3), T(3), T(2), T(2))
		testCases := []struct {
			name       string
			box1, box2 AABB[T]
			want       bool
		}{
			{name: "identical", box1: base, box2: base, want: true},
			{name: "overlapLeft", box1: base, box2: mk(T(2), T(3), T(2), T(2)), want: true},
			{name: "overlapRight", box1: base, box2: mk(T(4), T(3), T(2), T(2)), want: true},
			{name: "overlapTop", box1: base, box2: mk(T(3), T(2), T(2), T(2)), want: true},
			{name: "overlapBottom", box1: base, box2: mk(T(3), T(4), T(2), T(2)), want: true},
			{name: "containsOther", box1: base, box2: mk(T(3), T(3), T(1), T(1)), want: true},
			{name: "containedInOther", box1: base, box2: mk(T(2), T(2), T(4), T(4)), want: true},
			{name: "touchLeftEdge", box1: base, box2: mk(T(1), T(3), T(2), T(2)), want: true},
			{name: "touchRightEdge", box1: base, box2: mk(T(5), T(3), T(2), T(2)), want: true},
			{name: "touchTopEdge", box1: base, box2: mk(T(3), T(1), T(2), T(2)), want: true},
			{name: "touchBottomEdge", box1: base, box2: mk(T(3), T(5), T(2), T(2)), want: true},
			{name: "touchTopLeftCorner", box1: base, box2: mk(T(1), T(1), T(2), T(2)), want: true},
			{name: "touchTopRightCorner", box1: base, box2: mk(T(5), T(1), T(2), T(2)), want: true},
			{name: "touchBottomLeftCorner", box1: base, box2: mk(T(1), T(5), T(2), T(2)), want: true},
			{name: "touchBottomRightCorner", box1: base, box2: mk(T(5), T(5), T(2), T(2)), want: true},
			{name: "separateLeft", box1: base, box2: mk(T(0), T(3), T(2), T(2)), want: false},
			{name: "separateRight", box1: base, box2: mk(T(6), T(3), T(2), T(2)), want: false},
			{name: "separateAbove", box1: base, box2: mk(T(3), T(0), T(2), T(2)), want: false},
			{name: "separateBelow", box1: base, box2: mk(T(3), T(6), T(2), T(2)), want: false},
			{name: "separateDiagonal", box1: base, box2: mk(T(0), T(0), T(2), T(2)), want: false},
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
	})
}

func TestSortAABBsBy(t *testing.T) {
	runSortAABBsByTest[int](t, "int")
	runSortAABBsByTest[uint32](t, "uint32")
	runSortAABBsByTest[float64](t, "float64")
}

func runSortAABBsByTest[T Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		mk := func(x1, y1, x2, y2 T) AABB[T] {
			return NewAABB(NewVec(x1, y1), NewVec(x2, y2))
		}

		byLeftX := func(box AABB[T]) T { return box.TopLeft.X }
		byLeftY := func(box AABB[T]) T { return box.TopLeft.Y }
		byWidth := func(box AABB[T]) T { return box.BottomRight.X - box.TopLeft.X }

		testCases := []struct {
			name       string
			a, b       AABB[T]
			keyFns     []func(AABB[T]) T
			wantFirst  AABB[T]
			wantSecond AABB[T]
		}{
			{
				name:       "ordersByFirstKey",
				a:          mk(T(0), T(0), T(2), T(2)),
				b:          mk(T(5), T(0), T(7), T(2)),
				keyFns:     []func(AABB[T]) T{byLeftX},
				wantFirst:  mk(T(0), T(0), T(2), T(2)),
				wantSecond: mk(T(5), T(0), T(7), T(2)),
			},
			{
				name:       "reversesWhenFirstGreater",
				a:          mk(T(10), T(0), T(12), T(2)),
				b:          mk(T(3), T(0), T(5), T(2)),
				keyFns:     []func(AABB[T]) T{byLeftX},
				wantFirst:  mk(T(3), T(0), T(5), T(2)),
				wantSecond: mk(T(10), T(0), T(12), T(2)),
			},
			{
				name: "fallsBackToNextKey",
				a:    mk(T(1), T(5), T(3), T(7)),
				b:    mk(T(1), T(2), T(3), T(4)),
				keyFns: []func(AABB[T]) T{
					byLeftX,
					byLeftY,
				},
				wantFirst:  mk(T(1), T(2), T(3), T(4)),
				wantSecond: mk(T(1), T(5), T(3), T(7)),
			},
			{
				name: "keepsOriginalWhenAllKeysEqual",
				a:    mk(T(0), T(0), T(2), T(2)),
				b:    mk(T(0), T(0), T(2), T(2)),
				keyFns: []func(AABB[T]) T{
					byLeftX,
					byWidth,
				},
				wantFirst:  mk(T(0), T(0), T(2), T(2)),
				wantSecond: mk(T(0), T(0), T(2), T(2)),
			},
			{
				name:       "noKeysKeepsOriginal",
				a:          mk(T(2), T(2), T(4), T(4)),
				b:          mk(T(1), T(1), T(3), T(3)),
				wantFirst:  mk(T(2), T(2), T(4), T(4)),
				wantSecond: mk(T(1), T(1), T(3), T(3)),
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
	})
}

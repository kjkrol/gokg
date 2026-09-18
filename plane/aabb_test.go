package plane

import (
	"testing"

	"github.com/kjkrol/gokg/geom"
)

func TestAABB_NewAABB(t *testing.T) {
	runAABBNewTest[int](t, "int")
	runAABBNewTest[uint32](t, "uint32")
	runAABBNewTest[float64](t, "float64")
}

func runAABBNewTest[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		aabb := NewAABB[T](vec[T](0, 0), T(10), T(10))
		expected := vec[T](10, 10)
		if aabb.BottomRight != expected {
			t.Errorf("center %v not equal to expected %v", aabb.BottomRight, expected)
		}
	})
}

func TestAABB_IntersectsIncludingFrags(t *testing.T) {
	runAABBIntersectsIncludingFragsTest[int](t, "int")
	runAABBIntersectsIncludingFragsTest[uint32](t, "uint32")
	runAABBIntersectsIncludingFragsTest[float64](t, "float64")
}

func runAABBIntersectsIncludingFragsTest[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		testCases := []struct {
			name  string
			aabb1 AABB[T]
			aabb2 AABB[T]
			// overhang is how far aabb2 runs past the far edges; the fragments
			// noted beside each case follow from it.
			overhang geom.Vec[T]
			want     bool
		}{
			{
				name:  "returnsTrueWhenAnyFragmentsIntersect",
				aabb1: NewAABB(geom.NewVec(T(0), T(0)), T(2), T(2)),
				aabb2: NewAABB(geom.NewVec(T(4), T(4)), T(1), T(1)),
				// right (0,4)..(1,5), bottom (4,0)..(5,1), corner (0,0)..(1,1)
				overhang: geom.NewVec(T(1), T(1)),
				want:     true,
			},
			{
				name:  "returnsFalseWhenNoFragmentsIntersect",
				aabb1: NewAABB(geom.NewVec(T(0), T(0)), T(2), T(2)),
				aabb2: NewAABB(geom.NewVec(T(4), T(4)), T(2), T(2)),
				// right (0,4)..(1,6) only — it does not reach the bottom edge
				overhang: geom.NewVec(T(1), T(0)),
				want:     false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tc.aabb2.Overhang = tc.overhang

				if got := tc.aabb1.IntersectsWithFrags(tc.aabb2); got != tc.want {
					t.Errorf("expected Intersects to return %v, but got %v", tc.want, got)
				}
			})
		}
	})
}

func TestAABB_Contains(t *testing.T) {
	runAABBContainsTest[int](t, "int")
	runAABBContainsTest[uint32](t, "uint32")
	runAABBContainsTest[float64](t, "float64")
}

func runAABBContainsTest[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		outer := NewAABB(geom.NewVec(T(0), T(0)), T(10), T(10))
		testCases := []struct {
			name   string
			target AABB[T]
			want   bool
		}{
			{name: "containsInner", target: NewAABB(geom.NewVec(T(2), T(2)), T(6), T(6)), want: true},
			{name: "rejectsBoxStartingOutside", target: NewAABB(geom.NewVec(T(11), T(11)), T(1), T(1)), want: false},
			{name: "rejectsBoxExtendingBeyondBounds", target: NewAABB(geom.NewVec(T(5), T(5)), T(7), T(7)), want: false},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				if got := outer.ContainsWithFrags(tc.target); got != tc.want {
					t.Errorf("expected Contains to return %v, got %v", tc.want, got)
				}
			})
		}
	})
}

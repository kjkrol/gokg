package plane

import (
	"testing"

	"github.com/kjkrol/gokg/pkg/geom"
)

func TestAABB_NewAABB(t *testing.T) {
	runAABBNewTest[int](t, "int")
	runAABBNewTest[uint32](t, "uint32")
	runAABBNewTest[float64](t, "float64")
}

func runAABBNewTest[T geom.Numeric](t *testing.T, name string) {
	t.Run(name, func(t *testing.T) {
		aabb := newAABB[T](vec[T](0, 0), T(10), T(10))
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
			frags map[FragPosition]geom.AABB[T]
			want  bool
		}{
			{
				name:  "returnsTrueWhenAnyFragmentsIntersect",
				aabb1: newAABB(geom.NewVec(T(0), T(0)), T(2), T(2)),
				aabb2: newAABB(geom.NewVec(T(4), T(4)), T(1), T(1)),
				frags: map[FragPosition]geom.AABB[T]{
					FRAG_RIGHT:        geom.NewAABB(geom.NewVec(T(0), T(4)), geom.NewVec(T(1), T(5))),
					FRAG_BOTTOM:       geom.NewAABB(geom.NewVec(T(4), T(0)), geom.NewVec(T(5), T(1))),
					FRAG_BOTTOM_RIGHT: geom.NewAABB(geom.NewVec(T(0), T(0)), geom.NewVec(T(1), T(1))),
				},
				want: true,
			},
			{
				name:  "returnsFalseWhenNoFragmentsIntersect",
				aabb1: newAABB(geom.NewVec(T(0), T(0)), T(2), T(2)),
				aabb2: newAABB(geom.NewVec(T(4), T(4)), T(2), T(2)),
				frags: map[FragPosition]geom.AABB[T]{
					FRAG_RIGHT: geom.NewAABB(geom.NewVec(T(0), T(4)), geom.NewVec(T(1), T(6))),
				},
				want: false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				for pos, coords := range tc.frags {
					tc.aabb2.setFragment(pos, coords)
				}

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
		outer := newAABB(geom.NewVec(T(0), T(0)), T(10), T(10))
		testCases := []struct {
			name   string
			target AABB[T]
			want   bool
		}{
			{name: "containsInner", target: newAABB(geom.NewVec(T(2), T(2)), T(6), T(6)), want: true},
			{name: "rejectsBoxStartingOutside", target: newAABB(geom.NewVec(T(11), T(11)), T(1), T(1)), want: false},
			{name: "rejectsBoxExtendingBeyondBounds", target: newAABB(geom.NewVec(T(5), T(5)), T(7), T(7)), want: false},
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

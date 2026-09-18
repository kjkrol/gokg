package plane

import (
	"testing"

	"github.com/kjkrol/gokg/geom"
)

func TestAABB_NewAABB(t *testing.T) {
	aabb := NewAABB(vec(0, 0), 10, 10)
	expected := vec(10, 10)
	if aabb.BottomRight != expected {
		t.Errorf("center %v not equal to expected %v", aabb.BottomRight, expected)
	}
}

func TestAABB_IntersectsIncludingFrags(t *testing.T) {
	testCases := []struct {
		name  string
		aabb1 AABB
		aabb2 AABB
		// overhang is how far aabb2 runs past the far edges; the fragments
		// noted beside each case follow from it.
		overhang geom.Vec
		want     bool
	}{
		{
			name:  "returnsTrueWhenAnyFragmentsIntersect",
			aabb1: NewAABB(geom.NewVec(0, 0), 2, 2),
			aabb2: NewAABB(geom.NewVec(4, 4), 1, 1),
			// right (0,4)..(1,5), bottom (4,0)..(5,1), corner (0,0)..(1,1)
			overhang: geom.NewVec(1, 1),
			want:     true,
		},
		{
			name:  "returnsFalseWhenNoFragmentsIntersect",
			aabb1: NewAABB(geom.NewVec(0, 0), 2, 2),
			aabb2: NewAABB(geom.NewVec(4, 4), 2, 2),
			// right (0,4)..(1,6) only — it does not reach the bottom edge
			overhang: geom.NewVec(1, 0),
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
}

func TestAABB_Contains(t *testing.T) {
	outer := NewAABB(geom.NewVec(0, 0), 10, 10)
	testCases := []struct {
		name   string
		target AABB
		want   bool
	}{
		{name: "containsInner", target: NewAABB(geom.NewVec(2, 2), 6, 6), want: true},
		{name: "rejectsBoxStartingOutside", target: NewAABB(geom.NewVec(11, 11), 1, 1), want: false},
		{name: "rejectsBoxExtendingBeyondBounds", target: NewAABB(geom.NewVec(5, 5), 7, 7), want: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := outer.ContainsWithFrags(tc.target); got != tc.want {
				t.Errorf("expected Contains to return %v, got %v", tc.want, got)
			}
		})
	}
}

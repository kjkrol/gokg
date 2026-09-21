package plane

import (
	"testing"

	"github.com/kjkrol/aabbworld/geom"
)

func TestHasFragments_FollowsTheOverhang(t *testing.T) {
	for name, tc := range map[string]struct {
		overhang geom.Vec
		want     bool
	}{
		"clear of every edge": {geom.Vec{}, false},
		"off the right edge":  {geom.NewVec(5, 0), true},
		"off the bottom edge": {geom.NewVec(0, 5), true},
		"off the corner":      {geom.NewVec(5, 5), true},
	} {
		box := NewAABB(geom.NewVec(0, 0), 10, 10)
		box.Overhang = tc.overhang
		if got := box.hasFragments(); got != tc.want {
			t.Errorf("%s: hasFragments = %v, want %v", name, got, tc.want)
		}
	}
}

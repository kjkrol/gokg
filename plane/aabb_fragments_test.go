package plane

import (
	"testing"
	"unsafe"

	"github.com/kjkrol/gokg/geom"
)

func box(x0, y0, x1, y1 uint32) geom.AABB[uint32] {
	return geom.NewAABB(geom.NewVec(x0, y0), geom.NewVec(x1, y1))
}

// The fragments are rebuilt from the overhang rather than stored, so they have
// to come out exactly where they used to. The expectations below are written
// out by hand, not derived from the code under test.
func TestVisitFragments_RebuildsEveryPieceWhereItBelongs(t *testing.T) {
	space := NewToroidal2D[uint32](100, 100)

	cases := map[string]struct {
		x, y, w, h uint32
		want       map[FragPosition]geom.AABB[uint32]
	}{
		"clear of every edge": {
			10, 10, 20, 20,
			nil,
		},
		"past the right edge": {
			95, 10, 20, 20,
			map[FragPosition]geom.AABB[uint32]{FRAG_RIGHT: box(0, 10, 15, 30)},
		},
		"past the bottom edge": {
			10, 95, 20, 20,
			map[FragPosition]geom.AABB[uint32]{FRAG_BOTTOM: box(10, 0, 30, 15)},
		},
		"through the corner": {
			95, 95, 20, 20,
			map[FragPosition]geom.AABB[uint32]{
				FRAG_RIGHT:        box(0, 95, 15, 100),
				FRAG_BOTTOM:       box(95, 0, 100, 15),
				FRAG_BOTTOM_RIGHT: box(0, 0, 15, 15),
			},
		},
		"wider than the world itself": {
			10, 10, 150, 150,
			map[FragPosition]geom.AABB[uint32]{
				FRAG_RIGHT:        box(0, 10, 60, 100),
				FRAG_BOTTOM:       box(10, 0, 100, 60),
				FRAG_BOTTOM_RIGHT: box(0, 0, 60, 60),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			wrapped := space.WrapAABB(geom.NewAABBAt(geom.NewVec(tc.x, tc.y), tc.w, tc.h))

			got := map[FragPosition]geom.AABB[uint32]{}
			wrapped.VisitFragments(func(pos FragPosition, frag geom.AABB[uint32]) bool {
				got[pos] = frag
				return true
			})

			if len(got) != len(tc.want) {
				t.Fatalf("got %d fragments %v, want %d %v", len(got), got, len(tc.want), tc.want)
			}
			for pos, want := range tc.want {
				if got[pos] != want {
					t.Errorf("fragment %d = %v, want %v", pos, got[pos], want)
				}
			}
		})
	}
}

// Stopping early has to work on the rebuilt walk too — the sprite batch and the
// contact solver both lean on it.
func TestVisitFragments_StopsWhenAsked(t *testing.T) {
	space := NewToroidal2D[uint32](100, 100)
	wrapped := space.WrapAABB(geom.NewAABBAt(geom.NewVec[uint32](95, 95), 20, 20))

	seen := 0
	wrapped.VisitFragments(func(FragPosition, geom.AABB[uint32]) bool {
		seen++
		return false
	})
	if seen != 1 {
		t.Errorf("visited %d fragments after returning false, want 1", seen)
	}
}

// The overhang exists to keep AABB small; a cache creeping back would undo the
// point of it, and AABB is the hottest struct any Space-backed game carries.
func TestAABB_StaysSmall(t *testing.T) {
	if got := unsafe.Sizeof(AABB[uint32]{}); got != 32 {
		t.Errorf("AABB[uint32] is %d bytes, want 32 — corners, size and overhang and nothing else", got)
	}
	if got := unsafe.Sizeof(AABB[float64]{}); got != 64 {
		t.Errorf("AABB[float64] is %d bytes, want 64", got)
	}
}

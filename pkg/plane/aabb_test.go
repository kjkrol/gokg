package plane

import (
	"testing"

	"github.com/kjkrol/gokg/pkg/geom"
)

func TestAABB_NewAABB(t *testing.T) {
	aabb := newAABB(geom.NewVec(0, 0), 10, 10)
	expected := geom.NewVec(10, 10)
	if aabb.BottomRight != expected {
		t.Errorf("center %v not equal to expected %v", aabb.BottomRight, expected)
	}
}

func TestAABB_IntersectsIncludingFrags_ReturnsTrue(t *testing.T) {
	aabb1 := newAABB(geom.NewVec(0, 0), 2, 2)
	aabb2 := newAABB(geom.NewVec(4, 4), 1, 1)
	aabb2.frags[FRAG_RIGHT] = geom.NewAABB(geom.NewVec(0, 4), geom.NewVec(1, 5))
	aabb2.frags[FRAG_BOTTOM] = geom.NewAABB(geom.NewVec(4, 0), geom.NewVec(5, 1))
	aabb2.frags[FRAG_BOTTOM_RIGHT] = geom.NewAABB(geom.NewVec(0, 0), geom.NewVec(1, 1))

	if !aabb1.Intersects(aabb2) {
		t.Errorf("expected IntersectsAny to return true, but got false")
	}
}

func TestAABB_IntersectsIncludingFrags_ReturnsFalse(t *testing.T) {
	aabb1 := newAABB(geom.NewVec(0, 0), 2, 2)
	aabb2 := newAABB(geom.NewVec(4, 4), 2, 2)
	aabb2.frags[FRAG_RIGHT] = geom.NewAABB(geom.NewVec(0, 4), geom.NewVec(1, 6))

	if aabb1.Intersects(aabb2) {
		t.Errorf("expected IntersectsAny to return false, but got true")
	}
}

func TestAABB_Contains(t *testing.T) {
	outer := newAABB(geom.NewVec(0, 0), 10, 10)
	inner := newAABB(geom.NewVec(2, 2), 6, 6)
	onlyTopLeftInside := newAABB(geom.NewVec(-1, -1), 4, 4)
	onlyBottomRightInside := newAABB(geom.NewVec(5, 5), 7, 7)

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

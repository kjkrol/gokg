package plane

import (
	"testing"

	"github.com/kjkrol/gokg/pkg/geom"
)

func TestCartesian_Expand(t *testing.T) {
	cartesian := NewCartesian(10, 10)
	aabb := newAABB(geom.NewVec(2, 3), 3, 4)
	cartesian.Expand(&aabb, 2)
	expectAABBState(t, aabb, geom.NewVec(0, 1), geom.NewVec(7, 9), map[FragPosition][2]geom.Vec[int]{})
}

func TestCartesian_Expand_CornerCase(t *testing.T) {
	cartesian := NewCartesian(10, 10)
	aabb := newAABB(geom.NewVec(0, 0), 2, 2)
	cartesian.Expand(&aabb, 2)
	expectAABBState(t, aabb, geom.NewVec(0, 0), geom.NewVec(4, 4), map[FragPosition][2]geom.Vec[int]{})
}

func TestTorus_Expand_CornerCase(t *testing.T) {
	torus := NewTorus(10, 10)
	aabb := newAABB(geom.NewVec(0, 0), 2, 2)
	torus.Expand(&aabb, 2)
	expectAABBState(t, aabb, geom.NewVec(8, 8), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec[int]{
		FRAG_RIGHT:        {geom.NewVec(0, 8), geom.NewVec(4, 10)},
		FRAG_BOTTOM:       {geom.NewVec(8, 0), geom.NewVec(10, 4)},
		FRAG_BOTTOM_RIGHT: {geom.NewVec(0, 0), geom.NewVec(4, 4)},
	})
}

func TestTorus_Exapnad_ThenIntersects(t *testing.T) {
	torus := NewTorus(100, 100)

	aabb1 := newAABB(geom.NewVec(5, 5), 10, 10)
	aabb2 := newAABB(geom.NewVec(96, 96), 10, 10)

	torus.Expand(&aabb2, 0)

	if !aabb1.Intersects(aabb2) {
		t.Errorf("rect1 %v should intersect with rect2 %v", aabb1, aabb2)
	}

}

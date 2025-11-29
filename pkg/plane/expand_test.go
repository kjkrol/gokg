package plane

import (
	"testing"

	"github.com/kjkrol/gokg/pkg/geom"
)

func TestPlane_Expand_On_BoundedPlane(t *testing.T) {
	plane := NewCartesian(10, 10)
	planeBox := newAABB(geom.NewVec(2, 3), 3, 4)
	plane.Expand(&planeBox, 2)
	expectPlaneBoxState(t, planeBox, geom.NewVec(0, 1), geom.NewVec(7, 9), map[FragPosition][2]geom.Vec[int]{})
}

func TestPlane_Expand_On_BoundedPlane_CornerCase(t *testing.T) {
	plane := NewCartesian(10, 10)
	planeBox := newAABB(geom.NewVec(0, 0), 2, 2)
	plane.Expand(&planeBox, 2)
	expectPlaneBoxState(t, planeBox, geom.NewVec(0, 0), geom.NewVec(4, 4), map[FragPosition][2]geom.Vec[int]{})
}

func TestPlane_Expand_On_CyclicPlane_CornerCase(t *testing.T) {
	plane := NewTorus(10, 10)
	planeBox := newAABB(geom.NewVec(0, 0), 2, 2)
	plane.Expand(&planeBox, 2)
	expectPlaneBoxState(t, planeBox, geom.NewVec(8, 8), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec[int]{
		FRAG_RIGHT:        {geom.NewVec(0, 8), geom.NewVec(4, 10)},
		FRAG_BOTTOM:       {geom.NewVec(8, 0), geom.NewVec(10, 4)},
		FRAG_BOTTOM_RIGHT: {geom.NewVec(0, 0), geom.NewVec(4, 4)},
	})
}

func TestPlane_Exapnad_ThenIntersects(t *testing.T) {
	plane := NewTorus(100, 100)

	rect1 := newAABB(geom.NewVec(5, 5), 10, 10)
	rect2 := newAABB(geom.NewVec(96, 96), 10, 10)

	plane.Expand(&rect2, 0)

	if !rect1.Intersects(rect2) {
		t.Errorf("rect1 %v should intersect with rect2 %v", rect1, rect2)
	}

}

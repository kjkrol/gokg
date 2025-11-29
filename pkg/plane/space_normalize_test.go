package plane

import (
	"testing"

	"github.com/kjkrol/gokg/pkg/geom"
)

func TestPlaneBounded_normalizeBox(t *testing.T) {
	cartesian := NewCartesian(10, 10)

	aabb := newAABB(geom.NewVec(-2, -2), 2, 2)
	cartesian.(space2d[int]).normalizeBox(&aabb)
	expectPlaneBoxState(t, aabb, geom.NewVec(0, 0), geom.NewVec(0, 0), map[FragPosition][2]geom.Vec[int]{})

	aabb = newAABB(geom.NewVec(-1, -1), 2, 2)
	cartesian.(space2d[int]).normalizeBox(&aabb)
	expectPlaneBoxState(t, aabb, geom.NewVec(0, 0), geom.NewVec(1, 1), map[FragPosition][2]geom.Vec[int]{})

	aabb = newAABB(geom.NewVec(9, 9), 2, 2)
	cartesian.(space2d[int]).normalizeBox(&aabb)
	expectPlaneBoxState(t, aabb, geom.NewVec(9, 9), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec[int]{})

	aabb = newAABB(geom.NewVec(-11, -11), 2, 2)
	cartesian.(space2d[int]).normalizeBox(&aabb)
	expectPlaneBoxState(t, aabb, geom.NewVec(0, 0), geom.NewVec(0, 0), map[FragPosition][2]geom.Vec[int]{})

	aabb = newAABB(geom.NewVec(19, 19), 2, 2)
	cartesian.(space2d[int]).normalizeBox(&aabb)
	expectPlaneBoxState(t, aabb, geom.NewVec(10, 10), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec[int]{})
}

func TestPlaneCyclic_normalizeBox(t *testing.T) {
	torus := NewTorus(10, 10)

	aabb := newAABB(geom.NewVec(-2, -2), 2, 2)
	torus.(space2d[int]).normalizeBox(&aabb)
	expectPlaneBoxState(t, aabb, geom.NewVec(8, 8), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec[int]{})

	aabb = newAABB(geom.NewVec(-1, -1), 2, 2)
	torus.(space2d[int]).normalizeBox(&aabb)
	expectPlaneBoxState(t, aabb, geom.NewVec(9, 9), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec[int]{
		FRAG_RIGHT:        {geom.NewVec(0, 9), geom.NewVec(1, 10)},
		FRAG_BOTTOM:       {geom.NewVec(9, 0), geom.NewVec(10, 1)},
		FRAG_BOTTOM_RIGHT: {geom.NewVec(0, 0), geom.NewVec(1, 1)},
	})

	aabb = newAABB(geom.NewVec(9, 9), 2, 2)
	torus.(space2d[int]).normalizeBox(&aabb)
	expectPlaneBoxState(t, aabb, geom.NewVec(9, 9), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec[int]{
		FRAG_RIGHT:        {geom.NewVec(0, 9), geom.NewVec(1, 10)},
		FRAG_BOTTOM:       {geom.NewVec(9, 0), geom.NewVec(10, 1)},
		FRAG_BOTTOM_RIGHT: {geom.NewVec(0, 0), geom.NewVec(1, 1)},
	})

	aabb = newAABB(geom.NewVec(-11, -11), 2, 2)
	torus.(space2d[int]).normalizeBox(&aabb)
	expectPlaneBoxState(t, aabb, geom.NewVec(9, 9), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec[int]{
		FRAG_RIGHT:        {geom.NewVec(0, 9), geom.NewVec(1, 10)},
		FRAG_BOTTOM:       {geom.NewVec(9, 0), geom.NewVec(10, 1)},
		FRAG_BOTTOM_RIGHT: {geom.NewVec(0, 0), geom.NewVec(1, 1)},
	})

	aabb = newAABB(geom.NewVec(19, 19), 2, 2)
	torus.(space2d[int]).normalizeBox(&aabb)
	expectPlaneBoxState(t, aabb, geom.NewVec(9, 9), geom.NewVec(10, 10), map[FragPosition][2]geom.Vec[int]{
		FRAG_RIGHT:        {geom.NewVec(0, 9), geom.NewVec(1, 10)},
		FRAG_BOTTOM:       {geom.NewVec(9, 0), geom.NewVec(10, 1)},
		FRAG_BOTTOM_RIGHT: {geom.NewVec(0, 0), geom.NewVec(1, 1)},
	})
}

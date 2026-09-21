package plane

import (
	"github.com/kjkrol/aabbworld/plane"
	"testing"

	"github.com/kjkrol/aabbworld/geom"
)

func TestEuclidean2DExpand(t *testing.T) {
	euclidean := NewEuclidean2D(10, 10)
	aabb := plane.NewAABB(vec(2, 3), 3, 4)
	euclidean.Expand(&aabb, 2)
	expectAABBState(t, aabb, vec(0, 1), vec(7, 9), map[plane.FragPosition][2]geom.Vec{})
}

func TestEuclidean2DExpandCornerCase(t *testing.T) {
	euclidean := NewEuclidean2D(10, 10)
	aabb := plane.NewAABB(vec(0, 0), 2, 2)
	euclidean.Expand(&aabb, 2)
	expectAABBState(t, aabb, vec(0, 0), vec(4, 4), map[plane.FragPosition][2]geom.Vec{})
}

func TestToroidal2DExpandCornerCase(t *testing.T) {
	toroidal := NewToroidal2D(10, 10)
	aabb := plane.NewAABB(vec(0, 0), 2, 2)
	toroidal.Expand(&aabb, 2)
	expectAABBState(t, aabb, vec(8, 8), vec(10, 10), convertFragments(map[plane.FragPosition][2]geom.Vec{
		plane.FRAG_RIGHT:        {geom.NewVec(0, 8), geom.NewVec(4, 10)},
		plane.FRAG_BOTTOM:       {geom.NewVec(8, 0), geom.NewVec(10, 4)},
		plane.FRAG_BOTTOM_RIGHT: {geom.NewVec(0, 0), geom.NewVec(4, 4)},
	}))
}

func TestToroidal2DExpandThenIntersects(t *testing.T) {
	toroidal := NewToroidal2D(100, 100)

	aabb1 := plane.NewAABB(vec(5, 5), 10, 10)
	aabb2 := plane.NewAABB(vec(96, 96), 10, 10)

	toroidal.Expand(&aabb2, 0)

	if intersects := meetAcrossSeams(aabb1, aabb2); intersects != true {
		t.Errorf("unexpected intersection result. got %t, want %t for boxes %v and %v", intersects, true, aabb1, aabb2)
	}
}

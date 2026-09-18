package plane

import (
	"testing"

	"github.com/kjkrol/gokg/geom"
)

func TestAxisDistance(t *testing.T) {
	aa := NewAABB(geom.NewVec(0, 0), 2, 2)
	bb := NewAABB(geom.NewVec(5, 0), 2, 2)

	dx := aa.AxisDistanceX(bb.AABB)
	dy := aa.AxisDistanceY(bb.AABB)

	if dx != 3 {
		t.Errorf("expected dx=3, got %v", dx)
	}
	if dy != 0 {
		t.Errorf("expected dy=0, got %v", dy)
	}
}

func TestAABBDistance_Euclidean2DSpace(t *testing.T) {
	rectA := NewAABB(vec(0, 0), 2, 2)
	rectB := NewAABB(vec(4, 5), 2, 2)

	euclidean := NewEuclidean2D(20, 20)
	distance := euclidean.AABBDistance()(rectA.AABB, rectB.AABB)

	expected := euclidean.(*euclidean2d).metric(vec(2, 3), geom.NewVec(0, 0))
	if distance != expected {
		t.Errorf("expected distance %v, got %v", expected, distance)
	}
}

func TestAABBDistance_ReturnsZeroOnIntersection(t *testing.T) {
	rectA := NewAABB(vec(0, 0), 4, 4)
	rectB := NewAABB(vec(2, 2), 4, 4)

	euclidean := NewEuclidean2D(20, 20)
	distance := euclidean.AABBDistance()(rectA.AABB, rectB.AABB)
	if distance != 0 {
		t.Errorf("expected distance 0 for intersecting rectangles, got %v", distance)
	}
}

func TestAABBDistance_To_Vector_Euclidean2D(t *testing.T) {
	testCases := []struct {
		name   string
		first  geom.AABB
		second geom.AABB
		delta  geom.Vec
	}{
		{
			name:   "vectorAsPointLeft",
			first:  geom.NewAABBAt(vec(0, 0), 0, 0),
			second: NewAABB(vec(4, 0), 2, 2).AABB,
			delta:  vec(4, 0),
		},
		{
			name:   "vectorAsPointRight",
			first:  NewAABB(vec(0, 0), 2, 2).AABB,
			second: geom.NewAABBAt(vec(5, 6), 0, 0),
			delta:  vec(3, 4),
		},
	}

	euclidean := NewEuclidean2D(100, 100)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			distance := euclidean.AABBDistance()(tc.first, tc.second)
			expected := euclidean.(*euclidean2d).metric(tc.delta, geom.NewVec(0, 0))
			if distance != expected {
				t.Errorf("expected distance %v, got %v", expected, distance)
			}
		})
	}
}

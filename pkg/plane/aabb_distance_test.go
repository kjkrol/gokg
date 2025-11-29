package plane

import (
	"testing"

	"github.com/kjkrol/gokg/pkg/geom"
)

func TestAxisDistance(t *testing.T) {
	aa := newAABB(geom.NewVec(0, 0), 2, 2)
	bb := newAABB(geom.NewVec(5, 0), 2, 2)

	dx := aa.AxisDistanceTo(bb.AABB, func(v geom.Vec[int]) int { return v.X })
	dy := aa.AxisDistanceTo(bb.AABB, func(v geom.Vec[int]) int { return v.Y })

	if dx != 3 {
		t.Errorf("expected dx=3, got %d", dx)
	}
	if dy != 0 {
		t.Errorf("expected dy=0, got %d", dy)
	}
}

func TestAABBDistance_CartesianSpace(t *testing.T) {
	rectA := newAABB(geom.NewVec(0, 0), 2, 2)
	rectB := newAABB(geom.NewVec(4, 5), 2, 2)

	cartesian := NewCartesian(20, 20)
	distance := cartesian.AABBDistance()(rectA.AABB, rectB.AABB)

	expected := cartesian.(space2d[int]).metric(geom.Vec[int]{X: 2, Y: 3}, geom.NewVec(0, 0))
	if distance != expected {
		t.Errorf("expected distance %v, got %v", expected, distance)
	}
}

func TestAABBDistance_ReturnsZeroOnIntersection(t *testing.T) {
	rectA := newAABB(geom.NewVec(0, 0), 4, 4)
	rectB := newAABB(geom.NewVec(2, 2), 4, 4)

	cartesian := NewCartesian(20, 20)
	distance := cartesian.AABBDistance()(rectA.AABB, rectB.AABB)
	if distance != 0 {
		t.Errorf("expected distance 0 for intersecting rectangles, got %v", distance)
	}
}

func TestAABBDistance_To_Vector_Cartesian(t *testing.T) {
	vector := geom.NewVec(0, 0)
	aabbOfVector := geom.NewAABBAt(vector, 0, 0)
	aabb := newAABB(geom.NewVec(4, 0), 2, 2)

	cartesian := NewCartesian(100, 100)
	distance := cartesian.AABBDistance()(aabbOfVector, aabb.AABB)
	expected := cartesian.(space2d[int]).metric(geom.NewVec(4, 0), geom.NewVec(0, 0))
	if distance != expected {
		t.Errorf("expected distance %v, got %v", expected, distance)
	}
}

func TestAABBDistance_To_Vector_Cartesian_float(t *testing.T) {
	aabb := newAABB(geom.NewVec(0.0, 0), 2, 2)
	vector := geom.NewVec(5.0, 6.0)
	aabbOfVector := geom.NewAABBAt(vector, 0, 0)

	cartesian := NewCartesian(100.0, 100.0)
	distance := cartesian.AABBDistance()(aabb.AABB, aabbOfVector)
	expected := cartesian.(space2d[float64]).metric(geom.Vec[float64]{X: 3, Y: 4}, geom.NewVec(0.0, 0.0))
	if distance != expected {
		t.Errorf("expected distance %v, got %v", expected, distance)
	}
}

package geometry

import (
	"testing"
)

func TestRectanglesAxisDistance(t *testing.T) {
	aa := NewAABB(NewVec(0, 0), 2, 2)
	bb := NewAABB(NewVec(5, 0), 2, 2)

	dx := aa.axisDistanceTo(bb, func(v Vec[int]) int { return v.X })
	dy := aa.axisDistanceTo(bb, func(v Vec[int]) int { return v.Y })

	if dx != 3 {
		t.Errorf("expected dx=3, got %d", dx)
	}
	if dy != 0 {
		t.Errorf("expected dy=0, got %d", dy)
	}
}

func TestBoundingBoxDistance_ForBoundedPlane(t *testing.T) {
	rectA := NewAABB(NewVec(0, 0), 2, 2)
	rectB := NewAABB(NewVec(4, 5), 2, 2)

	plane := NewBoundedPlane(20, 20)
	distance := BoundingBoxDistance(plane.Metric)(rectA, rectB)

	expected := plane.Metric(Vec[int]{X: 2, Y: 3}, ZERO_INT_VEC)
	if distance != expected {
		t.Errorf("expected distance %v, got %v", expected, distance)
	}
}

func TestBoundingBoxDistance_ReturnsZeroOnIntersection(t *testing.T) {
	rectA := NewAABB(NewVec(0, 0), 4, 4)
	rectB := NewAABB(NewVec(2, 2), 4, 4)

	plane := NewBoundedPlane(20, 20)
	distance := BoundingBoxDistance(plane.Metric)(rectA, rectB)
	if distance != 0 {
		t.Errorf("expected distance 0 for intersecting rectangles, got %v", distance)
	}
}

func TestBoundingBoxDistance_VectorToAABB(t *testing.T) {
	vector := NewVec(0, 0)
	rect := NewAABB(NewVec(4, 0), 2, 2)
	plane := NewBoundedPlane(100, 100)

	distance := BoundingBoxDistanceForPlane(plane)(vector.Bounds(), rect)
	expected := plane.Metric(NewVec(4, 0), ZERO_INT_VEC)
	if distance != expected {
		t.Errorf("expected distance %v, got %v", expected, distance)
	}
}

func TestBoundingBoxDistance_AABBToVector(t *testing.T) {
	rect := NewAABB(NewVec(0.0, 0), 2, 2)
	vector := NewVec(5.0, 6.0)
	plane := NewBoundedPlane(100.0, 100.0)

	distance := BoundingBoxDistanceForPlane(plane)(rect, vector.Bounds())
	expected := plane.Metric(Vec[float64]{X: 3, Y: 4}, ZERO_FLOAT_VEC)
	if distance != expected {
		t.Errorf("expected distance %v, got %v", expected, distance)
	}
}

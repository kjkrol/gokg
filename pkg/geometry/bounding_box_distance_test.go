package geometry

import (
	"testing"
)

func TestRectanglesAxisDistance(t *testing.T) {
	aa := NewPlaneBox(NewVec(0, 0), 2, 2)
	bb := NewPlaneBox(NewVec(5, 0), 2, 2)

	dx := aa.axisDistanceTo(bb.BoundingBox, func(v Vec[int]) int { return v.X })
	dy := aa.axisDistanceTo(bb.BoundingBox, func(v Vec[int]) int { return v.Y })

	if dx != 3 {
		t.Errorf("expected dx=3, got %d", dx)
	}
	if dy != 0 {
		t.Errorf("expected dy=0, got %d", dy)
	}
}

func TestBoundingBoxDistance_ForBoundedPlane(t *testing.T) {
	rectA := NewPlaneBox(NewVec(0, 0), 2, 2)
	rectB := NewPlaneBox(NewVec(4, 5), 2, 2)

	plane := NewBoundedPlane(20, 20)
	distance := BoundingBoxDistance(plane.Metric)(rectA.BoundingBox, rectB.BoundingBox)

	expected := plane.Metric(Vec[int]{X: 2, Y: 3}, ZERO_INT_VEC)
	if distance != expected {
		t.Errorf("expected distance %v, got %v", expected, distance)
	}
}

func TestBoundingBoxDistance_ReturnsZeroOnIntersection(t *testing.T) {
	rectA := NewPlaneBox(NewVec(0, 0), 4, 4)
	rectB := NewPlaneBox(NewVec(2, 2), 4, 4)

	plane := NewBoundedPlane(20, 20)
	distance := BoundingBoxDistance(plane.Metric)(rectA.BoundingBox, rectB.BoundingBox)
	if distance != 0 {
		t.Errorf("expected distance 0 for intersecting rectangles, got %v", distance)
	}
}

func TestBoundingBoxDistance_VectorToAABB(t *testing.T) {
	vector := NewVec(0, 0)
	rect := NewPlaneBox(NewVec(4, 0), 2, 2)
	plane := NewBoundedPlane(100, 100)
	vectorBox := NewBoundingBoxAt(vector, 0, 0)

	distance := BoundingBoxDistanceForPlane(plane)(vectorBox, rect.BoundingBox)
	expected := plane.Metric(NewVec(4, 0), ZERO_INT_VEC)
	if distance != expected {
		t.Errorf("expected distance %v, got %v", expected, distance)
	}
}

func TestBoundingBoxDistance_AABBToVector(t *testing.T) {
	rect := NewPlaneBox(NewVec(0.0, 0), 2, 2)
	vector := NewVec(5.0, 6.0)
	plane := NewBoundedPlane(100.0, 100.0)
	vectorBox := NewBoundingBoxAt(vector, 0, 0)

	distance := BoundingBoxDistanceForPlane(plane)(rect.BoundingBox, vectorBox)
	expected := plane.Metric(Vec[float64]{X: 3, Y: 4}, ZERO_FLOAT_VEC)
	if distance != expected {
		t.Errorf("expected distance %v, got %v", expected, distance)
	}
}

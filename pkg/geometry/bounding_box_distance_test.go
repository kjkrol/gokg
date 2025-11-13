package geometry

import (
	"testing"
)

func TestRectanglesAxisDistance(t *testing.T) {
	aa := newPlaneBox(NewVec(0, 0), 2, 2)
	bb := newPlaneBox(NewVec(5, 0), 2, 2)

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
	rectA := newPlaneBox(NewVec(0, 0), 2, 2)
	rectB := newPlaneBox(NewVec(4, 5), 2, 2)

	plane := NewBoundedPlane(20, 20)
	distance := plane.BoundingBoxDistance()(rectA.BoundingBox, rectB.BoundingBox)

	expected := plane.metric(Vec[int]{X: 2, Y: 3}, ZERO_INT_VEC)
	if distance != expected {
		t.Errorf("expected distance %v, got %v", expected, distance)
	}
}

func TestBoundingBoxDistance_ReturnsZeroOnIntersection(t *testing.T) {
	rectA := newPlaneBox(NewVec(0, 0), 4, 4)
	rectB := newPlaneBox(NewVec(2, 2), 4, 4)

	plane := NewBoundedPlane(20, 20)
	distance := plane.BoundingBoxDistance()(rectA.BoundingBox, rectB.BoundingBox)
	if distance != 0 {
		t.Errorf("expected distance 0 for intersecting rectangles, got %v", distance)
	}
}

func TestBoundingBoxDistance_VectorToBox(t *testing.T) {
	vector := NewVec(0, 0)
	rect := newPlaneBox(NewVec(4, 0), 2, 2)
	plane := NewBoundedPlane(100, 100)
	vectorBox := NewBoundingBoxAt(vector, 0, 0)

	distance := plane.BoundingBoxDistance()(vectorBox, rect.BoundingBox)
	expected := plane.metric(NewVec(4, 0), ZERO_INT_VEC)
	if distance != expected {
		t.Errorf("expected distance %v, got %v", expected, distance)
	}
}

func TestBoundingBoxDistance_BoxToVector(t *testing.T) {
	rect := newPlaneBox(NewVec(0.0, 0), 2, 2)
	vector := NewVec(5.0, 6.0)
	plane := NewBoundedPlane(100.0, 100.0)
	vectorBox := NewBoundingBoxAt(vector, 0, 0)

	distance := plane.BoundingBoxDistance()(rect.BoundingBox, vectorBox)
	expected := plane.metric(Vec[float64]{X: 3, Y: 4}, ZERO_FLOAT_VEC)
	if distance != expected {
		t.Errorf("expected distance %v, got %v", expected, distance)
	}
}

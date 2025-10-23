package geometry

import "testing"

func TestRectanglesAxisDistance(t *testing.T) {
	aa := NewRectangle(Vec[int]{X: 0, Y: 0}, Vec[int]{X: 2, Y: 2})
	bb := NewRectangle(Vec[int]{X: 5, Y: 0}, Vec[int]{X: 7, Y: 2})

	dx := rectanglesAxisDistance(aa, bb, func(v Vec[int]) int { return v.X })
	dy := rectanglesAxisDistance(aa, bb, func(v Vec[int]) int { return v.Y })

	if dx != 3 {
		t.Errorf("expected dx=3, got %d", dx)
	}
	if dy != 0 {
		t.Errorf("expected dy=0, got %d", dy)
	}
}

func TestBoundingBoxDistance(t *testing.T) {
	rectA := NewRectangle(Vec[int]{X: 0, Y: 0}, Vec[int]{X: 2, Y: 2})
	rectB := NewRectangle(Vec[int]{X: 4, Y: 5}, Vec[int]{X: 6, Y: 7})

	plane := NewBoundedPlane(20, 20)
	distance := BoundingBoxDistance[int](plane.Metric)(&rectA, &rectB)

	expected := plane.Metric(Vec[int]{X: 2, Y: 3}, ZERO_INT_VEC)
	if distance != expected {
		t.Errorf("expected distance %v, got %v", expected, distance)
	}
}

func TestBoundingBoxDistance_ReturnsZeroOnIntersection(t *testing.T) {
	rectA := NewRectangle(Vec[int]{X: 0, Y: 0}, Vec[int]{X: 4, Y: 4})
	rectB := NewRectangle(Vec[int]{X: 2, Y: 2}, Vec[int]{X: 6, Y: 6})

	plane := NewBoundedPlane(20, 20)
	distance := BoundingBoxDistance[int](plane.Metric)(&rectA, &rectB)
	if distance != 0 {
		t.Errorf("expected distance 0 for intersecting rectangles, got %v", distance)
	}
}

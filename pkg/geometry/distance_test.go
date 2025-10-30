package geometry

import (
	"testing"
)

func TestBoundingBoxDistance_ForBoundedPlane_LineToLine(t *testing.T) {
	plane := NewBoundedPlane(10, 10)
	distanceFun := BoundingBoxDistanceForPlane(plane)

	line := NewLine(Vec[int]{X: 1, Y: 1}, Vec[int]{X: 4, Y: 3})
	other := NewLine(Vec[int]{X: 8, Y: 8}, Vec[int]{X: 9, Y: 9})
	distance := distanceFun(&line, &other)
	expectedDistance := plane.Metric(Vec[int]{X: 4, Y: 5}, ZERO_INT_VEC)
	if distance != expectedDistance {
		t.Errorf("expected distance %v, got %v", expectedDistance, distance)
	}
}

func TestBoundingBoxDistance_ForBoundedPlane_PolygonToPolygon(t *testing.T) {
	plane := NewBoundedPlane(10.0, 10.0)
	distanceFun := BoundingBoxDistanceForPlane(plane)

	polyA := NewPolygon(
		Vec[float64]{X: 0, Y: 0},
		Vec[float64]{X: 2, Y: 0},
		Vec[float64]{X: 1, Y: 2},
	)
	polyB := NewPolygon(
		Vec[float64]{X: 5, Y: 5},
		Vec[float64]{X: 7, Y: 5},
		Vec[float64]{X: 6, Y: 7},
	)

	distance := distanceFun(&polyA, &polyB)

	expected := plane.Metric(
		Vec[float64]{X: 3, Y: 3},
		Vec[float64]{},
	)
	if distance != expected {
		t.Errorf("expected distance %v, got %v", expected, distance)
	}
}

func TestRectanglesAxisDistance(t *testing.T) {
	aa := NewAABB(Vec[int]{X: 0, Y: 0}, Vec[int]{X: 2, Y: 2})
	bb := NewAABB(Vec[int]{X: 5, Y: 0}, Vec[int]{X: 7, Y: 2})

	dx := aa.AxisDistanceTo(bb, func(v Vec[int]) int { return v.X })
	dy := aa.AxisDistanceTo(bb, func(v Vec[int]) int { return v.Y })

	if dx != 3 {
		t.Errorf("expected dx=3, got %d", dx)
	}
	if dy != 0 {
		t.Errorf("expected dy=0, got %d", dy)
	}
}

func TestBoundingBoxDistance_ForBoundedPlane(t *testing.T) {
	rectA := NewPolygon(NewVec(0, 0), NewVec(2, 0), NewVec(2, 2), NewVec(0, 2))
	rectB := NewPolygon(NewVec(4, 5), NewVec(6, 5), NewVec(6, 7), NewVec(4, 7))

	plane := NewBoundedPlane(20, 20)
	distance := BoundingBoxDistance(plane.Metric)(&rectA, &rectB)

	expected := plane.Metric(Vec[int]{X: 2, Y: 3}, ZERO_INT_VEC)
	if distance != expected {
		t.Errorf("expected distance %v, got %v", expected, distance)
	}
}

func TestBoundingBoxDistance_ReturnsZeroOnIntersection(t *testing.T) {
	rectA := NewPolygon(NewVec(0, 0), NewVec(4, 0), NewVec(4, 4), NewVec(0, 4))
	rectB := NewPolygon(NewVec(2, 2), NewVec(6, 2), NewVec(6, 6), NewVec(2, 6))

	plane := NewBoundedPlane(20, 20)
	distance := BoundingBoxDistance(plane.Metric)(&rectA, &rectB)
	if distance != 0 {
		t.Errorf("expected distance 0 for intersecting rectangles, got %v", distance)
	}
}

func TestBoundingBoxDistance_VectorToAABB(t *testing.T) {
	vector := NewVec(0, 0)
	rect := NewPolygon(NewVec(4, 0), NewVec(6, 0), NewVec(6, 2), NewVec(4, 2))
	plane := NewBoundedPlane(100, 100)

	distance := BoundingBoxDistanceForPlane(plane)(&vector, &rect)
	expected := plane.Metric(NewVec(4, 0), ZERO_INT_VEC)
	if distance != expected {
		t.Errorf("expected distance %v, got %v", expected, distance)
	}
}

func TestBoundingBoxDistance_RectangleToVector(t *testing.T) {
	rect := NewPolygon(NewVec(0.0, 0.0), NewVec(2.0, 0.0), NewVec(2.0, 2.0), NewVec(0.0, 2.0))
	vector := NewVec(5.0, 6.0)
	plane := NewBoundedPlane(100.0, 100.0)

	distance := BoundingBoxDistanceForPlane(plane)(&rect, &vector)
	expected := plane.Metric(Vec[float64]{X: 3, Y: 4}, ZERO_FLOAT_VEC)
	if distance != expected {
		t.Errorf("expected distance %v, got %v", expected, distance)
	}
}

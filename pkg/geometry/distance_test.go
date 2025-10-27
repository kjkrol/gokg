package geometry

import (
	"testing"

	s "github.com/kjkrol/gokg/pkg/geometry/spatial"
)

func TestBoundingBoxDistance_ForBoundedPlane_LineToLine(t *testing.T) {
	plane := NewBoundedPlane(10, 10)
	distanceFun := BoundingBoxDistanceForPlane(plane)

	line := s.NewLine(vec[int]{X: 1, Y: 1}, vec[int]{X: 4, Y: 3})
	other := s.NewLine(vec[int]{X: 8, Y: 8}, vec[int]{X: 9, Y: 9})
	distance := distanceFun(&line, &other)
	expectedDistance := plane.Metric(vec[int]{X: 4, Y: 5}, s.ZERO_INT_VEC)
	if distance != expectedDistance {
		t.Errorf("expected distance %v, got %v", expectedDistance, distance)
	}
}

func TestBoundingBoxDistance_ForBoundedPlane_PolygonToPolygon(t *testing.T) {
	plane := NewBoundedPlane(10.0, 10.0)
	distanceFun := BoundingBoxDistanceForPlane(plane)

	polyA := s.NewPolygon(
		vec[float64]{X: 0, Y: 0},
		vec[float64]{X: 2, Y: 0},
		vec[float64]{X: 1, Y: 2},
	)
	polyB := s.NewPolygon(
		vec[float64]{X: 5, Y: 5},
		vec[float64]{X: 7, Y: 5},
		vec[float64]{X: 6, Y: 7},
	)

	distance := distanceFun(&polyA, &polyB)

	expected := plane.Metric(
		vec[float64]{X: 3, Y: 3},
		vec[float64]{},
	)
	if distance != expected {
		t.Errorf("expected distance %v, got %v", expected, distance)
	}
}

func TestRectanglesAxisDistance(t *testing.T) {
	aa := s.NewRectangle(vec[int]{X: 0, Y: 0}, vec[int]{X: 2, Y: 2})
	bb := s.NewRectangle(vec[int]{X: 5, Y: 0}, vec[int]{X: 7, Y: 2})

	dx := aa.AxisDistanceTo(bb, func(v vec[int]) int { return v.X })
	dy := aa.AxisDistanceTo(bb, func(v vec[int]) int { return v.Y })

	if dx != 3 {
		t.Errorf("expected dx=3, got %d", dx)
	}
	if dy != 0 {
		t.Errorf("expected dy=0, got %d", dy)
	}
}

func TestBoundingBoxDistance_ForBoundedPlane_RectangleToRectanle(t *testing.T) {
	rectA := s.NewRectangle(vec[int]{X: 0, Y: 0}, vec[int]{X: 2, Y: 2})
	rectB := s.NewRectangle(vec[int]{X: 4, Y: 5}, vec[int]{X: 6, Y: 7})

	plane := NewBoundedPlane(20, 20)
	distance := BoundingBoxDistance(plane.Metric)(&rectA, &rectB)

	expected := plane.Metric(vec[int]{X: 2, Y: 3}, s.ZERO_INT_VEC)
	if distance != expected {
		t.Errorf("expected distance %v, got %v", expected, distance)
	}
}

func TestBoundingBoxDistance_ReturnsZeroOnIntersection(t *testing.T) {
	rectA := s.NewRectangle(vec[int]{X: 0, Y: 0}, vec[int]{X: 4, Y: 4})
	rectB := s.NewRectangle(vec[int]{X: 2, Y: 2}, vec[int]{X: 6, Y: 6})

	plane := NewBoundedPlane(20, 20)
	distance := BoundingBoxDistance(plane.Metric)(&rectA, &rectB)
	if distance != 0 {
		t.Errorf("expected distance 0 for intersecting rectangles, got %v", distance)
	}
}

func TestBoundingBoxDistance_VectorToRectangle(t *testing.T) {
	vector := s.NewVec(0, 0)
	rect := s.NewRectangle(s.NewVec(4, 0), s.NewVec(6, 2))
	plane := NewBoundedPlane(100, 100)

	distance := BoundingBoxDistanceForPlane(plane)(&vector, &rect)
	expected := plane.Metric(s.NewVec(4, 0), s.ZERO_INT_VEC)
	if distance != expected {
		t.Errorf("expected distance %v, got %v", expected, distance)
	}
}

func TestBoundingBoxDistance_RectangleToVector(t *testing.T) {
	rect := s.NewRectangle(vec[float64]{X: 0, Y: 0}, vec[float64]{X: 2, Y: 2})
	vector := vec[float64]{X: 5, Y: 6}
	plane := NewBoundedPlane(100.0, 100.0)

	distance := BoundingBoxDistanceForPlane(plane)(&rect, &vector)
	expected := plane.Metric(vec[float64]{X: 3, Y: 4}, s.ZERO_FLOAT_VEC)
	if distance != expected {
		t.Errorf("expected distance %v, got %v", expected, distance)
	}
}

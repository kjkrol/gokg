package geometry

import "testing"

func TestWrapSpatialFragments_PolygonCrossesRightEdge(t *testing.T) {
	size := Vec[int]{X: 10, Y: 10}
	vecMath := VectorMathByType[int]()
	shape := NewPolygon(
		Vec[int]{X: 8, Y: 4},
		Vec[int]{X: 12, Y: 4},
		Vec[int]{X: 12, Y: 6},
		Vec[int]{X: 8, Y: 6},
	)
	fragments := wrapSpatialFragments(&shape, size, vecMath)
	if len(fragments) != 1 {
		t.Fatalf("expected 1 fragment, got %d", len(fragments))
	}
	expectedFirstPoints := map[[2]int]struct{}{
		{-2, 4}: {},
	}
	for _, fragment := range fragments {
		poly, ok := fragment.(*Polygon[int])
		if !ok {
			t.Fatalf("expected polygon fragment, got %T", fragment)
		}
		points := poly.Points()
		if len(points) == 0 {
			t.Fatal("polygon fragment has no points")
		}
		key := [2]int{points[0].X, points[0].Y}
		if _, ok := expectedFirstPoints[key]; !ok {
			t.Fatalf("unexpected fragment starting point %+v", key)
		}
		delete(expectedFirstPoints, key)
	}
	if len(expectedFirstPoints) != 0 {
		t.Fatalf("missing expected fragments: %+v", expectedFirstPoints)
	}
}

func TestWrapSpatialFragments_PolygonCrossesCorner(t *testing.T) {
	size := Vec[int]{X: 10, Y: 10}
	vecMath := VectorMathByType[int]()
	shape := NewPolygon(
		Vec[int]{X: 9, Y: 9},
		Vec[int]{X: 12, Y: 9},
		Vec[int]{X: 12, Y: 12},
		Vec[int]{X: 9, Y: 12},
	)
	fragments := wrapSpatialFragments(&shape, size, vecMath)
	if len(fragments) != 3 {
		t.Fatalf("expected 3 fragments, got %d", len(fragments))
	}
	expectedFirstPoints := map[[2]int]struct{}{
		{-1, 9}:  {},
		{9, -1}:  {},
		{-1, -1}: {},
	}
	for _, fragment := range fragments {
		poly, ok := fragment.(*Polygon[int])
		if !ok {
			t.Fatalf("expected polygon fragment, got %T", fragment)
		}
		points := poly.Points()
		if len(points) == 0 {
			t.Fatal("polygon fragment has no points")
		}
		key := [2]int{points[0].X, points[0].Y}
		if _, ok := expectedFirstPoints[key]; !ok {
			t.Fatalf("unexpected fragment starting point %+v", key)
		}
		delete(expectedFirstPoints, key)
	}
	if len(expectedFirstPoints) != 0 {
		t.Fatalf("missing expected fragments: %+v", expectedFirstPoints)
	}
}

func TestTranslateSpatial_SetsFragmentsOnSpatial(t *testing.T) {
	plane := NewCyclicBoundedPlane(10, 10)
	poly := NewPolygon(
		Vec[int]{X: 9, Y: 9},
		Vec[int]{X: 11, Y: 9},
		Vec[int]{X: 11, Y: 11},
		Vec[int]{X: 9, Y: 11},
	)
	plane.TranslateSpatial(&poly, Vec[int]{X: 0, Y: 0})
	fragments := poly.Fragments()
	if len(fragments) != 3 {
		t.Fatalf("expected 3 fragments set on polygon, got %d", len(fragments))
	}
}

func TestTranslateSpatial_ClearsFragmentsWhenNotWrapping(t *testing.T) {
	plane := NewCyclicBoundedPlane(10, 10)
	poly := NewPolygon(
		Vec[int]{X: 1, Y: 1},
		Vec[int]{X: 2, Y: 1},
		Vec[int]{X: 2, Y: 2},
		Vec[int]{X: 1, Y: 2},
	)
	plane.TranslateSpatial(&poly, Vec[int]{X: 1, Y: 0})
	if frags := poly.Fragments(); frags != nil {
		t.Fatalf("expected no fragments for polygon inside bounds, got %d", len(frags))
	}
}

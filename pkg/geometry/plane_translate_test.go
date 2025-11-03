package geometry

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

func TestWrapSpatialFragments_PolygonCrossesRightEdge(t *testing.T) {
	plane := NewCyclicBoundedPlane(10, 10)
	shape := NewPolygonBuilder[int]().
		Add(8, 4).
		Add(12, 4).
		Add(12, 6).
		Add(8, 6).
		Build()
	plane.createShapeFragments(&shape)
	if len(shape.fragments) != 2 {
		t.Fatalf("expected 2 fragments, got %d", len(shape.fragments))
	}
	expectedFragments := map[string]struct{}{
		polygonKey([]Vec[int]{
			{X: 8, Y: 4},
			{X: 9, Y: 4},
			{X: 9, Y: 6},
			{X: 8, Y: 6},
		}): {},
		polygonKey([]Vec[int]{
			{X: 0, Y: 4},
			{X: 2, Y: 4},
			{X: 2, Y: 6},
			{X: 0, Y: 6},
		}): {},
	}
	for _, fragment := range shape.fragments {
		poly, ok := fragment.(*Polygon[int])
		if !ok {
			t.Fatalf("expected polygon fragment, got %T", fragment)
		}
		points := poly.Points()
		if len(points) == 0 {
			t.Fatal("polygon fragment has no points")
		}
		key := polygonKey(points)
		if _, ok := expectedFragments[key]; !ok {
			t.Fatalf("unexpected fragment points %+v", points)
		}
		delete(expectedFragments, key)
	}
	if len(expectedFragments) != 0 {
		t.Fatalf("missing expected fragments: %+v", expectedFragments)
	}
}

func TestWrapFragments_PolygonCrossesCorner(t *testing.T) {
	plane := NewCyclicBoundedPlane(10, 10)
	shape := NewPolygon(
		Vec[int]{X: 9, Y: 9},
		Vec[int]{X: 12, Y: 9},
		Vec[int]{X: 12, Y: 12},
		Vec[int]{X: 9, Y: 12},
	)
	plane.createShapeFragments(&shape)
	if len(shape.fragments) != 4 {
		t.Fatalf("expected 4 fragments, got %d", len(shape.fragments))
	}
	expectedFragments := map[string]struct{}{
		polygonKey([]Vec[int]{
			{X: 0, Y: 9},
			{X: 2, Y: 9},
			{X: 2, Y: 9},
			{X: 0, Y: 9},
		}): {},
		polygonKey([]Vec[int]{
			{X: 9, Y: 9},
			{X: 9, Y: 9},
			{X: 9, Y: 9},
			{X: 9, Y: 9},
		}): {},
		polygonKey([]Vec[int]{
			{X: 9, Y: 0},
			{X: 9, Y: 0},
			{X: 9, Y: 2},
			{X: 9, Y: 2},
		}): {},
		polygonKey([]Vec[int]{
			{X: 0, Y: 0},
			{X: 2, Y: 0},
			{X: 2, Y: 2},
			{X: 0, Y: 2},
		}): {},
	}
	for _, fragment := range shape.fragments {
		poly, ok := fragment.(*Polygon[int])
		if !ok {
			t.Fatalf("expected polygon fragment, got %T", fragment)
		}
		points := poly.Points()
		if len(points) == 0 {
			t.Fatal("polygon fragment has no points")
		}
		key := polygonKey(points)
		if _, ok := expectedFragments[key]; !ok {
			t.Fatalf("unexpected fragment points %+v", points)
		}
		delete(expectedFragments, key)
	}
	if len(expectedFragments) != 0 {
		t.Fatalf("missing expected fragments: %+v", expectedFragments)
	}
}

func polygonKey[T SupportedNumeric](points []Vec[T]) string {
	if len(points) == 0 {
		return ""
	}

	sorted := make([]Vec[T], len(points))
	copy(sorted, points)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].X == sorted[j].X {
			return sorted[i].Y < sorted[j].Y
		}
		return sorted[i].X < sorted[j].X
	})

	keys := make([]string, len(sorted))
	for i, p := range sorted {
		keys[i] = fmt.Sprintf("(%v,%v)", p.X, p.Y)
	}
	return strings.Join(keys, ";")
}

func TestTranslate_SetsFragmentsOnSpatial(t *testing.T) {
	plane := NewCyclicBoundedPlane(10, 10)
	poly := NewPolygon(
		Vec[int]{X: 9, Y: 9},
		Vec[int]{X: 11, Y: 9},
		Vec[int]{X: 11, Y: 11},
		Vec[int]{X: 9, Y: 11},
	)
	plane.Translate(&poly, Vec[int]{X: 0, Y: 0})
	fragments := poly.Fragments()
	if len(fragments) != 4 {
		t.Fatalf("expected 3 fragments set on polygon, got %d", len(fragments))
	}
}

func TestTranslate_ClearsFragmentsWhenNotWrapping(t *testing.T) {
	plane := NewCyclicBoundedPlane(10, 10)
	poly := NewPolygon(
		Vec[int]{X: 1, Y: 1},
		Vec[int]{X: 2, Y: 1},
		Vec[int]{X: 2, Y: 2},
		Vec[int]{X: 1, Y: 2},
	)
	plane.Translate(&poly, Vec[int]{X: 1, Y: 0})
	if frags := poly.Fragments(); len(frags) > 0 {
		t.Fatalf("expected no fragments for polygon inside bounds, got %d", len(frags))
	}
}

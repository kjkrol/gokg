package geometry

import "testing"

func TestNewPolygon_Bounds(t *testing.T) {
	poly := NewPolygon(
		Vec[int]{X: 5, Y: 5},
		Vec[int]{X: 1, Y: 1},
		Vec[int]{X: 3, Y: 6},
	)

	bounds := poly.Bounds()
	expectedTopLeft := Vec[int]{X: 1, Y: 1}
	expectedBottomRight := Vec[int]{X: 5, Y: 6}

	if bounds.TopLeft != expectedTopLeft {
		t.Errorf("expected top-left %v, got %v", expectedTopLeft, bounds.TopLeft)
	}
	if bounds.BottomRight != expectedBottomRight {
		t.Errorf("expected bottom-right %v, got %v", expectedBottomRight, bounds.BottomRight)
	}
}

func TestPolygon_Probe(t *testing.T) {
	poly := NewPolygon(
		Vec[int]{X: 0, Y: 0},
		Vec[int]{X: 2, Y: 0},
		Vec[int]{X: 1, Y: 2},
	)
	plane := NewBoundedPlane(100, 100)

	probes := poly.Probe(1, plane)
	if len(probes) != 1 {
		t.Fatalf("expected single probe rectangle, got %d", len(probes))
	}

	expectedTopLeft := Vec[int]{X: -1, Y: -1}
	expectedBottomRight := Vec[int]{X: 3, Y: 3}
	if probes[0].TopLeft != expectedTopLeft {
		t.Errorf("expected probe top-left %v, got %v", expectedTopLeft, probes[0].TopLeft)
	}
	if probes[0].BottomRight != expectedBottomRight {
		t.Errorf("expected probe bottom-right %v, got %v", expectedBottomRight, probes[0].BottomRight)
	}
}

func TestPolygon_DistanceTo(t *testing.T) {
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

	plane := NewBoundedPlane(10.0, 10.0)
	dist := polyA.DistanceTo(&polyB, BoundingBoxDistanceForPlane(plane))

	expected := plane.Metric(
		Vec[float64]{X: 3, Y: 3},
		Vec[float64]{},
	)
	if dist != expected {
		t.Errorf("expected distance %v, got %v", expected, dist)
	}
}

func TestNewPolygon_PanicOnInvalidInput(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic for invalid polygon, got none")
		}
	}()

	NewPolygon[int](Vec[int]{X: 0, Y: 0}, Vec[int]{X: 1, Y: 1})
}

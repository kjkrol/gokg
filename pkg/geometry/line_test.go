package geometry

import "testing"

func TestLine_SpatialMethods(t *testing.T) {
	line := NewLine(Vec[int]{X: 1, Y: 1}, Vec[int]{X: 4, Y: 3})
	bounds := line.Bounds()
	expectedBounds := NewRectangle(Vec[int]{X: 1, Y: 1}, Vec[int]{X: 4, Y: 3})
	if !rectEquals(bounds, expectedBounds) {
		t.Errorf("expected bounds %v, got %v", expectedBounds, bounds)
	}

	plane := NewBoundedPlane(10, 10)
	probes := line.Probe(1, plane)
	if len(probes) != 1 {
		t.Fatalf("expected a single probe rectangle, got %d", len(probes))
	}
	expectedProbe := NewRectangle(Vec[int]{X: 0, Y: 0}, Vec[int]{X: 5, Y: 4})
	if !rectEquals(probes[0], expectedProbe) {
		t.Errorf("expected expanded rectangle %v, got %v", expectedProbe, probes[0])
	}

	other := NewLine(Vec[int]{X: 8, Y: 8}, Vec[int]{X: 9, Y: 9})
	distance := line.DistanceTo(&other, plane.Metric)
	expectedDistance := plane.Metric(Vec[int]{X: 4, Y: 5}, ZERO_INT_VEC)
	if distance != expectedDistance {
		t.Errorf("expected distance %v, got %v", expectedDistance, distance)
	}
}

func TestLine_ProbeCyclicWrap(t *testing.T) {
	line := NewLine(Vec[int]{X: 9, Y: 9}, Vec[int]{X: 11, Y: 10})
	plane := NewCyclicBoundedPlane(10, 10)

	probes := line.Probe(0, plane)
	if len(probes) <= 1 {
		t.Fatalf("expected wrapped probe rectangles, got %d", len(probes))
	}
}

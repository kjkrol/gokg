package geometry

import "testing"

func TestVec_Probe_WrapsOnCyclicPlane(t *testing.T) {
	vec := Vec[int]{X: 9, Y: 9}
	plane := NewCyclicBoundedPlane(10, 10)

	probes := vec.Probe(1, plane)
	if len(probes) <= 1 {
		t.Fatalf("expected additional wrapped probes, got %d", len(probes))
	}

	wrappedFound := false
	expectedTopLeft := Vec[int]{X: 8, Y: 8}
	expectedBottomRight := Vec[int]{X: 10, Y: 10}
	for _, p := range probes {
		if p.TopLeft != expectedTopLeft || p.BottomRight != expectedBottomRight {
			wrappedFound = true
		}
	}
	if !wrappedFound {
		t.Errorf("expected wrapped rectangle in probes %v", probes)
	}
}

func TestVec_DistanceTo_Rectangle(t *testing.T) {
	vec := Vec[int]{X: 0, Y: 0}
	rect := NewRectangle(Vec[int]{X: 4, Y: 0}, Vec[int]{X: 6, Y: 2})
	plane := NewBoundedPlane(100, 100)

	distance := vec.DistanceTo(&rect, BoundingBoxDistanceForPlane(plane))
	expected := plane.Metric(Vec[int]{X: 4, Y: 0}, ZERO_INT_VEC)
	if distance != expected {
		t.Errorf("expected distance %v, got %v", expected, distance)
	}
}

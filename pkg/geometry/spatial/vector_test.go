package spatial

// func TestVec_Probe_WrapsOnCyclicPlane(t *testing.T) {
// 	vec := Vec[int]{X: 9, Y: 9}
// 	plane := NewCyclicBoundedPlane(10, 10)

// 	probes := vec.Probe(1, plane)
// 	if len(probes) <= 1 {
// 		t.Fatalf("expected additional wrapped probes, got %d", len(probes))
// 	}

// 	wrappedFound := false
// 	expectedTopLeft := Vec[int]{X: 8, Y: 8}
// 	expectedBottomRight := Vec[int]{X: 10, Y: 10}
// 	for _, p := range probes {
// 		if p.TopLeft != expectedTopLeft || p.BottomRight != expectedBottomRight {
// 			wrappedFound = true
// 		}
// 	}
// 	if !wrappedFound {
// 		t.Errorf("expected wrapped rectangle in probes %v", probes)
// 	}
// }

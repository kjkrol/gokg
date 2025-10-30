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

func TestNewPolygon_PanicOnInvalidInput(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic for invalid polygon, got none")
		}
	}()

	NewPolygon(Vec[int]{X: 0, Y: 0}, Vec[int]{X: 1, Y: 1})
}

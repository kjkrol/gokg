package spatial

import "testing"

func TestLine_SpatialMethods(t *testing.T) {
	line := NewLine(Vec[int]{X: 1, Y: 1}, Vec[int]{X: 4, Y: 3})
	bounds := line.Bounds()
	expectedBounds := NewRectangle(Vec[int]{X: 1, Y: 1}, Vec[int]{X: 4, Y: 3})
	if !bounds.Equals(expectedBounds) {
		t.Errorf("expected bounds %v, got %v", expectedBounds, bounds)
	}
}

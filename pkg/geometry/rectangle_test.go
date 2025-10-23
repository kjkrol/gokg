package geometry

import "testing"

func TestRectangle_NewRectangle(t *testing.T) {
	rect := NewRectangle(ZERO_INT_VEC, Vec[int]{X: 10, Y: 10})
	expected := Vec[int]{X: 5, Y: 5}
	if rect.Center != expected {
		t.Errorf("center %v not equal to expected %v", rect.Center, expected)
	}
}

func TestRectangle_BuildRectangle(t *testing.T) {
	center := Vec[int]{X: 5, Y: 5}
	rect := BuildRectangle(center, 2)
	expectedTopLeft := Vec[int]{X: 3, Y: 3}
	expectedBottomRight := Vec[int]{X: 7, Y: 7}
	if rect.TopLeft != expectedTopLeft {
		t.Errorf("topLeft %v not equal to expected %v", rect.TopLeft, expectedTopLeft)
	}
	if rect.BottomRight != expectedBottomRight {
		t.Errorf("bottomRight %v not equal to expected %v", rect.BottomRight, expectedBottomRight)
	}
}

func TestRectangle_Split(t *testing.T) {
	parent := NewRectangle(ZERO_INT_VEC, Vec[int]{X: 10, Y: 10})
	splitted := parent.Split()

	expected := [4]Rectangle[int]{
		NewRectangle(ZERO_INT_VEC, Vec[int]{X: 5, Y: 5}),
		NewRectangle(Vec[int]{X: 5, Y: 0}, Vec[int]{X: 10, Y: 5}),
		NewRectangle(Vec[int]{X: 0, Y: 5}, Vec[int]{X: 5, Y: 10}),
		NewRectangle(Vec[int]{X: 5, Y: 5}, Vec[int]{X: 10, Y: 10}),
	}
	for i := 0; i < 4; i++ {
		if !rectEquals(splitted[i], expected[i]) {
			t.Errorf("split %v not equal to expected %v", splitted[i], expected[i])
		}
	}
}

func TestRectangle_Intersects(t *testing.T) {
	intersects := []struct{ rect1, rect2 Rectangle[int] }{
		{
			rect1: NewRectangle(Vec[int]{X: 3, Y: 3}, Vec[int]{X: 5, Y: 5}),
			rect2: NewRectangle(Vec[int]{X: 2, Y: 2}, Vec[int]{X: 4, Y: 4}),
		},
		{
			rect1: NewRectangle(Vec[int]{X: 3, Y: 3}, Vec[int]{X: 5, Y: 5}),
			rect2: NewRectangle(Vec[int]{X: 4, Y: 2}, Vec[int]{X: 6, Y: 4}),
		},
		{
			rect1: NewRectangle(Vec[int]{X: 3, Y: 3}, Vec[int]{X: 5, Y: 5}),
			rect2: NewRectangle(Vec[int]{X: 2, Y: 4}, Vec[int]{X: 4, Y: 6}),
		},
		{
			rect1: NewRectangle(Vec[int]{X: 3, Y: 3}, Vec[int]{X: 5, Y: 5}),
			rect2: NewRectangle(Vec[int]{X: 4, Y: 4}, Vec[int]{X: 6, Y: 6}),
		},
		{
			rect1: NewRectangle(Vec[int]{X: 3, Y: 3}, Vec[int]{X: 5, Y: 5}),
			rect2: NewRectangle(Vec[int]{X: 2, Y: 2}, Vec[int]{X: 6, Y: 6}),
		},
		{
			rect1: NewRectangle(Vec[int]{X: 3, Y: 3}, Vec[int]{X: 5, Y: 5}),
			rect2: NewRectangle(Vec[int]{X: 2, Y: 2}, Vec[int]{X: 3, Y: 4}),
		},
		{
			rect1: NewRectangle(Vec[int]{X: 3, Y: 3}, Vec[int]{X: 5, Y: 5}),
			rect2: NewRectangle(Vec[int]{X: 2, Y: 4}, Vec[int]{X: 3, Y: 6}),
		},
		{
			rect1: NewRectangle(Vec[int]{X: 3, Y: 3}, Vec[int]{X: 5, Y: 5}),
			rect2: NewRectangle(Vec[int]{X: 2, Y: 2}, Vec[int]{X: 6, Y: 3}),
		},
		{
			rect1: NewRectangle(Vec[int]{X: 3, Y: 3}, Vec[int]{X: 5, Y: 5}),
			rect2: NewRectangle(Vec[int]{X: 2, Y: 5}, Vec[int]{X: 6, Y: 6}),
		},
		{
			rect1: NewRectangle(Vec[int]{X: 3, Y: 3}, Vec[int]{X: 5, Y: 5}),
			rect2: NewRectangle(Vec[int]{X: 5, Y: 2}, Vec[int]{X: 6, Y: 6}),
		},
		{
			rect1: NewRectangle(Vec[int]{X: 3, Y: 3}, Vec[int]{X: 5, Y: 5}),
			rect2: NewRectangle(Vec[int]{X: 5, Y: 2}, Vec[int]{X: 6, Y: 6}),
		},
	}

	for _, intersection := range intersects {
		if !intersection.rect1.Intersects(intersection.rect2) {
			t.Errorf("rect1 %v should intersect with rect2 %v", intersection.rect1, intersection.rect2)
		}
		if !intersection.rect2.Intersects(intersection.rect1) {
			t.Errorf("rect2 %v should intersect with rect1 %v", intersection.rect2, intersection.rect1)
		}
	}

	notIntersects := []struct{ rect1, rect2 Rectangle[int] }{
		{
			rect1: NewRectangle(Vec[int]{X: 3, Y: 3}, Vec[int]{X: 5, Y: 5}),
			rect2: NewRectangle(Vec[int]{X: 1, Y: 1}, Vec[int]{X: 2, Y: 2}),
		},
		{
			rect1: NewRectangle(Vec[int]{X: 3, Y: 3}, Vec[int]{X: 5, Y: 5}),
			rect2: NewRectangle(Vec[int]{X: 6, Y: 0}, Vec[int]{X: 9, Y: 9}),
		},
		{
			rect1: NewRectangle(Vec[int]{X: 3, Y: 3}, Vec[int]{X: 5, Y: 5}),
			rect2: NewRectangle(Vec[int]{X: 0, Y: 6}, Vec[int]{X: 9, Y: 9}),
		},
	}
	for _, intersection := range notIntersects {
		if intersection.rect1.Intersects(intersection.rect2) {
			t.Errorf("rect1 %v should not intersect with rect2 %v", intersection.rect1, intersection.rect2)
		}
		if intersection.rect2.Intersects(intersection.rect1) {
			t.Errorf("rect2 %v should not intersect with rect1 %v", intersection.rect2, intersection.rect1)
		}
	}
}

func TestRectangle_IntersectsCyclic(t *testing.T) {
	intersects := []struct{ rect1, rect2 Rectangle[int] }{
		{
			rect1: NewRectangle(Vec[int]{X: 5, Y: 5}, Vec[int]{X: 15, Y: 15}),
			rect2: NewRectangle(Vec[int]{X: 95, Y: 95}, Vec[int]{X: 105, Y: 105}),
		},
	}
	size := Vec[int]{X: 100, Y: 100}
	plane := NewBoundedPlane(size.X, size.Y)

	for _, intersection := range intersects {
		wrapped := WrapRectangleCyclic(intersection.rect2, size, plane.Contains)
		if !intersection.rect1.IntersectsAny(wrapped) {
			t.Errorf("rect1 %v should intersect with rect2 %v", intersection.rect1, intersection.rect2)
		}
	}
}

func TestRectangle_IntersectsAny_ReturnsFalse(t *testing.T) {
	base := Rectangle[int]{
		TopLeft:     Vec[int]{X: 0, Y: 0},
		BottomRight: Vec[int]{X: 10, Y: 10},
		Center:      Vec[int]{X: 5, Y: 5},
	}

	others := []Rectangle[int]{
		{TopLeft: Vec[int]{X: 20, Y: 20}, BottomRight: Vec[int]{X: 30, Y: 30}},
		{TopLeft: Vec[int]{X: 40, Y: 0}, BottomRight: Vec[int]{X: 50, Y: 10}},
		{TopLeft: Vec[int]{X: 0, Y: 40}, BottomRight: Vec[int]{X: 10, Y: 50}},
	}

	if base.IntersectsAny(others) {
		t.Errorf("expected IntersectsAny to return false, but got true")
	}
}

func TestRectangle_Probe_CyclicWrap(t *testing.T) {
	rect := NewRectangle(Vec[int]{X: 8, Y: 8}, Vec[int]{X: 10, Y: 10})
	plane := NewCyclicBoundedPlane(10, 10)

	probes := rect.Probe(0, plane)
	if len(probes) < 2 {
		t.Fatalf("expected wrapped probes, got %d", len(probes))
	}

	wrappedFound := false
	for _, p := range probes {
		if !rectEquals(p, rect) {
			wrappedFound = true
		}
	}
	if !wrappedFound {
		t.Errorf("expected wrapped rectangle in probes %v", probes)
	}
}

func TestRectangle_DistanceTo_Point(t *testing.T) {
	rect := NewRectangle(Vec[float64]{X: 0, Y: 0}, Vec[float64]{X: 2, Y: 2})
	point := Vec[float64]{X: 5, Y: 6}
	plane := NewBoundedPlane(100.0, 100.0)

	distance := rect.DistanceTo(&point, plane.Metric)
	expected := plane.Metric(Vec[float64]{X: 3, Y: 4}, ZERO_FLOAT_VEC)
	if distance != expected {
		t.Errorf("expected distance %v, got %v", expected, distance)
	}
}

func TestRectangle_Contains(t *testing.T) {
	outer := NewRectangle(Vec[int]{X: 0, Y: 0}, Vec[int]{X: 10, Y: 10})
	inner := NewRectangle(Vec[int]{X: 2, Y: 2}, Vec[int]{X: 8, Y: 8})
	onlyTopLeftInside := NewRectangle(Vec[int]{X: -1, Y: -1}, Vec[int]{X: 5, Y: 5})
	onlyBottomRightInside := NewRectangle(Vec[int]{X: 5, Y: 5}, Vec[int]{X: 12, Y: 12})

	if !outer.Contains(inner) {
		t.Errorf("expected outer to contain inner")
	}
	if outer.Contains(onlyTopLeftInside) {
		t.Errorf("expected outer not to contain rectangle with outside top-left")
	}
	if outer.Contains(onlyBottomRightInside) {
		t.Errorf("expected outer not to contain rectangle with outside bottom-right")
	}
}

func TestRectangle_Expand(t *testing.T) {
	rect := NewRectangle(Vec[int]{X: 2, Y: 3}, Vec[int]{X: 5, Y: 7})
	expanded := rect.Expand(2)
	expectedTopLeft := Vec[int]{X: 0, Y: 1}
	expectedBottomRight := Vec[int]{X: 7, Y: 9}
	if expanded.TopLeft != expectedTopLeft || expanded.BottomRight != expectedBottomRight {
		t.Errorf("unexpected expanded rectangle: %+v", expanded)
	}
}

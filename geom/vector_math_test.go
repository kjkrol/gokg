package geom

import (
	"math"
	"testing"
)

func TestLength(t *testing.T) {
	if got := Length(NewVec(3, 4)); got != 5 {
		t.Errorf("Length = %v, want 5", got)
	}
	if got := Length(NewVec(0, 0)); got != 0 {
		t.Errorf("Length of a zero vector = %v, want 0", got)
	}
	// A continuous world is the point of float64: lengths need not be whole.
	if got := Length(NewVec(1, 1)); math.Abs(got-math.Sqrt2) > 1e-12 {
		t.Errorf("Length = %v, want sqrt(2)", got)
	}
}

func TestClamp(t *testing.T) {
	bounds := NewVec(4, 6)

	cases := map[string]struct{ start, want Vec }{
		"inside is left alone": {NewVec(2, 3), NewVec(2, 3)},
		"past the far edge":    {NewVec(5, 7), NewVec(4, 6)},
		"below zero":           {NewVec(-3, -1), NewVec(0, 0)},
		"exactly on the edge":  {NewVec(4, 6), NewVec(4, 6)},
		"fractional stays put": {NewVec(1.5, 2.25), NewVec(1.5, 2.25)},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if got := Clamp(c.start, bounds); got != c.want {
				t.Errorf("Clamp(%v) = %v, want %v", c.start, got, c.want)
			}
		})
	}
}

func TestWrap(t *testing.T) {
	size := NewVec(4, 4)

	cases := map[string]struct{ start, want Vec }{
		"inside is left alone":     {NewVec(1, 3), NewVec(1, 3)},
		"one step past the edge":   {NewVec(5, 7), NewVec(1, 3)},
		"exactly on the far edge":  {NewVec(4, 4), NewVec(0, 0)},
		"just below zero":          {NewVec(-3, -1), NewVec(1, 3)},
		"far out, past the fast p": {NewVec(19, -13), NewVec(3, 3)},
		"fractional wraps too":     {NewVec(4.5, -0.5), NewVec(0.5, 3.5)},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got := Wrap(c.start, size)
			if math.Abs(got.X-c.want.X) > 1e-12 || math.Abs(got.Y-c.want.Y) > 1e-12 {
				t.Errorf("Wrap(%v) = %v, want %v", c.start, got, c.want)
			}
		})
	}
}

// A size of zero has nothing to wrap into, so the value passes through rather
// than dividing by zero.
func TestWrap_ZeroSizePassesThrough(t *testing.T) {
	if got := Wrap(NewVec(7, -7), NewVec(0, 0)); got != NewVec(7, -7) {
		t.Errorf("Wrap with a zero size = %v, want the input unchanged", got)
	}
}

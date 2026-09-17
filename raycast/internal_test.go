package raycast

import (
	"math"
	"testing"

	"github.com/kjkrol/gokg/geom"
)

func TestWrapToPi(t *testing.T) {
	cases := []struct{ in, want float64 }{
		{0, 0},
		{math.Pi / 2, math.Pi / 2},
		{math.Pi, math.Pi},
		{-math.Pi, math.Pi},
		{3 * math.Pi / 2, -math.Pi / 2},
		{-3 * math.Pi / 2, math.Pi / 2},
		{5 * math.Pi, math.Pi},
		{-5 * math.Pi, math.Pi},
	}
	for _, tc := range cases {
		if got := wrapToPi(tc.in); math.Abs(got-tc.want) > 1e-9 {
			t.Errorf("wrapToPi(%v) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestShortestDelta(t *testing.T) {
	cases := []struct{ d, size, want float64 }{
		{10, 100, 10},  // already the short way
		{60, 100, -40}, // shorter going backwards
		{-60, 100, 40}, // shorter going forwards
		{50, 100, 50},  // exactly half stays put
		{10, 0, 10},    // a zero-sized axis does not wrap
	}
	for _, tc := range cases {
		if got := shortestDelta(tc.d, tc.size); got != tc.want {
			t.Errorf("shortestDelta(%v, %v) = %v, want %v", tc.d, tc.size, got, tc.want)
		}
	}
}

func TestSplitAxis(t *testing.T) {
	cases := map[string]struct {
		lo, hi, size float64
		want         [][2]float64
	}{
		"inside":               {10, 40, 100, [][2]float64{{10, 40}}},
		"crosses the end":      {80, 120, 100, [][2]float64{{80, 100}, {0, 20}}},
		"starts negative":      {-20, 10, 100, [][2]float64{{80, 100}, {0, 10}}},
		"wider than the world": {-50, 200, 100, [][2]float64{{0, 100}}},
		"zero-sized axis":      {0, 10, 0, [][2]float64{{0, 0}}},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := splitAxis(tc.lo, tc.hi, tc.size)
			if len(got) != len(tc.want) {
				t.Fatalf("splitAxis(%v,%v,%v) = %v, want %v", tc.lo, tc.hi, tc.size, got, tc.want)
			}
			for i := range tc.want {
				if math.Abs(got[i][0]-tc.want[i][0]) > 1e-9 || math.Abs(got[i][1]-tc.want[i][1]) > 1e-9 {
					t.Errorf("splitAxis(%v,%v,%v) = %v, want %v", tc.lo, tc.hi, tc.size, got, tc.want)
				}
			}
		})
	}
}

func TestHitDistance(t *testing.T) {
	box := geom.NewAABB(geom.NewVec(10.0, 10.0), geom.NewVec(20.0, 20.0))
	origin := geom.NewVec(0.0, 15.0)

	cases := map[string]struct {
		origin, dir geom.Vec[float64]
		want        float64
		hit         bool
	}{
		"straight at it, parallel to Y":   {origin, geom.NewVec(1.0, 0.0), 10, true},
		"parallel but above the slab":     {geom.NewVec(0.0, 5.0), geom.NewVec(1.0, 0.0), 0, false},
		"parallel but below the slab":     {geom.NewVec(0.0, 25.0), geom.NewVec(1.0, 0.0), 0, false},
		"parallel to X, straight down":    {geom.NewVec(15.0, 0.0), geom.NewVec(0.0, 1.0), 10, true},
		"pointing away":                   {origin, geom.NewVec(-1.0, 0.0), 0, false},
		"diagonal miss":                   {geom.NewVec(0.0, 0.0), geom.NewVec(0.0, 1.0), 0, false},
		"starting inside hits at zero":    {geom.NewVec(15.0, 15.0), geom.NewVec(1.0, 0.0), 0, true},
		"diagonal hit through the corner": {geom.NewVec(0.0, 0.0), geom.NewVec(0.7071067811865476, 0.7071067811865476), 14.142135623730951, true},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, ok := hitDistance(tc.origin, tc.dir, box)
			if ok != tc.hit {
				t.Fatalf("hitDistance hit = %v, want %v", ok, tc.hit)
			}
			if ok && math.Abs(got-tc.want) > 1e-9 {
				t.Errorf("hitDistance = %v, want %v", got, tc.want)
			}
		})
	}
}

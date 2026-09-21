package raycast_test

import (
	"math"
	"math/rand/v2"
	"testing"

	"github.com/kjkrol/aabbworld/geom"
	"github.com/kjkrol/aabbworld/internal/raycast"
	"github.com/kjkrol/uid"
)

// nearestAlong is the reach of one ray, tested against every box directly.
func nearestAlong(boxes [][4]float64, origin geom.Vec, angle, radius float64) float64 {
	dx, dy := math.Cos(angle), math.Sin(angle)
	best := radius
	for _, b := range boxes {
		near, far := math.Inf(-1), math.Inf(1)
		for axis, ray := range [2][2]float64{{origin.X, dx}, {origin.Y, dy}} {
			o, d := ray[0], ray[1]
			lo, hi := b[axis], b[axis+2]
			if d == 0 {
				if o < lo || o > hi {
					near, far = 1, -1
				}
				continue
			}
			t1, t2 := (lo-o)/d, (hi-o)/d
			if t1 > t2 {
				t1, t2 = t2, t1
			}
			near, far = math.Max(near, t1), math.Min(far, t2)
		}
		if far >= near && far >= 0 {
			if h := math.Max(near, 0); h < best {
				best = h
			}
		}
	}
	return best
}

func TestDepths_MatchRaysCastDirectly(t *testing.T) {
	r := rand.New(rand.NewPCG(3, 5))
	for trial := range 10 {
		s := newFake(4000, 4000, false)
		s.put(eye, 2000, 2000, 10, 10)

		var boxes [][4]float64
		id := uid.UID64(100)
		for cx := range 6 {
			for cy := range 6 {
				if cx == 3 && cy == 3 {
					continue
				}
				x := float64(1700 + cx*100 + r.IntN(20))
				y := float64(1700 + cy*100 + r.IntN(20))
				w, h := float64(20+r.IntN(40)), float64(20+r.IntN(40))
				s.put(id, x, y, w, h)
				boxes = append(boxes, [4]float64{float64(x), float64(y), float64(x + w), float64(y + h)})
				id++
			}
		}

		coneDir := float64(trial)
		cone := raycast.Cone{
			Direction: geom.NewVec(math.Cos(coneDir), math.Sin(coneDir)),
			HalfAngle: math.Pi / 4,
			Radius:    500,
		}
		origin := geom.NewVec(2005.0, 2005.0)

		v := &raycast.View{}
		if !v.Scan(s, eye, cone) {
			t.Fatalf("trial %d: Scan failed", trial)
		}

		const k = 33
		got := v.Depths(k, nil)
		if len(got) != k {
			t.Fatalf("trial %d: got %d depths, want %d", trial, len(got), k)
		}
		step := 2 * cone.HalfAngle / float64(k-1)
		for i, d := range got {
			want := nearestAlong(boxes, origin, coneDir-cone.HalfAngle+float64(i)*step, cone.Radius)
			if math.Abs(float64(d)-want) > 1e-3 {
				t.Errorf("trial %d: depth %d = %.4f, casting that ray gives %.4f", trial, i, d, want)
			}
		}
	}
}

func TestDepths_TwoSamplesAreTheConeEdges(t *testing.T) {
	s := newFake(2000, 2000, false)
	s.put(eye, 1000, 1000, 10, 10)
	s.put(near, 1200, 1200, 30, 30)

	const half = math.Pi / 4
	v := &raycast.View{}
	if !v.Scan(s, eye, eastward(half, 600)) {
		t.Fatal("Scan failed")
	}

	got := v.Depths(2, nil)
	if len(got) != 2 {
		t.Fatalf("got %d depths, want 2", len(got))
	}
	origin := geom.NewVec(1005.0, 1005.0)
	boxes := [][4]float64{{1200, 1200, 1230, 1230}}
	for i, want := range []float64{
		nearestAlong(boxes, origin, -half, 600),
		nearestAlong(boxes, origin, half, 600),
	} {
		if math.Abs(float64(got[i])-want) > 1e-3 {
			t.Errorf("edge %d = %.4f, want %.4f", i, got[i], want)
		}
	}
}

func TestDepths_EmptySceneReachesFullRange(t *testing.T) {
	s := newFake(2000, 2000, false)
	s.put(eye, 1000, 1000, 10, 10)

	v := &raycast.View{}
	v.Scan(s, eye, eastward(math.Pi/4, 400))
	for i, d := range v.Depths(16, nil) {
		if math.Abs(float64(d)-400) > 1e-3 {
			t.Errorf("depth %d = %.3f with nothing in the way, want the full range 400", i, d)
		}
	}
}

func TestDepths_LeavesDstUntouchedWhenItCannotAnswer(t *testing.T) {
	s := newFake(2000, 2000, false)
	s.put(eye, 1000, 1000, 10, 10)
	dst := []float32{7, 8}

	cases := map[string]func(*raycast.View) []float32{
		"never scanned":   func(v *raycast.View) []float32 { return v.Depths(8, dst) },
		"failed scan":     func(v *raycast.View) []float32 { v.Scan(s, eye, eastward(0, 400)); return v.Depths(8, dst) },
		"k below two":     func(v *raycast.View) []float32 { v.Scan(s, eye, eastward(math.Pi/4, 400)); return v.Depths(1, dst) },
		"k below two too": func(v *raycast.View) []float32 { v.Scan(s, eye, eastward(math.Pi/4, 400)); return v.Depths(0, dst) },
	}
	for name, run := range cases {
		t.Run(name, func(t *testing.T) {
			if got := run(&raycast.View{}); len(got) != 2 || got[0] != 7 || got[1] != 8 {
				t.Errorf("Depths = %v, want dst unchanged", got)
			}
		})
	}
}

func TestDepths_DoesNotDisturbTheExactSweep(t *testing.T) {
	s := newFake(2000, 2000, false)
	s.put(eye, 1000, 1000, 10, 10)
	s.put(near, 1200, 960, 20, 40)
	s.put(far, 1400, 1100, 30, 30)
	cone := eastward(math.Pi/3, 700)

	clean := &raycast.View{}
	clean.Scan(s, eye, cone)
	var wantIDs []seenPair
	clean.Entities(func(id uid.UID64, d float64) { wantIDs = append(wantIDs, seenPair{id, d}) })
	wantPts := clean.Outline(0, nil)

	v := &raycast.View{}
	v.Scan(s, eye, cone)
	v.Depths(24, nil)

	var gotIDs []seenPair
	v.Entities(func(id uid.UID64, d float64) { gotIDs = append(gotIDs, seenPair{id, d}) })
	if len(gotIDs) != len(wantIDs) {
		t.Fatalf("Entities reports %d after Depths, want %d", len(gotIDs), len(wantIDs))
	}
	for i := range wantIDs {
		if gotIDs[i] != wantIDs[i] {
			t.Errorf("entity %d = %v after Depths, want %v", i, gotIDs[i], wantIDs[i])
		}
	}
	if gotPts := v.Outline(0, nil); len(gotPts) != len(wantPts) {
		t.Fatalf("Outline has %d points after Depths, want %d", len(gotPts), len(wantPts))
	}
}

func TestDepths_ReturnsExactlyKEvenBelowTheSweepsTolerance(t *testing.T) {
	s := newFake(2000, 2000, false)
	s.put(eye, 1000, 1000, 10, 10)

	v := &raycast.View{}
	if !v.Scan(s, eye, eastward(math.Pi/4, 10)) {
		t.Fatal("Scan failed")
	}
	const k = 200
	if got := v.Depths(k, nil); len(got) != k {
		t.Errorf("got %d depths, want exactly %d", len(got), k)
	}
}

package raycast_test

import (
	"math"
	"testing"

	"github.com/kjkrol/gokg/geom"
	"github.com/kjkrol/uid"
)

// reach is how far the outline point at index i sits from the observer.
func reach(pts []geom.Vec[float64], i int) float64 {
	return math.Hypot(pts[i].X-pts[0].X, pts[i].Y-pts[0].Y)
}

func TestOutline_EmptySceneIsAnArcAtFullRange(t *testing.T) {
	s := newFake(2000, 2000, false)
	s.put(eye, 1000, 1000, 10, 10)

	pts := outline(t, s, eye, eastward(math.Pi/4, 300), 0)
	if len(pts) < 10 {
		t.Fatalf("outline has %d points, expected the arc to be subdivided", len(pts))
	}
	for i := 1; i < len(pts); i++ {
		if d := reach(pts, i); math.Abs(d-300) > 1e-6 {
			t.Fatalf("point %d sits at %.3f, want the full range 300", i, d)
		}
	}
}

func TestOutline_StartsAtTheObserver(t *testing.T) {
	s := newFake(2000, 2000, false)
	s.put(eye, 1000, 1000, 10, 10)

	pts := outline(t, s, eye, eastward(math.Pi/4, 300), 0)
	if pts[0].X != 1005 || pts[0].Y != 1005 {
		t.Errorf("outline starts at %v, want the observer's centre (1005,1005)", pts[0])
	}
}

func TestOutline_IsPulledInBehindABlocker(t *testing.T) {
	s := newFake(2000, 2000, false)
	s.put(eye, 1000, 1000, 10, 10)
	s.put(near, 1200, 990, 10, 30) // straight ahead, spanning the axis

	pts := outline(t, s, eye, eastward(math.Pi/4, 600), 0)

	var shortest, longest float64 = math.Inf(1), 0
	for i := 1; i < len(pts); i++ {
		d := reach(pts, i)
		shortest = math.Min(shortest, d)
		longest = math.Max(longest, d)
	}
	// Some directions reach the blocker (~195), others run to the full range.
	if shortest > 250 {
		t.Errorf("shortest reach %.0f — the outline was never pulled in to the blocker", shortest)
	}
	if longest < 599 {
		t.Errorf("longest reach %.0f — the outline never reaches past the blocker", longest)
	}
}

func TestOutline_StaysWithinTheCone(t *testing.T) {
	s := newFake(2000, 2000, false)
	s.put(eye, 1000, 1000, 10, 10)
	s.put(near, 1200, 1200, 40, 40)

	const half = math.Pi / 6
	pts := outline(t, s, eye, eastward(half, 500), 0)
	for i := 1; i < len(pts); i++ {
		a := math.Atan2(pts[i].Y-pts[0].Y, pts[i].X-pts[0].X)
		if math.Abs(a) > half+1e-9 {
			t.Errorf("point %d lies at %.4f rad, outside the cone's ±%.4f", i, a, half)
		}
	}
}

// Whatever Visible reports must also show up on the outline, otherwise the
// drawn shape and the returned set disagree.
func TestOutline_AgreesWithVisible(t *testing.T) {
	s := newFake(2000, 2000, false)
	s.put(eye, 1000, 1000, 10, 10)
	s.put(near, 1200, 960, 20, 40)
	s.put(far, 1400, 1100, 30, 30)
	s.put(tail, 1600, 1000, 20, 20)

	cone := eastward(math.Pi/3, 900)

	var dists []float64
	scan(t, s, eye, cone).Entities(func(_ uid.UID64, d float64) { dists = append(dists, d) })
	if len(dists) == 0 {
		t.Fatal("nothing visible; the fixture proves nothing")
	}

	pts := outline(t, s, eye, cone, 0)
	for _, want := range dists {
		found := false
		for i := 1; i < len(pts); i++ {
			if math.Abs(reach(pts, i)-want) < 1e-6 {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Visible reports an entity at %.3f, but no outline point sits there", want)
		}
	}
}

func TestOutline_NilForADegenerateQuery(t *testing.T) {
	s := newFake(2000, 2000, false)
	s.put(eye, 1000, 1000, 10, 10)

	if pts := outline(t, s, uid.UID64(999), eastward(math.Pi/4, 300), 0); pts != nil {
		t.Errorf("outline = %v for an unknown observer, want nil", pts)
	}
	if pts := outline(t, s, eye, eastward(0, 300), 0); pts != nil {
		t.Errorf("outline = %v for a zero-width cone, want nil", pts)
	}
}

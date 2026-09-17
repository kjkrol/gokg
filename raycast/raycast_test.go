package raycast_test

import (
	"math"
	"testing"

	"github.com/kjkrol/gokg/geom"
	"github.com/kjkrol/gokg/plane"
	"github.com/kjkrol/gokg/raycast"
	"github.com/kjkrol/uid"
)

// fakeSpace is a QueryableSpace over a plain slice — enough to pin the
// visibility rules without building a real index.
type fakeSpace struct {
	boxes    map[uid.UID64]geom.AABB[uint32]
	w, h     uint32
	toroidal bool
	// twice makes Query report every entity again, as it does for a box
	// straddling the seam; phantom is reported but has no indexed box.
	twice   bool
	phantom uid.UID64
}

func newFake(w, h uint32, toroidal bool) *fakeSpace {
	return &fakeSpace{boxes: map[uid.UID64]geom.AABB[uint32]{}, w: w, h: h, toroidal: toroidal}
}

func (f *fakeSpace) put(id uid.UID64, x, y, w, h uint32) {
	f.boxes[id] = geom.NewAABBAt(geom.NewVec(x, y), w, h)
}

func (f *fakeSpace) Bounds() (uint32, uint32, bool) { return f.w, f.h, f.toroidal }

func (f *fakeSpace) EntryAABB(id uid.UID64) (geom.AABB[uint32], bool) {
	b, ok := f.boxes[id]
	return b, ok
}

func (f *fakeSpace) Query(area geom.AABB[uint32], fn func(uid.UID64, plane.FragPosition)) int {
	n := 0
	for id, b := range f.boxes {
		if !area.Intersects(b) {
			continue
		}
		fn(id, plane.FRAG_MAIN)
		n++
		if f.twice {
			fn(id, plane.FRAG_RIGHT)
			n++
		}
	}
	if f.phantom != 0 {
		fn(f.phantom, plane.FRAG_MAIN)
		n++
	}
	return n
}

// eastward looks along +X with the given half-angle and reach.
func eastward(halfAngle, radius float64) raycast.Cone {
	return raycast.Cone{Direction: geom.NewVec(1.0, 0.0), HalfAngle: halfAngle, Radius: radius}
}

// scan answers through a fresh View, which is how a test looks at one query.
func scan(t *testing.T, s raycast.QueryableSpace, observer uid.UID64, c raycast.Cone) *raycast.View {
	t.Helper()
	v := &raycast.View{}
	v.Scan(s, observer, c)
	return v
}

func visible(t *testing.T, s raycast.QueryableSpace, observer uid.UID64, c raycast.Cone) map[uid.UID64]bool {
	t.Helper()
	got := map[uid.UID64]bool{}
	n := scan(t, s, observer, c).Entities(func(id uid.UID64, _ float64) { got[id] = true })
	if n != len(got) {
		t.Errorf("Entities returned %d but reported %d distinct ids", n, len(got))
	}
	return got
}

func outline(t *testing.T, s raycast.QueryableSpace, observer uid.UID64, c raycast.Cone, step float64) []geom.Vec[float64] {
	t.Helper()
	return scan(t, s, observer, c).Outline(step, nil)
}

func assertVisible(t *testing.T, got map[uid.UID64]bool, want ...uid.UID64) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("visible = %v, want exactly %v", keys(got), want)
	}
	for _, id := range want {
		if !got[id] {
			t.Errorf("entity %v missing from %v", id, keys(got))
		}
	}
}

func keys(m map[uid.UID64]bool) []uid.UID64 {
	out := make([]uid.UID64, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

const (
	eye  = uid.UID64(1)
	near = uid.UID64(2)
	far  = uid.UID64(3)
	tail = uid.UID64(4)
)

func TestVisible_NearestOfAColumnHidesTheRest(t *testing.T) {
	s := newFake(1000, 1000, false)
	s.put(eye, 100, 100, 10, 10)
	s.put(near, 200, 100, 10, 10)
	s.put(far, 300, 100, 10, 10)
	s.put(tail, 400, 100, 10, 10)

	assertVisible(t, visible(t, s, eye, eastward(math.Pi/4, 500)), near)
}

func TestVisible_TargetPeekingPastTheBlockerIsSeen(t *testing.T) {
	s := newFake(1000, 1000, false)
	s.put(eye, 100, 100, 10, 10)
	s.put(near, 200, 100, 10, 10)
	// Wide enough that its angular span reaches past the blocker's.
	s.put(far, 300, 60, 10, 100)

	assertVisible(t, visible(t, s, eye, eastward(math.Pi/4, 500)), near, far)
}

func TestVisible_SideBySideEntitiesDoNotHideEachOther(t *testing.T) {
	s := newFake(1000, 1000, false)
	s.put(eye, 100, 100, 10, 10)
	s.put(near, 300, 60, 10, 10)
	s.put(far, 300, 140, 10, 10)

	assertVisible(t, visible(t, s, eye, eastward(math.Pi/3, 500)), near, far)
}

func TestVisible_ExcludesOutOfRangeOutOfConeAndSelf(t *testing.T) {
	// Every rejected entity sits inside the search square, so each subtest
	// isolates one reason rather than leaning on the broad phase.
	cases := map[string]struct {
		x, y uint32
		want []uid.UID64
	}{
		"in cone and in range":           {700, 500, []uid.UID64{near}},
		"inside range, outside the cone": {500, 300, nil},
		// Inside the 2R-square but past R: the corner region the square covers
		// and the circle does not.
		"inside the cone, beyond the radius": {796, 608, nil},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			s := newFake(2000, 2000, false)
			s.put(eye, 500, 500, 10, 10)
			s.put(near, tc.x, tc.y, 10, 10)

			got := visible(t, s, eye, eastward(math.Pi/8, 300))
			assertVisible(t, got, tc.want...)
		})
	}
}

func TestScan_RefusesAnUnknownObserver(t *testing.T) {
	s := newFake(1000, 1000, false)
	s.put(near, 200, 100, 10, 10)

	if (&raycast.View{}).Scan(s, eye, eastward(math.Pi/4, 500)) {
		t.Error("Scan succeeded for an observer that is not indexed")
	}
}

func TestScan_RefusesDegenerateCones(t *testing.T) {
	s := newFake(1000, 1000, false)
	s.put(eye, 100, 100, 10, 10)
	s.put(near, 200, 100, 10, 10)

	cases := map[string]raycast.Cone{
		"zero half-angle": eastward(0, 500),
		"half-angle pi":   eastward(math.Pi, 500),
		"zero radius":     eastward(math.Pi/4, 0),
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if (&raycast.View{}).Scan(s, eye, c) {
				t.Error("Scan succeeded, want a refusal")
			}
		})
	}
}

func TestVisible_ReportsNearestFirst(t *testing.T) {
	s := newFake(1000, 1000, false)
	s.put(eye, 100, 100, 10, 10)
	s.put(near, 300, 60, 10, 10)
	s.put(far, 500, 140, 10, 10)

	var order []float64
	scan(t, s, eye, eastward(math.Pi/3, 900)).Entities(func(_ uid.UID64, d float64) {
		order = append(order, d)
	})
	if len(order) != 2 {
		t.Fatalf("saw %d entities, want 2", len(order))
	}
	if order[0] > order[1] {
		t.Errorf("distances %v are not ascending", order)
	}
}

func TestVisible_SeesAcrossTheSeamInAToroidalWorld(t *testing.T) {
	s := newFake(1000, 1000, true)
	s.put(eye, 970, 500, 10, 10)
	s.put(near, 20, 500, 10, 10) // just past the seam, 40 units away

	got := visible(t, s, eye, eastward(math.Pi/4, 200))
	assertVisible(t, got, near)
}

// An entity straddling the seam is indexed once per fragment, so Query reports
// it several times; it must still count as one.
func TestVisible_ReportsAFragmentedEntityOnce(t *testing.T) {
	s := newFake(1000, 1000, false)
	s.twice = true
	s.put(eye, 100, 100, 10, 10)
	s.put(near, 300, 100, 10, 10)

	assertVisible(t, visible(t, s, eye, eastward(math.Pi/4, 500)), near)
}

// An id the index reports but can no longer resolve is skipped rather than
// crashing the sweep.
func TestVisible_SkipsCandidatesWithoutAnIndexedBox(t *testing.T) {
	s := newFake(1000, 1000, false)
	s.phantom = uid.UID64(99)
	s.put(eye, 100, 100, 10, 10)
	s.put(near, 300, 100, 10, 10)

	assertVisible(t, visible(t, s, eye, eastward(math.Pi/4, 500)), near)
}

// The cone's bounding square runs off the left edge, so the wrapped search area
// must still cover what lies just across the seam.
func TestVisible_SeesAcrossTheSeamWhenTheObserverSitsAtTheOrigin(t *testing.T) {
	s := newFake(1000, 1000, true)
	s.put(eye, 20, 500, 10, 10)
	s.put(near, 960, 500, 10, 10) // behind the observer, just across the seam

	cone := raycast.Cone{Direction: geom.NewVec(-1.0, 0.0), HalfAngle: math.Pi / 4, Radius: 200}
	assertVisible(t, visible(t, s, eye, cone), near)
}

// Sweeping per angle resolves what ordering by nearest distance could not.
//
// The wall's nearest corner is close, but its far end recedes: along the angles
// where the target sits, the wall's surface is ~682 away while the target is at
// ~500. Ordering boxes by their nearest corner used to put the wall first and
// let it swallow the target across its whole angular span; measuring depth at
// each angle sees the target in front, where it belongs.
func TestVisible_SeesATargetInFrontOfAReceedingWall(t *testing.T) {
	s := newFake(2000, 2000, false)
	s.put(eye, 0, 0, 10, 10)     // centre (5,5)
	s.put(near, 20, 20, 980, 10) // wall spanning x∈[20,1000], y∈[20,30]
	s.put(far, 500, 14, 10, 4)   // in front of the wall's far end

	assertVisible(t, visible(t, s, eye, eastward(math.Pi/3, 1200)), near, far)
}

// Two boxes lying on the same pair of sight lines subtend exactly the same
// angles, so the sweep meets each critical angle twice and must not cast twice.
func TestVisible_BoxesOnTheSameSightLines(t *testing.T) {
	s := newFake(2000, 2000, false)
	s.put(eye, 100, 100, 10, 10) // centre (105,105)
	s.put(near, 200, 100, 10, 10)
	// Corners placed on the rays through near's corners, at twice the distance.
	s.put(far, 295, 95, 20, 20)

	assertVisible(t, visible(t, s, eye, eastward(math.Pi/4, 900)), near)
}

// A cone wider than a half-turn is no longer the intersection of two
// half-planes, so the cheap corner rejection stands down and the angular span
// alone has to keep what lies behind the observer out.
func TestVisible_WideConeStillExcludesWhatIsBehind(t *testing.T) {
	s := newFake(2000, 2000, false)
	s.put(eye, 1000, 1000, 10, 10)
	s.put(near, 1300, 995, 20, 20) // ahead
	s.put(tail, 700, 995, 20, 20)  // behind

	var got []uid.UID64
	scan(t, s, eye, eastward(2*math.Pi/3, 600)).Entities(func(id uid.UID64, _ float64) {
		got = append(got, id)
	})

	if len(got) != 1 || got[0] != near {
		t.Errorf("visible = %v, want only the entity ahead (%d)", got, near)
	}
}

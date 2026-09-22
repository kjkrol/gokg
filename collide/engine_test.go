package collide_test

import (
	"slices"
	"testing"

	"github.com/kjkrol/aabbworld"
	"github.com/kjkrol/aabbworld/collide"
	"github.com/kjkrol/aabbworld/geom"
	"github.com/kjkrol/aabbworld/plane"
	"github.com/kjkrol/uid"
)

func torusIf(toroidal bool) aabbworld.Edges {
	if toroidal {
		return aabbworld.Torus
	}
	return 0
}

func engineSpace(t *testing.T, toroidal bool) *aabbworld.Space {
	t.Helper()
	space, err := aabbworld.NewSpace(aabbworld.Config{
		Width: 1000, Height: 1000, Edges: torusIf(toroidal), BucketSize: 64, BucketCapacity: 8,
	})
	if err != nil {
		t.Fatalf("NewSpace: %v", err)
	}
	return space
}

// scene is a few boxes kept beside the space they are inserted into, resolved for the Engine by id.
type scene struct {
	space  *aabbworld.Space
	boxes  map[uid.UID64]*plane.AABB
	static map[uid.UID64]bool
	sensor map[uid.UID64]bool
	asked  []idPair
}

func newScene(t *testing.T, toroidal bool) *scene {
	return &scene{space: engineSpace(t, toroidal), boxes: map[uid.UID64]*plane.AABB{}, static: map[uid.UID64]bool{}, sensor: map[uid.UID64]bool{}}
}

func (s *scene) put(id uid.UID64, x, y float64) *plane.AABB {
	box := plane.NewAABB(geom.NewVec(x, y), 10, 10)
	s.space.Insert(id, &box)
	s.boxes[id] = &box
	return &box
}

func (s *scene) resolve(a, b uid.UID64) (collide.Body, collide.Body, bool) {
	s.asked = append(s.asked, idPair{a, b})
	return collide.Body{Box: s.boxes[a], Static: s.static[a], Sensor: s.sensor[a]},
		collide.Body{Box: s.boxes[b], Static: s.static[b], Sensor: s.sensor[b]}, true
}

func (s *scene) tick(e *collide.Engine, touch collide.Touch, onContact func(i int, pen geom.Vec)) {
	s.tickFor(e, 16, touch, onContact)
}

func (s *scene) tickFor(e *collide.Engine, iterations int, touch collide.Touch, onContact func(i int, pen geom.Vec)) {
	s.space.Flush(nil)
	s.asked = s.asked[:0]
	e.Tick(s.space, 0.5, aabbworld.AnyCapability, iterations, s.resolve, touch, onContact)
	s.space.Flush(nil)
}

// where lists every id the index finds inside box, each once.
func (s *scene) where(box geom.AABB) []uid.UID64 {
	var ids []uid.UID64
	s.space.Query(box, aabbworld.AnyCapability, func(id uid.UID64) {
		if !slices.Contains(ids, id) {
			ids = append(ids, id)
		}
	})
	slices.Sort(ids)
	return ids
}

func TestTick_SeparatesWhatOverlaps_AndTheIndexFollows(t *testing.T) {
	s := newScene(t, false)
	a, b := s.put(1, 100, 100), s.put(2, 104, 100)
	s.put(3, 500, 500)

	var e collide.Engine
	var contacts []int
	s.tick(&e, nil, func(i int, pen geom.Vec) {
		contacts = append(contacts, i)
		if pen == (geom.Vec{}) {
			t.Error("a contact was reported with no penetration")
		}
	})

	if len(s.asked) != 1 || s.asked[0] != (idPair{1, 2}) {
		t.Fatalf("resolve asked about %v, want only the pair that can touch", s.asked)
	}
	if len(contacts) != 1 || contacts[0] != 0 {
		t.Errorf("contacts = %v, want the one overlapping pair, once", contacts)
	}
	if a.Penetration(b.AABB) != (geom.Vec{}) {
		t.Errorf("a and b still overlap after Tick: %v and %v", a, b)
	}
	if a.TopLeft.X >= 100 || b.TopLeft.X <= 104 {
		t.Errorf("a at x=%v and b at x=%v, want both pushed apart", a.TopLeft.X, b.TopLeft.X)
	}
	if got := len(e.Left()); got != 0 {
		t.Errorf("Left names %d entities in a closed world, want none", got)
	}
	if got := s.where(geom.NewAABB(geom.NewVec(90, 100), geom.NewVec(99, 110))); !slices.Equal(got, []uid.UID64{1}) {
		t.Errorf("the index finds %v where only a came to rest, want [1] with no step taken by the caller", got)
	}

	s.tick(&e, nil, func(int, geom.Vec) { t.Error("a settled scene reported a contact") })
}

func TestTick_SeparatesAcrossTheSeamOfAToroidalSpace(t *testing.T) {
	s := newScene(t, true)
	left, right := s.put(1, 0, 100), s.put(2, 994, 100)

	var e collide.Engine
	s.tick(&e, nil, nil)

	if _, overlapping := left.DeepestOverlapWith(right); overlapping {
		t.Errorf("the pair still overlaps across the seam: %v and %v", left, right)
	}
	if travelled := left.TopLeft.X; travelled > 10 {
		t.Errorf("left ended up at x=%v — it went the long way round instead of through the seam", travelled)
	}
}

func TestTick_RefusedPairsAreNeitherTestedNorMoved(t *testing.T) {
	s := newScene(t, false)
	a, b := s.put(1, 100, 100), s.put(2, 104, 100)
	start := *a

	var e collide.Engine
	s.space.Flush(nil)
	e.Tick(s.space, 0.5, aabbworld.AnyCapability, 16, func(uid.UID64, uid.UID64) (collide.Body, collide.Body, bool) {
		return collide.Body{}, collide.Body{}, false
	}, nil, func(int, geom.Vec) { t.Error("a refused pair reported a contact") })

	if *a != start || b.TopLeft.X != 104 {
		t.Errorf("boxes moved to %v and %v although the pair was refused", a, b)
	}
}

func TestTick_TouchVetoesAPair(t *testing.T) {
	s := newScene(t, false)
	a, b := s.put(1, 100, 100), s.put(2, 104, 100)
	startA, startB := *a, *b

	var e collide.Engine
	asked := 0
	s.tick(&e, func(i int, pen geom.Vec) (geom.Vec, bool) {
		asked++
		return pen, false
	}, func(int, geom.Vec) { t.Error("a pair the shapes do not touch in reported a contact") })

	if asked != 1 {
		t.Errorf("touch asked %d times, want once for the one overlapping pair", asked)
	}
	if *a != startA || *b != startB {
		t.Errorf("boxes moved to %v and %v although touch said no", a, b)
	}
}

func TestTick_TouchMayRefineThePenetration(t *testing.T) {
	s := newScene(t, false)
	a, b := s.put(1, 100, 100), s.put(2, 104, 100)

	var e collide.Engine
	var reported geom.Vec
	s.tickFor(&e, 1, func(_ int, pen geom.Vec) (geom.Vec, bool) {
		return geom.NewVec(pen.X/3, 0), true
	}, func(_ int, pen geom.Vec) { reported = pen })

	if reported != geom.NewVec(-2, 0) {
		t.Errorf("onContact got %v, want the refined penetration (-2,0)", reported)
	}
	if a.TopLeft.X != 99 || b.TopLeft.X != 105 {
		t.Errorf("a at x=%v, b at x=%v, want one pass to move each by half of the refined 2", a.TopLeft.X, b.TopLeft.X)
	}
}

func TestTick_ASensorPairIsReportedButNeverSeparated(t *testing.T) {
	s := newScene(t, false)
	a, b := s.put(1, 100, 100), s.put(2, 104, 100)
	s.sensor[2] = true
	startA, startB := *a, *b

	var e collide.Engine
	reported := 0
	s.tick(&e, nil, func(int, geom.Vec) { reported++ })

	if reported != 1 {
		t.Errorf("reported %d contacts, want 1", reported)
	}
	if *a != startA || *b != startB {
		t.Errorf("a sensor pair was separated: %v and %v", a, b)
	}
}

func TestTick_AStaticSideNeverMoves(t *testing.T) {
	s := newScene(t, false)
	a, b := s.put(1, 100, 100), s.put(2, 104, 100)
	s.static[1] = true
	startA := *a

	var e collide.Engine
	s.tick(&e, nil, nil)

	if *a != startA {
		t.Errorf("the static side moved from %v to %v", startA.AABB, a.AABB)
	}
	if a.Penetration(b.AABB) != (geom.Vec{}) {
		t.Errorf("still overlapping: %v and %v", a, b)
	}
}

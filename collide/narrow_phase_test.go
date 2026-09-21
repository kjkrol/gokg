package collide_test

import (
	"github.com/kjkrol/aabbworld/collide"
	"github.com/kjkrol/uid"
	"testing"

	"github.com/kjkrol/aabbworld"
	"github.com/kjkrol/aabbworld/geom"
)

func torusIf(toroidal bool) aabbworld.Edges {
	if toroidal {
		return aabbworld.Torus
	}
	return 0
}

func solverSpace(t *testing.T, toroidal bool) *aabbworld.Space {
	t.Helper()
	space, err := aabbworld.NewSpace(aabbworld.Config{
		Width: 1000, Height: 1000, Edges: torusIf(toroidal), BucketSize: 64, BucketCapacity: 8,
	})
	if err != nil {
		t.Fatalf("NewSpace: %v", err)
	}
	return space
}

func TestSeparate_SeparatesWhatOverlaps_AndTheIndexFollows(t *testing.T) {
	space := solverSpace(t, false)
	idA, idB, idApart := uid.UID64(1), uid.UID64(2), uid.UID64(3)
	a := space.WrapAABB(geom.NewAABBAt(geom.NewVec(100, 100), 10, 10))
	b := space.WrapAABB(geom.NewAABBAt(geom.NewVec(104, 100), 10, 10))
	apart := space.WrapAABB(geom.NewAABBAt(geom.NewVec(500, 500), 10, 10))
	space.Insert(idA, &a)
	space.Insert(idB, &b)
	space.Insert(idApart, &apart)
	space.Flush(nil)

	var s collide.NarrowPhase
	touching := s.Add(collide.Pair{A: &a, B: &b, KeyA: 1, KeyB: 2})
	s.Add(collide.Pair{A: &a, B: &apart, KeyA: 1, KeyB: 3})
	whose := func(i int) (uid.UID64, uid.UID64) {
		if i == touching {
			return idA, idB
		}
		return idA, idApart
	}

	var contacts []int
	s.Separate(space, 16, func(i int, pen geom.Vec) {
		contacts = append(contacts, i)
		if pen == (geom.Vec{}) {
			t.Error("a contact was reported with no penetration")
		}
	}, whose)
	space.Flush(nil)

	if len(contacts) != 1 || contacts[0] != touching {
		t.Errorf("contacts = %v, want the one overlapping pair (%d), once", contacts, touching)
	}
	if a.Penetration(b.AABB) != (geom.Vec{}) {
		t.Errorf("a and b still overlap after Separate: %v and %v", a, b)
	}
	if a.TopLeft.X >= 100 || b.TopLeft.X <= 104 {
		t.Errorf("a at x=%v and b at x=%v, want both pushed apart", a.TopLeft.X, b.TopLeft.X)
	}
	if got := len(s.Left()); got != 0 {
		t.Errorf("Left names %d entities in a closed world, want none", got)
	}

	vacatedByNobody := geom.NewAABB(geom.NewVec(90, 100), geom.NewVec(99, 110))
	var there []uid.UID64
	space.Query(vacatedByNobody, aabbworld.AnyCapability, func(id uid.UID64) { there = append(there, id) })
	if len(there) != 1 || there[0] != idA {
		t.Errorf("the index finds %v where only a came to rest, want %v with no step taken by the caller", there, idA)
	}

	s.Reset()
	s.Separate(space, 16, func(int, geom.Vec) { t.Error("a reset NarrowPhase reported a contact") }, whose)
}

func TestSeparate_SeparatesAcrossTheSeamOfAToroidalSpace(t *testing.T) {
	space := solverSpace(t, true)
	left := space.WrapAABB(geom.NewAABBAt(geom.NewVec(0, 100), 10, 10))
	right := space.WrapAABB(geom.NewAABBAt(geom.NewVec(994, 100), 10, 10))
	space.Insert(uid.UID64(1), &left)
	space.Insert(uid.UID64(2), &right)
	space.Flush(nil)

	var s collide.NarrowPhase
	s.Add(collide.Pair{A: &left, B: &right})
	s.Separate(space, 16, nil, func(int) (uid.UID64, uid.UID64) { return uid.UID64(1), uid.UID64(2) })

	if _, overlapping := left.DeepestOverlapWith(&right); overlapping {
		t.Errorf("the pair still overlaps across the seam: %v and %v", left, right)
	}
	if travelled := left.TopLeft.X; travelled > 10 {
		t.Errorf("left ended up at x=%v — it went the long way round instead of through the seam", travelled)
	}
}

func TestSeparate_WithNobodyNamedLeavesTheIndexAlone(t *testing.T) {
	space := solverSpace(t, false)
	id := uid.UID64(1)
	a := space.WrapAABB(geom.NewAABBAt(geom.NewVec(100, 100), 10, 10))
	b := space.WrapAABB(geom.NewAABBAt(geom.NewVec(104, 100), 10, 10))
	space.Insert(id, &a)
	space.Flush(nil)

	var s collide.NarrowPhase
	s.Add(collide.Pair{A: &a, B: &b})
	s.Separate(space, 16, nil, nil)
	space.Flush(nil)

	if a.TopLeft.X >= 100 {
		t.Fatalf("a at x=%v, want it pushed", a.TopLeft.X)
	}
	stale := geom.NewAABB(geom.NewVec(108, 100), geom.NewVec(110, 110))
	if n := space.Query(stale, aabbworld.AnyCapability, func(uid.UID64) {}); n != 1 {
		t.Errorf("the index finds %d pieces where a used to reach, want it still there: nobody named it", n)
	}
}

package spatial

import (
	"testing"

	"github.com/kjkrol/gokg/plane"
	"github.com/kjkrol/uid"
)

// capGrid is a 256x256 toroidal index holding three entities side by side.
func capGrid(t *testing.T) (*GridIndexManager, []uid.UID64) {
	t.Helper()
	space := plane.NewToroidal2D(256, 256)
	m, err := NewGridIndexManager(space, GridIndexConfig{
		Resolution: Size256x256, BucketResolution: Size32x32, BucketCapacity: 8,
	})
	if err != nil {
		t.Fatalf("NewGridIndexManager: %v", err)
	}
	ids := []uid.UID64{1, 2, 3}
	for i, id := range ids {
		m.QueueInsert(id, space.WrapAABB(NewAABBAt(NewVec(float64(100+i*12), 100), 8, 8)))
	}
	m.Flush(nil)
	return m, ids
}

func found(m *GridIndexManager, box AABB, want Capability) map[uid.UID64]bool {
	got := map[uid.UID64]bool{}
	m.QueryRangeWith(box, want, func(id uid.UID64, _ plane.FragPosition) { got[id] = true })
	return got
}

// Nobody has to know capabilities exist: entries never given any answer an
// ordinary query exactly as they did before capabilities existed. What they
// must not do is answer a query for somebody else's — an entry with every
// capability could never be filtered out of anything.
func TestCapabilities_AnUnassignedEntityIsPlainOnly(t *testing.T) {
	m, ids := capGrid(t)
	box := NewAABBAt(NewVec(90, 90), 60, 40)

	for _, want := range []Capability{AnyCapability, Plain} {
		got := found(m, box, want)
		for _, id := range ids {
			if !got[id] {
				t.Errorf("entity %v missing from a query for %#x, which every unassigned entity answers", id, want)
			}
		}
	}

	const somebodyElse Capability = 1 << 7
	if got := found(m, box, somebodyElse); len(got) != 0 {
		t.Errorf("query for %#x found %v, want nothing — no entity was put there", somebodyElse, got)
	}
}

func TestCapabilities_QuerySeesOnlyWhatItAskedFor(t *testing.T) {
	m, ids := capGrid(t)

	const solid, ghost Capability = 1 << 1, 1 << 2
	m.QueueSetCapabilities(ids[0], solid)
	m.QueueSetCapabilities(ids[1], ghost)
	m.QueueSetCapabilities(ids[2], solid|ghost)
	m.Flush(nil)

	box := NewAABBAt(NewVec(90, 90), 60, 40)

	got := found(m, box, solid)
	if !got[ids[0]] || got[ids[1]] || !got[ids[2]] {
		t.Errorf("query for solid found %v, want the solid one and the one in both", got)
	}

	got = found(m, box, ghost)
	if got[ids[0]] || !got[ids[1]] || !got[ids[2]] {
		t.Errorf("query for ghost found %v, want the ghost one and the one in both", got)
	}

	if got = found(m, box, AnyCapability); len(got) != 3 {
		t.Errorf("query for AnyCapability found %v, want all three", got)
	}
}

// Moving onto a seam grows wrapped pieces, which the manager handles by
// removing and reinserting the entity. That is bookkeeping, not death: it is
// the same entity and must still be what it was.
func TestCapabilities_SurviveRefragmentation(t *testing.T) {
	space := plane.NewToroidal2D(256, 256)
	m, err := NewGridIndexManager(space, GridIndexConfig{
		Resolution: Size256x256, BucketResolution: Size32x32, BucketCapacity: 8,
	})
	if err != nil {
		t.Fatalf("NewGridIndexManager: %v", err)
	}

	// Not Plain: an entity that forgot its capabilities falls back to it, so
	// a test expecting the fallback could not tell the two apart.
	const solid Capability = 1 << 1
	id := uid.UID64(1)
	m.QueueInsert(id, space.WrapAABB(NewAABBAt(NewVec(100, 100), 8, 8)))
	m.QueueSetCapabilities(id, solid)
	m.Flush(nil)

	// Onto the right seam: the entity now has a wrapped piece it did not have.
	m.QueueUpdate(id, space.WrapAABB(NewAABBAt(NewVec(252, 100), 8, 8)), false)
	m.Flush(nil)

	// Asking for solid proves little on its own — what a loss cannot survive
	// is being asked for a capability the entity was never given.
	const ghost Capability = 1 << 2
	for name, box := range map[string]AABB{
		"its body":         NewAABBAt(NewVec(250, 98), 8, 12),
		"its wrapped part": NewAABBAt(NewVec(0, 98), 8, 12),
	} {
		if got := found(m, box, solid); !got[id] {
			t.Errorf("%s no longer answers a query for the capability it has", name)
		}
		if got := found(m, box, ghost); got[id] {
			t.Errorf("%s answers a query for a capability it never had — they were forgotten when it grew a wrapped piece", name)
		}
	}
}

// An index is reused once its entity dies; the newcomer must not inherit what
// the previous occupant could do.
func TestCapabilities_AreNotInheritedByARecycledIndex(t *testing.T) {
	space := plane.NewToroidal2D(256, 256)
	m, err := NewGridIndexManager(space, GridIndexConfig{
		Resolution: Size256x256, BucketResolution: Size32x32, BucketCapacity: 8,
	})
	if err != nil {
		t.Fatalf("NewGridIndexManager: %v", err)
	}

	const ghost Capability = 1 << 1
	var pool uid.UID64Pool
	pool.Init(4, 4)

	first := pool.Next()
	m.QueueInsert(first, space.WrapAABB(NewAABBAt(NewVec(100, 100), 8, 8)))
	m.QueueSetCapabilities(first, ghost)
	m.Flush(nil)

	m.QueueRemove(first)
	m.Flush(nil)
	pool.Release(first)

	second := pool.Next()
	if second.Index() != first.Index() {
		t.Skipf("pool did not recycle the index, nothing to prove here")
	}
	m.QueueInsert(second, space.WrapAABB(NewAABBAt(NewVec(100, 100), 8, 8)))
	m.Flush(nil)

	if got := found(m, NewAABBAt(NewVec(98, 98), 12, 12), Plain); !got[second] {
		t.Error("the new entity is not Plain — it inherited the dead one's capabilities")
	}
}

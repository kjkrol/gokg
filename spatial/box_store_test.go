package spatial

import (
	"testing"

	"github.com/kjkrol/gokg/plane"
	"github.com/kjkrol/uid"
)

func TestBoxStore_KeepsMainBoxesAndFragmentsApart(t *testing.T) {
	s := newBoxStore(4)

	id := uid.UID64(7)
	main := NewAABBAt(NewVec(10, 10), 4, 4)
	frag := NewAABBAt(NewVec(0, 10), 2, 4)

	s.set(withFrag(id, uint8(plane.FRAG_MAIN)), main)
	s.set(withFrag(id, uint8(plane.FRAG_RIGHT)), frag)

	if got, ok := s.get(withFrag(id, uint8(plane.FRAG_MAIN))); !ok || got != main {
		t.Errorf("main box = %v (found=%v), want %v", got, ok, main)
	}
	if got, ok := s.get(withFrag(id, uint8(plane.FRAG_RIGHT))); !ok || got != frag {
		t.Errorf("right fragment = %v (found=%v), want %v", got, ok, frag)
	}
	if s.len() != 2 {
		t.Errorf("len = %d, want 2 — a fragment is its own entry", s.len())
	}
}

// The main slice is addressed by index alone, so a recycled index would hand
// out its predecessor's box if the id were not checked as well. A map simply
// failed to find the dead id; this has to refuse just as firmly.
func TestBoxStore_RefusesAnIdWhoseIndexWasRecycled(t *testing.T) {
	s := newBoxStore(4)

	var pool uid.UID64Pool
	pool.Init(4, 4)
	first := pool.Next()
	s.set(first, NewAABBAt(NewVec(10, 10), 4, 4))
	pool.Release(first)

	second := pool.Next()
	if second.Index() != first.Index() {
		t.Skipf("pool did not recycle the index (%d vs %d), nothing to prove here", second.Index(), first.Index())
	}
	s.remove(first)
	s.set(second, NewAABBAt(NewVec(50, 50), 4, 4))

	if got, ok := s.get(first); ok {
		t.Errorf("the dead id still resolves to %v — it shares an index with a live entity", got)
	}
}

func TestBoxStore_RemoveAndClearKeepTheCount(t *testing.T) {
	s := newBoxStore(4)

	for i := range 5 {
		s.set(uid.UID64(i), NewAABBAt(NewVec(float64(i), 0), 2, 2))
	}
	s.set(withFrag(uid.UID64(1), uint8(plane.FRAG_BOTTOM)), NewAABBAt(NewVec(0, 0), 2, 2))
	if s.len() != 6 {
		t.Fatalf("len = %d after six inserts, want 6", s.len())
	}

	s.remove(uid.UID64(2))
	s.remove(uid.UID64(2)) // removing twice must not double-count
	s.remove(withFrag(uid.UID64(1), uint8(plane.FRAG_BOTTOM)))
	if s.len() != 4 {
		t.Errorf("len = %d after two distinct removals, want 4", s.len())
	}

	s.clear()
	if s.len() != 0 {
		t.Errorf("len = %d after clear, want 0", s.len())
	}
	if _, ok := s.get(uid.UID64(0)); ok {
		t.Error("an entry survived clear")
	}
}

// An entity straddling a seam is indexed once per wrapped piece, and a query
// has to find it through whichever piece it meets — the main box now living in
// a slice and the pieces in a map must not change that.
func TestQueryRange_FindsAWrappedEntityThroughEveryPiece(t *testing.T) {
	space := plane.NewToroidal2D(256, 256)
	m, err := NewGridIndexManager(space, GridIndexConfig{
		Resolution:       Size256x256,
		BucketResolution: Size32x32,
		BucketCapacity:   8,
	})
	if err != nil {
		t.Fatalf("NewGridIndexManager: %v", err)
	}

	id := uid.UID64(1)
	// Straddling the bottom-right corner, so it has all three wrapped pieces.
	m.QueueInsert(id, space.WrapAABB(NewAABBAt(NewVec(250, 250), 12, 12)))
	m.Flush(nil)

	corners := map[string]AABB{
		"where its body is":         NewAABBAt(NewVec(248, 248), 6, 6),
		"across the right seam":     NewAABBAt(NewVec(0, 250), 4, 4),
		"across the bottom seam":    NewAABBAt(NewVec(250, 0), 4, 4),
		"across both at the corner": NewAABBAt(NewVec(0, 0), 4, 4),
	}
	for name, box := range corners {
		t.Run(name, func(t *testing.T) {
			found := false
			m.QueryRange(box, func(got uid.UID64, _ plane.FragPosition) {
				if got == id {
					found = true
				}
			})
			if !found {
				t.Errorf("query at %v did not find the wrapped entity", box)
			}
		})
	}
}

package spatial

import "github.com/kjkrol/uid"

// boxStore holds the box of every entry the grid indexes, keyed by entry id.
//
// Main boxes live in a slice indexed by uid.UID64.Index(), which uid documents
// as the "main sequence for array/pool offsets" — the ids a grid is given come
// from uid.UID64Pool, which hands out indices densely from zero. Keying that by
// a hash map meant every candidate a query examined paid for hashing and a
// random probe, on what is the hottest loop in the package.
//
// Wrapped fragments keep a map, because they are rare by construction: only an
// entity straddling a seam has any, and a euclidean space has none at all.
type boxStore struct {
	main  []mainBox
	frags map[uid.UID64]AABB
	count int
}

// mainBox is one slot of the main slice. The id is kept and compared, not just
// the liveness flag: an index is reused once its entity dies, and a stale id
// left in a bucket would otherwise read its successor's box as if it were its
// own — a map simply failed to find it.
type mainBox struct {
	aabb AABB
	id   uid.UID64
	live bool
}

func newBoxStore(capacity int) *boxStore {
	return &boxStore{
		main:  make([]mainBox, 0, capacity),
		frags: make(map[uid.UID64]AABB),
	}
}

// get returns the box stored for id, and whether there is one.
func (s *boxStore) get(id uid.UID64) (AABB, bool) {
	if fragOf(id) != 0 {
		aabb, ok := s.frags[id]
		return aabb, ok
	}
	i := int(id.Index())
	if i >= len(s.main) {
		return AABB{}, false
	}
	slot := &s.main[i]
	if !slot.live || slot.id != id {
		return AABB{}, false
	}
	return slot.aabb, true
}

// set records id's box, replacing whatever was there.
func (s *boxStore) set(id uid.UID64, aabb AABB) {
	if fragOf(id) != 0 {
		if _, held := s.frags[id]; !held {
			s.count++
		}
		s.frags[id] = aabb
		return
	}
	i := int(id.Index())
	for i >= len(s.main) {
		s.main = append(s.main, mainBox{})
	}
	slot := &s.main[i]
	if !slot.live {
		s.count++
	}
	slot.aabb, slot.id, slot.live = aabb, id, true
}

// remove drops id's box if it has one.
func (s *boxStore) remove(id uid.UID64) {
	if fragOf(id) != 0 {
		if _, held := s.frags[id]; held {
			delete(s.frags, id)
			s.count--
		}
		return
	}
	i := int(id.Index())
	if i >= len(s.main) {
		return
	}
	slot := &s.main[i]
	if slot.live && slot.id == id {
		slot.live = false
		s.count--
	}
}

// len reports how many entries are stored.
func (s *boxStore) len() int { return s.count }

// clear empties the store, keeping the memory it has already grown into.
func (s *boxStore) clear() {
	s.main = s.main[:0]
	clear(s.frags)
	s.count = 0
}

package spatial

import "github.com/kjkrol/uid"

// boxStore holds the box and capabilities of every entry the grid indexes.
//
// Main boxes are addressed by uid.UID64.Index(), which uid.UID64Pool hands out
// densely from zero; wrapped fragments keep a map, being rare by construction.
type boxStore struct {
	main  []mainBox
	frags map[uid.UID64]AABB
	count int

	// caps is addressed like main but kept apart: a query tests it before it
	// has any reason to touch a box. Capabilities belong to the entity, so
	// every fragment of one shares the entry.
	caps []Capability
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

// init readies the store for capacity entries, in place: bucketGrid holds it
// by value, so get costs no indirection in the grid's hottest loop.
func (s *boxStore) init(capacity int) {
	s.main = make([]mainBox, 0, capacity)
	s.frags = make(map[uid.UID64]AABB)
	s.caps = make([]Capability, 0, capacity)
	s.count = 0
}

// capsOf returns id's capabilities, Plain if it was never given any.
func (s *boxStore) capsOf(id uid.UID64) Capability {
	i := int(id.Index())
	if i >= len(s.caps) {
		return Plain
	}
	return s.caps[i]
}

// setCaps records id's capabilities.
func (s *boxStore) setCaps(id uid.UID64, c Capability) {
	i := int(id.Index())
	for i >= len(s.caps) {
		s.caps = append(s.caps, Plain)
	}
	s.caps[i] = c
}

// clearCaps drops id's capabilities, so a recycled index starts from Plain
// rather than inheriting its predecessor's.
func (s *boxStore) clearCaps(id uid.UID64) {
	i := int(id.Index())
	if i < len(s.caps) {
		s.caps[i] = Plain
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
	s.caps = s.caps[:0]
	clear(s.frags)
	s.count = 0
}

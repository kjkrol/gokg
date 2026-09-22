package spatial

import "github.com/kjkrol/uid"

// boxStore holds the box and capabilities of every entry the grid indexes.
type boxStore struct {
	main  []mainBox
	frags map[uid.UID64]AABB
	count int

	// caps is addressed like main; every fragment of an entity shares the entry.
	caps []Capability

	// sizes is addressed like main and holds the whole entity's extent.
	sizes []Vec
}

// mainBox is one slot of the main slice; the id tells a live entry from a recycled index.
// mask names the pieces the entity is indexed in, tl and br the cells its main box spans.
type mainBox struct {
	aabb   AABB
	id     uid.UID64
	tl, br int
	mask   uint8
	live   bool
}

// init readies the store for capacity entries.
func (s *boxStore) init(capacity int) {
	s.main = make([]mainBox, 0, capacity)
	s.frags = make(map[uid.UID64]AABB)
	s.caps = make([]Capability, 0, capacity)
	s.sizes = make([]Vec, 0, capacity)
	s.count = 0
}

// setSize records the whole extent of the entity id names.
func (s *boxStore) setSize(id uid.UID64, size Vec) {
	i := int(id.Index())
	for i >= len(s.sizes) {
		s.sizes = append(s.sizes, Vec{})
	}
	s.sizes[i] = size
}

// sizeOf returns the whole extent of the entity id names, or main's own if none was given.
func (s *boxStore) sizeOf(id uid.UID64, main AABB) Vec {
	if i := int(id.Index()); i < len(s.sizes) && s.sizes[i] != (Vec{}) {
		return s.sizes[i]
	}
	return NewVec(main.BottomRight.X-main.TopLeft.X, main.BottomRight.Y-main.TopLeft.Y)
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

// clearCaps drops id's capabilities.
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

// slot returns the main slot of the entity id names, if it is live.
func (s *boxStore) slot(id uid.UID64) (*mainBox, bool) {
	i := int(id.Index())
	if i >= len(s.main) {
		return nil, false
	}
	slot := &s.main[i]
	if !slot.live || slot.id != withoutFrag(id) {
		return nil, false
	}
	return slot, true
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
	s.sizes = s.sizes[:0]
	clear(s.frags)
	s.count = 0
}

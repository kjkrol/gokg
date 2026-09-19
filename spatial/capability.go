package spatial

// Capability is what may be done with an entry, for queries that only care
// about part of what the index holds. One index serves every system asking a
// spatial question — collision, selection, line of sight — and each wants a
// different subset.
//
// The eight bits are a set: an entry can have several capabilities, and a
// query matches when any bit is shared. What the bits mean is the caller's
// business; the index only compares them.
type Capability uint8

const (
	// Plain is an entry that is nothing but geometry — where one starts
	// until a caller says otherwise. It is a bit rather than zero so that a
	// query can still ask for it.
	Plain Capability = 1 << 0

	// AnyCapability is the query mask that accepts every entry.
	AnyCapability Capability = 0xFF
)

// matches reports whether an entry with capabilities c answers a query for want.
func (c Capability) matches(want Capability) bool { return c&want != 0 }

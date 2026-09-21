package spatial

// Capability is what may be done with an entry, as a set of up to eight bits.
// A query matches when any bit is shared; what the bits mean is the caller's business.
type Capability uint8

const (
	// Plain is an entry that is nothing but geometry — what one is until told otherwise.
	Plain Capability = 1 << 0

	// AnyCapability is the query mask that accepts every entry.
	AnyCapability Capability = 0xFF
)

// matches reports whether an entry with capabilities c answers a query for want.
func (c Capability) matches(want Capability) bool { return c&want != 0 }

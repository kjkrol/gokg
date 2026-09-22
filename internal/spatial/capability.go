package spatial

// Capability is what may be done with an entry, as a set of up to eight bits.
// A query matches when any bit is shared; what the bits mean is the caller's business.
type Capability uint8

const (
	// Plain is an entry that is nothing but geometry; the zero value.
	Plain Capability = 0
	// CanCollide is an entry that takes part in collisions.
	CanCollide Capability = 1 << 1
	// Static is an entry nothing shifts.
	Static Capability = 1 << 2
	// Sensor is an entry only ever detected.
	Sensor Capability = 1 << 3

	// AnyCapability is the query mask that accepts every entry.
	AnyCapability Capability = 0xFF
)

// matches reports whether an entry with capabilities c answers a query for want.
func (c Capability) matches(want Capability) bool { return want == AnyCapability || c&want != 0 }

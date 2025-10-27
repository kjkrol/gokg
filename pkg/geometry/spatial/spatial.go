package spatial

// Spatial represents an entity that occupies space on a 2D plane and exposes its
// axis-aligned bounding box (AABB) for neighbor lookups.
type Spatial[T SupportedNumeric] interface {
	Bounds() Rectangle[T]
	// Vektory tworzace / kluczowe (np dla Rectangle sa to top-left i bottom-right)
	Vertices() []*Vec[T]
	Fragments() []Spatial[T]
	SetFragments([]Spatial[T])
}

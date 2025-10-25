package geometry

// Spatial represents an entity that occupies space on a 2D plane and exposes its
// axis-aligned bounding box (AABB) for neighbor lookups.
type Spatial[T SupportedNumeric] interface {
	Bounds() Rectangle[T]
	Probe(margin T, plane Plane[T]) []Rectangle[T]
	DistanceTo(other Spatial[T], distance Distance[T]) T
	// Vektory tworzace / kluczowe (np dla Rectangle sa to top-left i bottom-right)
	Vertices() []*Vec[T]
	Fragments() []Spatial[T]
	SetFragments([]Spatial[T])
}

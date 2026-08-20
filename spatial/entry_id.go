package spatial

import "github.com/kjkrol/uid"

// fragSegment reserves 2 bits of uid.UID64's metadata byte for the fragment
// index (plane.FragPosition has 4 values: MAIN/RIGHT/BOTTOM/BOTTOM_RIGHT) —
// separate from the 32-bit index and 24-bit generation, so tagging a fragment
// never collides with entity identity.
var fragSegment = uid.NewMetaSegment(2, 0)

// withFrag returns id tagged with frag, for use as a bucket-grid entry key —
// a single (entity, fragment) pair can then live as one comparable value.
func withFrag(id uid.UID64, frag uint8) uid.UID64 {
	return id.WithMetaSegment(fragSegment, frag)
}

// fragOf extracts the fragment tag from an id previously tagged by withFrag.
func fragOf(id uid.UID64) uint8 {
	return id.MetaSegment(fragSegment)
}

// withoutFrag strips the fragment tag, recovering the plain entity id.
func withoutFrag(id uid.UID64) uid.UID64 {
	return id.WithMetaSegment(fragSegment, 0)
}

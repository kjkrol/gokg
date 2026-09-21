package spatial

import "github.com/kjkrol/uid"

// fragSegment is the two metadata bits of a uid.UID64 that carry the fragment index.
var fragSegment = uid.NewMetaSegment(2, 0)

// withFrag returns id tagged with frag, for use as a bucket-grid entry key.
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

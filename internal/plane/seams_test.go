package plane

import (
	"github.com/kjkrol/aabbworld/geom"
	"github.com/kjkrol/aabbworld/plane"
)

// imagesOf is a box's main rectangle and every piece it wraps into.
func imagesOf(ab plane.AABB) []geom.AABB {
	images := []geom.AABB{ab.AABB}
	ab.VisitFragments(func(_ plane.FragPosition, box geom.AABB) bool {
		images = append(images, box)
		return true
	})
	return images
}

// meetAcrossSeams reports whether any image of a touches any image of b.
func meetAcrossSeams(a, b plane.AABB) bool {
	for _, ia := range imagesOf(a) {
		for _, ib := range imagesOf(b) {
			if ia.Intersects(ib) {
				return true
			}
		}
	}
	return false
}

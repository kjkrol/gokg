package plane

import (
	"github.com/kjkrol/gokg/geom"
)

// FragPosition identifies a fragment's position relative to its parent AABB (axis-aligned bounding box).
// Names follow logical cardinal directions of the parent; depending on screen
// coordinates they may appear flipped (e.g. right on a Euclidean grid may render left in screen space).
type FragPosition int

const (
	FRAG_MAIN FragPosition = iota
	// FRAG_RIGHT is the fragment that spans the parent's right edge.
	FRAG_RIGHT
	// FRAG_BOTTOM is the fragment along the parent's bottom edge.
	FRAG_BOTTOM
	// FRAG_BOTTOM_RIGHT is the fragment covering the parent's bottom-right quadrant.
	FRAG_BOTTOM_RIGHT
)

// AABB extends geom.AABB with the size and wrap overhang Space normalisation
// needs. It is the Space-aware view of a AABB: Space keeps AABB instances
// canonical within its domain.
//
// Overhang is how far the box runs past the world's far edges once normalised,
// and zero on an axis it does not cross. It is everything a wrapped fragment
// needs: each piece follows from it and the main box, so the pieces are rebuilt
// on demand by VisitFragments rather than carried in every copy of the struct.
// Holding them instead would cost six times the memory in what is, for a moving
// entity, a cache recomputed every tick and read once.
//
// The field is exported (not hidden behind a BinaryMarshaler) so AABB stays a
// plain, recursively POD-encodable type for goke's persist — embedding a type
// with its own MarshalBinary would silently drop any sibling fields on the
// embedding struct that MarshalBinary doesn't know about (Go promotes the
// method to the whole outer type).
type AABB struct {
	geom.AABB
	Size     geom.Vec
	Overhang geom.Vec
}

// NewAABB builds a AABB at pos with the given size, not yet wrapped by any Space.
func NewAABB(pos geom.Vec, width, height float64) AABB {
	return AABB{
		AABB: geom.AABB{
			TopLeft:     pos,
			BottomRight: geom.NewVec(pos.X+width, pos.Y+height),
		},
		Size: geom.NewVec(width, height),
	}
}

// --------------------------------------------------------------------------

// String formats the aabb using its top-left and bottom-right corners.
func (ab AABB) String() string {
	return ab.AABB.String()
}

// Equals reports whether ab and other share the same corners.
func (ab AABB) Equals(other AABB) bool {
	return ab.AABB.Equals(other.AABB)
}

// ContainsWithFrags reports whether ab, or any of its wrapped fragments,
// contains other or any of its own.
func (ab AABB) ContainsWithFrags(other AABB) bool {
	return ab.overlapsWithFrags(other, geom.AABB.Contains)
}

// IntersectsWithFrags reports whether ab, or any of its wrapped fragments,
// intersects other or any of its own.
func (ab AABB) IntersectsWithFrags(other AABB) bool {
	return ab.overlapsWithFrags(other, geom.AABB.Intersects)
}

// overlapsWithFrags asks whether any image of ab meets any image of other. The
// two exported forms differ only in what "meets" means, so they share the walk:
// main against main, then main against the other's fragments, then each of ab's
// fragments against the other's main and fragments.
func (ab AABB) overlapsWithFrags(other AABB, meets func(a, b geom.AABB) bool) bool {
	if meets(ab.AABB, other.AABB) {
		return true
	}
	if !ab.hasFragments() && !other.hasFragments() {
		return false
	}

	met := false
	other.VisitFragments(func(_ FragPosition, frag geom.AABB) bool {
		met = meets(ab.AABB, frag)
		return !met
	})
	if met {
		return true
	}

	ab.VisitFragments(func(_ FragPosition, frag geom.AABB) bool {
		if meets(other.AABB, frag) {
			met = true
			return false
		}
		other.VisitFragments(func(_ FragPosition, otherFrag geom.AABB) bool {
			met = meets(frag, otherFrag)
			return !met
		})
		return !met
	})
	return met
}

// hasFragments reports whether the box reaches past a world edge at all.
func (ab AABB) hasFragments() bool {
	var none float64
	return ab.Overhang.X > none || ab.Overhang.Y > none
}

type FragVisitor func(pos FragPosition, box geom.AABB) bool

// VisitFragments calls fn for each piece the box wraps into, rebuilding it from
// the overhang: the part past the right edge reappears at the left, the part
// past the bottom at the top, and the corner where both happen at once.
// fn returning false stops the walk.
func (ab *AABB) VisitFragments(fn FragVisitor) {
	var none float64
	dx, dy := ab.Overhang.X, ab.Overhang.Y

	if dx > none {
		box := geom.NewAABB(geom.NewVec(none, ab.TopLeft.Y), geom.NewVec(dx, ab.BottomRight.Y))
		if !fn(FRAG_RIGHT, box) {
			return
		}
	}
	if dy > none {
		box := geom.NewAABB(geom.NewVec(ab.TopLeft.X, none), geom.NewVec(ab.BottomRight.X, dy))
		if !fn(FRAG_BOTTOM, box) {
			return
		}
	}
	if dx > none && dy > none {
		fn(FRAG_BOTTOM_RIGHT, geom.NewAABB(geom.NewVec(none, none), geom.NewVec(dx, dy)))
	}
}

// fragmentation records how far the box ran past the far edges; the pieces
// themselves are rebuilt from it whenever someone asks.
func (ab *AABB) fragmentation(dx, dy float64) {
	ab.Overhang = geom.NewVec(dx, dy)
}

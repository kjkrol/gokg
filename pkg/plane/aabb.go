package plane

import "github.com/kjkrol/gokg/pkg/geom"

// FragPosition identifies a fragment's position relative to its parent AABB (axis-aligned bounding box).
// Names follow logical cardinal directions of the parent; depending on screen
// coordinates they may appear flipped (e.g. right on a Euclidean grid may render left in screen space).
type FragPosition int

const (
	// FRAG_RIGHT is the fragment that spans the parent's right edge.
	FRAG_RIGHT FragPosition = iota
	// FRAG_BOTTOM is the fragment along the parent's bottom edge.
	FRAG_BOTTOM
	// FRAG_BOTTOM_RIGHT is the fragment covering the parent's bottom-right quadrant.
	FRAG_BOTTOM_RIGHT
)

// AABB extends geom.AABB with cached width, height, and boundary fragments used by Space normalisation.
// It is the Space-aware view of a AABB: Space keeps AABB instances canonical within its domain.
type AABB[T geom.Numeric] struct {
	geom.AABB[T]
	size     geom.Vec[T]
	frags    [3]geom.AABB[T]
	fragMask uint8
}

// newAABB builds a AABB at pos with the given size, priming fragment storage for Space operations.
func newAABB[T geom.Numeric](pos geom.Vec[T], width, height T) AABB[T] {
	return AABB[T]{
		AABB: geom.AABB[T]{
			TopLeft:     pos,
			BottomRight: geom.NewVec(pos.X+width, pos.Y+height),
		},
		size: geom.NewVec(width, height),
	}
}

// --------------------------------------------------------------------------

// String formats the aabb using its top-left and bottom-right corners.
func (ab AABB[T]) String() string {
	return ab.AABB.String()
}

// Equals reports whether ab and other share the same corners.
func (ab AABB[T]) Equals(other AABB[T]) bool {
	return ab.AABB.Equals(other.AABB)
}

func (ab AABB[T]) ContainsWithFrags(other AABB[T]) bool {
	base := ab.AABB
	otherBase := other.AABB

	if base.Contains(otherBase) {
		return true
	}

	abHasFrags := ab.fragMask != 0
	otherHasFrags := other.fragMask != 0

	if !abHasFrags && !otherHasFrags {
		return false
	}

	for idx := range len(other.frags) {
		if other.fragMask&(1<<idx) != 0 && base.Contains(other.frags[idx]) {
			return true
		}
	}

	// other contains any fragment of ab
	for idx := range len(ab.frags) {
		if ab.fragMask&(1<<idx) == 0 {
			continue
		}
		frag := ab.frags[idx]
		if otherBase.Contains(frag) {
			return true
		}
		if !otherHasFrags {
			continue
		}
		if fragContainsAny(frag, other) {
			return true
		}
	}

	return false
}

// IntersectsWithFrags reports whether ab IntersectsWithFrags other or any of its fragments.
func (ab AABB[T]) IntersectsWithFrags(other AABB[T]) bool {
	base := ab.AABB
	otherBase := other.AABB

	if base.Intersects(otherBase) {
		return true
	}

	abHasFrags := ab.fragMask != 0
	otherHasFrags := other.fragMask != 0

	if !abHasFrags && !otherHasFrags {
		return false
	}

	for idx := range len(other.frags) {
		if other.fragMask&(1<<idx) != 0 && base.Intersects(other.frags[idx]) {
			return true
		}
	}

	for idx := range len(ab.frags) {
		if ab.fragMask&(1<<idx) == 0 {
			continue
		}
		frag := ab.frags[idx]
		if otherBase.Intersects(frag) {
			return true
		}
		if otherHasFrags && fragIntersectsAny(frag, other) {
			return true
		}
	}

	return false
}

type FragVisitor[T geom.Numeric] func(pos FragPosition, box geom.AABB[T]) bool

func (ab *AABB[T]) VisitFragments(fn FragVisitor[T]) {
	for pos := range len(ab.frags) {
		if ab.fragMask&(1<<pos) == 0 {
			continue
		}
		if !fn(FragPosition(pos), ab.frags[pos]) {
			return
		}
	}
}

func (ab *AABB[T]) setFragment(pos FragPosition, box geom.AABB[T]) {
	ab.frags[pos] = box
	ab.fragMask |= 1 << pos
}

func (ab *AABB[T]) clearFragment(pos FragPosition) {
	ab.fragMask &^= 1 << pos
}

func (ab *AABB[T]) fragmentation(dx, dy T) {
	if dx > 0 {
		ab.setFragment(FRAG_RIGHT, geom.NewAABB(geom.NewVec(0, ab.TopLeft.Y), geom.NewVec(dx, ab.BottomRight.Y)))
	} else {
		ab.clearFragment(FRAG_RIGHT)
	}
	if dy > 0 {
		ab.setFragment(FRAG_BOTTOM, geom.NewAABB(geom.NewVec(ab.TopLeft.X, 0), geom.NewVec(ab.BottomRight.X, dy)))
	} else {
		ab.clearFragment(FRAG_BOTTOM)
	}
	if dx > 0 && dy > 0 {
		ab.setFragment(FRAG_BOTTOM_RIGHT, geom.NewAABB(geom.NewVec[T](0, 0), geom.NewVec(dx, dy)))
	} else {
		ab.clearFragment(FRAG_BOTTOM_RIGHT)
	}
}

func fragContainsAny[T geom.Numeric](frag geom.AABB[T], other AABB[T]) bool {
	for j := range len(other.frags) {
		if other.fragMask&(1<<j) != 0 && frag.Contains(other.frags[j]) {
			return true
		}
	}
	return false
}

func fragIntersectsAny[T geom.Numeric](frag geom.AABB[T], other AABB[T]) bool {
	for j := range len(other.frags) {
		if other.fragMask&(1<<j) != 0 && frag.Intersects(other.frags[j]) {
			return true
		}
	}
	return false
}

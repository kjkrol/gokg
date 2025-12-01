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
	size  geom.Vec[T]
	frags [3]geom.AABB[T]
	set   [3]bool
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

func (ab AABB[T]) Contains(other AABB[T]) bool {
	return ab.compareWithFrags(other, geom.AABB[T].Contains)
}

// Intersects reports whether ab Intersects other or any of its fragments.
func (ab AABB[T]) Intersects(other AABB[T]) bool {
	return ab.compareWithFrags(other, geom.AABB[T].Intersects)
}

type FragVisitor[T geom.Numeric] func(pos FragPosition, box geom.AABB[T]) bool

func (ab *AABB[T]) VisitFragments(fn FragVisitor[T]) {
	for pos, set := range ab.set {
		if !set {
			continue
		}
		if !fn(FragPosition(pos), ab.frags[pos]) {
			return
		}
	}
}

func (ab *AABB[T]) setFragment(pos FragPosition, box geom.AABB[T]) {
	ab.frags[pos] = box
	ab.set[pos] = true
}

func (ab *AABB[T]) clearFragment(pos FragPosition) {
	ab.set[pos] = false
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

func (ab AABB[T]) compareWithFrags(other AABB[T], compareFn func(geom.AABB[T], geom.AABB[T]) bool) bool {
	if compareFn(ab.AABB, other.AABB) {
		return true
	}
	found := false
	(&other).VisitFragments(func(_ FragPosition, frag geom.AABB[T]) bool {
		if compareFn(ab.AABB, frag) {
			found = true
			return false
		}
		return true
	})
	if found {
		return true
	}
	(&ab).VisitFragments(func(_ FragPosition, frag geom.AABB[T]) bool {
		if compareFn(other.AABB, frag) {
			found = true
			return false
		}
		(&other).VisitFragments(func(_ FragPosition, otherFrag geom.AABB[T]) bool {
			if compareFn(frag, otherFrag) {
				found = true
				return false
			}
			return true
		})
		return !found
	})

	return found
}

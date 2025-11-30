package plane

import "github.com/kjkrol/gokg/pkg/geom"

// FragPosition identifies a fragment's position relative to its parent AABB (align axis bounding-box).
// Names follow logical cardinal directions of the parent; depending on screen
// coordinates they may appear flipped (e.g. right on Cartesian may render left in screen space).
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
	frags map[FragPosition]geom.AABB[T]
}

// newAABB builds a AABB at pos with the given size, priming fragment storage for Space operations.
func newAABB[T geom.Numeric](pos geom.Vec[T], width, height T) AABB[T] {
	return AABB[T]{
		AABB: geom.AABB[T]{
			TopLeft:     pos,
			BottomRight: geom.NewVec(pos.X+width, pos.Y+height),
		},
		size:  geom.NewVec(width, height),
		frags: make(map[FragPosition]geom.AABB[T], 4),
	}
}

// --------------------------------------------------------------------------

// String formats the box using its top-left and bottom-right corners.
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

// Fragments returns lazily computed fragments keyed by bound position.
func (ab AABB[T]) Fragments() map[FragPosition]geom.AABB[T] { return ab.frags }

func (ab *AABB[T]) fragmentation(dx, dy T) {
	if dx < 0 {
		ab.frags[FRAG_RIGHT] = geom.NewAABB(geom.NewVec(0, ab.TopLeft.Y), geom.NewVec(-dx, ab.BottomRight.Y))
	} else {
		delete(ab.frags, FRAG_RIGHT)
	}
	if dy < 0 {
		ab.frags[FRAG_BOTTOM] = geom.NewAABB(geom.NewVec(ab.TopLeft.X, 0), geom.NewVec(ab.BottomRight.X, -dy))
	} else {
		delete(ab.frags, FRAG_BOTTOM)
	}
	if dx < 0 && dy < 0 {
		ab.frags[FRAG_BOTTOM_RIGHT] = geom.NewAABB(geom.NewVec[T](0, 0), geom.NewVec(-dx, -dy))
	} else {
		delete(ab.frags, FRAG_BOTTOM_RIGHT)
	}
}

func (ab AABB[T]) compareWithFrags(other AABB[T], compareFn func(geom.AABB[T], geom.AABB[T]) bool) bool {
	if compareFn(ab.AABB, other.AABB) {
		return true
	}
	for _, frag := range other.frags {
		if compareFn(ab.AABB, frag) {
			return true
		}
	}
	for _, frag := range ab.frags {
		if compareFn(other.AABB, frag) {
			return true
		}
	}
	for _, left := range ab.frags {
		for _, right := range other.frags {
			if compareFn(left, right) {
				return true
			}
		}
	}

	return false
}

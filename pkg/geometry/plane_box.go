package geometry

// FragPosition identifies a fragment's position relative to its parent PlaneBox (bounding-box).
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

// PlaneBox extends BoundingBox with cached width, height, and boundary fragments used by Plane normalisation.
// It is the Plane-aware view of a BoundingBox: Plane keeps PlaneBox instances canonical within its domain.
type PlaneBox[T SupportedNumeric] struct {
	BoundingBox[T]
	size  Vec[T]
	frags map[FragPosition]BoundingBox[T]
}

// newPlaneBox builds a PlaneBox at pos with the given size, priming fragment storage for Plane operations.
func newPlaneBox[T SupportedNumeric](pos Vec[T], width, height T) PlaneBox[T] {
	return PlaneBox[T]{
		BoundingBox: BoundingBox[T]{
			TopLeft:     pos,
			BottomRight: NewVec(pos.X+width, pos.Y+height),
		},
		size:  NewVec(width, height),
		frags: make(map[FragPosition]BoundingBox[T], 4),
	}
}

// --------------------------------------------------------------------------

// String formats the box using its top-left and bottom-right corners.
func (ab PlaneBox[T]) String() string {
	return ab.BoundingBox.String()
}

// Equals reports whether ab and other share the same corners.
func (ab PlaneBox[T]) Equals(other PlaneBox[T]) bool {
	return ab.BoundingBox.Equals(other.BoundingBox)
}

func (ab PlaneBox[T]) Contains(other PlaneBox[T]) bool {
	return ab.compareWithFrags(other, BoundingBox[T].Contains)
}

// Intersects reports whether ab Intersects other or any of its fragments.
func (ab PlaneBox[T]) Intersects(other PlaneBox[T]) bool {
	return ab.compareWithFrags(other, BoundingBox[T].Intersects)
}

// Fragments returns lazily computed fragments keyed by bound position.
func (ab PlaneBox[T]) Fragments() map[FragPosition]BoundingBox[T] { return ab.frags }

func (ab *PlaneBox[T]) fragmentation(dx, dy T) {
	if dx < 0 {
		ab.frags[FRAG_RIGHT] = NewBoundingBox(NewVec(0, ab.TopLeft.Y), NewVec(-dx, ab.BottomRight.Y))
	} else {
		delete(ab.frags, FRAG_RIGHT)
	}
	if dy < 0 {
		ab.frags[FRAG_BOTTOM] = NewBoundingBox(NewVec(ab.TopLeft.X, 0), NewVec(ab.BottomRight.X, -dy))
	} else {
		delete(ab.frags, FRAG_BOTTOM)
	}
	if dx < 0 && dy < 0 {
		ab.frags[FRAG_BOTTOM_RIGHT] = NewBoundingBox(NewVec[T](0, 0), NewVec(-dx, -dy))
	} else {
		delete(ab.frags, FRAG_BOTTOM_RIGHT)
	}
}

func (ab PlaneBox[T]) compareWithFrags(other PlaneBox[T], compareFn func(BoundingBox[T], BoundingBox[T]) bool) bool {
	if compareFn(ab.BoundingBox, other.BoundingBox) {
		return true
	}
	for _, frag := range other.frags {
		if compareFn(ab.BoundingBox, frag) {
			return true
		}
	}
	for _, frag := range ab.frags {
		if compareFn(other.BoundingBox, frag) {
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

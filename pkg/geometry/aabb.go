package geometry

import "fmt"

// BoundPosition identifies a fragment's position relative to its parent AABB.
// Names follow logical cardinal directions of the parent; depending on screen
// coordinates they may appear flipped (e.g. right on Cartesian may render left in screen space).
type BoundPosition int

const (
	// BOUND_RIGHT is the fragment that spans the parent's right edge.
	BOUND_RIGHT BoundPosition = iota
	// BOUND_BOTTOM is the fragment along the parent's bottom edge.
	BOUND_BOTTOM
	// BOUND_BOTTOM_RIGHT is the fragment covering the parent's bottom-right quadrant.
	BOUND_BOTTOM_RIGHT
)

// AABB represents an axis-aligned bounding box parameterized by numeric type T.
// It stores the top-left and bottom-right corners plus cached fragments for subdivision.
type AABB[T SupportedNumeric] struct {
	TopLeft     Vec[T]
	width       T
	height      T
	BottomRight Vec[T]
	frags       map[BoundPosition]AABB[T]
}

func NewAABB[T SupportedNumeric](pos Vec[T], width, height T) AABB[T] {
	return AABB[T]{
		TopLeft:     pos,
		width:       width,
		height:      height,
		BottomRight: NewVec(pos.X+width, pos.Y+height),
		frags:       make(map[BoundPosition]AABB[T], 4),
	}
}

// --------------------------------------------------------------------------

func BuildAABB[T SupportedNumeric](center Vec[T], d T) AABB[T] {
	topLeft := Vec[T]{X: center.X - d, Y: center.Y - d}
	return NewAABB(topLeft, 2*d, 2*d)
}

// Split divides the box into four equal quadrants around its center.
func (ab *AABB[T]) Split() [4]AABB[T] {
	center := ab.center()
	return [4]AABB[T]{
		NewAABB(ab.TopLeft, ab.width/2, ab.height/2),                     // top left
		NewAABB(NewVec(center.X, ab.TopLeft.Y), ab.width/2, ab.height/2), // top right
		NewAABB(NewVec(ab.TopLeft.X, center.Y), ab.width/2, ab.height/2), // bottom left
		NewAABB(center, ab.width/2, ab.height/2),                         // bottom right
	}
}

func (ab *AABB[T]) center() Vec[T] {
	centerX := (ab.TopLeft.X + ab.BottomRight.X) / 2
	centerY := (ab.TopLeft.Y + ab.BottomRight.Y) / 2
	return Vec[T]{X: centerX, Y: centerY}
}

// Contains reports whether other lies entirely within ab.
func (ab AABB[T]) Contains(other AABB[T]) bool {
	return other.TopLeft.X >= ab.TopLeft.X &&
		other.TopLeft.Y >= ab.TopLeft.Y &&
		other.BottomRight.X <= ab.BottomRight.X &&
		other.BottomRight.Y <= ab.BottomRight.Y
}

// SortRectanglesBy orders two boxes using the provided key functions and returns them as (min,max).
func SortRectanglesBy[T SupportedNumeric](a, b AABB[T], keyFns ...func(AABB[T]) T) (aa, bb AABB[T]) {
	for _, keyFn := range keyFns {
		av, bv := keyFn(a), keyFn(b)
		if av < bv {
			return a, b
		}
		if av > bv {
			return b, a
		}
	}
	return a, b
}

// Intersects reports whether this AABB overlaps another bounding box.
// It returns true both when the boxes share any interior volume
// and when they only touch along edges or vertices.
func (ab AABB[T]) Intersects(other AABB[T]) bool {
	// x axis check
	xIntersects := axisIntersection(ab, other, func(v Vec[T]) T { return v.X })
	if !xIntersects {
		return false
	}
	// y axis check
	yIntersects := axisIntersection(ab, other, func(v Vec[T]) T { return v.Y })
	return yIntersects
}

// IntersectsIncludingFrags reports whether ab intersects other or any of its fragments.
func (ab AABB[T]) IntersectsIncludingFrags(other AABB[T]) bool {
	if ab.Intersects(other) {
		return true
	}
	for _, frag := range other.frags {
		if ab.IntersectsIncludingFrags(frag) {
			return true
		}
	}
	for _, frag := range ab.frags {
		if other.IntersectsIncludingFrags(frag) {
			return true
		}
	}
	return false
}

//-------------------------------------------------------------------------

func axisIntersection[T SupportedNumeric](aa, bb AABB[T], axisValue func(Vec[T]) T) bool {
	aa, bb = SortRectanglesBy(
		aa, bb,
		func(r AABB[T]) T { return axisValue(r.TopLeft) },
		func(r AABB[T]) T { return axisValue(r.BottomRight) },
	)

	noIntersection := axisValue(aa.TopLeft) < axisValue(bb.BottomRight) &&
		axisValue(aa.BottomRight) < axisValue(bb.TopLeft)
	return !noIntersection
}

// String formats the box using its top-left and bottom-right corners.
func (ab AABB[T]) String() string {
	return fmt.Sprintf("{%v %v}", ab.TopLeft, ab.BottomRight)
}

// Equals reports whether ab and other share the same corners.
func (ab AABB[T]) Equals(other AABB[T]) bool {
	return ab.TopLeft == other.TopLeft && ab.BottomRight == other.BottomRight
}

// Fragments returns lazily computed fragments keyed by bound position.
func (ab AABB[T]) Fragments() map[BoundPosition]AABB[T] { return ab.frags }

func (ab *AABB[T]) fragmentation(dx, dy T) {
	if dx < 0 {
		ab.frags[BOUND_RIGHT] = newAABBFrag(NewVec(0, ab.TopLeft.Y), NewVec(-dx, ab.BottomRight.Y))
	} else {
		delete(ab.frags, BOUND_RIGHT)
	}
	if dy < 0 {
		ab.frags[BOUND_BOTTOM] = newAABBFrag(NewVec(ab.TopLeft.X, 0), NewVec(ab.BottomRight.X, -dy))
	} else {
		delete(ab.frags, BOUND_BOTTOM)
	}
	if dx < 0 && dy < 0 {
		ab.frags[BOUND_BOTTOM_RIGHT] = newAABBFrag(NewVec[T](0, 0), NewVec(-dx, -dy))
	} else {
		delete(ab.frags, BOUND_BOTTOM_RIGHT)
	}
}

func newAABBFrag[T SupportedNumeric](pos, bottomRight Vec[T]) AABB[T] {
	return AABB[T]{
		TopLeft:     pos,
		width:       0,
		height:      0,
		BottomRight: bottomRight,
		frags:       nil,
	}
}

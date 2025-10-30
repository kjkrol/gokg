package geometry

import (
	"fmt"
)

// AABB (Axis-Aligned Bounding Box) represents a bounding box aligned to the
// coordinate system axes. In 2D it is a rectangle, and in 3D a cuboid,
// whose edges are always parallel to the X, Y (and Z) axes.
//
// Unlike a general rectangle or box, an AABB is never rotated — it always
// remains axis-aligned ("standing straight" along the coordinate axes).
//
// It is typically defined by two points:
//   - Min: the lower corner (xMin, yMin[, zMin])
//   - Max: the upper corner (xMax, yMax[, zMax])
//
// Example in 2D:
//
//	Min = (2, 1)
//	Max = (5, 4)
//	This defines a rectangle spanning from (2,1) to (5,4).
//
// AABBs are commonly used in spatial algorithms, collision detection,
// and distance computations. Instead of comparing objects point-to-point,
// algorithms may first approximate them by their bounding boxes to
// quickly estimate overlap or gap between them.
type AABB[T SupportedNumeric] struct {
	TopLeft     Vec[T]
	BottomRight Vec[T]
	Center      Vec[T]
}

func NewAABB[T SupportedNumeric](topLeft Vec[T], bottomRight Vec[T]) AABB[T] {
	centerX := (topLeft.X + bottomRight.X) / 2
	centerY := (topLeft.Y + bottomRight.Y) / 2
	center := Vec[T]{X: centerX, Y: centerY}
	return AABB[T]{TopLeft: topLeft, BottomRight: bottomRight, Center: center}
}

func BuildAABB[T SupportedNumeric](center Vec[T], d T) AABB[T] {
	topLeft := Vec[T]{X: center.X - d, Y: center.Y - d}
	bottomRight := Vec[T]{X: center.X + d, Y: center.Y + d}
	return NewAABB(topLeft, bottomRight)
}

// Bounds returns the rectangle itself.
func (r AABB[T]) Bounds() AABB[T] {
	return r
}

// Vertices returns the corners that define the rectangle.
func (r *AABB[T]) Vertices() []*Vec[T] {
	return []*Vec[T]{&r.TopLeft, &r.BottomRight}
}

func (r *AABB[T]) Split() [4]AABB[T] {
	return [4]AABB[T]{
		NewAABB(r.TopLeft, r.Center), // top left
		NewAABB(Vec[T]{X: r.Center.X, Y: r.TopLeft.Y}, Vec[T]{X: r.BottomRight.X, Y: r.Center.Y}), // top right
		NewAABB(Vec[T]{X: r.TopLeft.X, Y: r.Center.Y}, Vec[T]{X: r.Center.X, Y: r.BottomRight.Y}), // bottom left
		NewAABB(r.Center, r.BottomRight), // bottom right
	}
}

func (r AABB[T]) Contains(other AABB[T]) bool {
	return other.TopLeft.X >= r.TopLeft.X &&
		other.TopLeft.Y >= r.TopLeft.Y &&
		other.BottomRight.X <= r.BottomRight.X &&
		other.BottomRight.Y <= r.BottomRight.Y
}

func (r AABB[T]) Expand(margin T) AABB[T] {
	return NewAABB(
		Vec[T]{X: r.TopLeft.X - margin, Y: r.TopLeft.Y - margin},
		Vec[T]{X: r.BottomRight.X + margin, Y: r.BottomRight.Y + margin},
	)
}

//-------------------------------------------------------------------------

func (r AABB[T]) Intersects(other AABB[T]) bool {
	// x axis check
	xIntersects := axisIntersection(r, other, func(v Vec[T]) T { return v.X })
	if !xIntersects {
		return false
	}
	// y axis check
	yIntersects := axisIntersection(r, other, func(v Vec[T]) T { return v.Y })
	return yIntersects
}

func axisIntersection[T SupportedNumeric](aa, bb AABB[T], axisValue func(Vec[T]) T) bool {
	aa, bb = sortRectanglesBy(aa, bb, axisValue)
	noIntersection := axisValue(aa.TopLeft) < axisValue(bb.BottomRight) && axisValue(aa.BottomRight) < axisValue(bb.TopLeft)
	return !noIntersection
}

func sortRectanglesBy[T SupportedNumeric](a, b AABB[T], axisValue func(Vec[T]) T) (aa, bb AABB[T]) {
	if axisValue(a.TopLeft) < axisValue(b.TopLeft) {
		aa = a
		bb = b
	} else {
		aa = b
		bb = a
	}
	return
}

//-------------------------------------------------------------------------

func (r AABB[T]) AxisDistanceTo(
	other AABB[T],
	axisValue func(Vec[T]) T,
) T {
	r, other = sortRectanglesBy(r, other, axisValue)

	if axisValue(r.BottomRight) >= axisValue(other.TopLeft) {
		return 0
	}
	return axisValue(other.TopLeft) - axisValue(r.BottomRight)
}

//-------------------------------------------------------------------------

func (r AABB[T]) IntersectsAny(others []AABB[T]) bool {
	for _, wrapped := range others {
		if r.Intersects(wrapped) {
			return true
		}
	}
	return false
}

//-------------------------------------------------------------------------

func (r AABB[T]) String() string {
	return fmt.Sprintf("{%v %v %v}", r.TopLeft, r.BottomRight, r.Center)
}

func (r AABB[T]) Equals(other AABB[T]) bool {
	return r.TopLeft == other.TopLeft && r.BottomRight == other.BottomRight && r.Center == other.Center
}

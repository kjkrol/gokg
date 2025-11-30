package geom

import "fmt"

// AABB is a minimal axis-aligned rectangle defined by its top-left and bottom-right corners.
type AABB[T Numeric] struct {
	TopLeft     Vec[T]
	BottomRight Vec[T]
}

// NewAABB constructs a axis-aligned bounding box from explicit corner vectors.
func NewAABB[T Numeric](topLeft, bottomRight Vec[T]) AABB[T] {
	return AABB[T]{
		TopLeft:     topLeft,
		BottomRight: bottomRight,
	}
}

// NewAABBAt builds a AABB starting at pos with the provided width and height.
func NewAABBAt[T Numeric](pos Vec[T], width, height T) AABB[T] {
	bottomRight := NewVec(pos.X+width, pos.Y+height)
	return NewAABB(pos, bottomRight)
}

// NewAABBAround creates a AABB centered at center with half-size d on each axis.
func NewAABBAround[T Numeric](center Vec[T], d T) AABB[T] {
	topLeft := Vec[T]{X: center.X - d, Y: center.Y - d}
	return NewAABBAt(topLeft, 2*d, 2*d)
}

// Equals reports whether ab and other share the same corners.
func (ab AABB[T]) Equals(other AABB[T]) bool {
	return ab.TopLeft == other.TopLeft && ab.BottomRight == other.BottomRight
}

// String formats the box using its top-left and bottom-right corners.
func (ab AABB[T]) String() string {
	return fmt.Sprintf("{%v %v}", ab.TopLeft, ab.BottomRight)
}

// Split divides the box into four equal quadrants around its center.
func (ab *AABB[T]) Split() [4]AABB[T] {
	center := ab.center()
	width_half := (ab.BottomRight.X - ab.TopLeft.X) / 2
	height_half := (ab.BottomRight.Y - ab.TopLeft.Y) / 2
	return [4]AABB[T]{
		NewAABBAt(ab.TopLeft, width_half, height_half),                     // top left
		NewAABBAt(NewVec(center.X, ab.TopLeft.Y), width_half, height_half), // top right
		NewAABBAt(NewVec(ab.TopLeft.X, center.Y), width_half, height_half), // bottom left
		NewAABBAt(center, width_half, height_half),                         // bottom right
	}
}

func (ab AABB[T]) center() Vec[T] {
	centerX := (ab.TopLeft.X + ab.BottomRight.X) / 2
	centerY := (ab.TopLeft.Y + ab.BottomRight.Y) / 2
	return Vec[T]{X: centerX, Y: centerY}
}

// Contains reports whether other lies entirely within axis-aligned bounding box.
func (ab AABB[T]) Contains(other AABB[T]) bool {
	return ab.TopLeft.X <= other.TopLeft.X &&
		ab.TopLeft.Y <= other.TopLeft.Y &&
		ab.BottomRight.X >= other.BottomRight.X &&
		ab.BottomRight.Y >= other.BottomRight.Y
}

func (ab AABB[T]) ContainsVec(vec Vec[T]) bool {
	return vec.X < ab.BottomRight.X && vec.X > ab.TopLeft.X &&
		vec.Y < ab.TopLeft.Y && vec.Y > ab.BottomRight.Y
}

// Intersects reports whether this AABB overlaps another axis-aligned bounding box.
// It returns true both when the AABBs share any interior volume
// and when they only touch along edges or vertices.
func (ab AABB[T]) Intersects(other AABB[T]) bool {
	// x axis check
	xIntersects := ab.axisIntersection(other, func(v Vec[T]) T { return v.X })
	if !xIntersects {
		return false
	}
	// y axis check
	yIntersects := ab.axisIntersection(other, func(v Vec[T]) T { return v.Y })
	return yIntersects
}

func (ab AABB[T]) axisIntersection(other AABB[T], axisValue func(Vec[T]) T) bool {
	d := ab.AxisDistanceTo(other, axisValue)
	return float64(d) <= eps
}

// AxisDistanceTo returns the gap between tow given AABBs on the axis selected by axisValue.
func (ab AABB[T]) AxisDistanceTo(
	other AABB[T],
	axisValue func(Vec[T]) T,
) T {
	ab, other = SortAABBsBy(
		ab, other,
		func(box AABB[T]) T { return axisValue(box.TopLeft) },
		func(box AABB[T]) T { return axisValue(box.BottomRight) },
	)

	if axisValue(ab.BottomRight) >= axisValue(other.TopLeft) {
		return 0
	}
	return axisValue(other.TopLeft) - axisValue(ab.BottomRight)
}

// SortAABBsBy orders two AABBs using the provided key functions and returns them as (min,max).
func SortAABBsBy[T Numeric](a, b AABB[T], keyFns ...func(AABB[T]) T) (aa, bb AABB[T]) {
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

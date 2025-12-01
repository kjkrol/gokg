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
func (ab AABB[T]) Split() [4]AABB[T] {
	w := ab.BottomRight.X - ab.TopLeft.X
	h := ab.BottomRight.Y - ab.TopLeft.Y
	hw := w / 2
	hh := h / 2

	midX := ab.TopLeft.X + hw
	midY := ab.TopLeft.Y + hh

	topLeft := ab.TopLeft
	topRightPos := NewVec(midX, ab.TopLeft.Y)
	botLeftPos := NewVec(ab.TopLeft.X, midY)
	botRightPos := NewVec(midX, midY)

	return [4]AABB[T]{
		NewAABBAt(topLeft, w/2, h/2),
		NewAABBAt(topRightPos, w/2, h/2),
		NewAABBAt(botLeftPos, w/2, h/2),
		NewAABBAt(botRightPos, w/2, h/2),
	}
}

// Contains reports whether other lies entirely within axis-aligned bounding box.
func (ab AABB[T]) Contains(other AABB[T]) bool {
	return ab.TopLeft.X <= other.TopLeft.X &&
		ab.TopLeft.Y <= other.TopLeft.Y &&
		ab.BottomRight.X >= other.BottomRight.X &&
		ab.BottomRight.Y >= other.BottomRight.Y
}

func (ab AABB[T]) ContainsVec(vec Vec[T]) bool {
	return vec.X > ab.TopLeft.X && vec.X < ab.BottomRight.X &&
		vec.Y > ab.TopLeft.Y && vec.Y < ab.BottomRight.Y
}

// Intersects reports whether this AABB overlaps another axis-aligned bounding box.
// It returns true both when the AABBs share any interior volume
// and when they only touch along edges or vertices.
// Intersects returns true when AABBs overlap or touch edges.
func (ab AABB[T]) Intersects(other AABB[T]) bool {
	return ab.TopLeft.X <= other.BottomRight.X &&
		ab.BottomRight.X >= other.TopLeft.X &&
		ab.TopLeft.Y <= other.BottomRight.Y &&
		ab.BottomRight.Y >= other.TopLeft.Y
}

// AxisDistanceTo returns the gap between tow given AABBs on the axis selected by axisValue.
func (ab AABB[T]) AxisDistanceX(other AABB[T]) T {
	return axisDistance1D(ab.TopLeft.X, ab.BottomRight.X, other.TopLeft.X, other.BottomRight.X)
}

func (ab AABB[T]) AxisDistanceY(other AABB[T]) T {
	return axisDistance1D(ab.TopLeft.Y, ab.BottomRight.Y, other.TopLeft.Y, other.BottomRight.Y)
}

func axisDistance1D[T Numeric](aMin, aMax, bMin, bMax T) T {
	// a: [aMin, aMax]
	// b: [bMin, bMax]

	// a całkowicie przed b
	if aMax < bMin {
		return bMin - aMax
	}

	// b całkowicie przed a
	if bMax < aMin {
		return aMin - bMax
	}

	// nachodzą na siebie lub stykają się
	return 0
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

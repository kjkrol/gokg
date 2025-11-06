package geometry

import "fmt"

// BoundingBox is a minimal axis-aligned rectangle defined by its top-left and bottom-right corners.
type BoundingBox[T SupportedNumeric] struct {
	TopLeft     Vec[T]
	BottomRight Vec[T]
}

// NewBoundingBox constructs a BoundingBox from explicit corner vectors.
func NewBoundingBox[T SupportedNumeric](topLeft, bottomRight Vec[T]) BoundingBox[T] {
	return BoundingBox[T]{
		TopLeft:     topLeft,
		BottomRight: bottomRight,
	}
}

// NewBoundingBoxAt builds a BoundingBox starting at pos with the provided width and height.
func NewBoundingBoxAt[T SupportedNumeric](pos Vec[T], width, height T) BoundingBox[T] {
	bottomRight := NewVec(pos.X+width, pos.Y+height)
	return NewBoundingBox(pos, bottomRight)
}

// NewBoundingBoxAround creates a BoundingBox centered at center with half-size d on each axis.
func NewBoundingBoxAround[T SupportedNumeric](center Vec[T], d T) BoundingBox[T] {
	topLeft := Vec[T]{X: center.X - d, Y: center.Y - d}
	return NewBoundingBoxAt(topLeft, 2*d, 2*d)
}

// Equals reports whether ab and other share the same corners.
func (b BoundingBox[T]) Equals(other BoundingBox[T]) bool {
	return b.TopLeft == other.TopLeft && b.BottomRight == other.BottomRight
}

// String formats the box using its top-left and bottom-right corners.
func (b BoundingBox[T]) String() string {
	return fmt.Sprintf("{%v %v}", b.TopLeft, b.BottomRight)
}

// Split divides the box into four equal quadrants around its center.
func (ab *BoundingBox[T]) Split() [4]BoundingBox[T] {
	center := ab.center()
	width_half := (ab.BottomRight.X - ab.TopLeft.X) / 2
	height_half := (ab.BottomRight.Y - ab.TopLeft.Y) / 2
	return [4]BoundingBox[T]{
		NewBoundingBoxAt(ab.TopLeft, width_half, height_half),                     // top left
		NewBoundingBoxAt(NewVec(center.X, ab.TopLeft.Y), width_half, height_half), // top right
		NewBoundingBoxAt(NewVec(ab.TopLeft.X, center.Y), width_half, height_half), // bottom left
		NewBoundingBoxAt(center, width_half, height_half),                         // bottom right
	}
}

func (b BoundingBox[T]) center() Vec[T] {
	centerX := (b.TopLeft.X + b.BottomRight.X) / 2
	centerY := (b.TopLeft.Y + b.BottomRight.Y) / 2
	return Vec[T]{X: centerX, Y: centerY}
}

// Contains reports whether other lies entirely within ab.
func (b BoundingBox[T]) Contains(other BoundingBox[T]) bool {
	return other.TopLeft.X >= b.TopLeft.X &&
		other.TopLeft.Y >= b.TopLeft.Y &&
		other.BottomRight.X <= b.BottomRight.X &&
		other.BottomRight.Y <= b.BottomRight.Y
}

// Intersects reports whether this AABB overlaps another bounding box.
// It returns true both when the boxes share any interior volume
// and when they only touch along edges or vertices.
func (b BoundingBox[T]) Intersects(other BoundingBox[T]) bool {
	// x axis check
	xIntersects := b.axisIntersection(other, func(v Vec[T]) T { return v.X })
	if !xIntersects {
		return false
	}
	// y axis check
	yIntersects := b.axisIntersection(other, func(v Vec[T]) T { return v.Y })
	return yIntersects
}

func (b BoundingBox[T]) axisIntersection(other BoundingBox[T], axisValue func(Vec[T]) T) bool {
	eps := VectorMathByType[T]().OverlapEpsilon()
	return b.axisDistanceTo(other, axisValue) <= eps
}

// AxisDistanceTo returns the gap between ab and other on the axis selected by axisValue.
func (b BoundingBox[T]) axisDistanceTo(
	other BoundingBox[T],
	axisValue func(Vec[T]) T,
) T {
	b, other = SortBoxesBy(
		b, other,
		func(box BoundingBox[T]) T { return axisValue(box.TopLeft) },
		func(box BoundingBox[T]) T { return axisValue(box.BottomRight) },
	)

	if axisValue(b.BottomRight) >= axisValue(other.TopLeft) {
		return 0
	}
	return axisValue(other.TopLeft) - axisValue(b.BottomRight)
}

// SortBoxesBy orders two boxes using the provided key functions and returns them as (min,max).
func SortBoxesBy[T SupportedNumeric](a, b BoundingBox[T], keyFns ...func(BoundingBox[T]) T) (aa, bb BoundingBox[T]) {
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

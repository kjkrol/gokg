package geometry

import "fmt"

// Rectangle represents an axis-aligned area in a 2D space defined by its top-left and bottom-right corners.
// It also includes the center point of the rectangle for convenience.
//
// Fields:
// - TopLeft: The top-left corner of the rectangle as a vector.
// - BottomRight: The bottom-right corner of the rectangle as a vector.
// - Center: The center point of the rectangle as a vector.
type Rectangle[T SupportedNumeric] struct {
	TopLeft     Vec[T]
	BottomRight Vec[T]
	Center      Vec[T]
	fragments   []Spatial[T]
}

func NewRectangle[T SupportedNumeric](topLeft Vec[T], bottomRight Vec[T]) Rectangle[T] {
	centerX := (topLeft.X + bottomRight.X) / 2
	centerY := (topLeft.Y + bottomRight.Y) / 2
	center := Vec[T]{X: centerX, Y: centerY}
	return Rectangle[T]{TopLeft: topLeft, BottomRight: bottomRight, Center: center}
}

func BuildRectangle[T SupportedNumeric](center Vec[T], d T) Rectangle[T] {
	topLeft := Vec[T]{X: center.X - d, Y: center.Y - d}
	bottomRight := Vec[T]{X: center.X + d, Y: center.Y + d}
	return NewRectangle(topLeft, bottomRight)
}

// Bounds returns the rectangle itself.
func (r Rectangle[T]) Bounds() Rectangle[T] {
	return r
}

// Probe expands the rectangle by the given margin and wraps it for cyclic planes.
func (r Rectangle[T]) Probe(margin T, plane Plane[T]) []Rectangle[T] {
	probe := r.Expand(margin)
	rectangles := []Rectangle[T]{probe}
	if plane.Name() == "cyclic" {
		rectangles = append(rectangles, WrapRectangleCyclic(probe, plane.Size(), plane.Contains)...)
	}
	return rectangles
}

// DistanceTo computes the distance to another spatial object using their bounds.
func (r Rectangle[T]) DistanceTo(other Spatial[T], metric func(Vec[T], Vec[T]) T) T {
	return aabbDistance(&r, other, metric)
}

// Vertices returns the corners that define the rectangle.
func (r *Rectangle[T]) Vertices() []*Vec[T] {
	return []*Vec[T]{&r.TopLeft, &r.BottomRight}
}

func (r Rectangle[T]) Fragments() []Spatial[T] { return r.fragments }

func (r *Rectangle[T]) SetFragments(f []Spatial[T]) { r.fragments = f }

func (r *Rectangle[T]) Split() [4]Rectangle[T] {
	return [4]Rectangle[T]{
		NewRectangle(r.TopLeft, r.Center), // top left
		NewRectangle(Vec[T]{X: r.Center.X, Y: r.TopLeft.Y}, Vec[T]{X: r.BottomRight.X, Y: r.Center.Y}), // top right
		NewRectangle(Vec[T]{X: r.TopLeft.X, Y: r.Center.Y}, Vec[T]{X: r.Center.X, Y: r.BottomRight.Y}), // bottom left
		NewRectangle(r.Center, r.BottomRight), // bottom right
	}
}

func (r Rectangle[T]) Contains(other Rectangle[T]) bool {
	return other.TopLeft.X >= r.TopLeft.X &&
		other.TopLeft.Y >= r.TopLeft.Y &&
		other.BottomRight.X <= r.BottomRight.X &&
		other.BottomRight.Y <= r.BottomRight.Y
}

func (r Rectangle[T]) Expand(margin T) Rectangle[T] {
	return NewRectangle(
		Vec[T]{X: r.TopLeft.X - margin, Y: r.TopLeft.Y - margin},
		Vec[T]{X: r.BottomRight.X + margin, Y: r.BottomRight.Y + margin},
	)
}

//-------------------------------------------------------------------------

func (r Rectangle[T]) Intersects(other Rectangle[T]) bool {
	// x axis check
	xIntersects := axisIntersection(r, other, func(v Vec[T]) T { return v.X })
	if !xIntersects {
		return false
	}
	// y axis check
	yIntersects := axisIntersection(r, other, func(v Vec[T]) T { return v.Y })
	return yIntersects
}

func axisIntersection[T SupportedNumeric](aa, bb Rectangle[T], axisValue func(Vec[T]) T) bool {
	aa, bb = SortRectanglesBy(aa, bb, axisValue)
	noIntersection := axisValue(aa.TopLeft) < axisValue(bb.BottomRight) && axisValue(aa.BottomRight) < axisValue(bb.TopLeft)
	return !noIntersection
}

func SortRectanglesBy[T SupportedNumeric](a, b Rectangle[T], axisValue func(Vec[T]) T) (aa, bb Rectangle[T]) {
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

func (r Rectangle[T]) IntersectsAny(others []Rectangle[T]) bool {
	for _, wrapped := range others {
		if r.Intersects(wrapped) {
			return true
		}
	}
	return false
}

func WrapRectangleCyclic[T SupportedNumeric](
	r Rectangle[T],
	size Vec[T],
	contains func(Vec[T]) bool,
) []Rectangle[T] {
	var wrappedRectangles []Rectangle[T]

	// Predefined offset values for wrapping
	offsets := []Vec[T]{
		{X: size.X, Y: 0},      // Shift right
		{X: 0, Y: size.Y},      // Shift down
		{X: size.X, Y: size.Y}, // Shift right-down
	}

	vecMath := VectorMathByType[T]()

	// Generate wrapped versions for each offset
	for _, offset := range offsets {
		wrapped := Rectangle[T]{
			TopLeft:     r.TopLeft,
			BottomRight: r.BottomRight,
			Center:      r.Center,
		}
		vecMath.Wrap(&wrapped.TopLeft, offset)
		vecMath.Wrap(&wrapped.BottomRight, offset)
		vecMath.Wrap(&wrapped.Center, offset)

		if contains(wrapped.TopLeft) || contains(wrapped.BottomRight) {
			wrappedRectangles = append(wrappedRectangles, wrapped)
		}
	}

	return wrappedRectangles
}

//-------------------------------------------------------------------------

func (r Rectangle[T]) String() string {
	return fmt.Sprintf("{%v %v %v}", r.TopLeft, r.BottomRight, r.Center)
}

package geometry

const (
	BOUNDED = "bounded"
	CYCLIC  = "cyclic"
)

// Plane encapsulates a 2D surface with its own metric and boundary behaviour.
type Plane[T SupportedNumeric] struct {
	size       Vec[T]
	vectorMath VectorMath[T]
	normalize  func(*Vec[T])
	metric     func(v1, v2 Vec[T]) T
	name       string
	viewport   AABB[T]
}

// -----------------------------------------------------------------------------

// Size returns the plane width and height as a vector.
func (p Plane[T]) Size() Vec[T] { return p.size }

// Metric measures the distance between v1 and v2 using the plane-specific metric.
func (p Plane[T]) Metric(v1, v2 Vec[T]) T { return p.metric(v1, v2) }

// Contains reports whether vec lies within the plane's bounds.
func (p Plane[T]) Contains(vec Vec[T]) bool {
	return vec.X >= 0 && vec.X < p.size.X && vec.Y >= 0 && vec.Y < p.size.Y
}

// Expand grows the bounding box by margin and normalises it to the plane.
func (p Plane[T]) Expand(ab *AABB[T], margin T) {
	ab.TopLeft.AddMutable(NewVec(-margin, -margin))
	ab.width = ab.width + 2*margin
	ab.height = ab.height + 2*margin

	p.NormalizeAABB(ab)
}

// Translate shifts the bounding box by delta and normalises it to the plane.
func (p Plane[T]) Translate(ab *AABB[T], delta Vec[T]) {
	ab.TopLeft.AddMutable(delta)

	p.NormalizeAABB(ab)
}

// Normalize maps vec into the plane domain according to its boundary rules.
func (p Plane[T]) Normalize(vec *Vec[T]) { p.normalize(vec) }

// NormalizeAABB maps ab into the plane domain, adjusting fragments for wrap-around.
func (p Plane[T]) NormalizeAABB(ab *AABB[T]) {

	switch p.name {
	case BOUNDED:
		prev := ab.TopLeft
		p.normalize(&ab.TopLeft)
		dx := prev.X - ab.TopLeft.X
		dy := prev.Y - ab.TopLeft.Y
		ab.BottomRight = NewVec(ab.TopLeft.X+ab.width+dx, ab.TopLeft.Y+ab.height+dy)
	case CYCLIC:
		p.normalize(&ab.TopLeft)
		ab.BottomRight = NewVec(ab.TopLeft.X+ab.width, ab.TopLeft.Y+ab.height)
		dx := p.size.X - ab.TopLeft.X - ab.width
		dy := p.size.Y - ab.TopLeft.Y - ab.height

		if dx < 0 {
			ab.BottomRight.X = p.size.X
		}
		if dy < 0 {
			ab.BottomRight.Y = p.size.Y
		}
		ab.fragmentation(dx, dy)
	}
}

func (p Plane[T]) relativeMetric(v1, v2 Vec[T]) T {
	delta := v1.Sub(v2)
	p.normalize(&delta)
	return p.vectorMath.Length(delta)
}

// Name reports the plane mode (bounded or cyclic).
func (p Plane[T]) Name() string { return p.name }

// Viewport returns the canonical AABB covering the entire plane.
func (p Plane[T]) Viewport() AABB[T] { return p.viewport }

// -----------------------------------------------------------------------------

// NewBoundedPlane constructs a plane that clamps vectors to the given width and height.
func NewBoundedPlane[T SupportedNumeric](sizeX, sizeY T) Plane[T] {
	plane := Plane[T]{
		name:       BOUNDED,
		size:       NewVec(sizeX, sizeY),
		vectorMath: VectorMathByType[T](),
	}
	plane.viewport = NewAABB(NewVec[T](0, 0), plane.size.X, plane.size.Y)
	plane.normalize = func(v *Vec[T]) { plane.vectorMath.Clamp(v, plane.size) }
	plane.metric = func(v1, v2 Vec[T]) T { return max(plane.relativeMetric(v1, v2), plane.relativeMetric(v2, v1)) }
	return plane
}

// -----------------------------------------------------------------------------

// NewCyclicBoundedPlane constructs a plane with wrap-around behaviour on both axes.
func NewCyclicBoundedPlane[T SupportedNumeric](sizeX, sizeY T) Plane[T] {
	plane := Plane[T]{
		name:       CYCLIC,
		size:       NewVec(sizeX, sizeY),
		vectorMath: VectorMathByType[T](),
	}
	plane.viewport = NewAABB(NewVec[T](0, 0), plane.size.X, plane.size.Y)
	plane.normalize = func(v *Vec[T]) { plane.vectorMath.Wrap(v, plane.size) }
	plane.metric = func(v1, v2 Vec[T]) T { return min(plane.relativeMetric(v1, v2), plane.relativeMetric(v2, v1)) }
	return plane
}

// -----------------------------------------------------------------------------

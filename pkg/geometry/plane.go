package geometry

const (
	BOUNDED = "bounded"
	CYCLIC  = "cyclic"
)

// Plane encapsulates a 2D surface with its own metric and boundary behaviour.
type Plane[T SupportedNumeric] struct {
	size         Vec[T]
	vectorMath   VectorMath[T]
	normalize    func(*Vec[T])
	normalizeBox func(*PlaneBox[T])
	metric       func(v1, v2 Vec[T]) T
	name         string
	viewport     BoundingBox[T]
}

// -----------------------------------------------------------------------------

// Size returns the plane width and height as a vector.
func (p Plane[T]) Size() Vec[T] { return p.size }

// Metric measures the distance between v1 and v2 using the plane-specific metric.
func (p Plane[T]) Metric(v1, v2 Vec[T]) T { return p.metric(v1, v2) }

// Expand grows the bounding box by margin and normalises it to the plane.
func (p Plane[T]) Expand(ab *PlaneBox[T], margin T) {
	ab.TopLeft.AddMutable(NewVec(-margin, -margin))
	ab.size.AddMutable(NewVec(2*margin, 2*margin))
	p.normalizeBox(ab)
}

// Translate shifts the bounding box by delta and normalises it to the plane.
func (p Plane[T]) Translate(ab *PlaneBox[T], delta Vec[T]) {
	ab.TopLeft.AddMutable(delta)
	p.normalizeBox(ab)
}

// WrapBoundingBox converts a world-space BoundingBox into a PlaneBox normalized to this Plane.
func (p Plane[T]) WrapBoundingBox(box BoundingBox[T]) PlaneBox[T] {
	width := box.BottomRight.X - box.TopLeft.X
	height := box.BottomRight.Y - box.TopLeft.Y
	planeBox := newPlaneBox(box.TopLeft, width, height)
	p.normalizeBox(&planeBox)
	return planeBox
}

// WrapVec treats the point as a zero-area box and returns its Plane-normalized PlaneBox representation.
func (p Plane[T]) WrapVec(vec Vec[T]) PlaneBox[T] {
	box := NewBoundingBoxAt(vec, 0, 0)
	return p.WrapBoundingBox(box)
}

// Name reports the plane mode (bounded or cyclic).
func (p Plane[T]) Name() string { return p.name }

// Viewport returns the canonical PlaneBox (bounding-box) covering the entire plane.
func (p Plane[T]) Viewport() BoundingBox[T] { return p.viewport }

// -----------------------------------------------------------------------------

// NewBoundedPlane constructs a plane that clamps vectors to the given width and height.
func NewBoundedPlane[T SupportedNumeric](sizeX, sizeY T) Plane[T] {
	plane := Plane[T]{
		name:       BOUNDED,
		size:       NewVec(sizeX, sizeY),
		vectorMath: VectorMathByType[T](),
	}
	plane.viewport = NewBoundingBoxAt(NewVec[T](0, 0), plane.size.X, plane.size.Y)
	plane.normalize = func(v *Vec[T]) { plane.vectorMath.Clamp(v, plane.size) }
	plane.metric = func(v1, v2 Vec[T]) T { return max(plane.relativeMetric(v1, v2), plane.relativeMetric(v2, v1)) }
	plane.normalizeBox = func(pb *PlaneBox[T]) {
		pb.BottomRight = pb.TopLeft.Add(pb.size)
		plane.vectorMath.Clamp(&pb.BottomRight, plane.size)
		plane.normalize(&pb.TopLeft)
	}
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
	plane.viewport = NewBoundingBoxAt(NewVec[T](0, 0), plane.size.X, plane.size.Y)
	plane.normalize = func(v *Vec[T]) { plane.vectorMath.Wrap(v, plane.size) }
	plane.metric = func(v1, v2 Vec[T]) T { return min(plane.relativeMetric(v1, v2), plane.relativeMetric(v2, v1)) }
	plane.normalizeBox = func(pb *PlaneBox[T]) {
		plane.normalize(&pb.TopLeft)

		pb.BottomRight = pb.TopLeft.Add(pb.size)
		plane.vectorMath.Clamp(&pb.BottomRight, plane.size)

		d := plane.size.Sub(pb.TopLeft).Sub(pb.size)
		pb.fragmentation(d.X, d.Y)
	}
	return plane
}

// -----------------------------------------------------------------------------

func (p Plane[T]) relativeMetric(v1, v2 Vec[T]) T {
	delta := v1.Sub(v2)
	p.normalize(&delta)
	return p.vectorMath.Length(delta)
}

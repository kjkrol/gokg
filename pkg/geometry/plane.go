package geometry

const (
	BOUNDED = "bounded"
	CYCLIC  = "cyclic"
)

type Metric[T SupportedNumeric] func(v1, v2 Vec[T]) T

// Plane encapsulates a 2D surface with its own metric and boundary behaviour.
type Plane[T SupportedNumeric] struct {
	size                Vec[T]
	vectorMath          VectorMath[T]
	normalizeVec        func(*Vec[T])
	normalizeBox        func(*PlaneBox[T])
	metric              Metric[T]
	name                string
	viewport            BoundingBox[T]
	boundingBoxDistance BoundingBoxDistance[T]
}

// -----------------------------------------------------------------------------

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

// BoundingBoxDistance measures the distance between aa and bb using the plane-specific metric.
func (p Plane[T]) BoundingBoxDistance(aa, bb BoundingBox[T]) T {
	return p.boundingBoxDistance(aa, bb)
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

// Size returns the plane width and height as a vector.
func (p Plane[T]) Size() Vec[T] { return p.size }

// Viewport returns the canonical PlaneBox (bounding-box) covering the entire plane.
func (p Plane[T]) Viewport() BoundingBox[T] { return p.viewport }

// -----------------------------------------------------------------------------

// NewBoundedPlane constructs a plane that clamps vectors to the given width and height.
func NewBoundedPlane[T SupportedNumeric](sizeX, sizeY T) Plane[T] {
	return newPlane(BOUNDED, sizeX, sizeY, func(p *Plane[T]) {
		p.normalizeVec = func(v *Vec[T]) { p.vectorMath.Clamp(v, p.size) }
		p.normalizeBox = func(pb *PlaneBox[T]) {
			p.normalizePlaneBoxBottomRight(pb)
			p.normalizePlaneBoxTopLeft(pb)
		}
		p.metric = func(v1, v2 Vec[T]) T { return max(p.relativeMetric(v1, v2), p.relativeMetric(v2, v1)) }
	})
}

// -----------------------------------------------------------------------------

// NewCyclicBoundedPlane constructs a plane with wrap-around behaviour on both axes.
func NewCyclicBoundedPlane[T SupportedNumeric](sizeX, sizeY T) Plane[T] {
	return newPlane(CYCLIC, sizeX, sizeY, func(p *Plane[T]) {
		p.normalizeVec = func(v *Vec[T]) { p.vectorMath.Wrap(v, p.size) }
		p.normalizeBox = func(pb *PlaneBox[T]) {
			p.normalizePlaneBoxTopLeft(pb)
			if p.normalizePlaneBoxBottomRight(pb) {
				d := p.size.Sub(pb.TopLeft).Sub(pb.size)
				pb.fragmentation(d.X, d.Y)
			}
		}
		p.metric = func(v1, v2 Vec[T]) T { return min(p.relativeMetric(v1, v2), p.relativeMetric(v2, v1)) }
	})
}

// -----------------------------------------------------------------------------

func newPlane[T SupportedNumeric](name string, sizeX, sizeY T, setup func(p *Plane[T])) Plane[T] {
	plane := Plane[T]{
		name:       name,
		size:       NewVec(sizeX, sizeY),
		vectorMath: VectorMathByType[T](),
		viewport:   NewBoundingBoxAt(NewVec[T](0, 0), sizeX, sizeY),
	}
	setup(&plane)
	plane.boundingBoxDistance = newBoundingBoxDistance(plane.metric)
	return plane
}

func (p Plane[T]) relativeMetric(v1, v2 Vec[T]) T {
	delta := v1.Sub(v2)
	p.normalizeVec(&delta)
	return p.vectorMath.Length(delta)
}

func (p Plane[T]) normalizePlaneBoxTopLeft(pb *PlaneBox[T]) {
	if !p.Viewport().ContainsVec(pb.BoundingBox.TopLeft) {
		p.normalizeVec(&pb.TopLeft)
	}
}

func (p Plane[T]) normalizePlaneBoxBottomRight(pb *PlaneBox[T]) bool {
	pb.BottomRight = pb.TopLeft.Add(pb.size)
	if !p.Viewport().ContainsVec(pb.BoundingBox.BottomRight) {
		p.vectorMath.Clamp(&pb.BottomRight, p.size)
		return true
	}
	return false
}

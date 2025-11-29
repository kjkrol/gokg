package plane

import "github.com/kjkrol/gokg/pkg/geom"

const (
	modeCartesian = "cartesian"
	modeTorus     = "torus"
)

type (
	Space[T geom.Numeric] interface {
		WrapAABB(box geom.AABB[T]) AABB[T]
		WrapVec(vec geom.Vec[T]) AABB[T]
		Expand(ab *AABB[T], margin T)
		Translate(ab *AABB[T], delta geom.Vec[T])
		AABBDistance() AABBDistance[T]
		Name() string
		Viewport() geom.AABB[T]
	}

	Metric[T geom.Numeric] func(v1, v2 geom.Vec[T]) T
)

// -----------------------------------------------------------------------------

// NewCartesian constructs a plane that clamps vectors to the given width and height.
func NewCartesian[T geom.Numeric](sizeX, sizeY T) Space[T] {
	return newSpace2d(modeCartesian, sizeX, sizeY, func(p *space2d[T]) {
		p.normalizeVec = func(v *geom.Vec[T]) { p.vectorMath.Clamp(v, p.size) }
		p.normalizeBox = func(pb *AABB[T]) {
			p.normalizePlaneBoxBottomRight(pb)
			p.normalizePlaneBoxTopLeft(pb)
		}
		p.metric = func(v1, v2 geom.Vec[T]) T { return max(p.relativeMetric(v1, v2), p.relativeMetric(v2, v1)) }
	})
}

// -----------------------------------------------------------------------------

// NewTorus constructs a plane with wrap-around behaviour on both axes.
func NewTorus[T geom.Numeric](sizeX, sizeY T) Space[T] {
	return newSpace2d(string(modeTorus), sizeX, sizeY, func(p *space2d[T]) {
		p.normalizeVec = func(v *geom.Vec[T]) { p.vectorMath.Wrap(v, p.size) }
		p.normalizeBox = func(pb *AABB[T]) {
			p.normalizePlaneBoxTopLeft(pb)
			if p.normalizePlaneBoxBottomRight(pb) {
				step1 := p.vectorMath.Sub(p.size, pb.TopLeft)
				step2 := p.vectorMath.Sub(step1, pb.size)

				d := step2
				pb.fragmentation(d.X, d.Y)
			}
		}
		p.metric = func(v1, v2 geom.Vec[T]) T { return min(p.relativeMetric(v1, v2), p.relativeMetric(v2, v1)) }
	})
}

// -----------------------------------------------------------------------------

// space2d encapsulates a 2D surface with its own metric and boundary behaviour.
type space2d[T geom.Numeric] struct {
	size         geom.Vec[T]
	vectorMath   geom.VectorMath[T]
	normalizeVec func(*geom.Vec[T])
	normalizeBox func(*AABB[T])
	metric       Metric[T]
	name         string
	viewport     geom.AABB[T]
}

// WrapAABB converts a world-space AABB into a AABB normalized to this Plane.
func (p space2d[T]) WrapAABB(box geom.AABB[T]) AABB[T] {
	width := box.BottomRight.X - box.TopLeft.X
	height := box.BottomRight.Y - box.TopLeft.Y
	planeBox := newAABB(box.TopLeft, width, height)
	p.normalizeBox(&planeBox)
	return planeBox
}

// WrapVec treats the point as a zero-area box and returns its Plane-normalized PlaneBox representation.
func (p space2d[T]) WrapVec(vec geom.Vec[T]) AABB[T] {
	box := geom.NewAABBAt(vec, 0, 0)
	return p.WrapAABB(box)
}

// Expand grows the bounding box by margin and normalises it to the plane.
func (p space2d[T]) Expand(ab *AABB[T], margin T) {
	ab.TopLeft.AddMutable(geom.NewVec(-margin, -margin))
	ab.size.AddMutable(geom.NewVec(2*margin, 2*margin))
	p.normalizeBox(ab)
}

// Translate shifts the bounding box by delta and normalises it to the plane.
func (p space2d[T]) Translate(ab *AABB[T], delta geom.Vec[T]) {
	ab.TopLeft.AddMutable(delta)
	p.normalizeBox(ab)
}

// AABBDistance measures the distance between aa and bb using the plane-specific metric.
func (p space2d[T]) AABBDistance() AABBDistance[T] {
	return newAABBDistance(p.metric)
}

// Name reports the plane mode (bounded or cyclic).
func (p space2d[T]) Name() string { return p.name }

// Viewport returns the canonical PlaneBox (bounding-box) covering the entire plane.
func (p space2d[T]) Viewport() geom.AABB[T] { return p.viewport }

func newSpace2d[T geom.Numeric](name string, sizeX, sizeY T, setup func(p *space2d[T])) space2d[T] {
	plane := space2d[T]{
		name:       name,
		size:       geom.NewVec(sizeX, sizeY),
		vectorMath: geom.VectorMathByType[T](),
		viewport:   geom.NewAABBAt(geom.NewVec[T](0, 0), sizeX, sizeY),
	}
	setup(&plane)
	return plane
}

func (p space2d[T]) relativeMetric(v1, v2 geom.Vec[T]) T {
	delta := p.vectorMath.Sub(v1, v2)
	p.normalizeVec(&delta)
	return p.vectorMath.Length(delta)
}

func (p space2d[T]) normalizePlaneBoxTopLeft(pb *AABB[T]) {
	if !p.Viewport().ContainsVec(pb.AABB.TopLeft) {
		p.normalizeVec(&pb.TopLeft)
	}
}

func (p space2d[T]) normalizePlaneBoxBottomRight(pb *AABB[T]) bool {
	pb.BottomRight = pb.TopLeft.Add(pb.size)
	if !p.Viewport().ContainsVec(pb.AABB.BottomRight) {
		p.vectorMath.Clamp(&pb.BottomRight, p.size)
		return true
	}
	return false
}

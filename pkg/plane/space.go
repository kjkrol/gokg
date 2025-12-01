package plane

import "github.com/kjkrol/gokg/pkg/geom"

const (
	modeEuclidean2D = "Euclidean2D"
	modeToroidal2D  = "Toroidal2D"
)

type (
	Space2D[T geom.Numeric] interface {
		WrapAABB(aabb geom.AABB[T]) AABB[T]
		WrapVec(vec geom.Vec[T]) AABB[T]
		Expand(aabb *AABB[T], margin T)
		Translate(aabb *AABB[T], delta geom.Vec[T])
		AABBDistance() AABBDistance[T]
		Name() string
		Viewport() geom.AABB[T]
	}

	Metric[T geom.Numeric] func(vec1, vec2 geom.Vec[T]) T
)

// -----------------------------------------------------------------------------

// NewEuclidean2D constructs a 2D space that clamps vectors to the given width and height.
func NewEuclidean2D[T geom.Numeric](sizeX, sizeY T) Space2D[T] {
	return newSpace2d(modeEuclidean2D, sizeX, sizeY, func(space *space2d[T]) {
		space.normalizeVec = func(vec *geom.Vec[T]) { space.vectorMath.Clamp(vec, space.size) }
		space.normalizeAABB = func(aabb *AABB[T]) {
			space.normalizeAABBBottomRight(aabb)
			space.normalizeAABBTopLeft(aabb)
		}
		space.metric = func(vec1, vec2 geom.Vec[T]) T {
			return max(space.relativeMetric(vec1, vec2), space.relativeMetric(vec2, vec1))
		}
	})
}

// -----------------------------------------------------------------------------

// NewToroidal2D constructs a 2D space with wrap-around behaviour on both axes.
func NewToroidal2D[T geom.Numeric](sizeX, sizeY T) Space2D[T] {
	return newSpace2d(modeToroidal2D, sizeX, sizeY, func(space *space2d[T]) {
		space.normalizeVec = func(vec *geom.Vec[T]) { space.vectorMath.Wrap(vec, space.size) }
		space.normalizeAABB = func(aabb *AABB[T]) {
			space.normalizeAABBTopLeft(aabb)
			dx, dy := space.normalizeAABBBottomRight(aabb)
			aabb.fragmentation(dx, dy)
		}
		space.metric = func(vec1, vec2 geom.Vec[T]) T {
			return min(space.relativeMetric(vec1, vec2), space.relativeMetric(vec2, vec1))
		}
	})
}

// -----------------------------------------------------------------------------

// space2d encapsulates a 2D surface with its own metric and boundary behaviour.
type space2d[T geom.Numeric] struct {
	size          geom.Vec[T]
	vectorMath    geom.VectorMath[T]
	normalizeVec  func(*geom.Vec[T])
	normalizeAABB func(*AABB[T])
	metric        Metric[T]
	name          string
	viewport      geom.AABB[T]
}

// WrapAABB converts a world-space AABB into a AABB normalized to this Space.
func (s space2d[T]) WrapAABB(aabb geom.AABB[T]) AABB[T] {
	width := aabb.BottomRight.X - aabb.TopLeft.X
	height := aabb.BottomRight.Y - aabb.TopLeft.Y
	wrappedAABB := newAABB(aabb.TopLeft, width, height)
	s.normalizeAABB(&wrappedAABB)
	return wrappedAABB
}

// WrapVec treats the point as a zero-area box and returns its Space-normalized AABB representation.
func (s space2d[T]) WrapVec(vec geom.Vec[T]) AABB[T] {
	aabb := geom.NewAABBAt(vec, 0, 0)
	return s.WrapAABB(aabb)
}

// Expand grows the axis-aligned bounding box by margin and normalises it to the plane.
func (s space2d[T]) Expand(aabb *AABB[T], margin T) {
	aabb.TopLeft.AddMutable(geom.NewVec(-margin, -margin))
	aabb.size.AddMutable(geom.NewVec(2*margin, 2*margin))
	s.normalizeAABB(aabb)
}

// Translate shifts the axis-aligned bounding box by delta and normalises it to the plane.
func (s space2d[T]) Translate(aabb *AABB[T], delta geom.Vec[T]) {
	aabb.TopLeft.AddMutable(delta)
	s.normalizeAABB(aabb)
}

// AABBDistance measures the distance between aa and bb using the plane-specific metric.
func (s space2d[T]) AABBDistance() AABBDistance[T] {
	return newAABBDistance(s.metric)
}

// Name reports the space mode (Euclidean2D or Toroidal2D).
func (s space2d[T]) Name() string { return s.name }

// Viewport returns the canonical axis-aligned bounding box covering the entire plane.
func (s space2d[T]) Viewport() geom.AABB[T] { return s.viewport }

func newSpace2d[T geom.Numeric](name string, sizeX, sizeY T, setup func(p *space2d[T])) space2d[T] {
	space := space2d[T]{
		name:       name,
		size:       geom.NewVec(sizeX, sizeY),
		vectorMath: geom.VectorMathByType[T](),
		viewport:   geom.NewAABBAt(geom.NewVec[T](0, 0), sizeX, sizeY),
	}
	setup(&space)
	return space
}

func (s space2d[T]) relativeMetric(vec1, vec2 geom.Vec[T]) T {
	delta := s.vectorMath.Sub(vec1, vec2)
	s.normalizeVec(&delta)
	return s.vectorMath.Length(delta)
}

func (s space2d[T]) normalizeAABBTopLeft(aabb *AABB[T]) {
	if !s.Viewport().ContainsVec(aabb.AABB.TopLeft) {
		s.normalizeVec(&aabb.TopLeft)
	}
}

func (s space2d[T]) normalizeAABBBottomRight(aabb *AABB[T]) (dx T, dy T) {
	aabb.BottomRight = aabb.TopLeft.Add(aabb.size)
	dx = T(0)
	dy = T(0)
	if aabb.BottomRight.X > s.size.X {
		dx = aabb.BottomRight.X - s.size.X
	}
	if aabb.BottomRight.Y > s.size.Y {
		dy = aabb.BottomRight.Y - s.size.Y
	}
	s.vectorMath.Clamp(&aabb.BottomRight, s.size)
	return
}

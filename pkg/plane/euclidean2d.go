package plane

import "github.com/kjkrol/gokg/pkg/geom"

// NewEuclidean2D constructs a 2D space that clamps vectors to the given width and height.
func NewEuclidean2D[T geom.Numeric](sizeX, sizeY T) Space2D[T] {
	return &euclidean2d[T]{
		space2d: space2d[T]{
			size:       geom.NewVec(sizeX, sizeY),
			vectorMath: geom.VectorMathByType[T](),
			viewport:   geom.NewAABBAt(geom.NewVec[T](0, 0), sizeX, sizeY),
		},
	}
}

type euclidean2d[T geom.Numeric] struct{ space2d[T] }

func (s euclidean2d[T]) Name() string { return modeEuclidean2D }

func (s euclidean2d[T]) Viewport() geom.AABB[T] { return s.viewport }

func (s euclidean2d[T]) WrapAABB(aabb geom.AABB[T]) AABB[T] {
	width := aabb.BottomRight.X - aabb.TopLeft.X
	height := aabb.BottomRight.Y - aabb.TopLeft.Y
	wrappedAABB := newAABB(aabb.TopLeft, width, height)
	s.normalizeAABB(&wrappedAABB)
	return wrappedAABB
}

func (s euclidean2d[T]) WrapVec(vec geom.Vec[T]) AABB[T] {
	aabb := geom.NewAABBAt(vec, 0, 0)
	return s.WrapAABB(aabb)
}

func (s euclidean2d[T]) Expand(aabb *AABB[T], margin T) {
	aabb.TopLeft.AddMutable(geom.NewVec(-margin, -margin))
	aabb.size.AddMutable(geom.NewVec(2*margin, 2*margin))
	s.normalizeAABB(aabb)
}

func (s euclidean2d[T]) Translate(aabb *AABB[T], delta geom.Vec[T]) {
	aabb.TopLeft.AddMutable(delta)
	s.normalizeAABB(aabb)
}

func (s euclidean2d[T]) AABBDistance() AABBDistance[T] {
	return newAABBDistance(s.metric)
}

func (s euclidean2d[T]) normalizeVec(vec geom.Vec[T]) geom.Vec[T] {
	return s.vectorMath.Clamp(vec, s.size)
}

func (s euclidean2d[T]) normalizeAABB(aabb *AABB[T]) {
	s.normalizeAABBBottomRight(aabb)
	s.normalizeAABBTopLeft(aabb)
}

func (s euclidean2d[T]) normalizeAABBTopLeft(aabb *AABB[T]) {
	if !s.Viewport().ContainsVec(aabb.AABB.TopLeft) {
		aabb.TopLeft = s.normalizeVec(aabb.TopLeft)
	}
}

func (s euclidean2d[T]) normalizeAABBBottomRight(aabb *AABB[T]) (dx T, dy T) {
	aabb.BottomRight = aabb.TopLeft.Add(aabb.size)
	dx = T(0)
	dy = T(0)
	if aabb.BottomRight.X > s.size.X {
		dx = aabb.BottomRight.X - s.size.X
	}
	if aabb.BottomRight.Y > s.size.Y {
		dy = aabb.BottomRight.Y - s.size.Y
	}
	aabb.BottomRight = s.vectorMath.Clamp(aabb.BottomRight, s.size)
	return
}

func (s euclidean2d[T]) metric(vec1, vec2 geom.Vec[T]) T {
	dx := vec1.X
	if vec2.X > dx {
		dx = vec2.X - dx
	} else {
		dx = dx - vec2.X
	}
	dy := vec1.Y
	if vec2.Y > dy {
		dy = vec2.Y - dy
	} else {
		dy = dy - vec2.Y
	}
	delta := geom.NewVec(dx, dy)
	return s.vectorMath.Length(s.vectorMath.Clamp(delta, s.size))
}

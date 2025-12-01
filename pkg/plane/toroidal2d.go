package plane

import "github.com/kjkrol/gokg/pkg/geom"

// NewToroidal2D constructs a 2D space with wrap-around behaviour on both axes.
func NewToroidal2D[T geom.Numeric](sizeX, sizeY T) Space2D[T] {
	return &toroidal2d[T]{
		space2d: space2d[T]{
			size:       geom.NewVec(sizeX, sizeY),
			vectorMath: geom.VectorMathByType[T](),
			viewport:   geom.NewAABBAt(geom.NewVec[T](0, 0), sizeX, sizeY),
		},
	}
}

type toroidal2d[T geom.Numeric] struct{ space2d[T] }

func (s toroidal2d[T]) Name() string { return modeToroidal2D }

func (s toroidal2d[T]) Viewport() geom.AABB[T] { return s.viewport }

func (s toroidal2d[T]) WrapAABB(aabb geom.AABB[T]) AABB[T] {
	width := aabb.BottomRight.X - aabb.TopLeft.X
	height := aabb.BottomRight.Y - aabb.TopLeft.Y
	wrappedAABB := newAABB(aabb.TopLeft, width, height)
	s.normalizeAABB(&wrappedAABB)
	return wrappedAABB
}

func (s toroidal2d[T]) WrapVec(vec geom.Vec[T]) AABB[T] {
	aabb := geom.NewAABBAt(vec, 0, 0)
	return s.WrapAABB(aabb)
}

func (s toroidal2d[T]) Expand(aabb *AABB[T], margin T) {
	aabb.TopLeft.AddMutable(geom.NewVec(-margin, -margin))
	aabb.size.AddMutable(geom.NewVec(2*margin, 2*margin))
	s.normalizeAABB(aabb)
}

func (s toroidal2d[T]) Translate(aabb *AABB[T], delta geom.Vec[T]) {
	aabb.TopLeft.AddMutable(delta)
	s.normalizeAABB(aabb)
}

func (s toroidal2d[T]) AABBDistance() AABBDistance[T] {
	return newAABBDistance(s.metric)
}

func (s toroidal2d[T]) normalizeVec(vec geom.Vec[T]) geom.Vec[T] {
	return s.vectorMath.Wrap(vec, s.size)
}

func (s toroidal2d[T]) normalizeAABB(aabb *AABB[T]) {
	s.normalizeAABBTopLeft(aabb)
	dx, dy := s.normalizeAABBBottomRight(aabb)
	aabb.fragmentation(dx, dy)
}

func (s toroidal2d[T]) normalizeAABBTopLeft(aabb *AABB[T]) {
	aabb.TopLeft = s.normalizeVec(aabb.TopLeft)
}

func (s toroidal2d[T]) normalizeAABBBottomRight(aabb *AABB[T]) (dx T, dy T) {
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

func (s toroidal2d[T]) metric(vec1, vec2 geom.Vec[T]) T {
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

	delta := s.vectorMath.Wrap(geom.NewVec(dx, dy), s.size)

	if s.size.X != 0 {
		alt := s.size.X - delta.X
		if alt < delta.X {
			delta.X = alt
		}
	}
	if s.size.Y != 0 {
		alt := s.size.Y - delta.Y
		if alt < delta.Y {
			delta.Y = alt
		}
	}

	return s.vectorMath.Length(delta)
}

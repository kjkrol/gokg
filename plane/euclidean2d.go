package plane

import "github.com/kjkrol/gokg/geom"

// NewEuclidean2D constructs a 2D space that clamps vectors to the given width and height.
func NewEuclidean2D(sizeX, sizeY float64) Space2D {
	return &euclidean2d{
		space2d: space2d{
			size:     geom.NewVec(sizeX, sizeY),
			viewport: geom.NewAABBAt(geom.NewVec(0, 0), sizeX, sizeY),
		},
	}
}

type euclidean2d struct{ space2d }

func (s euclidean2d) Normalize(aabb geom.AABB) geom.AABB {
	wrappedAABB := s.WrapAABB(aabb)
	s.normalizeAABB(&wrappedAABB)
	return wrappedAABB.AABB
}

func (s euclidean2d) Name() string { return modeEuclidean2D }

func (s euclidean2d) Viewport() geom.AABB { return s.viewport }

func (s euclidean2d) WrapAABB(aabb geom.AABB) AABB {
	width := aabb.BottomRight.X - aabb.TopLeft.X
	height := aabb.BottomRight.Y - aabb.TopLeft.Y
	wrappedAABB := NewAABB(aabb.TopLeft, width, height)
	s.normalizeAABB(&wrappedAABB)
	return wrappedAABB
}

func (s euclidean2d) WrapVec(vec geom.Vec) AABB {
	aabb := geom.NewAABBAt(vec, 0, 0)
	return s.WrapAABB(aabb)
}

func (s euclidean2d) Expand(aabb *AABB, margin float64) {
	aabb.TopLeft.AddMutable(geom.NewVec(-margin, -margin))
	aabb.Size.AddMutable(geom.NewVec(2*margin, 2*margin))
	s.normalizeAABB(aabb)
}

func (s euclidean2d) Translate(aabb *AABB, delta geom.Vec) {
	aabb.TopLeft.AddMutable(delta)
	s.normalizeAABB(aabb)
}

// Reposition shifts aabb by delta, clamping its position so it stays
// within bounds without clipping its size — see Space2D.Reposition.
func (s euclidean2d) Reposition(aabb *AABB, delta geom.Vec) {
	aabb.TopLeft = aabb.TopLeft.Add(delta)
	aabb.TopLeft = s.clampPositionForSize(aabb.TopLeft, aabb.Size)
	aabb.BottomRight = aabb.TopLeft.Add(aabb.Size)
}

// clampPositionForSize clamps pos so a box of the given size stays
// within [0, s.size] without shrinking — a size at least as large as
// the space on an axis is pinned to 0 on that axis.
func (s euclidean2d) clampPositionForSize(pos, size geom.Vec) geom.Vec {
	var zero float64
	maxX, maxY := zero, zero
	if size.X < s.size.X {
		maxX = s.size.X - size.X
	}
	if size.Y < s.size.Y {
		maxY = s.size.Y - size.Y
	}
	return geom.Clamp(pos, geom.NewVec(maxX, maxY))
}

func (s euclidean2d) AABBDistance() AABBDistance {
	return newAABBDistance(s.metric)
}

func (s euclidean2d) normalizeVec(vec geom.Vec) geom.Vec {
	return geom.Clamp(vec, s.size)
}

func (s euclidean2d) normalizeAABB(aabb *AABB) {
	s.normalizeAABBBottomRight(aabb)
	s.normalizeAABBTopLeft(aabb)
}

func (s euclidean2d) normalizeAABBTopLeft(aabb *AABB) {
	aabb.TopLeft = s.normalizeVec(aabb.TopLeft)
}

func (s euclidean2d) normalizeAABBBottomRight(aabb *AABB) {
	aabb.BottomRight = aabb.TopLeft.Add(aabb.Size)
	aabb.BottomRight = geom.Clamp(aabb.BottomRight, s.size)
}

func (s euclidean2d) metric(vec1, vec2 geom.Vec) float64 {
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
	return geom.Length(geom.Clamp(delta, s.size))
}

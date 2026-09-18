package plane

import "github.com/kjkrol/gokg/geom"

// NewToroidal2D constructs a 2D space with wrap-around behaviour on both axes.
func NewToroidal2D(sizeX, sizeY float64) Space2D {
	return &toroidal2d{
		space2d: space2d{
			size:     geom.NewVec(sizeX, sizeY),
			viewport: geom.NewAABBAt(geom.NewVec(0, 0), sizeX, sizeY),
		},
	}
}

type toroidal2d struct{ space2d }

func (s toroidal2d) Normalize(aabb geom.AABB) geom.AABB {
	wrappedAABB := s.WrapAABB(aabb)
	s.normalizeAABB(&wrappedAABB)
	return wrappedAABB.AABB
}

func (s toroidal2d) Name() string { return modeToroidal2D }

func (s toroidal2d) Viewport() geom.AABB { return s.viewport }

func (s toroidal2d) WrapAABB(aabb geom.AABB) AABB {
	width := aabb.BottomRight.X - aabb.TopLeft.X
	height := aabb.BottomRight.Y - aabb.TopLeft.Y
	wrappedAABB := NewAABB(aabb.TopLeft, width, height)
	s.normalizeAABB(&wrappedAABB)
	return wrappedAABB
}

func (s toroidal2d) WrapVec(vec geom.Vec) AABB {
	aabb := geom.NewAABBAt(vec, 0, 0)
	return s.WrapAABB(aabb)
}

func (s toroidal2d) Expand(aabb *AABB, margin float64) {
	aabb.TopLeft.AddMutable(geom.NewVec(-margin, -margin))
	aabb.Size.AddMutable(geom.NewVec(2*margin, 2*margin))
	s.normalizeAABB(aabb)
}

func (s toroidal2d) Translate(aabb *AABB, delta geom.Vec) {
	aabb.TopLeft.AddMutable(delta)
	s.normalizeAABB(aabb)
}

// Reposition delegates to Translate — toroidal wrapping already keeps size fixed.
func (s toroidal2d) Reposition(aabb *AABB, delta geom.Vec) {
	s.Translate(aabb, delta)
}

func (s toroidal2d) AABBDistance() AABBDistance {
	return newAABBDistance(s.metric)
}

func (s toroidal2d) normalizeVec(vec geom.Vec) geom.Vec {
	return geom.Wrap(vec, s.size)
}

func (s toroidal2d) normalizeAABB(aabb *AABB) {
	s.normalizeAABBTopLeft(aabb)
	dx, dy := s.normalizeAABBBottomRight(aabb)
	aabb.fragmentation(dx, dy)
}

func (s toroidal2d) normalizeAABBTopLeft(aabb *AABB) {
	aabb.TopLeft = s.normalizeVec(aabb.TopLeft)
}

func (s toroidal2d) normalizeAABBBottomRight(aabb *AABB) (dx float64, dy float64) {
	aabb.BottomRight = aabb.TopLeft.Add(aabb.Size)
	dx = float64(0)
	dy = float64(0)
	if aabb.BottomRight.X > s.size.X {
		dx = aabb.BottomRight.X - s.size.X
	}
	if aabb.BottomRight.Y > s.size.Y {
		dy = aabb.BottomRight.Y - s.size.Y
	}
	aabb.BottomRight = geom.Clamp(aabb.BottomRight, s.size)
	return
}

func (s toroidal2d) metric(vec1, vec2 geom.Vec) float64 {
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

	delta := geom.Wrap(geom.NewVec(dx, dy), s.size)

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

	return geom.Length(delta)
}

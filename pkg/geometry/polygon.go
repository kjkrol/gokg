package geometry

import "strings"

// Polygon represents a simple polygon defined by an ordered list of vertices.
// The polygon is assumed to be closed, i.e. the last vertex connects back to the first one.
type Polygon[T SupportedNumeric] struct {
	points    []Vec[T] //TODO: oj oj, wygodniej byloby trzymac liste pointerow
	fragments map[OffsetRelativPos]Shape[T]
	bounds    AABB[T]
}

// NewPolygon constructs a polygon from a sequence of vertices.
// A minimum of three vertices is required to form a polygon.
func NewPolygon[T SupportedNumeric](vertices ...Vec[T]) Polygon[T] {
	if len(vertices) < 3 {
		panic("geometry.NewPolygon requires at least three vertices")
	}

	copyVertices := make([]Vec[T], len(vertices))
	copy(copyVertices, vertices)

	bounds := computeBounds(copyVertices)

	return Polygon[T]{
		points:    copyVertices,
		bounds:    bounds,
		fragments: make(map[OffsetRelativPos]Shape[T], 4),
	}
}

// ------- Builder -------------------------------

type PolygonBuilder[T SupportedNumeric] struct {
	points []Vec[T]
}

func NewPolygonBuilder[T SupportedNumeric]() *PolygonBuilder[T] {
	return &PolygonBuilder[T]{
		points: make([]Vec[T], 0),
	}
}

func (pb *PolygonBuilder[T]) Add(x, y T) *PolygonBuilder[T] {
	pb.points = append(pb.points, NewVec(x, y))
	return pb
}

func (pb *PolygonBuilder[T]) Build() Polygon[T] {
	return NewPolygon(pb.points...)
}

// -----------------------------------------------

func NewRect[T SupportedNumeric](topLeft Vec[T], witdh T, height T) Polygon[T] {
	topRight := topLeft.Add(NewVec(witdh, 0))
	bottomRight := topLeft.Add(NewVec(witdh, height))
	bottomLeft := topLeft.Add(NewVec(0, height))
	return NewPolygon(topLeft, topRight, bottomRight, bottomLeft)
}

func computeBounds[T SupportedNumeric](vertices []Vec[T]) AABB[T] {
	topLeft := vertices[0]
	bottomRight := vertices[0]

	for _, v := range vertices[1:] {
		if v.X < topLeft.X {
			topLeft.X = v.X
		}
		if v.Y < topLeft.Y {
			topLeft.Y = v.Y
		}
		if v.X > bottomRight.X {
			bottomRight.X = v.X
		}
		if v.Y > bottomRight.Y {
			bottomRight.Y = v.Y
		}
	}

	return NewAABB(topLeft, bottomRight)
}

// Bounds returns the axis-aligned bounding rectangle that contains the polygon.
func (p Polygon[T]) Bounds() AABB[T] {
	return p.bounds
}

func (p *Polygon[T]) UpdateBounds() {
	p.bounds = computeBounds(p.points)
}

// Points returns a copy of the polygon vertices.
func (p Polygon[T]) Points() []Vec[T] {
	pts := make([]Vec[T], len(p.points))
	copy(pts, p.points)
	return pts
}

// Vertices returns pointers to the underlying vertices so callers can mutate them.
func (p *Polygon[T]) Vertices() []*Vec[T] {
	ptrs := make([]*Vec[T], len(p.points))
	for i := range p.points {
		ptrs[i] = &p.points[i]
	}
	return ptrs
}

func (p Polygon[T]) Fragments() map[OffsetRelativPos]Shape[T] { return p.fragments }

func (p Polygon[T]) Clone() Shape[T] {
	pointsCopy := make([]Vec[T], len(p.points))
	copy(pointsCopy, p.points)
	clone := Polygon[T]{
		points: pointsCopy,
		bounds: p.bounds,
	}
	return &clone
}

func (p Polygon[T]) String() string {
	vertices := p.Vertices()
	sb := strings.Builder{}
	sb.WriteString("[")
	for i, v := range vertices {
		sb.WriteString(v.String())
		if i < len(vertices)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("]")
	return sb.String()
}

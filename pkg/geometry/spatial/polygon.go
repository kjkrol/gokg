package spatial

// Polygon represents a simple polygon defined by an ordered list of vertices.
// The polygon is assumed to be closed, i.e. the last vertex connects back to the first one.
type Polygon[T SupportedNumeric] struct {
	points    []Vec[T]
	bounds    Rectangle[T]
	fragments []Spatial[T]
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
		points: copyVertices,
		bounds: bounds,
	}
}

func computeBounds[T SupportedNumeric](vertices []Vec[T]) Rectangle[T] {
	minX, maxX := vertices[0].X, vertices[0].X
	minY, maxY := vertices[0].Y, vertices[0].Y
	for _, v := range vertices[1:] {
		if v.X < minX {
			minX = v.X
		}
		if v.X > maxX {
			maxX = v.X
		}
		if v.Y < minY {
			minY = v.Y
		}
		if v.Y > maxY {
			maxY = v.Y
		}
	}
	return NewRectangle(
		Vec[T]{X: minX, Y: minY},
		Vec[T]{X: maxX, Y: maxY},
	)
}

// Bounds returns the axis-aligned bounding rectangle that contains the polygon.
func (p Polygon[T]) Bounds() Rectangle[T] {
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

func (p Polygon[T]) Fragments() []Spatial[T] { return p.fragments }

func (p *Polygon[T]) SetFragments(f []Spatial[T]) { p.fragments = f }

func (p *Polygon[T]) Clone() Polygon[T] {
	pointsCopy := make([]Vec[T], len(p.points))
	copy(pointsCopy, p.points)
	return Polygon[T]{
		points: pointsCopy,
		bounds: p.bounds,
		// fragments świadomie zostawiasz puste,
	}
}

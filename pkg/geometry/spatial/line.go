package spatial

// Line represents a line segment defined by two endpoints.
// The segment is closed (includes both endpoints).
type Line[T SupportedNumeric] struct {
	Start     Vec[T]
	End       Vec[T]
	fragments []Spatial[T]
}

// NewLine constructs a line segment from two endpoints.
func NewLine[T SupportedNumeric](start, end Vec[T]) Line[T] {
	return Line[T]{Start: start, End: end}
}

// Bounds returns the axis-aligned bounding rectangle enclosing the line segment.
func (l Line[T]) Bounds() Rectangle[T] {
	minX, maxX := l.Start.X, l.Start.X
	minY, maxY := l.Start.Y, l.Start.Y

	if l.End.X < minX {
		minX = l.End.X
	}
	if l.End.X > maxX {
		maxX = l.End.X
	}
	if l.End.Y < minY {
		minY = l.End.Y
	}
	if l.End.Y > maxY {
		maxY = l.End.Y
	}

	return NewRectangle(
		Vec[T]{X: minX, Y: minY},
		Vec[T]{X: maxX, Y: maxY},
	)
}

// Vertices returns the endpoints of the line.
func (l *Line[T]) Vertices() []*Vec[T] {
	if l == nil {
		return nil
	}
	return []*Vec[T]{&l.Start, &l.End}
}

func (l Line[T]) Fragments() []Spatial[T] { return l.fragments }

func (l *Line[T]) SetFragments(f []Spatial[T]) { l.fragments = f }

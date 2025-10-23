package geometry

type Plane[T SupportedNumeric] struct {
	size       Vec[T]
	vectorMath VectorMath[T]
	normalize  func(*Vec[T])
	metric     func(v1, v2 Vec[T]) T
	name       string
}

// -----------------------------------------------------------------------------

func (p Plane[T]) Size() Vec[T] { return p.size }

func (p Plane[T]) Translate(vec *Vec[T], delta Vec[T]) {
	vec.AddMutable(delta)
	p.normalize(vec)
}

func (p Plane[T]) TranslateSpatial(spatial Spatial[T], delta Vec[T]) {
	if spatial == nil {
		return
	}

	switch s := spatial.(type) {
	case *Vec[T]:
		p.Translate(s, delta)
		spatial.SetFragments(nil)
		return
	case *Rectangle[T]:
		s.TopLeft.AddMutable(delta)
		s.BottomRight.AddMutable(delta)
		s.Center.AddMutable(delta)
		if p.name != "cyclic" {
			p.normalize(&s.TopLeft)
			p.normalize(&s.BottomRight)
			s.Center = Vec[T]{
				X: (s.TopLeft.X + s.BottomRight.X) / 2,
				Y: (s.TopLeft.Y + s.BottomRight.Y) / 2,
			}
			spatial.SetFragments(nil)
			return
		}
	case *Line[T]:
		s.Start.AddMutable(delta)
		s.End.AddMutable(delta)
		if p.name != "cyclic" {
			p.normalize(&s.Start)
			p.normalize(&s.End)
			spatial.SetFragments(nil)
			return
		}
	case *Polygon[T]:
		for i := range s.points {
			s.points[i] = s.points[i].Add(delta)
		}
		s.bounds = computeBounds(s.points)
		if p.name != "cyclic" {
			for i := range s.points {
				p.normalize(&s.points[i])
			}
			s.bounds = computeBounds(s.points)
			spatial.SetFragments(nil)
			return
		}
	default:
		vertices := spatial.Vertices()
		if len(vertices) == 0 {
			spatial.SetFragments(nil)
			return
		}
		for _, v := range vertices {
			if v == nil {
				continue
			}
			v.AddMutable(delta)
			if p.name != "cyclic" {
				p.normalize(v)
			}
		}
		if p.name != "cyclic" {
			spatial.SetFragments(nil)
			return
		}
	}

	if p.name != "cyclic" {
		spatial.SetFragments(nil)
		return
	}

	spatial.SetFragments(wrapSpatialFragments(spatial, p.size, p.vectorMath))
}

func (p Plane[T]) Metric(v1, v2 Vec[T]) T { return p.metric(v1, v2) }

func (p Plane[T]) Contains(vec Vec[T]) bool {
	return vec.X >= 0 && vec.X < p.size.X && vec.Y >= 0 && vec.Y < p.size.Y
}

func (p Plane[T]) Normalize(vec *Vec[T]) { p.normalize(vec) }

func (p Plane[T]) relativeMetric(v1, v2 Vec[T]) T {
	delta := v1.Sub(v2)
	p.normalize(&delta)
	return p.vectorMath.Length(delta)
}

func (p Plane[T]) Name() string { return p.name }

// -----------------------------------------------------------------------------

func NewBoundedPlane[T SupportedNumeric](sizeX, sizeY T) Plane[T] {
	plane := Plane[T]{
		name:       "bounded",
		size:       Vec[T]{sizeX, sizeY},
		vectorMath: VectorMathByType[T](),
	}
	plane.normalize = func(v *Vec[T]) { plane.vectorMath.Clamp(v, plane.size) }
	plane.metric = func(v1, v2 Vec[T]) T { return max(plane.relativeMetric(v1, v2), plane.relativeMetric(v2, v1)) }
	return plane
}

// -----------------------------------------------------------------------------

func NewCyclicBoundedPlane[T SupportedNumeric](sizeX, sizeY T) Plane[T] {
	plane := Plane[T]{
		name:       "cyclic",
		size:       Vec[T]{sizeX, sizeY},
		vectorMath: VectorMathByType[T](),
	}
	plane.normalize = func(v *Vec[T]) { plane.vectorMath.Wrap(v, plane.size) }
	plane.metric = func(v1, v2 Vec[T]) T { return min(plane.relativeMetric(v1, v2), plane.relativeMetric(v2, v1)) }
	return plane
}

// -----------------------------------------------------------------------------

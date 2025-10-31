package geometry

func translate[T SupportedNumeric](shape Shape[T], delta Vec[T], plane Plane[T]) {
	if shape == nil {
		return
	}

	translateInPlace(shape, delta)

	if plane.name == BOUNDED {
		for _, v := range shape.Vertices() {
			if v != nil {
				plane.normalize(v)
			}
		}
		shape.SetFragments(nil)
		return
	}

	switch s := shape.(type) {
	case *Vec[T]:
		plane.normalize(s)
	default:
		fragments := plane.createShapeFragmentsIfNeeded(shape)
		shape.SetFragments(fragments)
	}

}

func translateInPlace[T SupportedNumeric](shape Shape[T], delta Vec[T]) {
	switch s := shape.(type) {
	case *Vec[T]:
		s.AddMutable(delta)

	case *Line[T]:
		s.Start.AddMutable(delta)
		s.End.AddMutable(delta)

	case *Polygon[T]:
		for _, v := range s.Vertices() {
			if v != nil {
				v.AddMutable(delta)
			}
		}
		s.UpdateBounds()

	default:
		for _, v := range shape.Vertices() {
			if v != nil {
				v.AddMutable(delta)
			}
		}
	}
}

func (p Plane[T]) createShapeFragmentsIfNeeded(shape Shape[T]) []Shape[T] {
	vertices := shape.Vertices()
	if len(vertices) == 0 {
		return nil
	}
	base := *vertices[0]
	return GenerateBoundaryFragments(base, p, func(offset Vec[T]) (Shape[T], AABB[T], bool) {
		clone := shape.Clone()
		if clone == nil {
			return nil, AABB[T]{}, false
		}
		translateInPlace(clone, offset)
		return clone, clone.Bounds(), true
	})
}

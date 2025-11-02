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
	viewport := NewAABB(Vec[T]{X: 0, Y: 0}, Vec[T]{X: p.size.X, Y: p.size.Y})
	clipper := NewSutherlandHodgmanClipper(viewport)

	return GenerateBoundaryFragments(base, p, func(offset Vec[T]) (Shape[T], AABB[T], bool) {
		clone := shape.Clone()
		if clone == nil {
			return nil, AABB[T]{}, false
		}
		translateInPlace(clone, offset)
		preBounds := clone.Bounds()
		if !preBounds.Intersects(viewport) {
			return nil, AABB[T]{}, false
		}

		if polygon, ok := clone.(*Polygon[T]); ok {
			return buildPolygonFragment(polygon, clipper)
		}

		return normalizeFragment(p, clone, preBounds)
	})

}

func buildPolygonFragment[T SupportedNumeric](polygon *Polygon[T], clipper SutherlandHodgmanClipper[T]) (Shape[T], AABB[T], bool) {
	clipped := clipper.Clip(polygon.Points())
	if len(clipped) == 0 {
		return nil, AABB[T]{}, false
	}

	clippedPoly := NewPolygon(clipped...)
	return &clippedPoly, clippedPoly.Bounds(), true
}

func normalizeFragment[T SupportedNumeric](plane Plane[T], fragment Shape[T], originalBounds AABB[T]) (Shape[T], AABB[T], bool) {
	for _, v := range fragment.Vertices() {
		if v != nil {
			plane.normalize(v)
		}
	}
	if updater, ok := fragment.(interface{ UpdateBounds() }); ok {
		updater.UpdateBounds()
	}
	return fragment, originalBounds, true
}

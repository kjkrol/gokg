package geometry

func translate[T SupportedNumeric](shape Shape[T], delta Vec[T], plane Plane[T]) {
	if shape == nil {
		return
	}

	if plane.name == BOUNDED {
		translateInPlace(shape, delta)
		if !plane.ContainsShape(shape) {
			transformShape(shape, func(v *Vec[T]) { plane.vectorMath.Clamp(v, plane.size) })
		}
		return
	}

	switch s := shape.(type) {
	case *Vec[T]:
		translateInPlace(shape, delta)
		plane.normalize(s)
	default:
		translateCyclic(shape, delta, plane, nil, 0)
	}

}

func translateCyclic[T SupportedNumeric](shape Shape[T], delta Vec[T], plane Plane[T], parent Shape[T], relativePos OffsetRelativPos) {

	if len(shape.Fragments()) > 0 {

		for key, fragment := range shape.Fragments() {
			translateCyclic(fragment, delta, plane, shape, key)
		}
		if len(shape.Fragments()) == 1 {
			for key, fragment := range shape.Fragments() {
				shape = fragment
				delete(shape.Fragments(), key)
			}
		}
		return
	}

	translateInPlace(shape, delta)

	if parent != nil && !plane.viewport.Intersects(shape.Bounds()) {
		delete(parent.Fragments(), relativePos)
		return
	}

	if !plane.ContainsShape(shape) {
		plane.createShapeFragments(shape)
	}

}

func transformShape[T SupportedNumeric](shape Shape[T], transform func(*Vec[T])) {
	for _, v := range shape.Vertices() {
		if v != nil {
			transform(v)
		}
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
		// do nothing
	}
}

func (p Plane[T]) createShapeFragments(shape Shape[T]) {
	for offsetRelativePos, offsetVec := range p.offsets {
		temp := shape.Clone()
		translateInPlace(temp, offsetVec)
		if p.viewport.Intersects(temp.Bounds()) {
			transformShape(temp, func(v *Vec[T]) { p.vectorMath.Clamp(v, p.size) })
			shape.Fragments()[offsetRelativePos] = temp
			if polygon, ok := temp.(interface{ UpdateBounds() }); ok {
				polygon.UpdateBounds()
			}
		}
	}
}

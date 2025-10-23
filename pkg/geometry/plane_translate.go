package geometry

import "fmt"

func wrapSpatialFragments[T SupportedNumeric](spatial Spatial[T], size Vec[T], vecMath VectorMath[T]) []Spatial[T] {
	vertices := spatial.Vertices()
	if len(vertices) == 0 {
		return nil
	}
	reference := *vertices[0]
	wrappedRef := reference
	vecMath.Wrap(&wrappedRef, size)
	baseShift := wrappedRef.Sub(reference)

	offsets := candidateOffsets(baseShift, size)
	offsets = dedupeOffsets(offsets)
	viewport := NewRectangle(Vec[T]{X: 0, Y: 0}, Vec[T]{X: size.X, Y: size.Y})
	fragments := make([]Spatial[T], 0, len(offsets))

	for _, offset := range offsets {
		clone := cloneSpatialWithOffset(spatial, offset)
		if clone == nil {
			continue
		}
		if bounds := clone.Bounds(); bounds.Intersects(viewport) {
			fragments = append(fragments, clone)
		}
	}

	if len(fragments) == 0 {
		return nil
	}
	return fragments
}

func candidateOffsets[T SupportedNumeric](base Vec[T], size Vec[T]) []Vec[T] {
	offsets := make([]Vec[T], 0, 9)
	for _, mx := range [...]int{-1, 0, 1} {
		for _, my := range [...]int{-1, 0, 1} {
			offset := Vec[T]{
				X: base.X + T(mx)*size.X,
				Y: base.Y + T(my)*size.Y,
			}
			offsets = append(offsets, offset)
		}
	}
	return offsets
}

func dedupeOffsets[T SupportedNumeric](offsets []Vec[T]) []Vec[T] {
	seen := make(map[string]struct{}, len(offsets))
	result := make([]Vec[T], 0, len(offsets))
	for _, off := range offsets {
		key := fmt.Sprintf("%v:%v", off.X, off.Y)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, off)
	}
	return result
}

func cloneSpatialWithOffset[T SupportedNumeric](spatial Spatial[T], offset Vec[T]) Spatial[T] {
	switch s := spatial.(type) {
	case *Vec[T]:
		vecCopy := *s
		vecCopy.AddMutable(offset)
		return &vecCopy
	case *Rectangle[T]:
		rectCopy := *s
		rectCopy.TopLeft.AddMutable(offset)
		rectCopy.BottomRight.AddMutable(offset)
		rectCopy.Center.AddMutable(offset)
		return &rectCopy
	case *Line[T]:
		lineCopy := *s
		lineCopy.Start.AddMutable(offset)
		lineCopy.End.AddMutable(offset)
		return &lineCopy
	case *Polygon[T]:
		points := make([]Vec[T], len(s.points))
		copy(points, s.points)
		for i := range points {
			points[i] = points[i].Add(offset)
		}
		poly := NewPolygon(points...)
		return &poly
	default:
		vertices := spatial.Vertices()
		if len(vertices) == 0 {
			return nil
		}
		points := make([]Vec[T], 0, len(vertices))
		for _, v := range vertices {
			if v == nil {
				continue
			}
			points = append(points, v.Add(offset))
		}
		if len(points) < 3 {
			return nil
		}
		poly := NewPolygon(points...)
		return &poly
	}
}

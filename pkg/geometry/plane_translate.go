package geometry

import (
	"fmt"
)

func translateInPlace[T SupportedNumeric](spatialItem Shape[T], delta Vec[T]) {
	switch s := spatialItem.(type) {
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
		for _, v := range spatialItem.Vertices() {
			if v != nil {
				v.AddMutable(delta)
			}
		}
	}
}

func wrapSpatialFragments[T SupportedNumeric](shape Shape[T], size Vec[T], vecMath VectorMath[T]) []Shape[T] {
	vertices := shape.Vertices()
	if len(vertices) == 0 {
		return nil
	}
	reference := *vertices[0]
	wrappedRef := reference
	vecMath.Wrap(&wrappedRef, size)
	baseShift := wrappedRef.Sub(reference)

	offsets := candidateOffsets(baseShift, size)
	offsets = dedupeOffsets(offsets)
	viewport := NewAABB(Vec[T]{X: 0, Y: 0}, Vec[T]{X: size.X, Y: size.Y})
	fragments := make([]Shape[T], 0, len(offsets))

	for _, offset := range offsets {
		if offset.X == 0 && offset.Y == 0 {
			continue
		}
		clone := cloneSpatialWithOffset(shape, offset)
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
	if len(offsets) == 0 {
		return offsets
	}

	result := offsets[:0]
	for _, off := range offsets {
		duplicate := false
		for _, existing := range result {
			if off == existing {
				duplicate = true
				break
			}
		}
		if !duplicate {
			result = append(result, off)
		}
	}
	return result
}

func cloneSpatialWithOffset[T SupportedNumeric](spatialItem Shape[T], offset Vec[T]) Shape[T] {
	switch s := spatialItem.(type) {
	case *Vec[T]:
		vecCopy := *s
		translateInPlace(&vecCopy, offset)
		return &vecCopy

	case *Line[T]:
		lineCopy := *s
		translateInPlace(&lineCopy, offset)
		return &lineCopy

	case *Polygon[T]:
		polyCopy := s.Clone()
		translateInPlace(&polyCopy, offset)
		return &polyCopy

	default:
		panic(fmt.Sprintf("cloneSpatialWithOffset: unsupported spatial type %T", spatialItem))
	}
}

package geometry

import (
	"fmt"

	s "github.com/kjkrol/gokg/pkg/geometry/spatial"
)

func translateInPlace[T supportedNumeric](spatialItem spatial[T], delta vec[T]) {
	switch s := spatialItem.(type) {
	case *vec[T]:
		s.AddMutable(delta)

	case *rectangle[T]:
		s.TopLeft.AddMutable(delta)
		s.BottomRight.AddMutable(delta)
		s.Center.AddMutable(delta)

	case *line[T]:
		s.Start.AddMutable(delta)
		s.End.AddMutable(delta)

	case *polygon[T]:
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

func wrapSpatialFragments[T supportedNumeric](spatialItem spatial[T], size vec[T], vecMath VectorMath[T]) []spatial[T] {
	vertices := spatialItem.Vertices()
	if len(vertices) == 0 {
		return nil
	}
	reference := *vertices[0]
	wrappedRef := reference
	vecMath.Wrap(&wrappedRef, size)
	baseShift := wrappedRef.Sub(reference)

	offsets := candidateOffsets(baseShift, size)
	offsets = dedupeOffsets(offsets)
	viewport := s.NewRectangle(vec[T]{X: 0, Y: 0}, vec[T]{X: size.X, Y: size.Y})
	fragments := make([]spatial[T], 0, len(offsets))

	for _, offset := range offsets {
		if offset.X == 0 && offset.Y == 0 {
			continue
		}
		clone := cloneSpatialWithOffset(spatialItem, offset)
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

func candidateOffsets[T supportedNumeric](base vec[T], size vec[T]) []vec[T] {
	offsets := make([]vec[T], 0, 9)
	for _, mx := range [...]int{-1, 0, 1} {
		for _, my := range [...]int{-1, 0, 1} {
			offset := vec[T]{
				X: base.X + T(mx)*size.X,
				Y: base.Y + T(my)*size.Y,
			}
			offsets = append(offsets, offset)
		}
	}
	return offsets
}

func dedupeOffsets[T supportedNumeric](offsets []vec[T]) []vec[T] {
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

func cloneSpatialWithOffset[T supportedNumeric](spatialItem spatial[T], offset vec[T]) spatial[T] {
	switch s := spatialItem.(type) {
	case *vec[T]:
		vecCopy := *s
		translateInPlace(&vecCopy, offset)
		return &vecCopy

	case *rectangle[T]:
		rectCopy := *s
		translateInPlace(&rectCopy, offset)
		return &rectCopy

	case *line[T]:
		lineCopy := *s
		translateInPlace(&lineCopy, offset)
		return &lineCopy

	case *polygon[T]:
		polyCopy := s.Clone()
		translateInPlace(&polyCopy, offset)
		return &polyCopy

	default:
		panic(fmt.Sprintf("cloneSpatialWithOffset: unsupported spatial type %T", spatialItem))
	}
}

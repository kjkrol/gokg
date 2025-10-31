package geometry

// GenerateBoundaryFragments generates wrapped copies based on the reference point and collects those intersecting the viewport.
func GenerateBoundaryFragments[T SupportedNumeric, R any](
	base Vec[T],
	size Vec[T],
	vecMath VectorMath[T],
	build func(Vec[T]) (R, AABB[T], bool),
) []R {
	offsets := wrapOffsets(base, size, vecMath)
	viewport := NewAABB(Vec[T]{X: 0, Y: 0}, Vec[T]{X: size.X, Y: size.Y})
	return collectWrapped(offsets, viewport, build)
}

// WrapOffsets returns the candidate wrap offsets for a reference point in the given plane size.
// The zero offset is included; callers may skip it when they require only wrapped copies.
func wrapOffsets[T SupportedNumeric](reference Vec[T], size Vec[T], vecMath VectorMath[T]) []Vec[T] {
	wrapped := reference
	vecMath.Wrap(&wrapped, size)
	baseShift := wrapped.Sub(reference)
	return dedupeOffsets(candidateOffsets(baseShift, size))
}

// CollectWrapped iterates over offsets, builds wrapped items and collects those intersecting the viewport.
// The builder returns the wrapped value, its bounds, and whether the value is valid; invalid entries are skipped.
func collectWrapped[T SupportedNumeric, R any](
	offsets []Vec[T],
	viewport AABB[T],
	build func(Vec[T]) (R, AABB[T], bool),
) []R {
	results := make([]R, 0, len(offsets))
	for _, offset := range offsets {
		if offset.X == 0 && offset.Y == 0 {
			continue
		}
		item, bounds, ok := build(offset)
		if !ok {
			continue
		}
		if bounds.Intersects(viewport) {
			results = append(results, item)
		}
	}
	if len(results) == 0 {
		return nil
	}
	return results
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

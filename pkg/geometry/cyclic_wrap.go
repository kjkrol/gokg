package geometry

// GenerateBoundaryFragments generates wrapped copies based on the reference point and collects those intersecting the viewport.
func GenerateBoundaryFragments[T SupportedNumeric, R any](
	base Vec[T],
	plane Plane[T],
	build func(Vec[T]) (R, AABB[T], bool),
) []R {
	size := plane.Size()
	offsets := wrapOffsets(base, size, plane.normalize)
	viewport := NewAABB(Vec[T]{X: 0, Y: 0}, Vec[T]{X: size.X, Y: size.Y})
	return collectWrapped(offsets, viewport, build)
}

// WrapOffsets returns the candidate wrap offsets for a reference point in the given plane size.
// The zero offset is included; callers may skip it when they require only wrapped copies.
func wrapOffsets[T SupportedNumeric](reference Vec[T], size Vec[T], normalize func(*Vec[T])) []Vec[T] {
	wrapped := reference
	normalize(&wrapped)
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

func dedupeOffsets[T comparable](offsets []T) []T {
	seen := make(map[T]struct{}, len(offsets))
	out := offsets[:0]
	for _, o := range offsets {
		if _, ok := seen[o]; !ok {
			seen[o] = struct{}{}
			out = append(out, o)
		}
	}
	return out
}

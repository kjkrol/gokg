package geometry

// Spatial represents an entity that occupies space on a 2D plane and exposes its
// axis-aligned bounding box (AABB) for neighbor lookups.
type Spatial[T SupportedNumeric] interface {
	Bounds() Rectangle[T]
	Probe(margin T, plane Plane[T]) []Rectangle[T]
	DistanceTo(other Spatial[T], metric func(Vec[T], Vec[T]) T) T
	Vertices() []*Vec[T]
	Fragments() []Spatial[T]
	SetFragments([]Spatial[T])
}

func aabbDistance[T SupportedNumeric](
	a Spatial[T],
	b Spatial[T],
	metric func(Vec[T], Vec[T]) T,
) T {
	boundsA := a.Bounds()
	boundsB := b.Bounds()

	if boundsA.Intersects(boundsB) {
		return 0
	}

	dx := rectanglesAxisDistance(boundsA, boundsB, func(v Vec[T]) T { return v.X })
	dy := rectanglesAxisDistance(boundsA, boundsB, func(v Vec[T]) T { return v.Y })

	return metric(Vec[T]{X: dx, Y: dy}, Vec[T]{X: 0, Y: 0})
}

func rectanglesAxisDistance[T SupportedNumeric](
	aa, bb Rectangle[T],
	axisValue func(Vec[T]) T,
) T {
	aa, bb = SortRectanglesBy(aa, bb, axisValue)

	if axisValue(aa.BottomRight) >= axisValue(bb.TopLeft) {
		return 0
	}
	return axisValue(bb.TopLeft) - axisValue(aa.BottomRight)
}

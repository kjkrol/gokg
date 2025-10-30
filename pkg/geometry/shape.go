package geometry

type Shape[T SupportedNumeric] interface {
	Bounds() AABB[T]
	Vertices() []*Vec[T]
	Fragments() []Shape[T]
	SetFragments([]Shape[T])
}

package geometry

type Shape[T SupportedNumeric] interface {
	Bounds() AABB[T]
	Vertices() []*Vec[T]
	String() string
	Fragments() map[OffsetRelativPos]Shape[T]
	Clone() Shape[T]
}

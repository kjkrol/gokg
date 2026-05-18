package geom

import "testing"

var (
	aabbBoolSink bool
)

func Benchmark_AABB_Contains(b *testing.B) {
	outer := NewAABB(NewVec(0, 0), NewVec(10, 10))
	inner := NewAABB(NewVec(2, 2), NewVec(5, 5))

	b.ReportAllocs()
	for b.Loop() {
		aabbBoolSink = outer.Contains(inner)
	}
}

func Benchmark_AABB_Intersects(b *testing.B) {
	a := NewAABB(NewVec(0, 0), NewVec(5, 5))
	bb := NewAABB(NewVec(4, 4), NewVec(7, 7))

	b.ReportAllocs()
	for b.Loop() {
		aabbBoolSink = a.Intersects(bb)
	}
}

package geom

import "testing"

var (
	aabbBoolSink bool
)

func BenchmarkAABBContains(b *testing.B) {
	outer := NewAABB(NewVec(0, 0), NewVec(10, 10))
	inner := NewAABB(NewVec(2, 2), NewVec(5, 5))

	b.ReportAllocs()
	for b.Loop() {
		aabbBoolSink = outer.Contains(inner)
	}
}

func BenchmarkAABBIntersects(b *testing.B) {
	a := NewAABB(NewVec(0, 0), NewVec(5, 5))
	bb := NewAABB(NewVec(4, 4), NewVec(7, 7))

	b.ReportAllocs()
	for b.Loop() {
		aabbBoolSink = a.Intersects(bb)
	}
}

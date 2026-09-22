package bench_test

import (
	"testing"

	"github.com/kjkrol/aabbworld/geom"
)

var boolSink bool

// Benchmark_AABB_Contains tests one box inside another.
func Benchmark_AABB_Contains(b *testing.B) {
	outer := geom.NewAABB(geom.NewVec(0, 0), geom.NewVec(10, 10))
	inner := geom.NewAABB(geom.NewVec(2, 2), geom.NewVec(5, 5))
	b.ReportAllocs()
	for b.Loop() {
		boolSink = outer.Contains(inner)
	}
}

// Benchmark_AABB_Intersects tests two overlapping boxes.
func Benchmark_AABB_Intersects(b *testing.B) {
	a := geom.NewAABB(geom.NewVec(0, 0), geom.NewVec(5, 5))
	bb := geom.NewAABB(geom.NewVec(4, 4), geom.NewVec(7, 7))
	b.ReportAllocs()
	for b.Loop() {
		boolSink = a.Intersects(bb)
	}
}

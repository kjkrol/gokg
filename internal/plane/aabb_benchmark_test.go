package plane

import (
	"github.com/kjkrol/aabbworld/plane"
	"testing"

	"github.com/kjkrol/aabbworld/geom"
)

var (
	intersectsSink bool
	containsSink   bool
)

func Benchmark_AABB_Intersects(b *testing.B) {
	base := plane.NewAABB(geom.NewVec(0, 0), 5, 5)
	target := plane.NewAABB(geom.NewVec(4, 4), 3, 3)

	b.ReportAllocs()
	for b.Loop() {
		intersectsSink = base.Intersects(target.AABB)
	}
}

func Benchmark_AABB_Contains(b *testing.B) {
	outer := plane.NewAABB(geom.NewVec(0, 0), 10, 10)
	inner := plane.NewAABB(geom.NewVec(3, 3), 2, 2)

	b.ReportAllocs()
	for b.Loop() {
		containsSink = outer.Contains(inner.AABB)
	}
}

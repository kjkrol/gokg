package plane

import (
	"testing"

	"github.com/kjkrol/gokg/geom"
)

var (
	intersectsSink bool
	containsSink   bool
)

func Benchmark_AABB_Intersects(b *testing.B) {
	base := NewAABB(geom.NewVec(0, 0), 5, 5)
	target := NewAABB(geom.NewVec(4, 4), 3, 3)

	b.ReportAllocs()
	for b.Loop() {
		intersectsSink = base.Intersects(target.AABB)
	}
}

func Benchmark_AABB_IntersectsWithFrags(b *testing.B) {
	base := NewAABB(geom.NewVec(0, 0), 5, 5)
	target := NewAABB(geom.NewVec(4, 4), 3, 3)
	target.Overhang = geom.NewVec(2, 2) // hangs 2 past each far edge, so all three fragments exist

	b.ReportAllocs()
	for b.Loop() {
		intersectsSink = base.IntersectsWithFrags(target)
	}
}

func Benchmark_AABB_Contains(b *testing.B) {
	outer := NewAABB(geom.NewVec(0, 0), 10, 10)
	inner := NewAABB(geom.NewVec(3, 3), 2, 2)

	b.ReportAllocs()
	for b.Loop() {
		containsSink = outer.Contains(inner.AABB)
	}
}

func Benchmark_AABB_ContainsWithFrags(b *testing.B) {
	base := NewAABB(geom.NewVec(0, 0), 5, 5)
	target := NewAABB(geom.NewVec(4, 4), 3, 3)
	target.Overhang = geom.NewVec(2, 2) // hangs 2 past each far edge, so all three fragments exist

	b.ReportAllocs()
	for b.Loop() {
		intersectsSink = base.ContainsWithFrags(target)
	}
}

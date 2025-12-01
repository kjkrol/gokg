package plane

import (
	"testing"

	"github.com/kjkrol/gokg/pkg/geom"
)

var (
	intersectsSink bool
	containsSink   bool
)

func Benchmark_AABB_Intersects(b *testing.B) {
	base := newAABB(geom.NewVec(0, 0), 5, 5)
	target := newAABB(geom.NewVec(4, 4), 3, 3)

	b.ReportAllocs()
	for b.Loop() {
		intersectsSink = base.Intersects(target.AABB)
	}
}

func Benchmark_AABB_IntersectsWithFrags(b *testing.B) {
	base := newAABB(geom.NewVec(0, 0), 5, 5)
	target := newAABB(geom.NewVec(4, 4), 3, 3)
	target.setFragment(FRAG_RIGHT, geom.NewAABB(geom.NewVec(0, 4), geom.NewVec(2, 7)))
	target.setFragment(FRAG_BOTTOM, geom.NewAABB(geom.NewVec(4, 0), geom.NewVec(7, 2)))
	target.setFragment(FRAG_BOTTOM_RIGHT, geom.NewAABB(geom.NewVec(0, 0), geom.NewVec(2, 2)))

	b.ReportAllocs()
	for b.Loop() {
		intersectsSink = base.IntersectsWithFrags(target)
	}
}

func Benchmark_AABB_Contains(b *testing.B) {
	outer := newAABB(geom.NewVec(0, 0), 10, 10)
	inner := newAABB(geom.NewVec(3, 3), 2, 2)

	b.ReportAllocs()
	for b.Loop() {
		containsSink = outer.Contains(inner.AABB)
	}
}

func Benchmark_AABB_ContainsWithFrags(b *testing.B) {
	base := newAABB(geom.NewVec(0, 0), 5, 5)
	target := newAABB(geom.NewVec(4, 4), 3, 3)
	target.setFragment(FRAG_RIGHT, geom.NewAABB(geom.NewVec(0, 4), geom.NewVec(2, 7)))
	target.setFragment(FRAG_BOTTOM, geom.NewAABB(geom.NewVec(4, 0), geom.NewVec(7, 2)))
	target.setFragment(FRAG_BOTTOM_RIGHT, geom.NewAABB(geom.NewVec(0, 0), geom.NewVec(2, 2)))

	b.ReportAllocs()
	for b.Loop() {
		intersectsSink = base.ContainsWithFrags(target)
	}
}

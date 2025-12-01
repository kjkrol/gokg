package plane

import (
	"testing"

	"github.com/kjkrol/gokg/pkg/geom"
)

var (
	intersectsSink bool
	containsSink   bool
)

func BenchmarkAABBIntersects(b *testing.B) {
	baseTemplate := newAABB(geom.NewVec(0, 0), 5, 5)
	targetTemplate := newAABB(geom.NewVec(4, 4), 3, 3)

	b.ReportAllocs()
	var base, target AABB[int]
	for i := 0; i < b.N; i++ {
		base = baseTemplate
		target = targetTemplate
		intersectsSink = base.Intersects(target)
	}
}

func BenchmarkAABBIntersectsWithFrags(b *testing.B) {
	baseTemplate := newAABB(geom.NewVec(0, 0), 5, 5)
	targetTemplate := newAABB(geom.NewVec(4, 4), 3, 3)
	targetTemplate.setFragment(FRAG_RIGHT, geom.NewAABB(geom.NewVec(0, 4), geom.NewVec(2, 7)))
	targetTemplate.setFragment(FRAG_BOTTOM, geom.NewAABB(geom.NewVec(4, 0), geom.NewVec(7, 2)))
	targetTemplate.setFragment(FRAG_BOTTOM_RIGHT, geom.NewAABB(geom.NewVec(0, 0), geom.NewVec(2, 2)))

	b.ReportAllocs()
	var base, target AABB[int]
	for i := 0; i < b.N; i++ {
		base = baseTemplate
		target = targetTemplate
		intersectsSink = base.Intersects(target)
	}
}

func BenchmarkAABBContains(b *testing.B) {
	outerTemplate := newAABB(geom.NewVec(0, 0), 10, 10)
	innerTemplate := newAABB(geom.NewVec(3, 3), 2, 2)

	b.ReportAllocs()
	var outer, inner AABB[int]
	for i := 0; i < b.N; i++ {
		outer = outerTemplate
		inner = innerTemplate
		containsSink = outer.Contains(inner)
	}
}

package plane

import (
	"testing"

	"github.com/kjkrol/gokg/pkg/geom"
)

var (
	boolSink bool
	intSink  int
)

func BenchmarkEuclideanNormalizeAABB(b *testing.B) {
	space := NewEuclidean2D(10, 10).(*euclidean2d[int])
	template := newAABB(geom.NewVec(9, 9), 2, 2)

	b.ReportAllocs()
	var aabb AABB[int]
	for i := 0; i < b.N; i++ {
		aabb = template
		space.normalizeAABB(&aabb)
	}
	boolSink = aabb.BottomRight.Equals(geom.NewVec(10, 10))
}

func BenchmarkToroidalNormalizeAABB(b *testing.B) {
	space := NewToroidal2D(10, 10).(*toroidal2d[int])
	template := newAABB(geom.NewVec(9, 9), 2, 2)

	b.ReportAllocs()
	var aabb AABB[int]
	for i := 0; i < b.N; i++ {
		aabb = template
		space.normalizeAABB(&aabb)
	}
	boolSink = aabb.BottomRight.Equals(geom.NewVec(10, 10))
}

func BenchmarkEuclideanMetric(b *testing.B) {
	space := NewEuclidean2D(100, 100).(*euclidean2d[int])
	v1 := geom.NewVec(12, 34)
	v2 := geom.NewVec(78, 90)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		intSink = space.metric(v1, v2)
	}
}

func BenchmarkToroidalMetric(b *testing.B) {
	space := NewToroidal2D(100, 100).(*toroidal2d[int])
	v1 := geom.NewVec(12, 34)
	v2 := geom.NewVec(78, 90)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		intSink = space.metric(v1, v2)
	}
}

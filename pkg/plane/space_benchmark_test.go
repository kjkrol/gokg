package plane

import (
	"testing"

	"github.com/kjkrol/gokg/pkg/geom"
)

var (
	boolSink bool
	intSink  int
	aabbSink AABB[int]
)

func Benchmark_Euclidean_NormalizeAABB(b *testing.B) {
	space := NewEuclidean2D(10, 10).(*euclidean2d[int])
	template := newAABB(geom.NewVec(9, 9), 2, 2)
	aabbs := make([]AABB[int], b.N)
	for i := range aabbs {
		aabbs[i] = template
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := range aabbs {
		space.normalizeAABB(&aabbs[i])
	}
	boolSink = aabbs[len(aabbs)-1].BottomRight.Equals(geom.NewVec(10, 10))
}

func Benchmark_Toroidal_NormalizeAABB(b *testing.B) {
	space := NewToroidal2D(10, 10).(*toroidal2d[int])
	template := newAABB(geom.NewVec(9, 9), 2, 2)
	aabbs := make([]AABB[int], b.N)
	for i := range aabbs {
		aabbs[i] = template
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := range aabbs {
		space.normalizeAABB(&aabbs[i])
	}
	boolSink = aabbs[len(aabbs)-1].BottomRight.Equals(geom.NewVec(10, 10))
}

func Benchmark_Euclidean_Metric(b *testing.B) {
	space := NewEuclidean2D(100, 100).(*euclidean2d[int])
	v1 := geom.NewVec(12, 34)
	v2 := geom.NewVec(78, 90)

	b.ReportAllocs()
	for b.Loop() {
		intSink = space.metric(v1, v2)
	}
}

func Benchmark_Toroidal_Metric(b *testing.B) {
	space := NewToroidal2D(100, 100).(*toroidal2d[int])
	v1 := geom.NewVec(12, 34)
	v2 := geom.NewVec(78, 90)

	b.ReportAllocs()
	for b.Loop() {
		intSink = space.metric(v1, v2)
	}
}

func Benchmark_Euclidean_Expand(b *testing.B) {
	space := NewEuclidean2D(100, 100).(*euclidean2d[int])
	template := newAABB(geom.NewVec(90, 90), 8, 8)

	b.ReportAllocs()
	for b.Loop() {
		box := template
		space.Expand(&box, 5)
		aabbSink = box
	}
}

func Benchmark_Toroidal_Expand(b *testing.B) {
	space := NewToroidal2D(100, 100).(*toroidal2d[int])
	template := newAABB(geom.NewVec(90, 90), 8, 8)

	b.ReportAllocs()
	for b.Loop() {
		box := template
		space.Expand(&box, 5)
		aabbSink = box
	}
}

func Benchmark_Euclidean_Translate(b *testing.B) {
	space := NewEuclidean2D(100, 100).(*euclidean2d[int])
	template := newAABB(geom.NewVec(80, 80), 15, 10)
	delta := geom.NewVec(12, -18)

	b.ReportAllocs()
	for b.Loop() {
		box := template
		space.Translate(&box, delta)
		aabbSink = box
	}
}

func Benchmark_Toroidal_Translate(b *testing.B) {
	space := NewToroidal2D(100, 100).(*toroidal2d[int])
	template := newAABB(geom.NewVec(80, 80), 15, 10)
	delta := geom.NewVec(12, -18)

	b.ReportAllocs()
	for b.Loop() {
		box := template
		space.Translate(&box, delta)
		aabbSink = box
	}
}

func Benchmark_Euclidean_AABBDistance(b *testing.B) {
	space := NewEuclidean2D(200, 200).(*euclidean2d[int])
	distance := space.AABBDistance()
	rectA := newAABB(geom.NewVec(10, 10), 10, 10)
	rectB := newAABB(geom.NewVec(150, 160), 12, 12)

	b.ReportAllocs()
	for b.Loop() {
		intSink = distance(rectA.AABB, rectB.AABB)
	}
}

func Benchmark_Toroidal_AABBDistance(b *testing.B) {
	space := NewToroidal2D(200, 200).(*toroidal2d[int])
	distance := space.AABBDistance()
	rectA := newAABB(geom.NewVec(10, 10), 10, 10)
	rectB := newAABB(geom.NewVec(150, 160), 12, 12)

	b.ReportAllocs()
	for b.Loop() {
		intSink = distance(rectA.AABB, rectB.AABB)
	}
}

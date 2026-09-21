package plane

import (
	"github.com/kjkrol/aabbworld/plane"
	"testing"

	"github.com/kjkrol/aabbworld/geom"
)

var (
	boolSink bool
	intSink  float64
	aabbSink plane.AABB
)

func Benchmark_Euclidean_NormalizeAABB(b *testing.B) {
	space := NewEuclidean2D(10, 10)
	template := plane.NewAABB(geom.NewVec(9, 9), 2, 2)
	aabbs := make([]plane.AABB, b.N)
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
	space := NewToroidal2D(10, 10)
	template := plane.NewAABB(geom.NewVec(9, 9), 2, 2)
	aabbs := make([]plane.AABB, b.N)
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

func Benchmark_Euclidean_Expand(b *testing.B) {
	space := NewEuclidean2D(100, 100)
	template := plane.NewAABB(geom.NewVec(90, 90), 8, 8)

	b.ReportAllocs()
	for b.Loop() {
		box := template
		space.Expand(&box, 5)
		aabbSink = box
	}
}

func Benchmark_Toroidal_Expand(b *testing.B) {
	space := NewToroidal2D(100, 100)
	template := plane.NewAABB(geom.NewVec(90, 90), 8, 8)

	b.ReportAllocs()
	for b.Loop() {
		box := template
		space.Expand(&box, 5)
		aabbSink = box
	}
}

func Benchmark_Euclidean_Translate(b *testing.B) {
	space := NewEuclidean2D(100, 100)
	template := plane.NewAABB(geom.NewVec(80, 80), 15, 10)
	delta := geom.NewVec(12, -18)

	b.ReportAllocs()
	for b.Loop() {
		box := template
		space.Translate(&box, delta)
		aabbSink = box
	}
}

func Benchmark_Toroidal_Translate(b *testing.B) {
	space := NewToroidal2D(100, 100)
	template := plane.NewAABB(geom.NewVec(80, 80), 15, 10)
	delta := geom.NewVec(12, -18)

	b.ReportAllocs()
	for b.Loop() {
		box := template
		space.Translate(&box, delta)
		aabbSink = box
	}
}

package aabbworld_test

import (
	"fmt"
	"math"
	"math/rand/v2"
	"testing"

	"github.com/kjkrol/aabbworld"
	"github.com/kjkrol/aabbworld/geom"
	"github.com/kjkrol/aabbworld/plane"
	"github.com/kjkrol/uid"
)

func torusIf(toroidal bool) aabbworld.Edges {
	if toroidal {
		return aabbworld.Torus
	}
	return 0
}

// realSpace scatters n entities through a bucket-indexed Space and returns it with an observer.
func realSpace(b *testing.B, n int, toroidal bool) (*aabbworld.Space, uid.UID64) {
	b.Helper()
	space, err := aabbworld.NewSpace(aabbworld.Config{
		Width: 4000, Height: 4000, Edges: torusIf(toroidal),
		BucketSize: 256, BucketCapacity: 16,
	})
	if err != nil {
		b.Fatal(err)
	}
	observer := uid.UID64(1)
	space.Insert(observer, ptr(plane.NewAABB(geom.NewVec(2000, 2000), 10, 10)))
	r := rand.New(rand.NewPCG(1, 2))
	for i := range n {
		space.Insert(uid.UID64(100+i), ptr(plane.NewAABB(geom.NewVec(float64(r.IntN(4000)), float64(r.IntN(4000))), 20, 20)))
	}
	space.Flush(nil)
	return space, observer
}

func Benchmark_Real_Visible(b *testing.B) {
	for _, n := range []int{100, 1000} {
		b.Run(name(n), func(b *testing.B) {
			space, observer := realSpace(b, n, false)
			cone := aabbworld.Cone{Direction: geom.NewVec(1.0, 0.0), HalfAngle: math.Pi / 4, Radius: 800}
			v := &aabbworld.View{}
			b.ReportAllocs()
			for b.Loop() {
				if space.Scan(observer, cone, v) {
					v.Entities(func(uid.UID64, float64) {})
				}
			}
		})
	}
}

func Benchmark_Real_Outline(b *testing.B) {
	for _, n := range []int{100, 1000} {
		b.Run(name(n), func(b *testing.B) {
			space, observer := realSpace(b, n, false)
			cone := aabbworld.Cone{Direction: geom.NewVec(1.0, 0.0), HalfAngle: math.Pi / 4, Radius: 800}
			v := &aabbworld.View{}
			var fog []geom.Vec
			b.ReportAllocs()
			for b.Loop() {
				if space.Scan(observer, cone, v) {
					fog = v.Outline(0, fog[:0])
				}
			}
		})
	}
}

func Benchmark_Real_View(b *testing.B) {
	for _, n := range []int{100, 1000} {
		b.Run(name(n), func(b *testing.B) {
			space, observer := realSpace(b, n, false)
			cone := aabbworld.Cone{Direction: geom.NewVec(1.0, 0.0), HalfAngle: math.Pi / 4, Radius: 800}
			v := &aabbworld.View{}
			var fog []geom.Vec
			b.ReportAllocs()
			for b.Loop() {
				if space.Scan(observer, cone, v) {
					v.Entities(func(uid.UID64, float64) {})
					fog = v.Outline(0, fog[:0])
				}
			}
		})
	}
}

func Benchmark_Real_Depths(b *testing.B) {
	for _, n := range []int{100, 1000} {
		b.Run(name(n), func(b *testing.B) {
			space, observer := realSpace(b, n, false)
			cone := aabbworld.Cone{Direction: geom.NewVec(1.0, 0.0), HalfAngle: math.Pi / 4, Radius: 800}
			v := &aabbworld.View{}
			var fog []float32
			b.ReportAllocs()
			for b.Loop() {
				if space.Scan(observer, cone, v) {
					fog = v.Depths(63, fog[:0])
				}
			}
		})
	}
}

func name(n int) string { return fmt.Sprintf("entities=%d", n) }

package bench_test

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

// realSpace scatters n entities through a bucket-indexed Space and returns it with an observer.
func realSpace(b *testing.B, n int) (*aabbworld.Space, uid.UID64) {
	b.Helper()
	space, err := aabbworld.NewSpace(aabbworld.Config{
		Width: 4000, Height: 4000,
		BucketSize: 256,
	})
	if err != nil {
		b.Fatal(err)
	}
	observer := uid.UID64(1)
	items := []aabbworld.Item{{ID: observer, Box: plane.NewAABB(geom.NewVec(2000, 2000), 10, 10)}}
	r := rand.New(rand.NewPCG(1, 2))
	for i := range n {
		items = append(items, aabbworld.Item{ID: uid.UID64(100 + i), Box: plane.NewAABB(geom.NewVec(float64(r.IntN(4000)), float64(r.IntN(4000))), 20, 20)})
	}
	space.Rebuild(items)
	return space, observer
}

func entities(n int) string { return fmt.Sprintf("entities=%d", n) }

var viewCone = aabbworld.Cone{Direction: geom.NewVec(1.0, 0.0), HalfAngle: math.Pi / 4, Radius: 800}

// Benchmark_View_Visible scans a cone and lists what it sees, nearest first.
func Benchmark_View_Visible(b *testing.B) {
	for _, n := range []int{100, 1000} {
		b.Run(entities(n), func(b *testing.B) {
			space, observer := realSpace(b, n)
			v := &aabbworld.View{}
			b.ReportAllocs()
			for b.Loop() {
				if space.Scan(observer, viewCone, v) {
					v.Entities(func(uid.UID64, float64) {})
				}
			}
		})
	}
}

// Benchmark_View_Outline scans a cone and traces the lit region as a fan of points.
func Benchmark_View_Outline(b *testing.B) {
	for _, n := range []int{100, 1000} {
		b.Run(entities(n), func(b *testing.B) {
			space, observer := realSpace(b, n)
			v := &aabbworld.View{}
			var fog []geom.Vec
			b.ReportAllocs()
			for b.Loop() {
				if space.Scan(observer, viewCone, v) {
					fog = v.Outline(0, fog[:0])
				}
			}
		})
	}
}

// Benchmark_View_Both scans once and reads both the entities and the outline.
func Benchmark_View_Both(b *testing.B) {
	for _, n := range []int{100, 1000} {
		b.Run(entities(n), func(b *testing.B) {
			space, observer := realSpace(b, n)
			v := &aabbworld.View{}
			var fog []geom.Vec
			b.ReportAllocs()
			for b.Loop() {
				if space.Scan(observer, viewCone, v) {
					v.Entities(func(uid.UID64, float64) {})
					fog = v.Outline(0, fog[:0])
				}
			}
		})
	}
}

// Benchmark_View_Depths scans a cone and samples its reach at 63 angles.
func Benchmark_View_Depths(b *testing.B) {
	for _, n := range []int{100, 1000} {
		b.Run(entities(n), func(b *testing.B) {
			space, observer := realSpace(b, n)
			v := &aabbworld.View{}
			var depths []float32
			b.ReportAllocs()
			for b.Loop() {
				if space.Scan(observer, viewCone, v) {
					depths = v.Depths(63, depths[:0])
				}
			}
		})
	}
}

// Benchmark_View_Translucent is Benchmark_View_Both with three entities in ten see-through at τ = 0.5.
func Benchmark_View_Translucent(b *testing.B) {
	translucent := func(id uid.UID64) float64 {
		if id%10 < 3 {
			return 0.5
		}
		return 0
	}
	cone := viewCone
	cone.Transparency = translucent
	for _, n := range []int{100, 1000} {
		b.Run(entities(n), func(b *testing.B) {
			space, observer := realSpace(b, n)
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

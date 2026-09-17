package raycast_test

import (
	"math"
	"math/rand/v2"
	"testing"

	"github.com/kjkrol/gokg"
	"github.com/kjkrol/gokg/geom"
	"github.com/kjkrol/gokg/plane"
	"github.com/kjkrol/gokg/raycast"
	"github.com/kjkrol/gokg/spatial"
	"github.com/kjkrol/uid"
)

// realSpace scatters n entities through an actual bucket-indexed Space, so the
// benchmark measures raycast rather than a stand-in's map scan.
func realSpace(b *testing.B, n int, toroidal bool) (*gokg.Space, uid.UID64) {
	b.Helper()
	space, err := gokg.NewSpace(gokg.Config{
		Width: 4000, Height: 4000, Toroidal: toroidal,
		BucketSize: spatial.Size256x256, BucketCapacity: 16,
	})
	if err != nil {
		b.Fatal(err)
	}
	observer := uid.UID64(1)
	space.Insert(observer, plane.NewAABB(geom.NewVec[uint32](2000, 2000), 10, 10))
	r := rand.New(rand.NewPCG(1, 2))
	for i := range n {
		space.Insert(uid.UID64(100+i),
			plane.NewAABB(geom.NewVec(uint32(r.IntN(4000)), uint32(r.IntN(4000))), 20, 20))
	}
	space.Flush(nil)
	return space, observer
}

func Benchmark_Real_Visible(b *testing.B) {
	for _, n := range []int{100, 1000} {
		b.Run(name(n), func(b *testing.B) {
			space, observer := realSpace(b, n, false)
			cone := raycast.Cone{Direction: geom.NewVec(1.0, 0.0), HalfAngle: math.Pi / 4, Radius: 800}
			v := &raycast.View{}
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
			cone := raycast.Cone{Direction: geom.NewVec(1.0, 0.0), HalfAngle: math.Pi / 4, Radius: 800}
			v := &raycast.View{}
			var fog []geom.Vec[float64]
			b.ReportAllocs()
			for b.Loop() {
				if space.Scan(observer, cone, v) {
					fog = v.Outline(0, fog[:0])
				}
			}
		})
	}
}

// Both answers off one scan, through a View the caller keeps — what a game does
// per tick for an observer it also draws.
func Benchmark_Real_View(b *testing.B) {
	for _, n := range []int{100, 1000} {
		b.Run(name(n), func(b *testing.B) {
			space, observer := realSpace(b, n, false)
			cone := raycast.Cone{Direction: geom.NewVec(1.0, 0.0), HalfAngle: math.Pi / 4, Radius: 800}
			v := &raycast.View{}
			var fog []geom.Vec[float64]
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

// Depths is the drawing path a game takes per tick: scan once, then sample the
// cone at a fixed resolution into a buffer it keeps.
func Benchmark_Real_Depths(b *testing.B) {
	for _, n := range []int{100, 1000} {
		b.Run(name(n), func(b *testing.B) {
			space, observer := realSpace(b, n, false)
			cone := raycast.Cone{Direction: geom.NewVec(1.0, 0.0), HalfAngle: math.Pi / 4, Radius: 800}
			v := &raycast.View{}
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

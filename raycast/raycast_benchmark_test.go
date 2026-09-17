package raycast_test

import (
	"math"
	"math/rand/v2"
	"testing"

	"github.com/kjkrol/gokg/geom"
	"github.com/kjkrol/gokg/raycast"
	"github.com/kjkrol/uid"
)

func scatter(n int, toroidal bool) (*fakeSpace, uid.UID64) {
	s := newFake(4000, 4000, toroidal)
	s.put(eye, 2000, 2000, 10, 10)
	r := rand.New(rand.NewPCG(1, 2))
	for i := range n {
		s.put(uid.UID64(100+i), uint32(r.IntN(4000)), uint32(r.IntN(4000)), 20, 20)
	}
	return s, eye
}

func benchVisible(b *testing.B, n int, toroidal bool) {
	s, observer := scatter(n, toroidal)
	cone := raycast.Cone{Direction: geom.NewVec(1.0, 0.0), HalfAngle: math.Pi / 4, Radius: 800}
	v := &raycast.View{}
	b.ReportAllocs()
	for b.Loop() {
		if v.Scan(s, observer, cone) {
			v.Entities(func(uid.UID64, float64) {})
		}
	}
}

func Benchmark_Visible_Euclidean(b *testing.B) {
	for _, n := range []int{100, 1000} {
		b.Run(name(n), func(b *testing.B) { benchVisible(b, n, false) })
	}
}

func Benchmark_Visible_Toroidal(b *testing.B) {
	for _, n := range []int{100, 1000} {
		b.Run(name(n), func(b *testing.B) { benchVisible(b, n, true) })
	}
}

func name(n int) string {
	switch n {
	case 100:
		return "entities=100"
	default:
		return "entities=1000"
	}
}

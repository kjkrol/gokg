package geom

import "testing"

var (
	vecSink    Vec
	scalarSink float64
)

func Benchmark_Clamp(b *testing.B) {
	v := NewVec(12.5, -3.25)
	size := NewVec(10, 10)

	b.ReportAllocs()
	for b.Loop() {
		vecSink = Clamp(v, size)
	}
}

// The fast paths matter more than the general case: motion produces
// coordinates that are already inside or one step over an edge, and only
// something thrown far out reaches math.Mod.
func Benchmark_Wrap(b *testing.B) {
	size := NewVec(1024, 1024)

	cases := map[string]Vec{
		"already inside":  NewVec(512, 512),
		"one step over":   NewVec(1030, 1030),
		"just below zero": NewVec(-6, -6),
		"far out":         NewVec(99999, -99999),
	}
	for name, v := range cases {
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				vecSink = Wrap(v, size)
			}
		})
	}
}

func Benchmark_Length(b *testing.B) {
	v := NewVec(3.5, 4.5)

	b.ReportAllocs()
	for b.Loop() {
		scalarSink = Length(v)
	}
}

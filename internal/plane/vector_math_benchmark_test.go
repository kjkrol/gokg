package plane

import (
	"testing"

	"github.com/kjkrol/aabbworld/geom"
)

var (
	vecSink    geom.Vec
	scalarSink float64
)

func Benchmark_Clamp(b *testing.B) {
	v := geom.NewVec(12.5, -3.25)
	size := geom.NewVec(10, 10)

	b.ReportAllocs()
	for b.Loop() {
		vecSink = Clamp(v, size)
	}
}

func Benchmark_Wrap(b *testing.B) {
	size := geom.NewVec(1024, 1024)

	cases := map[string]geom.Vec{
		"already inside":  geom.NewVec(512, 512),
		"one step over":   geom.NewVec(1030, 1030),
		"just below zero": geom.NewVec(-6, -6),
		"far out":         geom.NewVec(99999, -99999),
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

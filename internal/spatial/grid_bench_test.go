package spatial

import (
	"math/rand/v2"
	"testing"

	"github.com/kjkrol/aabbworld/geom"
	iplane "github.com/kjkrol/aabbworld/internal/plane"
	"github.com/kjkrol/uid"
)

func BenchmarkRebuildPairsQuery(b *testing.B) {
	for _, cfg := range []struct {
		name  string
		n     int
		sides []float64
	}{
		{"n=2000,5x5", 2000, []float64{5}},
		{"n=8388,5x5", 8388, []float64{5}},
		{"n=8388,mixed", 8388, []float64{4, 5, 16, 60}},
	} {
		b.Run(cfg.name, func(b *testing.B) {
			rng := rand.New(rand.NewPCG(3, 5))
			sc := newScene(b, iplane.Torus, rng, cfg.n, cfg.sides)
			deltas := make([]geom.Vec, len(sc.items))
			for i, it := range sc.items {
				side := min(it.Box.Size.X, it.Box.Size.Y)
				deltas[i] = geom.NewVec((rng.Float64()*2-1)*side/2, (rng.Float64()*2-1)*side/2)
			}
			var probes []geom.AABB
			for range 8 {
				probes = append(probes, geom.NewAABBAt(geom.NewVec(rng.Float64()*900, rng.Float64()*900), 100, 100))
			}
			pairs, found := 0, 0
			onPair := func(_, _ uid.UID64) { pairs++ }
			onFound := func(uid.UID64) { found++ }
			b.ReportAllocs()
			for b.Loop() {
				for i := range sc.items {
					sc.surface.Translate(&sc.items[i].Box, deltas[i])
				}
				sc.grid.Rebuild(sc.items)
				sc.grid.Pairs(0.5, collides, onPair)
				for _, p := range probes {
					sc.grid.Query(p, collides, onFound)
				}
			}
		})
	}
}

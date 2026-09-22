package collide_test

import (
	"math/rand/v2"
	"testing"

	"github.com/kjkrol/aabbworld/geom"
	"github.com/kjkrol/aabbworld/internal/collide"
	iplane "github.com/kjkrol/aabbworld/internal/plane"
	"github.com/kjkrol/aabbworld/internal/spatial"
	"github.com/kjkrol/aabbworld/plane"
	"github.com/kjkrol/uid"
)

// benchBatch is a field of boxes standing clear of their neighbours, every 25th nudged into one.
type benchBatch struct {
	home  []plane.AABB
	pairs [][2]int
}

func newBenchBatch(cols, rows int) *benchBatch {
	b := &benchBatch{}
	for r := range rows {
		for c := range cols {
			x := float64(c) * 12
			if (r*cols+c)%25 == 0 {
				x += 4
			}
			b.home = append(b.home, plane.NewAABB(geom.NewVec(x, float64(r)*12), 10, 10))
		}
	}
	for r := range rows {
		for c := range cols {
			i := r*cols + c
			if c+1 < cols {
				b.pairs = append(b.pairs, [2]int{i, i + 1})
			}
			if r+1 < rows {
				b.pairs = append(b.pairs, [2]int{i, i + cols})
			}
		}
	}
	return b
}

func BenchmarkSolve(b *testing.B) {
	batch := newBenchBatch(100, 50)
	surface := iplane.NewEuclidean2D(2000, 2000)

	items := make([]spatial.Item, len(batch.home))
	var s collide.Solver
	for b.Loop() {
		for i, box := range batch.home {
			items[i] = spatial.Item{ID: uid.UID64(i + 1), Box: box}
		}
		s.Reset(len(items))
		for _, p := range batch.pairs {
			s.Add(collide.Pair{A: int32(p[0]), B: int32(p[1])})
		}
		s.Solve(items, surface, iterations, nil, nil)
	}
	b.ReportMetric(float64(len(batch.pairs)), "pairs")
}

// denseScene is the collision demo's shape: 5x5 boxes over a fifth of a torus, paired by the grid.
type denseScene struct {
	surface *iplane.Surface
	home    []spatial.Item
	pairs   []collide.Pair
}

func newDenseScene(n int) *denseScene {
	rng := rand.New(rand.NewPCG(3, 5))
	surface := iplane.NewToroidal2D(1024, 1024)
	sc := &denseScene{surface: surface}
	for i := range n {
		box := plane.NewAABB(geom.NewVec(rng.Float64()*1024, rng.Float64()*1024), 5, 5)
		surface.Translate(&box, geom.Vec{})
		sc.home = append(sc.home, spatial.Item{ID: uid.UID64(i + 1), Box: box, Caps: 1})
	}
	grid := spatial.NewGrid(surface, spatial.Size1024x1024, spatial.Size16x16)
	grid.Rebuild(sc.home)
	grid.Pairs(0.5, 1, func(a, b int32) {
		sc.pairs = append(sc.pairs, collide.Pair{A: a, B: b})
	})
	return sc
}

func BenchmarkSolve_Dense(b *testing.B) {
	sc := newDenseScene(8388)
	items := make([]spatial.Item, len(sc.home))
	var s collide.Solver
	contacts := 0
	onContact := func(int, geom.Vec) { contacts++ }
	for b.Loop() {
		copy(items, sc.home)
		s.Reset(len(items))
		for _, p := range sc.pairs {
			s.Add(p)
		}
		contacts = 0
		s.Solve(items, sc.surface, iterations, nil, onContact)
	}
	b.ReportMetric(float64(len(sc.pairs)), "pairs")
	b.ReportMetric(float64(contacts), "contacts")
}

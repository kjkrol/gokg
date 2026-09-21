package collide_test

import (
	iplane "github.com/kjkrol/aabbworld/internal/plane"
	"testing"

	"github.com/kjkrol/aabbworld/geom"
	"github.com/kjkrol/aabbworld/internal/collide"
	"github.com/kjkrol/aabbworld/plane"
)

// benchBatch is a field of boxes standing clear of their neighbours, every 25th nudged into one.
type benchBatch struct {
	home  []plane.AABB
	boxes []plane.AABB
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
	b.boxes = make([]plane.AABB, len(b.home))
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

	for _, keyed := range []bool{false, true} {
		name := "unkeyed"
		if keyed {
			name = "keyed"
		}
		b.Run(name, func(b *testing.B) {
			var s collide.Solver
			for b.Loop() {
				copy(batch.boxes, batch.home)
				s.Reset()
				for _, p := range batch.pairs {
					pair := collide.Pair{A: &batch.boxes[p[0]], B: &batch.boxes[p[1]]}
					if keyed {
						pair.KeyA, pair.KeyB = uint32(p[0]), uint32(p[1])
					}
					s.Add(pair)
				}
				s.Solve(surface, iterations, nil)
			}
			b.ReportMetric(float64(len(batch.pairs)), "pairs")
		})
	}
}

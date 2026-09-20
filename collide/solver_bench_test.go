package collide_test

import (
	"testing"

	"github.com/kjkrol/gokg/collide"
	"github.com/kjkrol/gokg/geom"
	"github.com/kjkrol/gokg/plane"
)

// benchBatch is a field of boxes standing clear of their neighbours, every
// 25th nudged into the one beside it — what a broad phase hands over: mostly
// candidates that turn out not to touch. Pairs name each box and its
// neighbours to the right and below.
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
	surface := plane.NewEuclidean2D(2000, 2000)

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

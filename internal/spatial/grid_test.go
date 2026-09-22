package spatial

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/kjkrol/aabbworld/geom"
	iplane "github.com/kjkrol/aabbworld/internal/plane"
	"github.com/kjkrol/aabbworld/plane"
	"github.com/kjkrol/uid"
)

const collides = Capability(1 << 1)

// scene is a population on a surface, held by a grid and by the slice the grid was built from.
type scene struct {
	surface *iplane.Surface
	grid    *Grid
	items   []Item
}

func newScene(t testing.TB, edges iplane.Edges, rng *rand.Rand, count int, sides []float64) *scene {
	t.Helper()
	surface := iplane.NewSurface(1024, 1024, edges)
	sc := &scene{surface: surface, grid: NewGrid(surface, Size1024x1024, Size16x16)}
	for i := range count {
		w, h := sides[rng.IntN(len(sides))], sides[rng.IntN(len(sides))]
		box := plane.NewAABB(geom.NewVec(rng.Float64()*1024-w/2, rng.Float64()*1024-h/2), w, h)
		surface.Translate(&box, geom.Vec{})
		if surface.Left(&box) {
			continue
		}
		caps := Plain
		if rng.IntN(3) != 0 {
			caps |= collides
		}
		sc.items = append(sc.items, Item{ID: uid.UID64(i + 1), Box: box, Caps: caps})
	}
	sc.grid.Rebuild(sc.items)
	return sc
}

// step moves every item by a small random delta and rebuilds.
func (sc *scene) step(rng *rand.Rand) {
	kept := sc.items[:0]
	for i := range sc.items {
		it := &sc.items[i]
		side := min(it.Box.Size.X, it.Box.Size.Y)
		sc.surface.Translate(&it.Box, geom.NewVec((rng.Float64()*2-1)*side/2, (rng.Float64()*2-1)*side/2))
		if !sc.surface.Left(&it.Box) {
			kept = append(kept, *it)
		}
	}
	sc.items = kept
	sc.grid.Rebuild(sc.items)
}

// inside lists, by brute force, every item sharing want whose main box or image intersects box.
func (sc *scene) inside(box geom.AABB, want Capability) (ids []uid.UID64, pieces int) {
	area, ok := sc.grid.clamp(box)
	if !ok {
		return nil, 0
	}
	for _, it := range sc.items {
		if !it.Caps.matches(want) {
			continue
		}
		hit := false
		for _, img := range images(it.Box) {
			if clamped, ok := sc.grid.clamp(img); ok && clamped.Intersects(area) {
				pieces++
				hit = true
			}
		}
		if hit {
			ids = append(ids, it.ID)
		}
	}
	slices.Sort(ids)
	return ids, pieces
}

func TestGrid_AgreesWithBruteForce(t *testing.T) {
	for name, edges := range map[string]iplane.Edges{"closed": 0, "torus": iplane.Torus, "wrapX openY": iplane.WrapX | iplane.OpenY} {
		t.Run(name, func(t *testing.T) {
			rng := rand.New(rand.NewPCG(7, 11))
			sc := newScene(t, edges, rng, 600, []float64{4, 5, 16, 60})
			for round := range 40 {
				for q := range 30 {
					w, h := 8+rng.Float64()*120, 8+rng.Float64()*120
					box := sc.surface.WrapAABB(geom.NewAABBAt(geom.NewVec(rng.Float64()*1024-20, rng.Float64()*1024-20), w, h))
					want := collides
					if rng.IntN(2) == 0 {
						want = AnyCapability
					}
					wantIDs, wantPieces := sc.inside(box.AABB, want)
					var got []uid.UID64
					n := sc.grid.Query(box.AABB, want, func(id uid.UID64) { got = append(got, id) })
					slices.Sort(got)
					got = slices.Compact(got)
					if n != wantPieces || !slices.Equal(got, wantIDs) {
						t.Fatalf("round %d query %d: grid %d %v, brute force %d %v", round, q, n, got, wantPieces, wantIDs)
					}
				}

				want := sc.near(0.5, collides)
				var got []idPair
				sc.grid.Pairs(0.5, collides, func(x, y uid.UID64) { got = append(got, idPair{x, y}) })
				sortPairs(got)
				if !slices.Equal(got, want) {
					t.Fatalf("round %d: grid found %d pairs, brute force %d", round, len(got), len(want))
				}

				for _, it := range sc.items {
					box, ok := sc.grid.EntryAABB(it.ID)
					if clamped, _ := sc.grid.clamp(it.Box.AABB); !ok || box != clamped {
						t.Fatalf("round %d: EntryAABB(%v) = %v/%v, want %v", round, it.ID, box, ok, clamped)
					}
				}
				if _, ok := sc.grid.EntryAABB(uid.UID64(99999)); ok {
					t.Fatal("grid knows an entity it was never told of")
				}
				sc.step(rng)
			}
		})
	}
}

func TestGrid_IsEmptyBeforeItsFirstRebuild(t *testing.T) {
	g := NewGrid(iplane.NewToroidal2D(256, 256), Size256x256, Size32x32)
	if n := g.Query(geom.NewAABBAt(geom.NewVec(0, 0), 256, 256), AnyCapability, func(uid.UID64) {}); n != 0 {
		t.Errorf("a fresh grid answered a query with %d pieces", n)
	}
	g.Pairs(0.5, AnyCapability, func(uid.UID64, uid.UID64) { t.Error("a fresh grid paired something") })
	if _, ok := g.EntryAABB(1); ok {
		t.Error("a fresh grid knows an entity")
	}
}

func TestGrid_ForgetsWhatTheLastRebuildLeftOut(t *testing.T) {
	surface := iplane.NewEuclidean2D(256, 256)
	g := NewGrid(surface, Size256x256, Size32x32)
	items := []Item{
		{ID: 1, Box: plane.NewAABB(geom.NewVec(10, 10), 8, 8), Caps: Plain},
		{ID: 2, Box: plane.NewAABB(geom.NewVec(100, 100), 8, 8), Caps: Plain},
	}
	g.Rebuild(items)
	g.Rebuild(items[1:])

	if _, ok := g.EntryAABB(1); ok {
		t.Error("an item left out of the last Rebuild is still known")
	}
	if n := g.Query(geom.NewAABBAt(geom.NewVec(0, 0), 32, 32), AnyCapability, func(uid.UID64) {}); n != 0 {
		t.Errorf("Query finds %d pieces where the left-out item was", n)
	}
	if n := g.Query(geom.NewAABBAt(geom.NewVec(96, 96), 32, 32), AnyCapability, func(uid.UID64) {}); n != 1 {
		t.Errorf("Query finds %d pieces where the kept item is, want 1", n)
	}
}

func TestGrid_RefusesCoordinatesWithNoCell(t *testing.T) {
	g := NewGrid(iplane.NewEuclidean2D(256, 256), Size256x256, Size32x32)
	for name, v := range map[string]geom.Vec{
		"negative":          geom.NewVec(-5, -5),
		"NaN":               geom.NewVec(math.NaN(), 10),
		"beyond the grid":   geom.NewVec(1e9, 10),
		"positive infinity": geom.NewVec(math.Inf(1), 10),
	} {
		t.Run(name, func(t *testing.T) {
			if c := g.cellOf(v.X); c != 0 && c != g.side-1 {
				t.Errorf("cellOf(%v) = %d, want it held at an edge cell", v.X, c)
			}
			if _, ok := g.clamp(geom.NewAABBAt(v, 4, 4)); ok && (v.X < 0 || math.IsNaN(v.X)) {
				box, _ := g.clamp(geom.NewAABBAt(v, 4, 4))
				if box.TopLeft.X != 0 {
					t.Errorf("clamp put a box with no cell at %v", box)
				}
			}
		})
	}
}

func TestGrid_HoldsItemsAtFractionalPositions(t *testing.T) {
	g := NewGrid(iplane.NewEuclidean2D(256, 256), Size256x256, Size32x32)
	g.Rebuild([]Item{{ID: 1, Box: plane.NewAABB(geom.NewVec(41.7, 43.2), 3.5, 3.5), Caps: Plain}})
	if n := g.Query(geom.NewAABBAt(geom.NewVec(40, 40), 10, 10), AnyCapability, func(uid.UID64) {}); n != 1 {
		t.Errorf("Query finds %d pieces, want the fractionally placed item", n)
	}
}

func TestGrid_AnItemSpanningFourCellsIsReportedOnce(t *testing.T) {
	g := NewGrid(iplane.NewEuclidean2D(256, 256), Size256x256, Size32x32)
	g.Rebuild([]Item{{ID: 1, Box: plane.NewAABB(geom.NewVec(28, 28), 8, 8), Caps: Plain}})
	seen := 0
	if n := g.Query(geom.NewAABBAt(geom.NewVec(0, 0), 64, 64), AnyCapability, func(uid.UID64) { seen++ }); n != 1 || seen != 1 {
		t.Errorf("a box across four cells was reported %d times (count %d), want once", seen, n)
	}
}

func TestGrid_QuerySeesOnlyWhatItAskedFor(t *testing.T) {
	g := NewGrid(iplane.NewToroidal2D(256, 256), Size256x256, Size32x32)
	g.Rebuild([]Item{
		{ID: 1, Box: plane.NewAABB(geom.NewVec(100, 100), 8, 8), Caps: Plain},
		{ID: 2, Box: plane.NewAABB(geom.NewVec(112, 100), 8, 8), Caps: Plain | collides},
	})
	box := geom.NewAABBAt(geom.NewVec(90, 90), 40, 40)
	var got []uid.UID64
	g.Query(box, collides, func(id uid.UID64) { got = append(got, id) })
	if !slices.Equal(got, []uid.UID64{2}) {
		t.Errorf("asking for collides found %v, want [2]", got)
	}
	got = got[:0]
	g.Query(box, AnyCapability, func(id uid.UID64) { got = append(got, id) })
	slices.Sort(got)
	if !slices.Equal(got, []uid.UID64{1, 2}) {
		t.Errorf("asking for anything found %v, want both", got)
	}
}

// near is every pair Pairs owes, worked out by comparing every two grown boxes, images and all.
func (sc *scene) near(reach float64, want Capability) []idPair {
	grown := make([]plane.AABB, len(sc.items))
	for i, it := range sc.items {
		grown[i] = it.Box
		sc.surface.Expand(&grown[i], reach*min(it.Box.Size.X, it.Box.Size.Y))
	}
	var pairs []idPair
	for i := range sc.items {
		if !sc.items[i].Caps.matches(want) {
			continue
		}
		for j := i + 1; j < len(sc.items); j++ {
			if !sc.items[j].Caps.matches(want) || !meet(grown[i], grown[j]) {
				continue
			}
			a, b := sc.items[i].ID, sc.items[j].ID
			if a.Index() > b.Index() {
				a, b = b, a
			}
			pairs = append(pairs, idPair{a, b})
		}
	}
	sortPairs(pairs)
	return pairs
}

func images(ab plane.AABB) []geom.AABB {
	out := []geom.AABB{ab.AABB}
	ab.VisitFragments(func(_ plane.FragPosition, box geom.AABB) bool {
		out = append(out, box)
		return true
	})
	return out
}

func meet(a, b plane.AABB) bool {
	for _, ia := range images(a) {
		for _, ib := range images(b) {
			if ia.Intersects(ib) {
				return true
			}
		}
	}
	return false
}

type idPair struct{ a, b uid.UID64 }

func sortPairs(pairs []idPair) {
	slices.SortFunc(pairs, func(l, r idPair) int {
		if l.a != r.a {
			if l.a < r.a {
				return -1
			}
			return 1
		}
		if l.b < r.b {
			return -1
		}
		if l.b > r.b {
			return 1
		}
		return 0
	})
}

package collide_test

import (
	"cmp"
	"github.com/kjkrol/aabbworld"
	"github.com/kjkrol/aabbworld/collide"
	"github.com/kjkrol/aabbworld/internal/core"
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/kjkrol/aabbworld/geom"
	"github.com/kjkrol/aabbworld/plane"
	"github.com/kjkrol/uid"
)

const pairsReach = 0.5

type idPair struct{ a, b uid.UID64 }

func sortPairs(pairs []idPair) {
	slices.SortFunc(pairs, func(l, r idPair) int {
		if c := cmp.Compare(l.a, r.a); c != 0 {
			return c
		}
		return cmp.Compare(l.b, r.b)
	})
}

// crowd is a population kept beside the space it was inserted into, so a test
// can work out for itself who is near whom.
type crowd struct {
	space *aabbworld.Space
	ids   []uid.UID64
	boxes []plane.AABB
	caps  []aabbworld.Capability
}

func newCrowd(t testing.TB, width, height uint32, edges aabbworld.Edges) *crowd {
	t.Helper()
	space, err := aabbworld.NewSpace(aabbworld.Config{
		Width: width, Height: height, Edges: edges,
		BucketSize: 64, BucketCapacity: 8,
	})
	if err != nil {
		t.Fatalf("NewSpace: %v", err)
	}
	return &crowd{space: space}
}

func (c *crowd) add(x, y, w, h float64, caps aabbworld.Capability) {
	id := uid.UID64(len(c.ids))
	box := c.space.WrapAABB(plane.NewAABB(geom.NewVec(x, y), w, h).AABB)
	c.space.Insert(id, &box)
	c.space.SetCapabilities(id, caps)
	c.ids, c.boxes, c.caps = append(c.ids, id), append(c.boxes, box), append(c.caps, caps)
}

// scatter fills the crowd with count boxes of the given sides, anywhere in the world.
func (c *crowd) scatter(rng *rand.Rand, count int, sides []float64, whole bool) {
	width, height, _ := c.space.Bounds()
	for range count {
		w, h := sides[rng.IntN(len(sides))], sides[rng.IntN(len(sides))]
		x, y := rng.Float64()*float64(width), rng.Float64()*float64(height)
		if !c.space.Edges.WrapsX() {
			x = rng.Float64() * (float64(width) - w)
		}
		if !c.space.Edges.WrapsY() {
			y = rng.Float64() * (float64(height) - h)
		}
		if whole {
			x, y = float64(int(x)), float64(int(y))
		}
		caps := aabbworld.CanCollide
		if rng.IntN(5) == 0 {
			caps = aabbworld.Plain
		}
		c.add(x, y, w, h, caps)
	}
	c.space.Flush(nil)
}

// near is every pair Pairs owes, worked out by comparing every two grown boxes.
func (c *crowd) near(want aabbworld.Capability) []idPair {
	reaches := make([]plane.AABB, len(c.boxes))
	for i, box := range c.boxes {
		reaches[i] = box
		core.Of(c.space).Surface.Expand(&reaches[i], pairsReach*min(box.Size.X, box.Size.Y))
	}
	var pairs []idPair
	for i := range reaches {
		for j := i + 1; j < len(reaches); j++ {
			if c.caps[i]&want == 0 || c.caps[j]&want == 0 {
				continue
			}
			if meetAcrossSeams(reaches[i], reaches[j]) {
				pairs = append(pairs, idPair{c.ids[i], c.ids[j]})
			}
		}
	}
	sortPairs(pairs)
	return pairs
}

// imagesOf is a box's main rectangle and every piece it wraps into.
func imagesOf(ab plane.AABB) []geom.AABB {
	images := []geom.AABB{ab.AABB}
	ab.VisitFragments(func(_ plane.FragPosition, box geom.AABB) bool {
		images = append(images, box)
		return true
	})
	return images
}

// meetAcrossSeams reports whether any image of a touches any image of b.
func meetAcrossSeams(a, b plane.AABB) bool {
	for _, ia := range imagesOf(a) {
		for _, ib := range imagesOf(b) {
			if ia.Intersects(ib) {
				return true
			}
		}
	}
	return false
}

// candidates lists the pairs Engine.Tick resolves, refusing each so nothing is tested or moved.
func candidates(space *aabbworld.Space, want aabbworld.Capability) []idPair {
	var pairs []idPair
	var e collide.Engine
	e.Tick(space, pairsReach, want, 0, func(a, b uid.UID64) (collide.Body, collide.Body, bool) {
		pairs = append(pairs, idPair{a, b})
		return collide.Body{}, collide.Body{}, false
	}, nil, nil)
	return pairs
}

func (c *crowd) pairs(want aabbworld.Capability) []idPair { return candidates(c.space, want) }

func TestBroadPhase_FindsExactlyWhoIsNear(t *testing.T) {
	sizes := map[string][]float64{
		"one size":            {10},
		"2, 16 and 100":       {2, 16, 100},
		"wider than a bucket": {40, 150},
	}
	for _, edges := range []aabbworld.Edges{0, aabbworld.Torus, aabbworld.WrapX, aabbworld.WrapY, aabbworld.WrapX | aabbworld.OpenY, aabbworld.OpenX | aabbworld.OpenY} {
		for name, sides := range sizes {
			for _, whole := range []bool{false, true} {
				rng := rand.New(rand.NewPCG(17, 23))
				for round := range 12 {
					c := newCrowd(t, 1000, 800, edges)
					c.scatter(rng, 400, sides, whole)

					want := c.near(aabbworld.CanCollide)
					got := c.pairs(aabbworld.CanCollide)
					for _, p := range got {
						if p.a.Index() >= p.b.Index() {
							t.Fatalf("edges=%04b %s whole=%v round %d: pair (%v, %v) is not lower index first", edges, name, whole, round, p.a, p.b)
						}
					}
					sortPairs(got)
					if !slices.Equal(got, want) {
						t.Fatalf("edges=%04b %s whole=%v round %d: Pairs found %d pairs, want %d\nmissing: %v\nextra:   %v",
							edges, name, whole, round, len(got), len(want), missing(want, got), missing(got, want))
					}
				}
			}
		}
	}
}

// grown is the whole of b, seam pieces and all, reaching margin further on every side.
func grown(b plane.AABB, margin float64) geom.AABB {
	return geom.NewAABB(
		geom.NewVec(b.TopLeft.X-margin, b.TopLeft.Y-margin),
		geom.NewVec(b.TopLeft.X+b.Size.X+margin, b.TopLeft.Y+b.Size.Y+margin),
	)
}

// missing is what from holds and in does not, both sorted.
func missing(from, in []idPair) []idPair {
	var out []idPair
	left := slices.Clone(in)
	for _, p := range from {
		if i := slices.Index(left, p); i >= 0 {
			left = slices.Delete(left, i, i+1)
			continue
		}
		out = append(out, p)
		if len(out) == 8 {
			break
		}
	}
	return out
}

func TestBroadPhase_AgreesWithAProbePerEntity(t *testing.T) {
	const side = 10
	for _, edges := range []aabbworld.Edges{0, aabbworld.Torus, aabbworld.WrapX, aabbworld.WrapY, aabbworld.WrapX | aabbworld.OpenY, aabbworld.OpenX | aabbworld.OpenY} {
		rng := rand.New(rand.NewPCG(5, 9))
		c := newCrowd(t, 1000, 800, edges)
		c.scatter(rng, 600, []float64{side}, true)

		var want []idPair
		for i, id := range c.ids {
			if c.caps[i]&aabbworld.CanCollide == 0 {
				continue
			}
			seen := map[uid.UID64]bool{}
			c.space.Query(grown(c.boxes[i], side), aabbworld.CanCollide, func(other uid.UID64) {
				if other.Index() > id.Index() && !seen[other] {
					seen[other] = true
					want = append(want, idPair{id, other})
				}
			})
		}
		sortPairs(want)

		got := c.pairs(aabbworld.CanCollide)
		sortPairs(got)
		if !slices.Equal(got, want) {
			t.Errorf("edges=%04b: Pairs found %d pairs, a probe per entity %d\nmissing: %v\nextra:   %v",
				edges, len(got), len(want), missing(want, got), missing(got, want))
		}
	}
}

func TestBroadPhase_SaysTheSameThingTwice(t *testing.T) {
	c := newCrowd(t, 1000, 800, aabbworld.Torus)
	c.scatter(rand.New(rand.NewPCG(1, 2)), 500, []float64{2, 16, 100}, false)

	first, second := c.pairs(aabbworld.CanCollide), c.pairs(aabbworld.CanCollide)
	if len(first) == 0 {
		t.Fatal("a crowd this dense has to hold pairs")
	}
	if !slices.Equal(first, second) {
		t.Error("two sweeps over the same space named their pairs in a different order")
	}
}

func TestBroadPhase_FollowsAnEntityThatMoved(t *testing.T) {
	c := newCrowd(t, 1000, 800, 0)
	c.add(100, 100, 10, 10, aabbworld.CanCollide)
	c.add(300, 100, 10, 10, aabbworld.CanCollide)
	c.space.Flush(nil)
	if got := c.pairs(aabbworld.CanCollide); len(got) != 0 {
		t.Fatalf("two boxes 190 apart are a pair: %v", got)
	}

	c.space.Translate(c.ids[1], &c.boxes[1], geom.NewVec(-185, 0))
	c.space.Flush(nil)
	if got := c.pairs(aabbworld.CanCollide); len(got) != 1 {
		t.Errorf("after closing to 5 apart Pairs found %v, want the one pair", got)
	}

	c.space.Remove(c.ids[0])
	c.space.Flush(nil)
	if got := c.pairs(aabbworld.CanCollide); len(got) != 0 {
		t.Errorf("a removed entity is still paired: %v", got)
	}
}

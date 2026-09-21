package gokg

import (
	"cmp"
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/kjkrol/gokg/geom"
	"github.com/kjkrol/gokg/plane"
	"github.com/kjkrol/gokg/spatial"
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
	space *Space
	ids   []uid.UID64
	boxes []plane.AABB
	caps  []Capability
}

func newCrowd(t testing.TB, width, height uint32, toroidal bool) *crowd {
	t.Helper()
	space, err := NewSpace(Config{
		Width: width, Height: height, Toroidal: toroidal,
		BucketSize: spatial.ResolutionFrom(64), BucketCapacity: 8, OpsBufferSize: 1 << 14,
	})
	if err != nil {
		t.Fatalf("NewSpace: %v", err)
	}
	return &crowd{space: space}
}

func (c *crowd) add(x, y, w, h float64, caps Capability) {
	id := uid.UID64(len(c.ids))
	box := c.space.WrapAABB(plane.NewAABB(geom.NewVec(x, y), w, h).AABB)
	c.space.Insert(id, box)
	c.space.SetCapabilities(id, caps)
	c.ids, c.boxes, c.caps = append(c.ids, id), append(c.boxes, box), append(c.caps, caps)
}

// scatter fills the crowd with count boxes of the given sides, anywhere in the
// world — edges, corners and seams included — on whole coordinates or not.
func (c *crowd) scatter(rng *rand.Rand, count int, sides []float64, whole bool) {
	width, height, _ := c.space.Bounds()
	for range count {
		w, h := sides[rng.IntN(len(sides))], sides[rng.IntN(len(sides))]
		x, y := rng.Float64()*float64(width), rng.Float64()*float64(height)
		if !c.space.Toroidal {
			x, y = rng.Float64()*(float64(width)-w), rng.Float64()*(float64(height)-h)
		}
		if whole {
			x, y = float64(int(x)), float64(int(y))
		}
		caps := CanCollide
		if rng.IntN(5) == 0 {
			caps = Plain
		}
		c.add(x, y, w, h, caps)
	}
	c.space.Flush(nil)
}

// near is every pair Pairs owes, worked out the slow way: grow each box by its
// own reach and compare every two, wrapped images and all.
func (c *crowd) near(want Capability) []idPair {
	reaches := make([]plane.AABB, len(c.boxes))
	for i, box := range c.boxes {
		reaches[i] = box
		c.space.ExpandOnly(&reaches[i], pairsReach*min(box.Size.X, box.Size.Y))
	}
	var pairs []idPair
	for i := range reaches {
		for j := i + 1; j < len(reaches); j++ {
			if c.caps[i]&want == 0 || c.caps[j]&want == 0 {
				continue
			}
			if reaches[i].IntersectsWithFrags(reaches[j]) {
				pairs = append(pairs, idPair{c.ids[i], c.ids[j]})
			}
		}
	}
	sortPairs(pairs)
	return pairs
}

func (c *crowd) pairs(want Capability) []idPair {
	var pairs []idPair
	c.space.Pairs(pairsReach, want, func(a, b uid.UID64) { pairs = append(pairs, idPair{a, b}) })
	return pairs
}

func TestPairs_FindsExactlyWhoIsNear(t *testing.T) {
	sizes := map[string][]float64{
		"one size":            {10},
		"2, 16 and 100":       {2, 16, 100},
		"wider than a bucket": {40, 150},
	}
	for _, toroidal := range []bool{false, true} {
		for name, sides := range sizes {
			for _, whole := range []bool{false, true} {
				rng := rand.New(rand.NewPCG(17, 23))
				for round := range 12 {
					c := newCrowd(t, 1000, 800, toroidal)
					c.scatter(rng, 400, sides, whole)

					want := c.near(CanCollide)
					got := c.pairs(CanCollide)
					for _, p := range got {
						if p.a.Index() >= p.b.Index() {
							t.Fatalf("toroidal=%v %s whole=%v round %d: pair (%v, %v) is not lower index first", toroidal, name, whole, round, p.a, p.b)
						}
					}
					sortPairs(got)
					if !slices.Equal(got, want) {
						t.Fatalf("toroidal=%v %s whole=%v round %d: Pairs found %d pairs, want %d\nmissing: %v\nextra:   %v",
							toroidal, name, whole, round, len(got), len(want), missing(want, got), missing(got, want))
					}
				}
			}
		}
	}
}

// missing is what from holds and in does not, both sorted; a pair listed twice
// in from and once in in counts as missing once.
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

// For a world of one size, reaching half a side each is the margin of one side
// a probe per entity used to be given — so the two must agree.
func TestPairs_AgreesWithAProbePerEntity(t *testing.T) {
	const side = 10
	for _, toroidal := range []bool{false, true} {
		rng := rand.New(rand.NewPCG(5, 9))
		c := newCrowd(t, 1000, 800, toroidal)
		c.scatter(rng, 600, []float64{side}, true)

		var want []idPair
		for i, id := range c.ids {
			if c.caps[i]&CanCollide == 0 {
				continue
			}
			probe := c.boxes[i]
			seen := map[uid.UID64]bool{}
			c.space.Neighbours(&probe, side, CanCollide, func(other uid.UID64, _ plane.FragPosition) {
				if other.Index() > id.Index() && !seen[other] {
					seen[other] = true
					want = append(want, idPair{id, other})
				}
			})
		}
		sortPairs(want)

		got := c.pairs(CanCollide)
		sortPairs(got)
		if !slices.Equal(got, want) {
			t.Errorf("toroidal=%v: Pairs found %d pairs, a probe per entity %d\nmissing: %v\nextra:   %v",
				toroidal, len(got), len(want), missing(want, got), missing(got, want))
		}
	}
}

func TestPairs_SaysTheSameThingTwice(t *testing.T) {
	c := newCrowd(t, 1000, 800, true)
	c.scatter(rand.New(rand.NewPCG(1, 2)), 500, []float64{2, 16, 100}, false)

	first, second := c.pairs(CanCollide), c.pairs(CanCollide)
	if len(first) == 0 {
		t.Fatal("a crowd this dense has to hold pairs")
	}
	if !slices.Equal(first, second) {
		t.Error("two sweeps over the same space named their pairs in a different order")
	}
}

func TestPairs_FollowsAnEntityThatMoved(t *testing.T) {
	c := newCrowd(t, 1000, 800, false)
	c.add(100, 100, 10, 10, CanCollide)
	c.add(300, 100, 10, 10, CanCollide)
	c.space.Flush(nil)
	if got := c.pairs(CanCollide); len(got) != 0 {
		t.Fatalf("two boxes 190 apart are a pair: %v", got)
	}

	c.space.Translate(c.ids[1], &c.boxes[1], geom.NewVec(-185, 0))
	c.space.Flush(nil)
	if got := c.pairs(CanCollide); len(got) != 1 {
		t.Errorf("after closing to 5 apart Pairs found %v, want the one pair", got)
	}

	c.space.Remove(c.ids[0])
	c.space.Flush(nil)
	if got := c.pairs(CanCollide); len(got) != 0 {
		t.Errorf("a removed entity is still paired: %v", got)
	}
}

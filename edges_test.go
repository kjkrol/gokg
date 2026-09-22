package aabbworld_test

import (
	"github.com/kjkrol/aabbworld/collide"
	"math"
	"slices"
	"testing"

	"github.com/kjkrol/aabbworld"
	"github.com/kjkrol/aabbworld/geom"
	"github.com/kjkrol/aabbworld/plane"
	"github.com/kjkrol/uid"
)

func edgedSpace(t *testing.T, edges aabbworld.Edges) *aabbworld.Space {
	t.Helper()
	space, err := aabbworld.NewSpace(aabbworld.Config{
		Width: 1000, Height: 1000, Edges: edges, BucketSize: 64, BucketCapacity: 8,
	})
	if err != nil {
		t.Fatalf("NewSpace: %v", err)
	}
	return space
}

// found is every id a query over the whole world reports, each once, in order.
func found(space *aabbworld.Space) []uid.UID64 {
	var ids []uid.UID64
	space.Query(geom.NewAABBAt(geom.NewVec(0, 0), 1000, 1000), aabbworld.AnyCapability, func(id uid.UID64) {
		if !slices.Contains(ids, id) {
			ids = append(ids, id)
		}
	})
	slices.Sort(ids)
	return ids
}

func paired(space *aabbworld.Space) int {
	n := 0
	var e collide.Engine
	e.Tick(space, 0.5, aabbworld.AnyCapability, 0, func(_, _ uid.UID64) (collide.Body, collide.Body, bool) {
		n++
		return collide.Body{}, collide.Body{}, false
	}, nil, nil)
	return n
}

func TestNewSpace_RefusesAnAxisThatBothWrapsAndIsOpen(t *testing.T) {
	for _, edges := range []aabbworld.Edges{aabbworld.WrapX | aabbworld.OpenX, aabbworld.WrapY | aabbworld.OpenY} {
		if _, err := aabbworld.NewSpace(aabbworld.Config{Width: 100, Height: 100, Edges: edges, BucketSize: 16}); err == nil {
			t.Errorf("edges %04b accepted, want an error", edges)
		}
	}
}

func TestTranslate_ClosedWorldKeepsTheBoxWholeAtEveryEdge(t *testing.T) {
	space := edgedSpace(t, 0)
	id := uid.UID64(1)
	box := plane.NewAABB(geom.NewVec(500, 500), 20, 20)
	space.Insert(id, &box)

	for _, push := range []geom.Vec{{X: -5000}, {X: 5000}, {Y: -5000}, {Y: 5000}, {X: 5000, Y: 5000}} {
		if !space.Translate(id, &box, push) {
			t.Fatalf("push %v: reported gone from a closed world", push)
		}
		if w, h := box.BottomRight.X-box.TopLeft.X, box.BottomRight.Y-box.TopLeft.Y; w != 20 || h != 20 {
			t.Fatalf("push %v: box is %vx%v, want it whole at 20x20", push, w, h)
		}
		if box.TopLeft.X < 0 || box.TopLeft.Y < 0 || box.BottomRight.X > 1000 || box.BottomRight.Y > 1000 {
			t.Fatalf("push %v: box %v reaches outside the world", push, box)
		}
	}
}

func TestTranslate_EntityOnTheSeamFallsThroughTheOpenBottom(t *testing.T) {
	space := edgedSpace(t, aabbworld.WrapX|aabbworld.OpenY)
	faller, stayer := uid.UID64(1), uid.UID64(2)

	box := plane.NewAABB(geom.NewVec(990, 970), 20, 20)
	near := plane.NewAABB(geom.NewVec(5, 975), 20, 20)
	space.Insert(faller, &box)
	space.Insert(stayer, &near)
	space.Translate(faller, &box, geom.Vec{})
	space.Flush(nil)

	if box.Overhang.X != 10 {
		t.Fatalf("overhang %v, want the box wrapped 10 past the right edge", box.Overhang)
	}
	if got := paired(space); got != 1 {
		t.Fatalf("Pairs named %d pairs across the seam, want 1", got)
	}

	if !space.Translate(faller, &box, geom.NewVec(0, 25)) {
		t.Fatal("reported gone while 5 units of it are still inside")
	}
	if space.Translate(faller, &box, geom.NewVec(0, 5)) {
		t.Fatal("not reported gone once wholly below the open bottom")
	}
	if space.Translate(faller, &box, geom.NewVec(0, 5)) {
		t.Fatal("reported back while still outside")
	}
	space.Flush(nil)

	if got := found(space); !slices.Equal(got, []uid.UID64{stayer}) {
		t.Errorf("Query finds %v after the fall, want only %v", got, stayer)
	}
	if got := paired(space); got != 0 {
		t.Errorf("Pairs still names %d pairs with the fallen entity", got)
	}

	back := plane.NewAABB(geom.NewVec(990, 970), 20, 20)
	space.Insert(faller, &back)
	space.Flush(nil)
	if got := found(space); !slices.Equal(got, []uid.UID64{faller, stayer}) {
		t.Errorf("Query finds %v after Insert, want both back", got)
	}
}

func TestTranslate_OpenEdgeLetsABoxComeBackWhereItWas(t *testing.T) {
	space := edgedSpace(t, aabbworld.OpenX|aabbworld.OpenY)
	id := uid.UID64(1)
	box := plane.NewAABB(geom.NewVec(5, 500), 10, 10)
	start := box
	space.Insert(id, &box)

	space.Translate(id, &box, geom.NewVec(-8, 0))
	space.Flush(nil)
	if got := found(space); len(got) != 1 {
		t.Fatalf("Query finds %v while 2 units of the box are inside, want it", got)
	}
	space.Translate(id, &box, geom.NewVec(8, 0))
	if box != start {
		t.Errorf("box came back as %+v, want exactly %+v", box, start)
	}
}

func TestTick_NamesABoxItPushedOutOfAnOpenEdge(t *testing.T) {
	space := edgedSpace(t, aabbworld.OpenX)
	wall, pushed := uid.UID64(1), uid.UID64(2)
	a := plane.NewAABB(geom.NewVec(0, 100), 10, 10)
	b := plane.NewAABB(geom.NewVec(-9, 100), 10, 10)
	space.Insert(wall, &a)
	space.Insert(pushed, &b)
	space.Flush(nil)

	var e collide.Engine
	e.Tick(space, 0.5, aabbworld.AnyCapability, 4, func(x, y uid.UID64) (collide.Body, collide.Body, bool) {
		if x != wall || y != pushed {
			t.Fatalf("resolve asked about (%v, %v), want (%v, %v)", x, y, wall, pushed)
		}
		return collide.Body{Box: &a, Static: true}, collide.Body{Box: &b}, true
	}, nil, nil)
	space.Flush(nil)

	if got := e.Left(); !slices.Equal(got, []uid.UID64{pushed}) {
		t.Errorf("Left = %v after the wall pushed %v past the open edge, want exactly it", got, pushed)
	}
	if got := found(space); !slices.Equal(got, []uid.UID64{wall}) {
		t.Errorf("Query finds %v, want only the wall", got)
	}
}

func TestMoveTo_PlacesTheBoxUnderTheEdgeRules(t *testing.T) {
	id := uid.UID64(1)

	t.Run("onto a wrapping seam", func(t *testing.T) {
		space := edgedSpace(t, aabbworld.Torus)
		box := plane.NewAABB(geom.NewVec(500, 500), 20, 20)
		space.Insert(id, &box)
		if !space.MoveTo(id, &box, geom.NewVec(990, 500)) {
			t.Fatal("reported gone from a torus")
		}
		space.Flush(nil)

		if box.Overhang.X != 10 {
			t.Errorf("overhang %v, want the box wrapped 10 past the right edge", box.Overhang)
		}
		for name, probe := range map[string]geom.AABB{
			"before the seam": geom.NewAABB(geom.NewVec(992, 505), geom.NewVec(998, 510)),
			"after the seam":  geom.NewAABB(geom.NewVec(2, 505), geom.NewVec(8, 510)),
		} {
			if n := space.Query(probe, aabbworld.AnyCapability, func(uid.UID64) {}); n != 1 {
				t.Errorf("%s: Query finds %d pieces, want the box", name, n)
			}
		}
		if n := space.Query(geom.NewAABB(geom.NewVec(495, 495), geom.NewVec(525, 525)), aabbworld.AnyCapability, func(uid.UID64) {}); n != 0 {
			t.Errorf("Query still finds %d pieces where the box used to be", n)
		}
	})

	t.Run("into a closed corner", func(t *testing.T) {
		space := edgedSpace(t, 0)
		box := plane.NewAABB(geom.NewVec(500, 500), 20, 20)
		space.Insert(id, &box)
		space.MoveTo(id, &box, geom.NewVec(5000, -5000))

		want := plane.NewAABB(geom.NewVec(980, 0), 20, 20)
		if box != want {
			t.Errorf("box %+v, want it whole in the corner at %+v", box, want)
		}
	})

	t.Run("past an open edge", func(t *testing.T) {
		space := edgedSpace(t, aabbworld.OpenX)
		box := plane.NewAABB(geom.NewVec(500, 500), 20, 20)
		space.Insert(id, &box)
		space.Flush(nil)
		if space.MoveTo(id, &box, geom.NewVec(1500, 500)) {
			t.Fatal("not reported gone once moved past the open edge")
		}
		space.Flush(nil)
		if got := found(space); len(got) != 0 {
			t.Errorf("Query finds %v, want nothing", got)
		}
	})
}

func TestScan_SeesThroughAWrappingSeamOnly(t *testing.T) {
	observer, target := uid.UID64(1), uid.UID64(2)
	for _, tc := range []struct {
		name      string
		edges     aabbworld.Edges
		eye, mark geom.Vec
		facing    geom.Vec
		wantSeen  bool
	}{
		{"west across a wrapping X seam", aabbworld.WrapX, geom.NewVec(20, 500), geom.NewVec(940, 500), geom.NewVec(-1, 0), true},
		{"north across a closed top", aabbworld.WrapX, geom.NewVec(500, 20), geom.NewVec(500, 940), geom.NewVec(0, -1), false},
		{"north across a wrapping Y seam", aabbworld.WrapY, geom.NewVec(500, 20), geom.NewVec(500, 940), geom.NewVec(0, -1), true},
		{"west across a closed side", aabbworld.WrapY, geom.NewVec(20, 500), geom.NewVec(940, 500), geom.NewVec(-1, 0), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			space := edgedSpace(t, tc.edges)
			space.Insert(observer, ptr(plane.NewAABB(tc.eye, 10, 10)))
			space.Insert(target, ptr(plane.NewAABB(tc.mark, 10, 10)))
			space.Flush(nil)

			var view aabbworld.View
			if !space.Scan(observer, aabbworld.Cone{Direction: tc.facing, HalfAngle: math.Pi / 6, Radius: 200}, &view) {
				t.Fatal("Scan refused")
			}
			seen := false
			view.Entities(func(id uid.UID64, _ float64) { seen = seen || id == target })
			if seen != tc.wantSeen {
				t.Errorf("target seen = %v, want %v", seen, tc.wantSeen)
			}
		})
	}
}

func TestQuery_CrossesAWrappingSeamOnly_AndFiltersByCapability(t *testing.T) {
	space := edgedSpace(t, aabbworld.WrapX)
	acrossX, acrossY, solid := uid.UID64(1), uid.UID64(2), uid.UID64(3)
	space.Insert(acrossX, ptr(plane.NewAABB(geom.NewVec(10, 20), 10, 10)))
	space.Insert(acrossY, ptr(plane.NewAABB(geom.NewVec(960, 980), 10, 10)))
	space.Insert(solid, ptr(plane.NewAABB(geom.NewVec(970, 20), 10, 10)))
	space.SetCapabilities(solid, aabbworld.Plain|aabbworld.CanCollide)
	space.Flush(nil)

	corner := geom.NewAABB(geom.NewVec(950, -30), geom.NewVec(1030, 40))
	ask := func(want aabbworld.Capability) []uid.UID64 {
		var got []uid.UID64
		n := space.Query(corner, want, func(id uid.UID64) { got = append(got, id) })
		if n != len(got) {
			t.Errorf("Query returned %d, called fn %d times", n, len(got))
		}
		slices.Sort(got)
		return got
	}

	if got := ask(aabbworld.AnyCapability); !slices.Equal(got, []uid.UID64{acrossX, solid}) {
		t.Errorf("found %v, want %v: through the X seam, never through the closed top", got, []uid.UID64{acrossX, solid})
	}
	if got := ask(aabbworld.CanCollide); !slices.Equal(got, []uid.UID64{solid}) {
		t.Errorf("found %v asking for CanCollide, want only %v", got, solid)
	}
}

func TestInsert_PlacesTheBoxUnderTheEdgeRules(t *testing.T) {
	id := uid.UID64(1)

	t.Run("on a wrapping seam", func(t *testing.T) {
		space := edgedSpace(t, aabbworld.Torus)
		box := plane.NewAABB(geom.NewVec(990, 500), 20, 20)
		if !space.Insert(id, &box) {
			t.Fatal("refused on a torus")
		}
		space.Flush(nil)

		if box.Overhang.X != 10 {
			t.Errorf("overhang %v, want the caller's box wrapped 10 past the right edge", box.Overhang)
		}
		across := geom.NewAABB(geom.NewVec(2, 505), geom.NewVec(8, 510))
		if n := space.Query(across, aabbworld.AnyCapability, func(uid.UID64) {}); n != 1 {
			t.Errorf("Query finds %d pieces across the seam, want the box", n)
		}
	})

	t.Run("past a closed corner", func(t *testing.T) {
		space := edgedSpace(t, 0)
		box := plane.NewAABB(geom.NewVec(990, -10), 20, 20)
		space.Insert(id, &box)
		space.Flush(nil)

		want := plane.NewAABB(geom.NewVec(980, 0), 20, 20)
		if box != want {
			t.Errorf("box %+v, want it whole in the corner at %+v", box, want)
		}
		if got := found(space); !slices.Equal(got, []uid.UID64{id}) {
			t.Errorf("Query finds %v, want the box", got)
		}
	})

	t.Run("past an open edge", func(t *testing.T) {
		space := edgedSpace(t, aabbworld.OpenX)
		box := plane.NewAABB(geom.NewVec(1500, 500), 20, 20)
		if space.Insert(id, &box) {
			t.Fatal("accepted a box lying wholly past the open edge")
		}
		space.Flush(nil)
		if got := found(space); len(got) != 0 {
			t.Errorf("Query finds %v, want nothing", got)
		}
	})
}

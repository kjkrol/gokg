package plane

import (
	"testing"

	"github.com/kjkrol/gokg/geom"
)

// boxAt spells a box as position and size; the sibling box() in this package
// spells one as two corners.
func boxAt(x, y, w, h float64) geom.AABB {
	return geom.NewAABB(geom.NewVec(x, y), geom.NewVec(x+w, y+h))
}

func TestPenetration(t *testing.T) {
	cases := map[string]struct {
		r1, r2 geom.AABB
		want   geom.Vec
	}{
		"nowhere near each other": {
			boxAt(0, 0, 10, 10), boxAt(100, 100, 10, 10), geom.Vec{},
		},
		"sharing an edge, not overlapping": {
			boxAt(0, 0, 10, 10), boxAt(10, 0, 10, 10), geom.Vec{},
		},
		"shallow across, deep down: push sideways": {
			boxAt(0, 0, 10, 10), boxAt(8, 0, 10, 10), geom.NewVec(-2, 0),
		},
		"shallow down, deep across: push up": {
			boxAt(0, 0, 10, 10), boxAt(0, 9, 10, 10), geom.NewVec(0, -1),
		},
		"one box swallowed by the other: out its nearest edge": {
			// r2 sits inside r1, closest to r1's right edge.
			boxAt(0, 0, 100, 100), boxAt(90, 40, 5, 5), geom.NewVec(-10, 0),
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := Penetration(tc.r1, tc.r2); got != tc.want {
				t.Errorf("Penetration = %v, want %v", got, tc.want)
			}
		})
	}
}

// Applying the penetration to r1 has to actually free it, whichever axis it
// chose — the property the table above spells out case by case.
func TestPenetration_ApplyingItSeparatesTheBoxes(t *testing.T) {
	r1, r2 := boxAt(0, 0, 10, 10), boxAt(8, 4, 10, 10)
	pen := Penetration(r1, r2)

	moved := geom.NewAABB(r1.TopLeft.Add(pen), r1.BottomRight.Add(pen))
	if Penetration(moved, r2) != (geom.Vec{}) {
		t.Errorf("after applying %v the boxes still overlap: %v vs %v", pen, moved, r2)
	}
}

// This is the combination collisions.findActiveCollision skipped, on the
// grounds that it "cannot occur". It can. A sits against the right edge near
// the top, so its wrapped part reappears at the left, still near the top.
// B sits against the bottom edge near the left, so its wrapped part reappears
// at the top, still near the left. The two reappearances meet close to the
// origin — and no other pair of their images meets at all, because A's body is
// far to the right and B's body is far down.
func TestDeepestOverlap_FragmentMeetsFragment(t *testing.T) {
	space := NewToroidal2D(100, 100)

	a := space.WrapAABB(boxAt(98, 1, 10, 10)) // overhangs right only
	b := space.WrapAABB(boxAt(3, 98, 10, 10)) // overhangs bottom only

	if a.Overhang.Y != 0 || b.Overhang.X != 0 {
		t.Fatalf("fixture is wrong: a overhang %v, b overhang %v — a must wrap only right, b only bottom", a.Overhang, b.Overhang)
	}
	// Neither main box reaches the other's, nor either main box the other's fragment.
	if a.AABB.Intersects(b.AABB) {
		t.Fatal("fixture is wrong: the main boxes already meet, so the case proves nothing")
	}

	got, ok := a.DeepestOverlapWith(&b)
	if !ok {
		t.Fatal("no overlap found, though a's wrapped part stands on b's wrapped part near the origin")
	}
	if got.Penetration.X == 0 && got.Penetration.Y == 0 {
		t.Errorf("Penetration = %v, want a non-zero push", got.Penetration)
	}
	// Both sides reported must be the wrapped images, not the main boxes.
	if got.BoxA.TopLeft.X != 0 {
		t.Errorf("BoxA = %v, want a's wrapped part starting at x=0", got.BoxA)
	}
	if got.BoxB.TopLeft.Y != 0 {
		t.Errorf("BoxB = %v, want b's wrapped part starting at y=0", got.BoxB)
	}
}

// IntersectsWithFrags has always walked this combination; the point here is
// that the two now agree, because they share one walk.
func TestDeepestOverlap_AgreesWithIntersectsWithFrags(t *testing.T) {
	space := NewToroidal2D(100, 100)

	cases := map[string]struct{ a, b AABB }{
		"clear of each other":     {space.WrapAABB(boxAt(10, 10, 10, 10)), space.WrapAABB(boxAt(50, 50, 10, 10))},
		"plain overlap":           {space.WrapAABB(boxAt(10, 10, 10, 10)), space.WrapAABB(boxAt(15, 10, 10, 10))},
		"across the right seam":   {space.WrapAABB(boxAt(96, 10, 10, 10)), space.WrapAABB(boxAt(2, 10, 10, 10))},
		"across the bottom seam":  {space.WrapAABB(boxAt(10, 96, 10, 10)), space.WrapAABB(boxAt(10, 2, 10, 10))},
		"fragment meets fragment": {space.WrapAABB(boxAt(98, 1, 10, 10)), space.WrapAABB(boxAt(3, 98, 10, 10))},
		"corner against corner":   {space.WrapAABB(boxAt(97, 97, 10, 10)), space.WrapAABB(boxAt(2, 2, 10, 10))},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			_, got := tc.a.DeepestOverlapWith(&tc.b)
			if want := tc.a.IntersectsWithFrags(tc.b); got != want {
				t.Errorf("DeepestOverlapWith says %v, IntersectsWithFrags says %v", got, want)
			}
		})
	}
}

// The deepest overlap is the point: a box may graze another's main body while
// running well into its wrapped part, and it is the wrapped part that has to
// be reported — resolving against the graze would push it the wrong way.
func TestDeepestOverlap_DeeperFragmentBeatsTheMainBoxes(t *testing.T) {
	space := NewToroidal2D(16, 16)
	// b runs off the right edge, so it also stands at 0..6. a reaches 13,
	// barely touching b's main box but well inside its wrapped part.
	a := space.WrapAABB(boxAt(0, 0, 13, 10))
	b := space.WrapAABB(boxAt(12, 0, 10, 10))

	got, ok := a.DeepestOverlapWith(&b)
	if !ok {
		t.Fatal("no overlap found")
	}
	if !got.BoxA.Equals(a.AABB) {
		t.Errorf("BoxA = %v, want a's main box — a never wraps here", got.BoxA)
	}
	if got.BoxB.TopLeft.X != 0 {
		t.Errorf("BoxB = %v, want b's wrapped part starting at x=0", got.BoxB)
	}
	if abs(got.Penetration.X) <= 1 {
		t.Errorf("penetration %v across, want the fragment's deeper overlap rather than the main boxes' single unit", got.Penetration.X)
	}
}

func TestDeepestOverlap_FragmentFindsWhatTheMainBoxesMiss(t *testing.T) {
	space := NewToroidal2D(20, 20)
	// a runs off the right edge and reappears at 0..5, where b stands. Their
	// main boxes sit at opposite ends and never touch.
	a := space.WrapAABB(boxAt(15, 0, 10, 10))
	b := space.WrapAABB(boxAt(0, 0, 10, 10))

	got, ok := a.DeepestOverlapWith(&b)
	if !ok {
		t.Fatal("no contact found, though a's wrapped part stands on top of b")
	}
	if got.BoxA.TopLeft.X != 0 {
		t.Errorf("BoxA = %v, want a's wrapped part starting at x=0", got.BoxA)
	}
	if !got.BoxB.Equals(b.AABB) {
		t.Errorf("BoxB = %v, want b's main box — b never wraps here", got.BoxB)
	}
}

func TestDeepestOverlap_UnwrappedBoxesReportTheirMainBoxes(t *testing.T) {
	space := NewToroidal2D(1000, 1000)
	a := space.WrapAABB(boxAt(0, 0, 10, 10))
	b := space.WrapAABB(boxAt(8, 0, 10, 10))

	got, ok := a.DeepestOverlapWith(&b)
	if !ok {
		t.Fatal("no overlap found")
	}
	if !got.BoxA.Equals(a.AABB) || !got.BoxB.Equals(b.AABB) {
		t.Errorf("reported %v / %v, want the main boxes when neither side wraps", got.BoxA, got.BoxB)
	}
	if got.Penetration != geom.NewVec(-2, 0) {
		t.Errorf("Penetration = %v, want {-2 0}", got.Penetration)
	}
}

func TestDeepestOverlap_ReportsNothingWhenNothingMeets(t *testing.T) {
	space := NewToroidal2D(1000, 1000)
	a := space.WrapAABB(boxAt(0, 0, 10, 10))
	b := space.WrapAABB(boxAt(500, 500, 10, 10))

	if got, ok := a.DeepestOverlapWith(&b); ok {
		t.Errorf("reported an overlap %v between boxes half a world apart", got.Penetration)
	}
}

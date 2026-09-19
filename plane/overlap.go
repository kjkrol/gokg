package plane

import "github.com/kjkrol/gokg/geom"

// Penetration is the shortest translation that separates r1 from r2, or a
// zero vector when they do not overlap. Applying it to r1 moves r1 clear of
// r2 along whichever axis needs the least travel.
func Penetration(r1, r2 geom.AABB) geom.Vec {
	leftX := max(r1.TopLeft.X, r2.TopLeft.X)
	rightX := min(r1.BottomRight.X, r2.BottomRight.X)
	if rightX-leftX <= 0 {
		return geom.Vec{}
	}

	topY := max(r1.TopLeft.Y, r2.TopLeft.Y)
	bottomY := min(r1.BottomRight.Y, r2.BottomRight.Y)
	if bottomY-topY <= 0 {
		return geom.Vec{}
	}

	pushRight := r2.BottomRight.X - r1.TopLeft.X
	pushLeft := r1.BottomRight.X - r2.TopLeft.X
	pushDown := r2.BottomRight.Y - r1.TopLeft.Y
	pushUp := r1.BottomRight.Y - r2.TopLeft.Y

	minPush := pushRight
	mtv := geom.Vec{X: pushRight, Y: 0}

	if pushLeft < minPush {
		minPush = pushLeft
		mtv = geom.Vec{X: -pushLeft, Y: 0}
	}
	if pushDown < minPush {
		minPush = pushDown
		mtv = geom.Vec{X: 0, Y: pushDown}
	}
	if pushUp < minPush {
		mtv = geom.Vec{X: 0, Y: -pushUp}
	}

	return mtv
}

// Overlap is where two wrapped boxes meet: which image of each was involved,
// and how deeply they interpenetrate.
type Overlap struct {
	BoxA, BoxB  geom.AABB
	Penetration geom.Vec
}

// DeepestOverlapWith reports where ab and other meet most deeply, considering
// every wrapped image of both, and whether they meet at all.
//
// Which images are in contact is not fixed for the duration of a solve: push
// two boxes apart and a different pair of their images can become the one that
// overlaps. That is why this recomputes from current geometry rather than
// being cached alongside the pair.
//
// Both sides are pointers, unlike the by-value Intersects/Contains siblings:
// this one runs once per candidate pair per solver pass, and an AABB is not
// small enough to copy twice that often for nothing.
func (ab *AABB) DeepestOverlapWith(other *AABB) (Overlap, bool) {
	// The overwhelmingly common case is two boxes nowhere near a world edge,
	// and it deserves to stay a straight line rather than a walk.
	if !ab.HasFragments() && !other.HasFragments() {
		pen := Penetration(ab.AABB, other.AABB)
		if pen.X == 0 && pen.Y == 0 {
			return Overlap{}, false
		}
		return Overlap{BoxA: ab.AABB, BoxB: other.AABB, Penetration: pen}, true
	}

	var best Overlap
	deepest := float64(0)
	ab.visitImagePairs(*other, func(a, b geom.AABB) bool {
		pen := Penetration(a, b)
		if pen.X == 0 && pen.Y == 0 {
			return true
		}
		if d := abs(pen.X) + abs(pen.Y); d > deepest {
			deepest = d
			best = Overlap{BoxA: a, BoxB: b, Penetration: pen}
		}
		return true
	})
	return best, deepest > 0
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

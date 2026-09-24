package raycast

import (
	"math"

	"github.com/kjkrol/aabbworld/geom"
)

// hitDistance reports where a unit ray enters box (zero from inside) and leaves it, if at all.
func hitDistance(origin, dir geom.Vec, box geom.AABB) (near, far float64, ok bool) {
	near, far = math.Inf(-1), math.Inf(1)

	if dir.X == 0 {
		if origin.X < box.TopLeft.X || origin.X > box.BottomRight.X {
			return 0, 0, false
		}
	} else {
		inv := 1 / dir.X
		t1 := (box.TopLeft.X - origin.X) * inv
		t2 := (box.BottomRight.X - origin.X) * inv
		if t1 > t2 {
			t1, t2 = t2, t1
		}
		near, far = t1, t2
	}

	if dir.Y == 0 {
		if origin.Y < box.TopLeft.Y || origin.Y > box.BottomRight.Y {
			return 0, 0, false
		}
	} else {
		inv := 1 / dir.Y
		t1 := (box.TopLeft.Y - origin.Y) * inv
		t2 := (box.BottomRight.Y - origin.Y) * inv
		if t1 > t2 {
			t1, t2 = t2, t1
		}
		if t1 > near {
			near = t1
		}
		if t2 < far {
			far = t2
		}
	}

	if far < near || far < 0 {
		return 0, 0, false
	}
	return math.Max(near, 0), far, true
}

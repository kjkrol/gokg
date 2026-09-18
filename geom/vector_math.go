package geom

import "math"

// Length returns the Euclidean magnitude of v.
func Length(v Vec) float64 { return math.Sqrt(v.X*v.X + v.Y*v.Y) }

// Clamp bounds v to the inclusive box [0,size].
func Clamp(v, size Vec) Vec {
	return Vec{X: clamp(v.X, size.X), Y: clamp(v.Y, size.Y)}
}

// Wrap folds v back into [0,size) on each axis.
func Wrap(v, size Vec) Vec {
	return Vec{X: wrap(v.X, size.X), Y: wrap(v.Y, size.Y)}
}

func clamp(val, max float64) float64 {
	if val > max {
		return max
	}
	if val < 0 {
		return 0
	}
	return val
}

// wrap folds val into [0,max). The two fast paths cover what motion actually
// produces — a coordinate already inside, or one that has just stepped over an
// edge — so math.Mod is reached only by something thrown far out of bounds.
func wrap(val, max float64) float64 {
	if max == 0 {
		return val
	}
	if max < 0 {
		max = -max
	}

	if val >= 0 && val < max {
		return val
	}
	if val >= max && val < max+max {
		return val - max
	}
	if val < 0 && val >= -max {
		return val + max
	}

	r := math.Mod(val, max)
	if r < 0 {
		r += max
	}
	return r
}

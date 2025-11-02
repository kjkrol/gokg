package geometry

import "math"

// SutherlandHodgmanClipper clips polygons against an axis-aligned viewport using the Sutherland–Hodgman algorithm.
type SutherlandHodgmanClipper[T SupportedNumeric] struct {
	viewport AABB[T]
}

// NewSutherlandHodgmanClipper constructs a clipper for the supplied viewport.
func NewSutherlandHodgmanClipper[T SupportedNumeric](viewport AABB[T]) SutherlandHodgmanClipper[T] {
	return SutherlandHodgmanClipper[T]{viewport: viewport}
}

// Clip returns the polygon points clipped to the viewport. When the polygon lies completely outside, nil is returned.
func (c SutherlandHodgmanClipper[T]) Clip(points []Vec[T]) []Vec[T] {
	if len(points) == 0 {
		return nil
	}

	minX := c.viewport.TopLeft.X
	maxX := c.viewport.BottomRight.X
	minY := c.viewport.TopLeft.Y
	maxY := c.viewport.BottomRight.Y

	type edge struct {
		inside       func(Vec[T]) bool
		intersection func(Vec[T], Vec[T]) Vec[T]
	}

	edges := []edge{
		{
			inside: func(v Vec[T]) bool { return v.X >= minX },
			intersection: func(start, end Vec[T]) Vec[T] {
				return intersectVertical(start, end, minX)
			},
		},
		{
			inside: func(v Vec[T]) bool { return v.X <= maxX },
			intersection: func(start, end Vec[T]) Vec[T] {
				return intersectVertical(start, end, maxX)
			},
		},
		{
			inside: func(v Vec[T]) bool { return v.Y >= minY },
			intersection: func(start, end Vec[T]) Vec[T] {
				return intersectHorizontal(start, end, minY)
			},
		},
		{
			inside: func(v Vec[T]) bool { return v.Y <= maxY },
			intersection: func(start, end Vec[T]) Vec[T] {
				return intersectHorizontal(start, end, maxY)
			},
		},
	}

	output := append([]Vec[T](nil), points...)
	for _, e := range edges {
		if len(output) == 0 {
			return nil
		}
		input := output
		output = make([]Vec[T], 0, len(input))

		prev := input[len(input)-1]
		for _, curr := range input {
			prevInside := e.inside(prev)
			currInside := e.inside(curr)

			switch {
			case prevInside && currInside:
				output = append(output, curr)
			case !prevInside && currInside:
				output = append(output, e.intersection(prev, curr))
				output = append(output, curr)
			case prevInside && !currInside:
				output = append(output, e.intersection(prev, curr))
			}

			prev = curr
		}
	}

	return dedupeVertices(output)
}

func intersectVertical[T SupportedNumeric](start, end Vec[T], xBound T) Vec[T] {
	if start.X == end.X {
		return Vec[T]{X: xBound, Y: start.Y}
	}

	t := toFloat64(xBound-start.X) / toFloat64(end.X-start.X)
	y := toFloat64(start.Y) + t*(toFloat64(end.Y-start.Y))
	return Vec[T]{X: xBound, Y: fromFloat64[T](y)}
}

func intersectHorizontal[T SupportedNumeric](start, end Vec[T], yBound T) Vec[T] {
	if start.Y == end.Y {
		return Vec[T]{X: start.X, Y: yBound}
	}

	t := toFloat64(yBound-start.Y) / toFloat64(end.Y-start.Y)
	x := toFloat64(start.X) + t*(toFloat64(end.X-start.X))
	return Vec[T]{X: fromFloat64[T](x), Y: yBound}
}

func dedupeVertices[T SupportedNumeric](points []Vec[T]) []Vec[T] {
	if len(points) == 0 {
		return points
	}

	out := make([]Vec[T], 0, len(points))
	for _, p := range points {
		if len(out) == 0 || !sameVertex(out[len(out)-1], p) {
			out = append(out, p)
		}
	}

	if len(out) > 1 && sameVertex(out[0], out[len(out)-1]) {
		out = out[:len(out)-1]
	}

	if len(out) < 3 {
		return nil
	}
	return out
}

func sameVertex[T SupportedNumeric](a, b Vec[T]) bool {
	switch ax := any(a.X).(type) {
	case int:
		return a == b
	case float64:
		const epsilon = 1e-9
		ay := any(a.Y).(float64)
		bx := any(b.X).(float64)
		by := any(b.Y).(float64)
		return math.Abs(ax-bx) < epsilon && math.Abs(ay-by) < epsilon
	default:
		panic("unsupported numeric type")
	}
}

func toFloat64[T SupportedNumeric](value T) float64 {
	switch v := any(value).(type) {
	case int:
		return float64(v)
	case float64:
		return v
	default:
		panic("unsupported numeric type")
	}
}

func fromFloat64[T SupportedNumeric](value float64) T {
	var zero T
	switch any(zero).(type) {
	case int:
		return T(int(math.Round(value)))
	case float64:
		return T(value)
	default:
		panic("unsupported numeric type")
	}
}

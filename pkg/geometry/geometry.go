package geometry

import "math"

type Geometry[T Number] interface {
	Length(v Vec[T]) T
	Distance(v1, v2 Vec[T]) T
	Mod(v1, v2 Vec[T]) Vec[T]
	ModMutable(v1 *Vec[T], v2 Vec[T])
}

// Floating points geometry

type Float64Geometry struct{}

func (g Float64Geometry) Length(v Vec[float64]) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (g Float64Geometry) Distance(v1, v2 Vec[float64]) float64 {
	delta := v1.Sub(v2)
	return g.Length(delta)
}

func (g Float64Geometry) ModMutable(v1 *Vec[float64], v2 Vec[float64]) {
	v1.X = math.Mod(v1.X, v2.X)
	v1.Y = math.Mod(v1.Y, v2.Y)
}

func (g Float64Geometry) Mod(v1, v2 Vec[float64]) Vec[float64] {
	return Vec[float64]{math.Mod(v1.X, v2.X), math.Mod(v1.Y, v2.Y)}
}

// Integral points geometry

type IntGeometry struct{}

func (g IntGeometry) Length(v Vec[int]) int {
	return int(math.Ceil(math.Sqrt(float64(v.X*v.X + v.Y*v.Y))))
}

func (g IntGeometry) Distance(v1, v2 Vec[int]) int {
	delta := v1.Sub(v2)
	return g.Length(delta)
}

func (g IntGeometry) ModMutable(v1 *Vec[int], v2 Vec[int]) { v1.X %= v2.X; v1.Y %= v2.Y }

func (g IntGeometry) Mod(v1, v2 Vec[int]) Vec[int] { return Vec[int]{v1.X % v2.X, v1.Y % v2.Y} }

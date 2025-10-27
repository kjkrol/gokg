package geometry

import s "github.com/kjkrol/gokg/pkg/geometry/spatial"

type (
	supportedNumeric              = s.SupportedNumeric
	vec[T supportedNumeric]       = s.Vec[T]
	spatial[T supportedNumeric]   = s.Spatial[T]
	rectangle[T supportedNumeric] = s.Rectangle[T]
	line[T supportedNumeric]      = s.Line[T]
	polygon[T supportedNumeric]   = s.Polygon[T]
)

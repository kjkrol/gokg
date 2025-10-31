package geometry

func translateInPlace[T SupportedNumeric](shape Shape[T], delta Vec[T]) {
	switch s := shape.(type) {
	case *Vec[T]:
		s.AddMutable(delta)

	case *Line[T]:
		s.Start.AddMutable(delta)
		s.End.AddMutable(delta)

	case *Polygon[T]:
		for _, v := range s.Vertices() {
			if v != nil {
				v.AddMutable(delta)
			}
		}
		s.UpdateBounds()

	default:
		for _, v := range shape.Vertices() {
			if v != nil {
				v.AddMutable(delta)
			}
		}
	}
}

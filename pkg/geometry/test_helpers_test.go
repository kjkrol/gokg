package geometry

func rectEquals[T SupportedNumeric](a, b Rectangle[T]) bool {
	return a.TopLeft == b.TopLeft && a.BottomRight == b.BottomRight && a.Center == b.Center
}

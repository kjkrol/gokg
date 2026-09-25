package raycast

// Shadow is a stretch of ground along one of the angles a View reads, From to To away from the
// observer, that the observer cannot see.
type Shadow struct {
	Sample   int
	From, To float32
}

// shadowing gathers the shadows of the angle being cast from the ground points it is shown in
// order: a run of hidden points is one shadow, reaching halfway to the lit points either side, or
// to the radius when none follows.
type shadowing struct {
	on     bool
	sample int
	out    []Shadow

	lastLit     float64
	firstHidden float64
	lastHidden  float64
	hiding      bool
}

func (s *shadowing) begin() { s.lastLit, s.hiding = 0, false }

// point takes the ground at d, visible or not.
func (s *shadowing) point(d float64, visible bool) {
	switch {
	case visible && s.hiding:
		s.emit((s.lastLit+s.firstHidden)/2, (s.lastHidden+d)/2)
		s.hiding, s.lastLit = false, d
	case visible:
		s.lastLit = d
	case !s.hiding:
		s.hiding, s.firstHidden, s.lastHidden = true, d, d
	default:
		s.lastHidden = d
	}
}

// end closes a run still hidden at the radius.
func (s *shadowing) end(radius float64) {
	if s.hiding {
		s.emit((s.lastLit+s.firstHidden)/2, radius)
		s.hiding = false
	}
}

func (s *shadowing) emit(from, to float64) {
	if to > from {
		s.out = append(s.out, Shadow{Sample: s.sample, From: float32(from), To: float32(to)})
	}
}

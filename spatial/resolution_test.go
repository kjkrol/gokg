package spatial

import "testing"

func TestResolutionFrom(t *testing.T) {
	cases := []struct {
		number uint32
		want   Resolution
	}{
		{0, Size1x1},
		{1, Size1x1},
		{2, Size2x2},
		{1000, Size1024x1024}, // smallest power of two >= 1000
		{1024, Size1024x1024}, // exact power of two must not double
		{1025, Size2048x2048},
		{2048, Size2048x2048},
	}
	for _, c := range cases {
		if got := ResolutionFrom(c.number); got != c.want {
			t.Errorf("ResolutionFrom(%d) = %v, want %v", c.number, got, c.want)
		}
	}
}

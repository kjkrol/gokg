package spatial

import "testing"

func BenchmarkNewMortonCode(b *testing.B) {
	x, y := uint32(12345), uint32(54321)

	for b.Loop() {
		_ = NewMortonCode(x, y)
	}
}

func BenchmarkMortonCodeDecode(b *testing.B) {
	code := NewMortonCode(12345, 54321)

	for b.Loop() {
		_, _ = code.Decode()
	}
}

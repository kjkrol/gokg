package spatial

import (
	"fmt"
	"math/bits"
)

// -----------------------------------------
// Enum where the name reflects total cell count (side length squared).
// -----------------------------------------
type Resolution uint8

const (
	Size1x1 Resolution = iota
	Size2x2
	Size4x4
	Size8x8
	Size16x16
	Size32x32
	Size64x64
	Size128x128
	Size256x256
	Size512x512
	Size1024x1024
	Size2048x2048
)

func NewResolution(cord uint8) Resolution {
	return Resolution(cord)
}

func ResolutionFrom(number uint32) Resolution {
	if number == 0 {
		return Size1x1
	}
	bits := bits.Len32(number)
	return NewResolution(uint8(bits))
}

// Side returns the side length of the square grid, e.g. Size8x8 → 8.
func (s Resolution) Side() uint32 {
	return uint32(1) << s
}

// MaxCoord returns the maximum coordinate on each axis, i.e. side - 1.
func (s Resolution) MaxCoord() uint32 {
	return (uint32(1) << s) - 1
}

// Cells returns the total number of cells in the grid (side * side).
// To jest 1 << (2*s).
func (s Resolution) Cells() int {
	return int(1) << (s << 1)
}

// String returns a human-readable name like "1x1", "1024x1024".
func (s Resolution) String() string {
	side := s.Side()
	return fmt.Sprintf("%dx%d", side, side)
}

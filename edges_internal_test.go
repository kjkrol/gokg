package aabbworld

import (
	"testing"

	iplane "github.com/kjkrol/aabbworld/internal/plane"
)

func TestEdges_AgreeWithTheSurfaceBitForBit(t *testing.T) {
	for name, pair := range map[string][2]uint8{
		"WrapX": {uint8(WrapX), uint8(iplane.WrapX)},
		"WrapY": {uint8(WrapY), uint8(iplane.WrapY)},
		"OpenX": {uint8(OpenX), uint8(iplane.OpenX)},
		"OpenY": {uint8(OpenY), uint8(iplane.OpenY)},
		"Torus": {uint8(Torus), uint8(iplane.Torus)},
	} {
		if pair[0] != pair[1] {
			t.Errorf("%s is %#b here and %#b in the surface", name, pair[0], pair[1])
		}
	}
}

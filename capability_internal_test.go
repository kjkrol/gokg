package aabbworld

import (
	"testing"

	"github.com/kjkrol/aabbworld/internal/spatial"
)

func TestCapability_AgreesWithTheIndexBitForBit(t *testing.T) {
	for name, pair := range map[string][2]uint8{
		"Plain":         {uint8(Plain), uint8(spatial.Plain)},
		"AnyCapability": {uint8(AnyCapability), uint8(spatial.AnyCapability)},
	} {
		if pair[0] != pair[1] {
			t.Errorf("%s is %#b here and %#b in the index", name, pair[0], pair[1])
		}
	}
	if CanCollide&Plain != 0 {
		t.Errorf("CanCollide (%#b) shares a bit with Plain (%#b)", CanCollide, Plain)
	}
}

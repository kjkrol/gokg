package aabbworld

import (
	"testing"

	"github.com/kjkrol/aabbworld/internal/spatial"
)

func TestCapability_AgreesWithTheIndexBitForBit(t *testing.T) {
	for name, pair := range map[string][2]uint8{
		"Plain":         {uint8(Plain), uint8(spatial.Plain)},
		"CanCollide":    {uint8(CanCollide), uint8(spatial.CanCollide)},
		"Static":        {uint8(Static), uint8(spatial.Static)},
		"Sensor":        {uint8(Sensor), uint8(spatial.Sensor)},
		"AnyCapability": {uint8(AnyCapability), uint8(spatial.AnyCapability)},
	} {
		if pair[0] != pair[1] {
			t.Errorf("%s is %#b here and %#b in the index", name, pair[0], pair[1])
		}
	}
	if Plain != 0 {
		t.Errorf("Plain is %#b, want the zero value so a bare Item is nothing but geometry", Plain)
	}
	for name, c := range map[string]Capability{"CanCollide": CanCollide, "Static": Static, "Sensor": Sensor} {
		if c == 0 || c&(c-1) != 0 {
			t.Errorf("%s is %#b, want a single bit", name, c)
		}
	}
}

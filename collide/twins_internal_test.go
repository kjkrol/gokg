package collide

import (
	"testing"
	"unsafe"

	icollide "github.com/kjkrol/aabbworld/internal/collide"
)

func TestPair_MatchesItsInternalTwinFieldForField(t *testing.T) {
	var pair Pair
	var ipair icollide.Pair
	if unsafe.Sizeof(pair) != unsafe.Sizeof(ipair) {
		t.Fatalf("sizes differ: Pair %d vs %d", unsafe.Sizeof(pair), unsafe.Sizeof(ipair))
	}
	for name, offsets := range map[string][2]uintptr{
		"Pair.B":       {unsafe.Offsetof(pair.B), unsafe.Offsetof(ipair.B)},
		"Pair.KeyA":    {unsafe.Offsetof(pair.KeyA), unsafe.Offsetof(ipair.KeyA)},
		"Pair.KeyB":    {unsafe.Offsetof(pair.KeyB), unsafe.Offsetof(ipair.KeyB)},
		"Pair.StaticA": {unsafe.Offsetof(pair.StaticA), unsafe.Offsetof(ipair.StaticA)},
		"Pair.StaticB": {unsafe.Offsetof(pair.StaticB), unsafe.Offsetof(ipair.StaticB)},
		"Pair.Sensor":  {unsafe.Offsetof(pair.Sensor), unsafe.Offsetof(ipair.Sensor)},
	} {
		if offsets[0] != offsets[1] {
			t.Errorf("%s sits at %d here and %d in the internal twin", name, offsets[0], offsets[1])
		}
	}
}

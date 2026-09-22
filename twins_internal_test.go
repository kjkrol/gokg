package aabbworld

import (
	"testing"
	"unsafe"

	iraycast "github.com/kjkrol/aabbworld/internal/raycast"
	"github.com/kjkrol/aabbworld/internal/spatial"
)

func TestCone_MatchesItsInternalTwinFieldForField(t *testing.T) {
	var cone Cone
	var icone iraycast.Cone
	if unsafe.Sizeof(cone) != unsafe.Sizeof(icone) {
		t.Fatalf("sizes differ: Cone %d vs %d", unsafe.Sizeof(cone), unsafe.Sizeof(icone))
	}
	for name, offsets := range map[string][2]uintptr{
		"Cone.HalfAngle": {unsafe.Offsetof(cone.HalfAngle), unsafe.Offsetof(icone.HalfAngle)},
		"Cone.Radius":    {unsafe.Offsetof(cone.Radius), unsafe.Offsetof(icone.Radius)},
	} {
		if offsets[0] != offsets[1] {
			t.Errorf("%s sits at %d here and %d in the internal twin", name, offsets[0], offsets[1])
		}
	}
}

func TestItem_MatchesItsInternalTwinFieldForField(t *testing.T) {
	var item Item
	var iitem spatial.Item
	if unsafe.Sizeof(item) != unsafe.Sizeof(iitem) {
		t.Fatalf("sizes differ: Item %d vs %d", unsafe.Sizeof(item), unsafe.Sizeof(iitem))
	}
	for name, offsets := range map[string][2]uintptr{
		"Item.ID":   {unsafe.Offsetof(item.ID), unsafe.Offsetof(iitem.ID)},
		"Item.Box":  {unsafe.Offsetof(item.Box), unsafe.Offsetof(iitem.Box)},
		"Item.Caps": {unsafe.Offsetof(item.Caps), unsafe.Offsetof(iitem.Caps)},
	} {
		if offsets[0] != offsets[1] {
			t.Errorf("%s sits at %d here and %d in the internal twin", name, offsets[0], offsets[1])
		}
	}
}

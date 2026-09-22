package collide_test

import (
	"testing"

	"github.com/kjkrol/aabbworld"
	"github.com/kjkrol/aabbworld/collide"
	"github.com/kjkrol/aabbworld/geom"
	"github.com/kjkrol/aabbworld/internal/core"
	"github.com/kjkrol/aabbworld/plane"
	"github.com/kjkrol/uid"
)

func TestTick_ReadsTheSpaceItIsHanded(t *testing.T) {
	build := func() *aabbworld.Space {
		space, err := aabbworld.NewSpace(aabbworld.Config{Width: 256, Height: 256, BucketSize: 64})
		if err != nil {
			t.Fatalf("NewSpace: %v", err)
		}
		return space
	}
	crowded, empty := build(), build()
	crowded.Insert(uid.UID64(1), ptr(plane.NewAABB(geom.NewVec(10, 10), 10, 10)))
	crowded.Insert(uid.UID64(2), ptr(plane.NewAABB(geom.NewVec(18, 10), 10, 10)))
	crowded.Flush(nil)

	count := func(space *aabbworld.Space) int {
		n := 0
		var e collide.Engine
		e.Tick(space, 0.5, aabbworld.AnyCapability, 0, func(_, _ uid.UID64) (collide.Body, collide.Body, bool) {
			n++
			return collide.Body{}, collide.Body{}, false
		}, nil, nil)
		return n
	}
	if got := count(crowded); got != 1 {
		t.Errorf("found %d pairs in the space holding two touching boxes, want 1", got)
	}
	if got := count(empty); got != 0 {
		t.Errorf("found %d pairs in an empty space, want 0", got)
	}
	if core.Of(crowded) != core.Of(crowded) || core.Of(crowded) == core.Of(empty) {
		t.Error("the bridge does not hand back each space's own inside")
	}
}

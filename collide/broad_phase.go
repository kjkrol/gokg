package collide

import (
	"github.com/kjkrol/aabbworld"
	"github.com/kjkrol/aabbworld/internal/core"
	"github.com/kjkrol/aabbworld/internal/spatial"
	"github.com/kjkrol/uid"
)

// BroadPhase calls fn once per pair sharing want that may touch in a step of reach × shorter side.
func BroadPhase(space *aabbworld.Space, reach float64, want aabbworld.Capability, fn func(a, b uid.UID64)) {
	core.Of(space).Index.Pairs(reach, spatial.Capability(want), fn)
}

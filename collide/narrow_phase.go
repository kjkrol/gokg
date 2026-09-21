package collide

import (
	"slices"

	"github.com/kjkrol/aabbworld"
	"github.com/kjkrol/aabbworld/geom"
	icollide "github.com/kjkrol/aabbworld/internal/collide"
	"github.com/kjkrol/aabbworld/internal/core"
	"github.com/kjkrol/aabbworld/plane"
	"github.com/kjkrol/uid"
)

// Pair is two boxes that may be in contact; the narrow phase moves them in place through A and B.
// A static side is never moved; a Sensor pair is reported but never separated.
type Pair struct {
	A, B             *plane.AABB
	KeyA, KeyB       uint32 // one key per box across the batch; zero measures every pass
	StaticA, StaticB bool
	Sensor           bool
}

// NarrowPhase tests a batch of candidate pairs exactly and separates the ones that overlap.
// The zero value is ready, buffers are reused between batches; one serves one goroutine.
type NarrowPhase struct {
	impl icollide.Solver
	left []uid.UID64
}

// Reset empties the batch, keeping the memory.
func (n *NarrowPhase) Reset() {
	n.impl.Reset()
	n.left = n.left[:0]
}

// Add enters a pair into the batch and returns the index onContact and whose name it by.
func (n *NarrowPhase) Add(p Pair) int { return n.impl.Add(icollide.Pair(p)) }

// Separate reports each overlap once, pushes it apart, and reindexes what moved as whose names it.
func (n *NarrowPhase) Separate(
	space *aabbworld.Space, iterations int,
	onContact func(i int, pen geom.Vec), whose func(i int) (a, b uid.UID64),
) {
	inside := core.Of(space)
	n.impl.Solve(inside.Surface, iterations, onContact)
	n.left = n.left[:0]
	if whose == nil {
		return
	}
	n.impl.VisitMoved(func(i int, movedA, movedB bool) {
		p := n.impl.Pair(i)
		idA, idB := whose(i)
		if movedA && !inside.Follow(idA, p.A) {
			n.leaves(idA)
		}
		if movedB && !inside.Follow(idB, p.B) {
			n.leaves(idB)
		}
	})
}

func (n *NarrowPhase) leaves(id uid.UID64) {
	if !slices.Contains(n.left, id) {
		n.left = append(n.left, id)
	}
}

// Left is who the last Separate pushed out through an open edge; good until the next Reset.
func (n *NarrowPhase) Left() []uid.UID64 { return n.left }

package collide

import (
	"slices"

	"github.com/kjkrol/aabbworld"
	"github.com/kjkrol/aabbworld/geom"
	icollide "github.com/kjkrol/aabbworld/internal/collide"
	"github.com/kjkrol/aabbworld/internal/core"
	"github.com/kjkrol/aabbworld/internal/spatial"
	"github.com/kjkrol/aabbworld/plane"
	"github.com/kjkrol/uid"
)

// Body is one side of a candidate pair as the caller knows it: its box, and how it takes part.
type Body struct {
	Box    *plane.AABB
	Static bool // never moved; the whole push goes to the other side
	Sensor bool // reported, never separated
}

// Resolve turns a candidate pair into its two bodies; false refuses the pair.
type Resolve func(a, b uid.UID64) (Body, Body, bool)

// Touch is asked once per pair, the first pass its boxes overlap: whether the shapes inside
// really touch, and the penetration to separate them by; false drops the pair for this tick.
type Touch func(i int, pen geom.Vec) (geom.Vec, bool)

// Engine finds and separates the overlapping boxes of a Space. The zero value is ready,
// buffers are reused between ticks; one Engine serves one goroutine.
type Engine struct {
	solver  icollide.Solver
	left    []uid.UID64
	resolve Resolve
	onPair  func(a, b uid.UID64)
}

// Tick resolves every pair sharing want that may touch within a step, separates the overlapping.
func (e *Engine) Tick(space *aabbworld.Space, reach float64, want aabbworld.Capability, iterations int,
	resolve Resolve, touch Touch, onContact func(i int, pen geom.Vec)) {
	inside := core.Of(space)
	e.solver.Reset()
	e.left = e.left[:0]
	e.resolve = resolve
	if e.onPair == nil {
		e.onPair = e.add
	}
	inside.Index.Pairs(reach, spatial.Capability(want), e.onPair)
	e.resolve = nil

	e.solver.Solve(inside.Surface, iterations, icollide.Touch(touch), onContact)
	e.solver.VisitMoved(func(id uid.UID64, box *plane.AABB) {
		if !inside.Follow(id, box) && !slices.Contains(e.left, id) {
			e.left = append(e.left, id)
		}
	})
}

// Left is who the last Tick pushed out through an open edge; good until the next Tick.
func (e *Engine) Left() []uid.UID64 { return e.left }

func (e *Engine) add(a, b uid.UID64) {
	bodyA, bodyB, ok := e.resolve(a, b)
	if !ok {
		return
	}
	e.solver.Add(icollide.Pair{
		A: bodyA.Box, B: bodyB.Box, IDA: a, IDB: b,
		Flags: flagsOf(bodyA, bodyB),
	})
}

func flagsOf(a, b Body) uint8 {
	var f uint8
	if a.Static {
		f |= icollide.StaticA
	}
	if b.Static {
		f |= icollide.StaticB
	}
	if a.Sensor || b.Sensor {
		f |= icollide.Sensor
	}
	return f
}

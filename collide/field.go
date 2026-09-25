package collide

import (
	"github.com/kjkrol/aabbworld/geom"
	"github.com/kjkrol/uid"
)

// SolidField is the solid ground of a world laid out on a grid — walls, rock, whatever no one
// passes — as an Engine asks it, instead of being boxes in the Space. Solid calls visit with every
// solid box of the ground that box of entity id may overlap, for that entity (a wall may let a
// flyer over); visit returning false ends the walk. Boxes lie in the frame of the box asked about:
// across a wrapping seam, shifted to its side.
type SolidField interface {
	Solid(id uid.UID64, box geom.AABB, visit func(FieldBox) bool)
}

// FieldBox is one solid box of a SolidField: where it is, which cell of the ground it is, for the
// handler, and which of its sides face open ground — the only ways a box inside it is pushed out,
// so a box sliding along a wall of many cells never catches on the seams between them.
type FieldBox struct {
	Box  geom.AABB
	Cell uint64
	Open Sides
}

// Sides is a set of the sides of a box.
type Sides uint8

const (
	Left Sides = 1 << iota
	Right
	Top
	Bottom
)

// FieldHandler is a Handler also told about the solid ground: once per entity and cell a tick,
// TouchField asks whether they really touch, and by how much (false: not this tick), and
// ContactField tells the contact, with the push that separates the entity. A Handler that is not
// a FieldHandler is told nothing of the ground; the Engine pushes all the same.
type FieldHandler interface {
	TouchField(id uid.UID64, cell uint64, pen geom.Vec) (geom.Vec, bool)
	ContactField(id uid.UID64, cell uint64, pen geom.Vec)
}

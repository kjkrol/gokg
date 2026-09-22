// Package aabbworld is a 2D world of axis-aligned boxes: a plane with its own edge rules, a spatial
// index over the boxes placed in it, and the questions a simulation asks of them — who is near
// whom, who sees what, what overlaps. Space is the whole of it; geom and plane are its vocabulary.
package aabbworld

import (
	"fmt"

	"github.com/kjkrol/aabbworld/geom"
	"github.com/kjkrol/aabbworld/internal/core"
	iplane "github.com/kjkrol/aabbworld/internal/plane"
	"github.com/kjkrol/aabbworld/internal/spatial"
	"github.com/kjkrol/aabbworld/plane"
	"github.com/kjkrol/uid"
)

// Space is a plane with boundary rules and a spatial index over the boxes placed in it.
// Changes are queued and take effect at the next Flush.
type Space struct {
	Config
	core core.Space
}

// Capability is what may be done with an entity, as a set of up to eight bits —
// a question names the bits it wants, and only entities sharing one are answered with.
type Capability uint8

const (
	// Plain is an entity that is nothing but geometry — what one is until SetCapabilities.
	Plain Capability = 1 << 0
	// CanCollide is an entity that takes part in collisions — what collide.BroadPhase asks for.
	CanCollide Capability = 1 << 1
	// AnyCapability accepts every entity.
	AnyCapability Capability = 0xFF
)

// Edges is what happens to a box at the edges of the world, per axis; zero stops it whole.
type Edges uint8

const (
	// WrapX returns a box leaving by the left or right edge through the opposite one.
	WrapX Edges = 1 << iota
	// WrapY returns a box leaving by the top or bottom edge through the opposite one.
	WrapY
	// OpenX lets a box leave by the left or right edge; wholly outside, it is dropped from the index.
	OpenX
	// OpenY lets a box leave by the top or bottom edge; wholly outside, it is dropped from the index.
	OpenY
	// Torus wraps on both axes.
	Torus = WrapX | WrapY
)

// WrapsX reports whether the left and right edges wrap.
func (e Edges) WrapsX() bool { return e&WrapX != 0 }

// WrapsY reports whether the top and bottom edges wrap.
func (e Edges) WrapsY() bool { return e&WrapY != 0 }

// Config sizes a Space.
type Config struct {
	// Width and Height of the world, in world units.
	Width, Height uint32
	// Edges is what happens to a box at the edges of the world; the zero value stops it whole.
	Edges Edges
	// BucketSize is the side of one index bucket in world units, rounded up to a power of two.
	BucketSize uint32
	// BucketCapacity is how many entities a bucket holds before it grows.
	BucketCapacity int
}

func init() {
	core.Of = func(space any) *core.Space { return &space.(*Space).core }
}

// NewSpace builds a Space, fitting the world into the nearest power-of-two grid.
func NewSpace(cfg Config) (*Space, error) {
	if cfg.Width == 0 || cfg.Height == 0 || cfg.BucketSize == 0 {
		return nil, fmt.Errorf("invalid dimensions")
	}
	if !iplane.Edges(cfg.Edges).Valid() {
		return nil, fmt.Errorf("an axis cannot both wrap and be open")
	}
	surface := iplane.NewSurface(float64(cfg.Width), float64(cfg.Height), iplane.Edges(cfg.Edges))

	index, err := spatial.NewGridIndexManager(surface, spatial.GridIndexConfig{
		Resolution:       spatial.ResolutionFrom(max(cfg.Width, cfg.Height)),
		BucketResolution: spatial.ResolutionFrom(cfg.BucketSize),
		BucketCapacity:   cfg.BucketCapacity,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create spatial index: %w", err)
	}
	return &Space{core: core.Space{Index: index, Surface: surface}, Config: cfg}, nil
}

// Insert places box under the edge rules and queues it for the index; false if past an open edge.
func (s *Space) Insert(id uid.UID64, box *plane.AABB) bool {
	s.core.Surface.Translate(box, geom.Vec{})
	if s.core.Surface.Left(box) {
		return false
	}
	s.core.Index.QueueInsert(id, *box)
	return true
}

// Remove queues id's removal from the index.
func (s *Space) Remove(id uid.UID64) { s.core.Index.QueueRemove(id) }

// SetCapabilities records what may be done with id, from the next Flush on.
func (s *Space) SetCapabilities(id uid.UID64, c Capability) {
	s.core.Index.QueueSetCapabilities(id, spatial.Capability(c))
}

// Flush applies every queued change, reporting each changed box to onDirty if given.
func (s *Space) Flush(onDirty func(geom.AABB)) { s.core.Index.Flush(onDirty) }

// Translate moves box by delta under the edge rules; false once it has left by an open edge.
func (s *Space) Translate(id uid.UID64, box *plane.AABB, delta geom.Vec) bool {
	s.core.Surface.Translate(box, delta)
	return s.core.Follow(id, box)
}

// MoveTo puts box with its top-left at to, under the edge rules; false once it has left that way.
func (s *Space) MoveTo(id uid.UID64, box *plane.AABB, to geom.Vec) bool {
	return s.Translate(id, box, to.Sub(box.TopLeft))
}

// WrapAABB folds any rectangle into the space: wrapped pieces on a wrapping axis, else clipped.
func (s *Space) WrapAABB(box geom.AABB) plane.AABB { return s.core.Surface.WrapAABB(box) }

// Bounds is the world's size and what happens at its edges.
func (s *Space) Bounds() (width, height uint32, edges Edges) {
	return s.Config.Width, s.Config.Height, s.Config.Edges
}

// Query calls fn for every indexed piece in box sharing want, seams included, and returns how many.
func (s *Space) Query(box geom.AABB, want Capability, fn func(id uid.UID64)) int {
	area := s.core.Surface.WrapAABB(box)
	caps := spatial.Capability(want)
	found := s.core.Index.QueryRangeWith(area.AABB, caps, fn)
	area.VisitFragments(func(_ plane.FragPosition, image geom.AABB) bool {
		found += s.core.Index.QueryRangeWith(image, caps, fn)
		return true
	})
	return found
}

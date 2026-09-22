package aabbworld

import (
	"fmt"
	"unsafe"

	"github.com/kjkrol/aabbworld/collide"
	"github.com/kjkrol/aabbworld/geom"
	icollide "github.com/kjkrol/aabbworld/internal/collide"
	iplane "github.com/kjkrol/aabbworld/internal/plane"
	"github.com/kjkrol/aabbworld/internal/spatial"
	"github.com/kjkrol/aabbworld/plane"
	"github.com/kjkrol/uid"
)

// Space is a plane with edge rules and a grid over the items of its last Rebuild, read as they
// are now: a Tick that pushes them is seen by the next Query without another Rebuild.
type Space struct {
	Config
	grid    *spatial.Grid
	surface *iplane.Surface
}

// Capability is what may be done with an entity, as a set of up to eight bits —
// a question names the bits it wants, and only entities sharing one are answered with.
type Capability uint8

const (
	// Plain is an entity that is nothing but geometry; the zero value.
	Plain Capability = 0
	// CanCollide is an entity that takes part in collisions — what a CollideEngine pairs up.
	CanCollide Capability = 1 << 1
	// Static is an entity nothing shifts; a contact pushes only the other side.
	Static Capability = 1 << 2
	// Sensor is an entity only ever detected; a contact with it separates nothing.
	Sensor Capability = 1 << 3
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
	// BucketSize is the side of one grid cell in world units, rounded up to a power of two.
	BucketSize uint32
}

// Item is one box a Rebuild hands the Space: whose it is, where it lies, what may be done with it.
type Item struct {
	ID   uid.UID64
	Box  plane.AABB
	Caps Capability
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
	grid := spatial.NewGrid(surface, spatial.ResolutionFrom(max(cfg.Width, cfg.Height)), spatial.ResolutionFrom(cfg.BucketSize))
	return &Space{Config: cfg, grid: grid, surface: surface}, nil
}

// Rebuild hands the Space items, which must stay put until the next Rebuild. Indexing waits for
// the first Query or Scan that needs it, so several Rebuilds in a row cost one indexing.
func (s *Space) Rebuild(items []Item) {
	if len(items) == 0 {
		s.grid.Rebuild(nil)
		return
	}
	s.grid.Rebuild(unsafe.Slice((*spatial.Item)(unsafe.Pointer(&items[0])), len(items)))
}

// Place puts a fresh box under the edge rules; false if it lies wholly past an open edge.
func (s *Space) Place(box *plane.AABB) bool {
	s.surface.Translate(box, geom.Vec{})
	return !s.surface.Left(box)
}

// Move moves box by delta under the edge rules; false once it has left by an open edge.
func (s *Space) Move(box *plane.AABB, delta geom.Vec) bool {
	s.surface.Translate(box, delta)
	return !s.surface.Left(box)
}

// MoveTo puts box with its top-left at to, under the edge rules; false once it has left that way.
func (s *Space) MoveTo(box *plane.AABB, to geom.Vec) bool {
	return s.Move(box, to.Sub(box.TopLeft))
}

// WrapAABB folds any rectangle into the space: wrapped pieces on a wrapping axis, else clipped.
func (s *Space) WrapAABB(box geom.AABB) plane.AABB { return s.surface.WrapAABB(box) }

// Bounds is the world's size and what happens at its edges.
func (s *Space) Bounds() (width, height uint32, edges Edges) {
	return s.Config.Width, s.Config.Height, s.Config.Edges
}

// Query calls fn for every indexed piece in box sharing want, seams included, and returns how many.
func (s *Space) Query(box geom.AABB, want Capability, fn func(id uid.UID64)) int {
	area := s.surface.WrapAABB(box)
	caps := spatial.Capability(want)
	found := s.grid.Query(area.AABB, caps, fn)
	area.VisitFragments(func(_ plane.FragPosition, image geom.AABB) bool {
		found += s.grid.Query(image, caps, fn)
		return true
	})
	return found
}

// CollideEngine builds an engine over the Space's CanCollide boxes, reporting to handler.
func (s *Space) CollideEngine(handler collide.Handler, cfg collide.Config) collide.Engine {
	return icollide.New(s.grid, s.surface, handler, cfg.Reach, cfg.Iterations)
}

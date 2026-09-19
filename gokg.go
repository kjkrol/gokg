package gokg

import (
	"fmt"
	"math"

	"github.com/kjkrol/gokg/collide"
	"github.com/kjkrol/gokg/geom"
	"github.com/kjkrol/gokg/plane"
	"github.com/kjkrol/gokg/raycast"
	"github.com/kjkrol/gokg/spatial"
	"github.com/kjkrol/uid"
)

// Space represents the main physical domain of the simulation.
// It acts as a facade that synchronizes boundary-aware geometry (plane.Space2D)
// with a highly optimized spatial index (spatial.GridIndexManager), providing
// a single, safe entry point for entity manipulation and querying.
type Space struct {
	Config
	surface      plane.Space2D
	spatialIndex *spatial.GridIndexManager
}

// Config defines the properties of the Space.
type Config struct {
	// Width of the physical simulation world.
	Width uint32
	// Height of the physical simulation world.
	Height uint32
	// Toroidal determines if the world wraps around its edges (true) or clamps them (false).
	Toroidal bool
	// BucketSize determines the resolution of a single cell in the spatial grid.
	BucketSize spatial.Resolution
	// BucketCapacity is the initial number of entities a single grid bucket can hold before allocating more memory.
	BucketCapacity int
	// OpsBufferSize is the size of the channel buffer used for queuing spatial index updates.
	OpsBufferSize int
	// TrackBucketDeltas records per-bucket entity movement for later
	// consumption. Off by default — see spatial.GridIndexConfig.
	TrackBucketDeltas bool
}

// NewSpace constructs a new Space with the given Config.
// It automatically handles asymmetric world dimensions by fitting them into
// the nearest power-of-two spatial grid internally, keeping the API simple and hiding complex topology.
func NewSpace(cfg Config) (*Space, error) {
	if cfg.Width == 0 || cfg.Height == 0 || cfg.BucketSize == 0 {
		return nil, fmt.Errorf("invalid dimensions")
	}
	var surface plane.Space2D
	if cfg.Toroidal {
		surface = plane.NewToroidal2D(float64(cfg.Width), float64(cfg.Height))
	} else {
		surface = plane.NewEuclidean2D(float64(cfg.Width), float64(cfg.Height))
	}

	bucketRes := cfg.BucketSize

	// Automatically calculate the required power-of-two grid resolution based on the longest dimension
	maxDim := uint32(math.Max(float64(cfg.Width), float64(cfg.Height)))
	gridRes := spatial.ResolutionFrom(maxDim)

	indexCfg := spatial.GridIndexConfig{
		Resolution:        gridRes,
		BucketResolution:  bucketRes,
		BucketCapacity:    cfg.BucketCapacity,
		OpsBufferSize:     cfg.OpsBufferSize,
		TrackBucketDeltas: cfg.TrackBucketDeltas,
	}

	spatialIndex, err := spatial.NewGridIndexManager(surface, indexCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create spatial index: %w", err)
	}

	return &Space{
		surface:      surface,
		spatialIndex: spatialIndex,
		Config:       cfg,
	}, nil
}

// Insert adds a new entity to the space.
// It first normalizes the AABB according to the Space topology (e.g., wraps it if Toroidal)
// and then queues it for insertion into the spatial grid.
func (w *Space) Insert(id uid.UID64, aabb plane.AABB) {
	w.surface.Normalize(aabb.AABB)
	w.spatialIndex.QueueInsert(id, aabb)
}

// Remove queues the entity with the given ID for removal from the spatial grid.
func (w *Space) Remove(id uid.UID64) {
	w.spatialIndex.QueueRemove(id)
}

// WrapAABB folds an arbitrary (possibly out-of-bounds) AABB into the
// space's canonical bounds — for a toroidal space this may split it into
// up to three additional wrapped fragments (Frags/FragMask), exactly as
// Insert/Translate already do internally for entity positions. Use this
// to build a valid Query box out of a rectangle that isn't already
// known to be canonical (e.g. one derived from screen coordinates).
func (w *Space) WrapAABB(aabb geom.AABB) plane.AABB {
	return w.surface.WrapAABB(aabb)
}

// Translate moves the given AABB by the specified delta, recalculates its fragments
// based on the boundary rules, and queues a spatial index update to reflect the new position.
func (w *Space) Translate(id uid.UID64, aabb *plane.AABB, delta geom.Vec) {
	w.surface.Translate(aabb, delta)
	w.spatialIndex.QueueUpdate(id, *aabb, true)
}

// Expand grows or shrinks the given AABB by the specified margin,
// and immediately queues an update to the spatial index.
func (w *Space) Expand(id uid.UID64, aabb *plane.AABB, margin float64) {
	w.surface.Expand(aabb, margin)
	w.spatialIndex.QueueUpdate(id, *aabb, true)
}

// ExpandOnly geometrically expands the AABB without updating the spatial index.
// This is useful for creating temporary probe boxes for broad-phase queries.
func (w *Space) ExpandOnly(aabb *plane.AABB, margin float64) {
	w.surface.Expand(aabb, margin)
}

// TranslateOnly moves the AABB under the space's boundary rules without
// queuing an index update — the counterpart of ExpandOnly. Use it when a box
// is moved repeatedly before it settles (an iterative contact solver, say):
// only where it ends up is worth indexing, and queuing every intermediate
// step can outrun OpsBufferSize. Pair it with Reindex once the box is final.
func (w *Space) TranslateOnly(aabb *plane.AABB, delta geom.Vec) {
	w.surface.Translate(aabb, delta)
}

// Reindex queues one spatial-index update for a box already moved by
// TranslateOnly, so the index catches up with the geometry at the next Flush.
func (w *Space) Reindex(id uid.UID64, aabb plane.AABB) {
	w.spatialIndex.QueueUpdate(id, aabb, true)
}

// Query searches the spatial grid for all entities intersecting the provided AABB.
// The collector function fn is called for every entity found, providing its ID and the exact fragment that was hit.
func (w *Space) Query(aabb geom.AABB, fn func(id uid.UID64, frag plane.FragPosition)) int {
	return w.spatialIndex.QueryRange(aabb, fn)
}

// Bounds reports the world's dimensions and whether it wraps at its edges.
func (w *Space) Bounds() (width, height uint32, toroidal bool) {
	return w.Config.Width, w.Config.Height, w.Config.Toroidal
}

// EntryAABB returns the indexed box of a single entity.
func (w *Space) EntryAABB(id uid.UID64) (geom.AABB, bool) {
	return w.spatialIndex.EntryAABB(id)
}

// Scan fills v with what observer sees through cone, and reports whether the
// query was answerable. Read the result off v as entities, as an outline, or as
// both — see [raycast.View]. Keeping v across ticks reuses its buffers.
func (w *Space) Scan(observer uid.UID64, cone raycast.Cone, v *raycast.View) bool {
	return v.Scan(w, observer, cone)
}

// Neighbours calls fn for every entity within margin of box, box's own
// wrapped images included — without which a probe straddling a seam silently
// misses whatever lies across it.
//
// box is the caller's scratch: it is expanded in place, so hand in a buffer
// held across ticks rather than a fresh copy each time. fn may see the same
// entity twice only if that entity is itself indexed in more than one image.
//
// fn takes the same arguments Query's collector does, so it can be handed
// straight through rather than wrapped — a wrapper here would be one more
// indirect call per candidate found, on the hottest path there is.
func (w *Space) Neighbours(box *plane.AABB, margin float64, fn func(id uid.UID64, frag plane.FragPosition)) {
	w.surface.Expand(box, margin)

	w.spatialIndex.QueryRange(box.AABB, fn)
	box.VisitFragments(func(_ plane.FragPosition, image geom.AABB) bool {
		w.spatialIndex.QueryRange(image, fn)
		return true
	})
}

// Resolve pushes the pairs gathered in s apart under this space's boundary
// rules — see [collide.Solver.Solve]. The index is not told anything: walk
// s.VisitMoved afterwards and Reindex what settled somewhere new.
func (w *Space) Resolve(s *collide.Solver, iterations int, onContact func(i int, pen geom.Vec)) {
	s.Solve(w.surface, iterations, onContact)
}

// Flush processes all pending queued operations (Insert, Remove, Translate, Expand)
// and applies them to the underlying bucket grid. The onDirty callback is invoked
// for every modified bucket area, which is useful for triggering visual redraws.
func (w *Space) Flush(onDirty func(geom.AABB)) {
	w.spatialIndex.Flush(onDirty)
}

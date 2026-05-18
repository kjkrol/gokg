package gokg

import (
	"fmt"
	"math"

	"github.com/kjkrol/gokg/geom"
	"github.com/kjkrol/gokg/plane"
	"github.com/kjkrol/gokg/spatial"
)

// Space represents the main physical domain of the simulation.
// It acts as a facade that synchronizes boundary-aware geometry (plane.Space2D)
// with a highly optimized spatial index (spatial.GridIndexManager), providing
// a single, safe entry point for entity manipulation and querying.
type Space struct {
	Config
	surface      plane.Space2D[uint32]
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
}

// NewSpace constructs a new Space with the given Config.
// It automatically handles asymmetric world dimensions by fitting them into
// the nearest power-of-two spatial grid internally, keeping the API simple and hiding complex topology.
func NewSpace(cfg Config) (*Space, error) {
	if cfg.Width == 0 || cfg.Height == 0 || cfg.BucketSize == 0 {
		return nil, fmt.Errorf("invalid dimensions")
	}
	var surface plane.Space2D[uint32]
	if cfg.Toroidal {
		surface = plane.NewToroidal2D(cfg.Width, cfg.Height)
	} else {
		surface = plane.NewEuclidean2D(cfg.Width, cfg.Height)
	}

	bucketRes := cfg.BucketSize

	// Automatically calculate the required power-of-two grid resolution based on the longest dimension
	maxDim := uint32(math.Max(float64(cfg.Width), float64(cfg.Height)))
	gridRes := spatial.ResolutionFrom(maxDim)

	indexCfg := spatial.GridIndexConfig{
		Resolution:       gridRes,
		BucketResolution: bucketRes,
		BucketCapacity:   cfg.BucketCapacity,
		OpsBufferSize:    cfg.OpsBufferSize,
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
func (w *Space) Insert(id uint64, aabb plane.AABB[uint32]) {
	w.surface.Normalize(aabb.AABB)
	w.spatialIndex.QueueInsert(id, aabb)
}

// Remove queues the entity with the given ID for removal from the spatial grid.
func (w *Space) Remove(id uint64) {
	w.spatialIndex.QueueRemove(id)
}

// Translate moves the given AABB by the specified delta, recalculates its fragments
// based on the boundary rules, and queues a spatial index update to reflect the new position.
func (w *Space) Translate(id uint64, aabb *plane.AABB[uint32], delta geom.Vec[uint32]) {
	w.surface.Translate(aabb, delta)
	w.spatialIndex.QueueUpdate(id, *aabb, true)
}

// Expand grows or shrinks the given AABB by the specified margin,
// and immediately queues an update to the spatial index.
func (w *Space) Expand(id uint64, aabb *plane.AABB[uint32], margin uint32) {
	w.surface.Expand(aabb, margin)
	w.spatialIndex.QueueUpdate(id, *aabb, true)
}

// ExpandOnly geometrically expands the AABB without updating the spatial index.
// This is useful for creating temporary probe boxes for broad-phase queries.
func (w *Space) ExpandOnly(aabb *plane.AABB[uint32], margin uint32) {
	w.surface.Expand(aabb, margin)
}

// Query searches the spatial grid for all entities intersecting the provided AABB.
// The collector function fn is called for every entity found, providing its ID and the exact fragment that was hit.
func (w *Space) Query(aabb geom.AABB[uint32], fn func(id uint64, frag plane.FragPosition)) int {
	return w.spatialIndex.QueryRange(aabb, fn)
}

// Flush processes all pending queued operations (Insert, Remove, Translate, Expand)
// and applies them to the underlying bucket grid. The onDirty callback is invoked
// for every modified bucket area, which is useful for triggering visual redraws.
func (w *Space) Flush(onDirty func(geom.AABB[uint32])) {
	w.spatialIndex.Flush(onDirty)
}

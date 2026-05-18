package gokg

import (
	"fmt"
	"math"

	"github.com/kjkrol/gokg/geom"
	"github.com/kjkrol/gokg/plane"
	"github.com/kjkrol/gokg/spatial"
)

type Space struct {
	Config
	surface      plane.Space2D[uint32]
	spatialIndex *spatial.GridIndexManager
}

type Config struct {
	Width          uint32
	Height         uint32
	Toroidal       bool
	BucketSize     spatial.Resolution
	BucketCapacity int
	OpsBufferSize  int
}

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

func (w *Space) Insert(id uint64, aabb plane.AABB[uint32]) {
	w.surface.Normalize(aabb.AABB)
	w.spatialIndex.QueueInsert(id, aabb)
}

func (w *Space) Remove(id uint64) {
	w.spatialIndex.QueueRemove(id)
}

func (w *Space) Translate(id uint64, aabb *plane.AABB[uint32], delta geom.Vec[uint32]) {
	w.surface.Translate(aabb, delta)
	w.spatialIndex.QueueUpdate(id, *aabb, true)
}

func (w *Space) Expand(id uint64, aabb *plane.AABB[uint32], margin uint32) {
	w.surface.Expand(aabb, margin)
	w.spatialIndex.QueueUpdate(id, *aabb, true)
}

func (w *Space) Query(aabb geom.AABB[uint32], fn func(uint64, plane.FragPosition)) int {
	return w.spatialIndex.QueryRange(aabb, fn)
}

func (w *Space) Flush(onDirty func(geom.AABB[uint32])) {
	w.spatialIndex.Flush(onDirty)
}

package gokg

import (
	"testing"

	"github.com/kjkrol/gokg/geom"
	"github.com/kjkrol/gokg/spatial"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSpace_WrapAABB_DelegatesAndFragments guards that Space.WrapAABB
// exposes the underlying Space2D.WrapAABB — a box overflowing a toroidal
// space's bounds comes back with non-empty Frags for the wrapped part.
func TestSpace_WrapAABB_DelegatesAndFragments(t *testing.T) {
	cfg := Config{
		Width: 100, Height: 100,
		Toroidal:       true,
		BucketSize:     spatial.Size8x8,
		BucketCapacity: 10,
	}
	space, err := NewSpace(cfg)
	require.NoError(t, err)

	// straddles the right/bottom edge of the 100x100 world
	overflowing := geom.AABB[uint32]{TopLeft: geom.NewVec[uint32](95, 95), BottomRight: geom.NewVec[uint32](105, 105)}

	wrapped := space.WrapAABB(overflowing)

	assert.NotEqual(t, uint8(0), wrapped.FragMask, "expected a box overflowing the world's edge to produce fragments")
}

// TestSpace_WrapAABB_EuclideanClips guards that a non-toroidal space
// clips (rather than fragments) an out-of-bounds box, matching
// Insert/Translate's own behavior there.
func TestSpace_WrapAABB_EuclideanClips(t *testing.T) {
	cfg := Config{
		Width: 100, Height: 100,
		Toroidal:       false,
		BucketSize:     spatial.Size8x8,
		BucketCapacity: 10,
	}
	space, err := NewSpace(cfg)
	require.NoError(t, err)

	overflowing := geom.AABB[uint32]{TopLeft: geom.NewVec[uint32](95, 95), BottomRight: geom.NewVec[uint32](105, 105)}

	wrapped := space.WrapAABB(overflowing)

	assert.Equal(t, uint32(100), wrapped.BottomRight.X, "expected BottomRight.X clipped to the world edge")
	assert.Equal(t, uint32(100), wrapped.BottomRight.Y, "expected BottomRight.Y clipped to the world edge")
}

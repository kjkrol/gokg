package gokg

import (
	"testing"

	"github.com/kjkrol/gokg/geom"
	"github.com/kjkrol/gokg/plane"
	"github.com/kjkrol/gokg/spatial"
	"github.com/stretchr/testify/assert"
)

func TestSpace_Lifecycle(t *testing.T) {
	// 1. Initialize the world
	cfg := Config{
		Width:          1000,
		Height:         1000,
		Toroidal:       true,
		BucketSize:     spatial.Size1024x1024,
		BucketCapacity: 10,
	}
	space, err := NewSpace(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, space)

	entityID := uint64(42)
	// A 20x20 object at position (10,10)
	box := plane.NewAABB(geom.NewVec[uint32](10, 10), 20, 20)

	// 2. Insert
	space.Insert(entityID, box)
	space.Flush(nil) // Remember to Flush so channel operations enter the buckets!

	// 3. Query (finds the object)
	foundIDs := []uint64{}
	queryBox := geom.NewAABBAt(geom.NewVec[uint32](15, 15), 5, 5)

	space.Query(queryBox, func(id uint64, frag plane.FragPosition) {
		foundIDs = append(foundIDs, id)
	})
	assert.Contains(t, foundIDs, entityID, "Object should be found at its initial position")

	// 4. Translate
	shift := geom.NewVec[uint32](100, 0)
	space.Translate(entityID, &box, shift) // box will be updated automatically
	space.Flush(nil)

	// 5. Query at the OLD position (it's no longer there)
	foundIDs = []uint64{}
	space.Query(queryBox, func(id uint64, frag plane.FragPosition) {
		foundIDs = append(foundIDs, id)
	})
	assert.NotContains(t, foundIDs, entityID, "Object should no longer be visible at the old position")

	// 6. Query at the NEW position (it should be there)
	foundIDs = []uint64{}
	queryBoxNew := geom.NewAABBAt(geom.NewVec[uint32](115, 15), 5, 5)
	space.Query(queryBoxNew, func(id uint64, frag plane.FragPosition) {
		foundIDs = append(foundIDs, id)
	})
	assert.Contains(t, foundIDs, entityID, "Object should be found at the new position")

	// 7. Remove
	space.Remove(entityID)
	space.Flush(nil)

	// 8. Final check (empty)
	foundIDs = []uint64{}
	space.Query(queryBoxNew, func(id uint64, frag plane.FragPosition) {
		foundIDs = append(foundIDs, id)
	})
	assert.Empty(t, foundIDs, "Object should be completely removed from the spatial grid")
}

func TestSpace_ToroidalWrap(t *testing.T) {
	// Create an asymmetric, wide world (e.g., 6 horizontal "screens")
	cfg := Config{
		Width:          6000,
		Height:         1000,
		Toroidal:       true,
		BucketSize:     spatial.Size512x512,
		BucketCapacity: 10,
	}
	space, err := NewSpace(cfg)
	assert.NoError(t, err)

	entityID := uint64(99)
	// Place the object near the right edge: X=5990
	box := plane.NewAABB(geom.NewVec[uint32](5990, 50), 20, 20)

	space.Insert(entityID, box)
	space.Flush(nil)

	// Shift it right by 30 pixels. It should cross the edge (6000)
	// and wrap around to the left side of the toroidal plane at X = 20.
	shift := geom.NewVec[uint32](30, 0)
	space.Translate(entityID, &box, shift)
	space.Flush(nil)

	// Check the new box position
	assert.Equal(t, uint32(20), box.TopLeft.X, "Object should physically wrap around to position X=20")

	// Query the object on the left side of the world (around X=20)
	foundIDs := []uint64{}
	queryBoxWrapped := geom.NewAABBAt(geom.NewVec[uint32](25, 55), 5, 5)

	space.Query(queryBoxWrapped, func(id uint64, frag plane.FragPosition) {
		foundIDs = append(foundIDs, id)
	})
	assert.Contains(t, foundIDs, entityID, "Object should be flawlessly queried on the left side of the plane after wrapping")
}

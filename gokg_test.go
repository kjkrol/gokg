package gokg

import (
	"math"
	"testing"

	"github.com/kjkrol/gokg/geom"
	"github.com/kjkrol/gokg/plane"
	"github.com/kjkrol/gokg/raycast"
	"github.com/kjkrol/gokg/spatial"
	"github.com/kjkrol/uid"
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

	entityID := uid.UID64(42)
	// A 20x20 object at position (10,10)
	box := plane.NewAABB(geom.NewVec(10, 10), 20, 20)

	// 2. Insert
	space.Insert(entityID, box)
	space.Flush(nil) // Remember to Flush so channel operations enter the buckets!

	// 3. Query (finds the object)
	foundIDs := []uid.UID64{}
	queryBox := geom.NewAABBAt(geom.NewVec(15, 15), 5, 5)

	space.Query(queryBox, func(id uid.UID64, frag plane.FragPosition) {
		foundIDs = append(foundIDs, id)
	})
	assert.Contains(t, foundIDs, entityID, "Object should be found at its initial position")

	// 4. Translate
	shift := geom.NewVec(100, 0)
	space.Translate(entityID, &box, shift) // box will be updated automatically
	space.Flush(nil)

	// 5. Query at the OLD position (it's no longer there)
	foundIDs = []uid.UID64{}
	space.Query(queryBox, func(id uid.UID64, frag plane.FragPosition) {
		foundIDs = append(foundIDs, id)
	})
	assert.NotContains(t, foundIDs, entityID, "Object should no longer be visible at the old position")

	// 6. Query at the NEW position (it should be there)
	foundIDs = []uid.UID64{}
	queryBoxNew := geom.NewAABBAt(geom.NewVec(115, 15), 5, 5)
	space.Query(queryBoxNew, func(id uid.UID64, frag plane.FragPosition) {
		foundIDs = append(foundIDs, id)
	})
	assert.Contains(t, foundIDs, entityID, "Object should be found at the new position")

	// 7. Remove
	space.Remove(entityID)
	space.Flush(nil)

	// 8. Final check (empty)
	foundIDs = []uid.UID64{}
	space.Query(queryBoxNew, func(id uid.UID64, frag plane.FragPosition) {
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

	entityID := uid.UID64(99)
	// Place the object near the right edge: X=5990
	box := plane.NewAABB(geom.NewVec(5990, 50), 20, 20)

	space.Insert(entityID, box)
	space.Flush(nil)

	// Shift it right by 30 pixels. It should cross the edge (6000)
	// and wrap around to the left side of the toroidal plane at X = 20.
	shift := geom.NewVec(30, 0)
	space.Translate(entityID, &box, shift)
	space.Flush(nil)

	// Check the new box position
	assert.Equal(t, float64(20), box.TopLeft.X, "Object should physically wrap around to position X=20")

	// Query the object on the left side of the world (around X=20)
	foundIDs := []uid.UID64{}
	queryBoxWrapped := geom.NewAABBAt(geom.NewVec(25, 55), 5, 5)

	space.Query(queryBoxWrapped, func(id uid.UID64, frag plane.FragPosition) {
		foundIDs = append(foundIDs, id)
	})
	assert.Contains(t, foundIDs, entityID, "Object should be flawlessly queried on the left side of the plane after wrapping")
}

func TestSpace_Visible(t *testing.T) {
	cfg := Config{
		Width:          2000,
		Height:         2000,
		Toroidal:       false,
		BucketSize:     spatial.Size256x256,
		BucketCapacity: 10,
	}
	space, err := NewSpace(cfg)
	assert.NoError(t, err)

	const (
		guard  = uid.UID64(1)
		inView = uid.UID64(2)
		behind = uid.UID64(3)
	)
	space.Insert(guard, plane.NewAABB(geom.NewVec(500, 500), 10, 10))
	space.Insert(inView, plane.NewAABB(geom.NewVec(700, 500), 10, 10))
	space.Insert(behind, plane.NewAABB(geom.NewVec(900, 500), 10, 10))
	space.Flush(nil)

	w, h, toroidal := space.Bounds()
	assert.Equal(t, cfg.Width, w)
	assert.Equal(t, cfg.Height, h)
	assert.Equal(t, cfg.Toroidal, toroidal)

	box, ok := space.EntryAABB(inView)
	assert.True(t, ok)
	assert.Equal(t, float64(700), box.TopLeft.X)

	_, ok = space.EntryAABB(uid.UID64(999))
	assert.False(t, ok)

	cone := raycast.Cone{
		Direction: geom.NewVec(1.0, 0.0),
		HalfAngle: math.Pi / 6,
		Radius:    600,
	}
	// One scan answers both what the guard sees and where the view ends.
	var view raycast.View
	assert.True(t, space.Scan(guard, cone, &view))

	var seen []uid.UID64
	n := view.Entities(func(id uid.UID64, _ float64) { seen = append(seen, id) })
	assert.Equal(t, 1, n)
	assert.Equal(t, []uid.UID64{inView}, seen, "behind must be hidden by inView")

	outline := view.Outline(0, nil)
	assert.NotEmpty(t, outline)
	assert.Equal(t, float64(505), outline[0].X, "the fan starts at the observer")

	assert.False(t, space.Scan(uid.UID64(999), cone, &view), "unknown observer")
}

// Expand reaches the spatial index; ExpandOnly changes the box and nothing
// else, which is what makes it usable for a throwaway probe.
func TestSpace_Expand(t *testing.T) {
	space, err := NewSpace(Config{
		Width: 1000, Height: 1000,
		BucketSize: spatial.Size256x256, BucketCapacity: 8,
	})
	assert.NoError(t, err)

	const id = uid.UID64(1)
	box := plane.NewAABB(geom.NewVec(100, 100), 20, 20)
	space.Insert(id, box)
	space.Flush(nil)

	space.Expand(id, &box, 10)
	space.Flush(nil)

	indexed, ok := space.EntryAABB(id)
	assert.True(t, ok)
	assert.Equal(t, float64(90), indexed.TopLeft.X, "the indexed box grew with the margin")

	probe := plane.NewAABB(geom.NewVec(500, 500), 20, 20)
	space.ExpandOnly(&probe, 10)
	assert.Equal(t, float64(490), probe.TopLeft.X)

	_, ok = space.EntryAABB(uid.UID64(2))
	assert.False(t, ok, "ExpandOnly must not put a probe into the index")
}

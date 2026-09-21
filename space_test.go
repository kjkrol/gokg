package aabbworld

import (
	"math"
	"testing"

	"github.com/kjkrol/aabbworld/geom"
	"github.com/kjkrol/aabbworld/plane"
	"github.com/kjkrol/uid"
	"github.com/stretchr/testify/assert"
)

func TestSpace_Lifecycle(t *testing.T) {
	cfg := Config{
		Width:          1000,
		Height:         1000,
		Edges:          Torus,
		BucketSize:     1024,
		BucketCapacity: 10,
	}
	space, err := NewSpace(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, space)

	entityID := uid.UID64(42)
	box := plane.NewAABB(geom.NewVec(10, 10), 20, 20)

	space.Insert(entityID, &box)
	space.Flush(nil)

	foundIDs := []uid.UID64{}
	queryBox := geom.NewAABBAt(geom.NewVec(15, 15), 5, 5)

	space.Query(queryBox, AnyCapability, func(id uid.UID64) {
		foundIDs = append(foundIDs, id)
	})
	assert.Contains(t, foundIDs, entityID, "Object should be found at its initial position")

	shift := geom.NewVec(100, 0)
	space.Translate(entityID, &box, shift)
	space.Flush(nil)

	foundIDs = []uid.UID64{}
	space.Query(queryBox, AnyCapability, func(id uid.UID64) {
		foundIDs = append(foundIDs, id)
	})
	assert.NotContains(t, foundIDs, entityID, "Object should no longer be visible at the old position")

	foundIDs = []uid.UID64{}
	queryBoxNew := geom.NewAABBAt(geom.NewVec(115, 15), 5, 5)
	space.Query(queryBoxNew, AnyCapability, func(id uid.UID64) {
		foundIDs = append(foundIDs, id)
	})
	assert.Contains(t, foundIDs, entityID, "Object should be found at the new position")

	space.Remove(entityID)
	space.Flush(nil)

	foundIDs = []uid.UID64{}
	space.Query(queryBoxNew, AnyCapability, func(id uid.UID64) {
		foundIDs = append(foundIDs, id)
	})
	assert.Empty(t, foundIDs, "Object should be completely removed from the spatial grid")
}

func TestSpace_ToroidalWrap(t *testing.T) {
	cfg := Config{
		Width:          6000,
		Height:         1000,
		Edges:          Torus,
		BucketSize:     512,
		BucketCapacity: 10,
	}
	space, err := NewSpace(cfg)
	assert.NoError(t, err)

	entityID := uid.UID64(99)
	box := plane.NewAABB(geom.NewVec(5990, 50), 20, 20)

	space.Insert(entityID, &box)
	space.Flush(nil)

	shift := geom.NewVec(30, 0)
	space.Translate(entityID, &box, shift)
	space.Flush(nil)

	assert.Equal(t, float64(20), box.TopLeft.X, "Object should physically wrap around to position X=20")

	foundIDs := []uid.UID64{}
	queryBoxWrapped := geom.NewAABBAt(geom.NewVec(25, 55), 5, 5)

	space.Query(queryBoxWrapped, AnyCapability, func(id uid.UID64) {
		foundIDs = append(foundIDs, id)
	})
	assert.Contains(t, foundIDs, entityID, "Object should be flawlessly queried on the left side of the plane after wrapping")
}

func TestSpace_Visible(t *testing.T) {
	cfg := Config{
		Width:          2000,
		Height:         2000,
		Edges:          0,
		BucketSize:     256,
		BucketCapacity: 10,
	}
	space, err := NewSpace(cfg)
	assert.NoError(t, err)

	const (
		guard  = uid.UID64(1)
		inView = uid.UID64(2)
		behind = uid.UID64(3)
	)
	space.Insert(guard, ptr(plane.NewAABB(geom.NewVec(500, 500), 10, 10)))
	space.Insert(inView, ptr(plane.NewAABB(geom.NewVec(700, 500), 10, 10)))
	space.Insert(behind, ptr(plane.NewAABB(geom.NewVec(900, 500), 10, 10)))
	space.Flush(nil)

	w, h, edges := space.Bounds()
	assert.Equal(t, cfg.Width, w)
	assert.Equal(t, cfg.Height, h)
	assert.Equal(t, cfg.Edges, edges)

	cone := Cone{
		Direction: geom.NewVec(1.0, 0.0),
		HalfAngle: math.Pi / 6,
		Radius:    600,
	}
	var view View
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

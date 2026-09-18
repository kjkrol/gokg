package raycast_test

import (
	"fmt"
	"math"

	"github.com/kjkrol/gokg"
	"github.com/kjkrol/gokg/geom"
	"github.com/kjkrol/gokg/plane"
	"github.com/kjkrol/gokg/raycast"
	"github.com/kjkrol/gokg/spatial"
	"github.com/kjkrol/uid"
)

// A guard looking east sees the nearest entity ahead of it; anything directly
// behind that one is hidden. Scanning once answers both what the guard can see
// and where the lit region ends, and a View kept across ticks — as a field of
// whatever owns the guard, never a local — stops allocating after the first.
func ExampleView() {
	space, err := gokg.NewSpace(gokg.Config{
		Width: 1000, Height: 1000,
		BucketSize: spatial.Size256x256, BucketCapacity: 8,
	})
	if err != nil {
		panic(err)
	}

	const guard = uid.UID64(1)
	space.Insert(guard, plane.NewAABB(geom.NewVec(100, 100), 10, 10))
	space.Insert(uid.UID64(2), plane.NewAABB(geom.NewVec(300, 100), 10, 10))
	space.Insert(uid.UID64(3), plane.NewAABB(geom.NewVec(500, 100), 10, 10))
	space.Flush(nil)

	cone := raycast.Cone{
		Direction: geom.NewVec(1.0, 0.0),
		HalfAngle: math.Pi / 6,
		Radius:    600,
	}

	var view raycast.View
	var fog []geom.Vec

	grew := 0
	for tick := 1; tick <= 3; tick++ {
		before := cap(fog)
		if !space.Scan(guard, cone, &view) {
			continue
		}

		// dist is measured from the guard's centre along the ray, to the point
		// where it meets the entity.
		view.Entities(func(id uid.UID64, dist float64) {
			if tick == 1 {
				fmt.Printf("sees %d at %.0f\n", id, dist)
			}
		})

		// fog[:0] keeps the array and drops only the length, so append reuses
		// it. The result must be assigned back, as with any append.
		fog = view.Outline(0, fog[:0])
		if cap(fog) != before {
			grew++
		}
	}
	fmt.Printf("outline buffer grew %d time(s) over 3 ticks\n", grew)
	// Output:
	// sees 2 at 195
	// outline buffer grew 1 time(s) over 3 ticks
}

package aabbworld_test

import (
	"fmt"

	"github.com/kjkrol/aabbworld"
	"github.com/kjkrol/aabbworld/geom"
	"github.com/kjkrol/aabbworld/plane"
	"github.com/kjkrol/uid"
)

func ExampleSpace_Translate() {
	space, err := aabbworld.NewSpace(aabbworld.Config{
		Width: 16, Height: 16, Edges: aabbworld.Torus, BucketSize: 4, BucketCapacity: 4,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	id := uid.UID64(1)
	box := space.WrapAABB(geom.NewAABBAt(geom.NewVec(0, 0), 2, 2))
	space.Insert(id, &box)
	space.Translate(id, &box, geom.NewVec(-1, -1))

	fmt.Printf("New position: %s\n", box)
	box.VisitFragments(func(pos plane.FragPosition, piece geom.AABB) bool {
		fmt.Printf("- Fragment %d: %s\n", pos, piece)
		return true
	})
	// Output:
	// New position: {(15,15) (16,16)}
	// - Fragment 1: {(0,15) (1,16)}
	// - Fragment 2: {(15,0) (16,1)}
	// - Fragment 3: {(0,0) (1,1)}
}

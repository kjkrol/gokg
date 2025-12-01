package plane_test

import (
	"fmt"

	"github.com/kjkrol/gokg/pkg/geom"
	"github.com/kjkrol/gokg/pkg/plane"
)

// Demonstrates how shifting a contiguous box beyond the toroidal plane boundary
// causes it to fragment into multiple wrapped pieces and prints those fragments.
func ExampleSpace2D_Translate() {
	cyclicPlane := plane.NewToroidal2D(10, 10)

	box := geom.NewAABBAt(geom.NewVec(0, 0), 2, 2)
	planeBox := cyclicPlane.WrapAABB(box)

	shift := geom.NewVec(-1, -1)
	cyclicPlane.Translate(&planeBox, shift)

	fragments := map[plane.FragPosition]geom.AABB[int]{}
	planeBox.VisitFragments(func(pos plane.FragPosition, box geom.AABB[int]) bool {
		fragments[pos] = box
		return true
	})
	if len(fragments) < 3 {
		fmt.Printf("Unexpected fragment count (%d)\n", len(fragments))
		return
	}

	fmt.Printf("New position: %s\n", planeBox)
	if fragment, ok := fragments[plane.FRAG_RIGHT]; ok {
		fmt.Printf("- Fragment %d: %s\n", plane.FRAG_RIGHT, fragment)
	}
	if fragment, ok := fragments[plane.FRAG_BOTTOM]; ok {
		fmt.Printf("- Fragment %d: %s\n", plane.FRAG_BOTTOM, fragment)
	}
	if fragment, ok := fragments[plane.FRAG_BOTTOM_RIGHT]; ok {
		fmt.Printf("- Fragment %d: %s\n", plane.FRAG_BOTTOM_RIGHT, fragment)
	}
	// Output:
	// New position: {(9,9) (10,10)}
	// - Fragment 0: {(0,9) (1,10)}
	// - Fragment 1: {(9,0) (10,1)}
	// - Fragment 2: {(0,0) (1,1)}
}

package geometry_test

import (
	"fmt"

	"github.com/kjkrol/gokg/pkg/geometry"
)

// Demonstrates how shifting a contiguous box beyond the cyclic plane boundary
// causes it to fragment into multiple wrapped pieces and prints those fragments.
func ExamplePlane_Translate() {
	cyclicPlane := geometry.NewCyclicBoundedPlane(10, 10)

	box := geometry.NewBoundingBoxAt(geometry.NewVec(0, 0), 2, 2)
	planeBox := cyclicPlane.WrapBoundingBox(box)

	shift := geometry.NewVec(-1, -1)
	cyclicPlane.Translate(&planeBox, shift)

	fragments := planeBox.Fragments()
	if len(fragments) < 3 {
		fmt.Printf("Unexpected fragment count (%d)\n", len(fragments))
		return
	}

	fmt.Printf("New position: %s\n", planeBox)
	if fragment, ok := fragments[geometry.FRAG_RIGHT]; ok {
		fmt.Printf("- Fragment %d: %s\n", geometry.FRAG_RIGHT, fragment)
	}
	if fragment, ok := fragments[geometry.FRAG_BOTTOM]; ok {
		fmt.Printf("- Fragment %d: %s\n", geometry.FRAG_BOTTOM, fragment)
	}
	if fragment, ok := fragments[geometry.FRAG_BOTTOM_RIGHT]; ok {
		fmt.Printf("- Fragment %d: %s\n", geometry.FRAG_BOTTOM_RIGHT, fragment)
	}
	// Output:
	// New position: {(9,9) (10,10)}
	// - Fragment 0: {(0,9) (1,10)}
	// - Fragment 1: {(9,0) (10,1)}
	// - Fragment 2: {(0,0) (1,1)}
}

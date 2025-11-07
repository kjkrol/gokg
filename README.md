# GOKG

*aka "Golang kjkrol Geometry"*

**GOKG** is a Go toolkit focused on practical 2D computational geometry.  
It centres on two core primitives — `Vec` for points/directions and axis-aligned bounding boxes (`BoundingBox`) — plus a plane-aware wrapper (`PlaneBox`) that keeps those boxes canonical to a selected plane. `PlaneBox` stores size, fragments, and helper logic so translations, wraps, and clamps obey the plane’s rules without forcing the caller to reimplement them.  
The library also models both finite and cyclic/“infinite” planes, providing operations that classify, intersect, and reconcile geometric entities under each boundary model.  
The focus remains purely on mathematical geometry; rendering or UI concerns live in neighbouring packages.

## PlaneBox boundary handling

- Finite planes clamp plane boxes to the viewport while keeping their size consistent, so
  expansions never bleed beyond the defined world.
- Cyclic planes automatically wrap plane boxes that cross an edge and split them into
  fragments (`Fragments()`) that continue on the opposite side, making toroidal
  worlds easy to model.
- The helper methods `Translate` and `Expand` renormalise plane boxes on every call,
  updating cached fragments and ensuring touch/collision queries remain accurate
  without extra bookkeeping.

## Usage example

```go
package main

import (
	"fmt"

	"github.com/kjkrol/gokg/pkg/geometry"
)

// Demonstrates how shifting a contiguous box beyond the cyclic plane boundary
// causes it to fragment into multiple wrapped pieces and prints those fragments.
func main() {
	cyclicPlane := geometry.NewCyclicBoundedPlane(10, 10)

	planeBox := geometry.NewPlaneBoxFromBox(
		geometry.NewBoundingBoxAt(geometry.NewVec(0, 0), 2, 2),
	)

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
```

For more scenarios, browse the example-based tests under `pkg/geometry`, which double as runnable documentation.

## Projects using GOKG

- `gokq` — quadtree utilities that rely on `Vec`, `BoundingBox`, and `PlaneBox` operations.
- `gokx` — graphical experiments that consume the geometric primitives from this package.

# GOKG

*aka "Golang kjkrol Geometry"*

**GOKG** is a Go toolkit focused on practical 2D computational geometry.  
It centres on two core primitives — `Vec` for points and directions, and axis-aligned
bounding boxes (`BoundingBox`) — together with a plane-aware overlay (`PlaneBox`) that
transforms, combines, and analyses them.  
The library also models both finite and infinite planes, providing operations that
help classify, intersect, and project geometric entities.  
The focus remains purely on mathematical geometry; rendering or UI concerns live in
neighbouring packages.

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

func main() {
	position := geometry.NewVec(3.0, 4.0)
	velocity := geometry.NewVec(1.5, -0.5)

	next := position.Add(velocity)

	region := geometry.NewPlaneBox(position, 2.0, 2.0)

	plane := geometry.NewBoundedPlane[float64](100, 100)
	plane.Translate(&region, velocity)

	fmt.Printf("next position: %v\n", next)
	pointRegion := geometry.NewPlaneBox(next, 0, 0)
	fmt.Printf("region contains next? %v\n", region.Contains(pointRegion))
}
```

For more scenarios, browse the example-based tests under `pkg/geometry`, which double as runnable documentation.

## Projects using GOKG

- `gokq` — quadtree utilities that rely on `Vec`, `BoundingBox`, and `PlaneBox` operations.
- `gokx` — graphical experiments that consume the geometric primitives from this package.

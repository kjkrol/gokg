package geometry_test

import (
	"fmt"

	"github.com/kjkrol/gokg/pkg/geometry"
	"github.com/kjkrol/gokg/pkg/geometry/spatial"
)

// ExampleNewBoundedPlane_int demonstrates how to use NewBoundedPlane with int type.
func ExampleNewBoundedPlane_int() {
	plane := geometry.NewBoundedPlane[int](10, 10)
	vec := spatial.Vec[int]{X: 5, Y: 5}
	delta := spatial.Vec[int]{X: 3, Y: 4}
	plane.Translate(&vec, delta)
	fmt.Println(vec)
	// Output: (8,9)
}

// ExampleNewBoundedPlane_float64 demonstrates how to use NewBoundedPlane with float64 type.
func ExampleNewBoundedPlane_float64() {
	plane := geometry.NewBoundedPlane[float64](10.0, 10.0)
	vec := spatial.Vec[float64]{X: 5.0, Y: 5.0}
	delta := spatial.Vec[float64]{X: 3.0, Y: 4.0}
	plane.Translate(&vec, delta)
	fmt.Println(vec)
	// Output: (8,9)
}

// ExampleNewCyclicBoundedPlane_int demonstrates how to use NewCyclicBoundedPlane with int type.
func ExampleNewCyclicBoundedPlane_int() {
	plane := geometry.NewCyclicBoundedPlane[int](10, 10)
	vec := spatial.Vec[int]{X: 9, Y: 9}
	delta := spatial.Vec[int]{X: 3, Y: 4}
	plane.Translate(&vec, delta)
	fmt.Println(vec)
	// Output: (2,3)
}

// ExampleNewCyclicBoundedPlane_float64 demonstrates how to use NewCyclicBoundedPlane with float64 type.
func ExampleNewCyclicBoundedPlane_float64() {
	plane := geometry.NewCyclicBoundedPlane[float64](10.0, 10.0)
	vec := spatial.Vec[float64]{X: 9.0, Y: 9.0}
	delta := spatial.Vec[float64]{X: 3.0, Y: 4.0}
	plane.Translate(&vec, delta)
	fmt.Println(vec)
	// Output: (2,3)
}

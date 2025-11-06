package geometry

import (
	"fmt"
)

// ExampleNewBoundedPlane_int demonstrates how to use NewBoundedPlane with int type.
func ExampleNewBoundedPlane_int() {
	plane := NewBoundedPlane(10, 10)
	vec := NewVec(5, 5)
	delta := NewVec(3, 4)
	vec.AddMutable(delta)
	plane.normalize(&vec)
	fmt.Println(vec)
	// Output: (8,9)
}

// ExampleNewBoundedPlane_float64 demonstrates how to use NewBoundedPlane with float64 type.
func ExampleNewBoundedPlane_float64() {
	plane := NewBoundedPlane(10.0, 10.0)
	vec := NewVec(5., 5)
	delta := NewVec(3., 4)
	vec.AddMutable(delta)
	plane.normalize(&vec)
	fmt.Println(vec)
	// Output: (8,9)
}

// ExampleNewCyclicBoundedPlane_int demonstrates how to use NewCyclicBoundedPlane with int type.
func ExampleNewCyclicBoundedPlane_int() {
	plane := NewCyclicBoundedPlane(10, 10)
	vec := NewVec(9, 9)
	delta := NewVec(3, 4)
	vec.AddMutable(delta)
	plane.normalize(&vec)
	fmt.Println(vec)
	// Output: (2,3)
}

// ExampleNewCyclicBoundedPlane_float64 demonstrates how to use NewCyclicBoundedPlane with float64 type.
func ExampleNewCyclicBoundedPlane_float64() {
	plane := NewCyclicBoundedPlane(10.0, 10.0)
	vec := NewVec(9., 9)
	delta := NewVec(3., 4)
	vec.AddMutable(delta)
	plane.normalize(&vec)
	fmt.Println(vec)
	// Output: (2,3)
}

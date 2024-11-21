package geometry_test

import (
	"fmt"

	"github.com/kjkrol/gokg/pkg/geometry"
)

// ExampleVec_Add demonstrates how to use the Add method.
func ExampleVec_Add() {
	// Create two vectors
	v1 := geometry.Vec[float64]{X: 1.0, Y: 2.0}
	v2 := geometry.Vec[float64]{X: 3.0, Y: 4.0}

	// Add the vectors
	result := v1.Add(v2)

	// Output the result
	fmt.Println(result)
	// Output: (4,6)
}

// ExampleVec_Sub demonstrates how to use the Sub method.
func ExampleVec_Sub() {
	// Create two vectors
	v1 := geometry.Vec[float64]{X: 5.0, Y: 7.0}
	v2 := geometry.Vec[float64]{X: 2.0, Y: 3.0}

	// Subtract the vectors
	result := v1.Sub(v2)

	// Output the result
	fmt.Println(result)
	// Output: (3,4)
}

// ExampleVec_AddMutable demonstrates how to use the AddMutable method.
func ExampleVec_AddMutable() {
	// Create two vectors
	v1 := geometry.Vec[float64]{X: 1.0, Y: 2.0}
	v2 := geometry.Vec[float64]{X: 3.0, Y: 4.0}

	// Add the vectors (mutating v1)
	v1.AddMutable(v2)

	// Output the result
	fmt.Println(v1)
	// Output: (4,6)
}

// ExampleVec_SubMutable demonstrates how to use the SubMutable method.
func ExampleVec_SubMutable() {
	// Create two vectors
	v1 := geometry.Vec[float64]{X: 5.0, Y: 7.0}
	v2 := geometry.Vec[float64]{X: 2.0, Y: 3.0}

	// Subtract the vectors (mutating v1)
	v1.SubMutable(v2)

	// Output the result
	fmt.Println(v1)
	// Output: (3,4)
}

// Example_float64Geometry_Lenght demonstrates how to use the Length method
// of the FLOAT_64_GEOMETRY type to calculate the length of a vector with
// float64 components. In this example, the vector (3.0, 4.0) is used, and
// the expected output is 5.
func Example_float64Geometry_Lenght() {
	floatGeometry := geometry.FLOAT_64_GEOMETRY
	v1 := geometry.Vec[float64]{X: 3.0, Y: 4.0}

	// Calculate length
	length := floatGeometry.Length(v1)
	fmt.Println(length)
	// Output: 5
}

// Example_float64Geometry_Distance demonstrates how to use the Distance method
// of the FLOAT_64_GEOMETRY type to calculate the distance between two vectors
// of type Vec[float64]. The example creates two vectors, v1 and v2, and calculates
// the distance between them, which is then printed to the console. The expected
// output is 5.
func Example_float64Geometry_Distance() {
	floatGeometry := geometry.FLOAT_64_GEOMETRY
	v1 := geometry.Vec[float64]{X: 3.0, Y: 4.0}
	v2 := geometry.Vec[float64]{X: 6.0, Y: 8.0}

	// Calculate distance
	distance := floatGeometry.Distance(v1, v2)
	fmt.Println(distance)
	// Output: 5
}

// Example_float64Geometry_Mod demonstrates the usage of the Mod method
// from the FLOAT_64_GEOMETRY type in the geometry package. It performs
// a modulus operation on two Vec[float64] vectors and prints the result.
func Example_float64Geometry_Mod() {
	floatGeometry := geometry.FLOAT_64_GEOMETRY
	v1 := geometry.Vec[float64]{X: 3.0, Y: 4.0}
	v2 := geometry.Vec[float64]{X: 6.0, Y: 8.0}

	// Modulus operation
	mod := floatGeometry.Mod(v1, v2)
	fmt.Println(mod)
	// Output: (3,4)
}

// Example_intGeometry_Length demonstrates how to use the Length method
// of the INT_GEOMETRY type from the geometry package to calculate the
// length of a vector with integer coordinates.
func Example_intGeometry_Length() {
	intGeometry := geometry.INT_GEOMETRY
	v1 := geometry.Vec[int]{X: 3, Y: 4}

	// Calculate length
	length := intGeometry.Length(v1)
	fmt.Println(length)
	// Output: 5
}

// Example_intGeometry_Distance demonstrates how to use the Distance method
// of the intGeometry type to calculate the distance between two integer vectors.
// It creates two vectors, v1 and v2, and calculates the distance between them,
// printing the result. The expected output is 5.
func Example_intGeometry_Distance() {
	intGeometry := geometry.INT_GEOMETRY
	v1 := geometry.Vec[int]{X: 3, Y: 4}
	v2 := geometry.Vec[int]{X: 6, Y: 8}

	// Calculate distance
	distance := intGeometry.Distance(v1, v2)
	fmt.Println(distance)
	// Output: 5
}

// Example_intGeometry_Mod demonstrates the usage of the Mod method from the
// INT_GEOMETRY type in the geometry package. It performs a modulus operation
// on two vectors v1 and v2, and prints the result.
func Example_intGeometry_Mod() {
	intGeometry := geometry.INT_GEOMETRY
	v1 := geometry.Vec[int]{X: 3, Y: 4}
	v2 := geometry.Vec[int]{X: 6, Y: 8}

	// Modulus operation
	mod := intGeometry.Mod(v1, v2)
	fmt.Println(mod)
	// Output: (3,4)
}

// ExampleBoundedPlane_Translate demonstrates how to use the BoundedPlane Translate method.
func ExampleBoundedPlane_Translate() {
	plane := geometry.NewDiscreteBoundedPlane(10, 10)
	v1 := geometry.Vec[int]{X: 3, Y: 4}
	delta := geometry.Vec[int]{X: 2, Y: 2}

	// Translate vector
	plane.Translate(&v1, delta)
	fmt.Println(v1)
	// Output: (5,6)
}

// ExampleBoundedPlane_Contains demonstrates how to use the BoundedPlane Contains method.
func ExampleBoundedPlane_Contains() {
	plane := geometry.NewDiscreteBoundedPlane(10, 10)
	v1 := geometry.Vec[int]{X: 5, Y: 6}

	// Check if vector is contained in the plane
	contains := plane.Contains(v1)
	fmt.Println(contains)
	// Output: true
}

// ExampleBoundedPlane_Max demonstrates how to use the BoundedPlane Max methods.
func ExampleBoundedPlane_Max() {
	plane := geometry.NewDiscreteBoundedPlane(10, 10)

	// Get max vector
	max := plane.Max()
	fmt.Println(max)
	// Output: (9,9)
}

// ExampleBoundedPlane_Min demonstrates how to use the BoundedPlane Min methods.
func ExampleBoundedPlane_Min() {
	plane := geometry.NewDiscreteBoundedPlane(10, 10)

	// Get min vector
	min := plane.Min()
	fmt.Println(min)
	// Output: (0,0)
}

// ExampleCyclicBoundedPlane_Translate demonstrates how to use the CyclicBoundedPlane Translate method.
func ExampleCyclicBoundedPlane_Translate() {
	plane := geometry.NewDiscreteCyclicBoundedPlane(10, 10)
	v1 := geometry.Vec[int]{X: 9, Y: 9}
	delta := geometry.Vec[int]{X: 2, Y: 2}

	// Translate vector
	plane.Translate(&v1, delta)
	fmt.Println(v1)
	// Output: (1,1)
}

// ExampleCyclicBoundedPlane_Contains demonstrates how to use the CyclicBoundedPlane Contains method.
func ExampleCyclicBoundedPlane_Contains() {
	plane := geometry.NewDiscreteCyclicBoundedPlane(10, 10)
	v1 := geometry.Vec[int]{X: 1, Y: 1}

	// Check if vector is contained in the plane
	contains := plane.Contains(v1)
	fmt.Println(contains)
	// Output: true
}

// ExampleCyclicBoundedPlane_Max demonstrates how to use the CyclicBoundedPlane Max methods.
func ExampleCyclicBoundedPlane_Max() {
	plane := geometry.NewDiscreteCyclicBoundedPlane(10, 10)

	// Get max vector
	max := plane.Max()
	fmt.Println(max)
	// Output: (9,9)
}

// ExampleCyclicBoundedPlane_Min demonstrates how to use the CyclicBoundedPlane Max methods.
func ExampleCyclicBoundedPlane_Min() {
	plane := geometry.NewDiscreteCyclicBoundedPlane(10, 10)

	// Get main vector
	min := plane.Min()
	fmt.Println(min)
	// Output: (0,0)
}

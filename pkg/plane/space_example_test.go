package plane

import (
	"fmt"

	"github.com/kjkrol/gokg/pkg/geom"
)

// ExampleNormalizeVecCartesian_int demonstrates how to use NewBoundedPlane with int type.
func ExampleNewCartesian_normalizeVec_int() {
	plane := NewCartesian(10, 10)
	vec := geom.NewVec(5, 5)
	delta := geom.NewVec(3, 4)
	vec.AddMutable(delta)
	plane.(space2d[int]).normalizeVec(&vec)
	fmt.Println(vec)
	// Output: (8,9)
}

// ExampleNewCartesianSpace_float64 demonstrates how to use NewBoundedPlane with float64 type.
func ExampleNewCartesian_normalizeVec_float64() {
	cartesian := NewCartesian(10.0, 10.0)
	vec := geom.NewVec(5., 5)
	delta := geom.NewVec(3., 4)
	vec.AddMutable(delta)
	cartesian.(space2d[float64]).normalizeVec(&vec)
	fmt.Println(vec)
	// Output: (8,9)
}

// ExampleNewTorus_normalizeVec_int demonstrates how to use NewCyclicBoundedPlane with int type.
func ExampleNewTorus_normalizeVec_int() {
	torus := NewTorus(10, 10)
	vec := geom.NewVec(9, 9)
	delta := geom.NewVec(3, 4)
	vec.AddMutable(delta)
	torus.(space2d[int]).normalizeVec(&vec)
	fmt.Println(vec)
	// Output: (2,3)
}

// ExampleNewTorus_normalizeVec_float64 demonstrates how to use NewCyclicBoundedPlane with float64 type.
func ExampleNewTorus_normalizeVec_float64() {
	torus := NewTorus(10.0, 10.0)
	vec := geom.NewVec(9., 9)
	delta := geom.NewVec(3., 4)
	vec.AddMutable(delta)
	torus.(space2d[float64]).normalizeVec(&vec)
	fmt.Println(vec)
	// Output: (2,3)
}

package plane

import (
	"fmt"

	"github.com/kjkrol/gokg/pkg/geom"
)

// ExampleNormalizeVecEuclidean2D_int demonstrates how to use NewBoundedPlane with int type.
func ExampleNewEuclidean2D_normalizeVec_int() {
	plane := NewEuclidean2D(10, 10)
	vec := geom.NewVec(5, 5)
	delta := geom.NewVec(3, 4)
	vec.AddMutable(delta)
	vec = plane.(*euclidean2d[int]).normalizeVec(vec)
	fmt.Println(vec)
	// Output: (8,9)
}

// ExampleNewEuclidean2D_normalizeVec_float64 demonstrates how to use NewBoundedPlane with float64 type.
func ExampleNewEuclidean2D_normalizeVec_float64() {
	euclidean := NewEuclidean2D(10.0, 10.0)
	vec := geom.NewVec(5., 5)
	delta := geom.NewVec(3., 4)
	vec.AddMutable(delta)
	vec = euclidean.(*euclidean2d[float64]).normalizeVec(vec)
	fmt.Println(vec)
	// Output: (8,9)
}

// ExampleNewToroidal2D_normalizeVec_int demonstrates how to use NewToroidal2D with int type.
func ExampleNewToroidal2D_normalizeVec_int() {
	toroidal := NewToroidal2D(10, 10)
	vec := geom.NewVec(9, 9)
	delta := geom.NewVec(3, 4)
	vec.AddMutable(delta)
	vec = toroidal.(*toroidal2d[int]).normalizeVec(vec)
	fmt.Println(vec)
	// Output: (2,3)
}

// ExampleNewToroidal2D_normalizeVec_float64 demonstrates how to use NewToroidal2D with float64 type.
func ExampleNewToroidal2D_normalizeVec_float64() {
	toroidal := NewToroidal2D(10.0, 10.0)
	vec := geom.NewVec(9., 9)
	delta := geom.NewVec(3., 4)
	vec.AddMutable(delta)
	vec = toroidal.(*toroidal2d[float64]).normalizeVec(vec)
	fmt.Println(vec)
	// Output: (2,3)
}

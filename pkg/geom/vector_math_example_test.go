package geom_test

import (
	"fmt"

	"github.com/kjkrol/gokg/pkg/geom"
)

// ExampleVectorMath_Length demonstrates how to use the Length method.
func ExampleVectorMath_Length() {
	v := geom.Vec[float64]{X: 3.0, Y: 4.0}
	vm := geom.VectorMathByType[float64]()
	length := vm.Length(v)
	fmt.Println(length)
	// Output: 5
}

// ExampleVectorMath_Clamp demonstrates how to use the Clamp method.
func ExampleVectorMath_Clamp() {
	v := geom.Vec[float64]{X: 5.5, Y: 7.5}
	size := geom.Vec[float64]{X: 5.0, Y: 5.0}
	vm := geom.VectorMathByType[float64]()
	vm.Clamp(&v, size)
	fmt.Println(v)
	// Output: (5,5)
}

// ExampleVectorMath_Wrap demonstrates how to use the Wrap method.
func ExampleVectorMath_Wrap() {
	v := geom.Vec[float64]{X: 5.5, Y: 9.5}
	size := geom.Vec[float64]{X: 5.0, Y: 5.0}
	vm := geom.VectorMathByType[float64]()
	vm.Wrap(&v, size)
	fmt.Println(v)
	// Output: (0.5,4.5)
}

// Example_intVectorMath_Length demonstrates how to use the Length method with int type.
func Example_intVectorMath_Length() {
	v := geom.Vec[int]{X: 3, Y: 4}
	vm := geom.VectorMathByType[int]()
	length := vm.Length(v)
	fmt.Println(length)
	// Output: 5
}

// Example_intVectorMath_Clamp demonstrates how to use the Clamp method with int type.
func Example_intVectorMath_Clamp() {
	v := geom.Vec[int]{X: 6, Y: 5}
	size := geom.Vec[int]{X: 5, Y: 5}
	vm := geom.VectorMathByType[int]()
	vm.Clamp(&v, size)
	fmt.Println(v)
	// Output: (5,5)
}

// Example_intVectorMath_Wrap demonstrates how to use the Wrap method with int type.
func Example_intVectorMath_Wrap() {
	v := geom.Vec[int]{X: 7, Y: 9}
	size := geom.Vec[int]{X: 5, Y: 5}
	vm := geom.VectorMathByType[int]()
	vm.Wrap(&v, size)
	fmt.Println(v)
	// Output: (2,4)
}

// Example_intVectorMath_Wrap demonstrates how to use the Wrap method with int type.
func Example_intVectorMath_Wrap_minus() {
	target := geom.NewVec(2, 2)
	v := geom.Vec[int]{X: -101, Y: -101}
	size := geom.Vec[int]{X: 10, Y: 10}
	vm := geom.VectorMathByType[int]()

	target.AddMutable(v)
	vm.Wrap(&target, size)

	fmt.Println(target)
	// Output: (1,1)
}

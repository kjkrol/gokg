package geometry_test

import (
	"fmt"

	"github.com/kjkrol/gokg/pkg/geometry"
	"github.com/kjkrol/gokg/pkg/geometry/spatial"
)

// ExampleVectorMath_Length demonstrates how to use the Length method.
func ExampleVectorMath_Length() {
	v := spatial.Vec[float64]{X: 3.0, Y: 4.0}
	vm := geometry.VectorMathByType[float64]()
	length := vm.Length(v)
	fmt.Println(length)
	// Output: 5
}

// ExampleVectorMath_Clamp demonstrates how to use the Clamp method.
func ExampleVectorMath_Clamp() {
	v := spatial.Vec[float64]{X: 5.5, Y: 7.5}
	size := spatial.Vec[float64]{X: 5.0, Y: 5.0}
	vm := geometry.VectorMathByType[float64]()
	vm.Clamp(&v, size)
	fmt.Println(v)
	// Output: (4.9999,4.9999)
}

// ExampleVectorMath_Wrap demonstrates how to use the Wrap method.
func ExampleVectorMath_Wrap() {
	v := spatial.Vec[float64]{X: 5.5, Y: 9.5}
	size := spatial.Vec[float64]{X: 5.0, Y: 5.0}
	vm := geometry.VectorMathByType[float64]()
	vm.Wrap(&v, size)
	fmt.Println(v)
	// Output: (0.5,4.5)
}

// Example_intVectorMath_Length demonstrates how to use the Length method with int type.
func Example_intVectorMath_Length() {
	v := spatial.Vec[int]{X: 3, Y: 4}
	vm := geometry.VectorMathByType[int]()
	length := vm.Length(v)
	fmt.Println(length)
	// Output: 5
}

// Example_intVectorMath_Clamp demonstrates how to use the Clamp method with int type.
func Example_intVectorMath_Clamp() {
	v := spatial.Vec[int]{X: 6, Y: 5}
	size := spatial.Vec[int]{X: 5, Y: 5}
	vm := geometry.VectorMathByType[int]()
	vm.Clamp(&v, size)
	fmt.Println(v)
	// Output: (4,4)
}

// Example_intVectorMath_Wrap demonstrates how to use the Wrap method with int type.
func Example_intVectorMath_Wrap() {
	v := spatial.Vec[int]{X: 7, Y: 9}
	size := spatial.Vec[int]{X: 5, Y: 5}
	vm := geometry.VectorMathByType[int]()
	vm.Wrap(&v, size)
	fmt.Println(v)
	// Output: (2,4)
}

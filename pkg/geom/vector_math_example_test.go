package geom_test

import (
	"fmt"

	"github.com/kjkrol/gokg/pkg/geom"
)

// Example_float64VectorMath_Length demonstrates how to use the Length method with float64 type.
func Example_float64VectorMath_Length() {
	v := geom.Vec[float64]{X: 3.0, Y: 4.0}
	vm := geom.FloatVectorMath[float64]{}
	exampleVectorMathLength(v, vm)
	// Output: 5
}

// Example_float64VectorMath_Clamp demonstrates how to use the Clamp method with float64 type.
func Example_float64VectorMath_Clamp() {
	v := geom.Vec[float64]{X: 5.5, Y: 7.5}
	size := geom.Vec[float64]{X: 5.0, Y: 5.0}
	vm := geom.FloatVectorMath[float64]{}
	exampleVectorMathClamp(v, size, vm)
	// Output: (5,5)
}

// Example_float64VectorMath_Wrap demonstrates how to use the Wrap method with float64 type.
func Example_float64VectorMath_Wrap() {
	v := geom.Vec[float64]{X: 5.5, Y: 9.5}
	size := geom.Vec[float64]{X: 5.0, Y: 5.0}
	vm := geom.FloatVectorMath[float64]{}
	exampleVectorMathWrap(v, size, vm)
	// Output: (0.5,4.5)
}

// Example_intVectorMath_Length demonstrates how to use the Length method with int type.
func Example_intVectorMath_Length() {
	v := geom.Vec[int]{X: 3, Y: 4}
	vm := geom.SignedIntVectorMath[int]{}
	exampleVectorMathLength(v, vm)
	// Output: 5
}

// Example_intVectorMath_Clamp demonstrates how to use the Clamp method with int type.
func Example_intVectorMath_Clamp() {
	v := geom.Vec[int]{X: 6, Y: 5}
	size := geom.Vec[int]{X: 5, Y: 5}
	vm := geom.SignedIntVectorMath[int]{}
	exampleVectorMathClamp(v, size, vm)
	// Output: (5,5)
}

// Example_intVectorMath_Wrap demonstrates how to use the Wrap method with int type.
func Example_intVectorMath_Wrap() {
	v := geom.Vec[int]{X: 7, Y: 9}
	size := geom.Vec[int]{X: 5, Y: 5}
	vm := geom.SignedIntVectorMath[int]{}
	exampleVectorMathWrap(v, size, vm)
	// Output: (2,4)
}

// Example_intVectorMath_Wrap demonstrates how to use the Wrap method with int type.
func Example_intVectorMath_Wrap_minus() {
	target := geom.NewVec(2, 2)
	v := geom.Vec[int]{X: -101, Y: -101}
	size := geom.Vec[int]{X: 10, Y: 10}
	vm := geom.SignedIntVectorMath[int]{}

	target.AddMutable(v)
	exampleVectorMathWrap(target, size, vm)
	// Output: (1,1)
}

// Example_uint32VectorMath_Length demonstrates how to use the Length method with uint32 type.
func Example_uint32VectorMath_Length() {
	v := geom.Vec[uint32]{X: 3, Y: 4}
	vm := geom.UnsignedIntVectorMath[uint32]{}
	exampleVectorMathLength(v, vm)
	// Output: 5
}

// Example_uint32VectorMath_Clamp demonstrates how to use the Clamp method with uint32 type.
func Example_uint32VectorMath_Clamp() {
	v := geom.Vec[uint32]{X: 6, Y: 5}
	size := geom.Vec[uint32]{X: 5, Y: 5}
	vm := geom.UnsignedIntVectorMath[uint32]{}
	exampleVectorMathClamp(v, size, vm)
	// Output: (5,5)
}

// Example_uint32VectorMath_Wrap demonstrates how to use the Wrap method with uint32 type.
func Example_uint32VectorMath_Wrap() {
	v := geom.Vec[uint32]{X: 7, Y: 9}
	size := geom.Vec[uint32]{X: 5, Y: 5}
	vm := geom.UnsignedIntVectorMath[uint32]{}
	exampleVectorMathWrap(v, size, vm)
	// Output: (2,4)
}

func exampleVectorMathLength[T geom.Numeric](v geom.Vec[T], vm geom.VectorMath[T]) {
	fmt.Println(vm.Length(v))
}

func exampleVectorMathClamp[T geom.Numeric](v, size geom.Vec[T], vm geom.VectorMath[T]) {
	vm.Clamp(&v, size)
	fmt.Println(v)
}

func exampleVectorMathWrap[T geom.Numeric](v, size geom.Vec[T], vm geom.VectorMath[T]) {
	vm.Wrap(&v, size)
	fmt.Println(v)
}

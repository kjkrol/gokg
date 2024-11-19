package geometry_test

import (
	"fmt"
	"testing"

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

func TestVecSub(t *testing.T) {
	for _, test := range []struct {
		arg1     geometry.Vec[int]
		arg2     geometry.Vec[int]
		expected geometry.Vec[int]
	}{
		{geometry.Vec[int]{2, 3}, geometry.Vec[int]{1, 2}, geometry.Vec[int]{1, 1}},
		{geometry.Vec[int]{1, 2}, geometry.Vec[int]{1, 2}, geometry.Vec[int]{0, 0}},
		{geometry.Vec[int]{0, 0}, geometry.Vec[int]{4, 4}, geometry.Vec[int]{-4, -4}},
		{geometry.Vec[int]{4, 0}, geometry.Vec[int]{1, 0}, geometry.Vec[int]{3, 0}},
		{geometry.Vec[int]{1, 0}, geometry.Vec[int]{4, 0}, geometry.Vec[int]{-3, 0}},
	} {
		if output := test.arg1.Sub(test.arg2); !output.Equals(test.expected) {
			t.Errorf("result %v not equal to expected %v", test.arg1, test.expected)
		}
	}
}

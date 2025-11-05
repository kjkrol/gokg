package geometry

import (
	"fmt"
)

// ExampleVec_Add demonstrates how to use the Add method.
func ExampleVec_Add() {
	v1 := NewVec(1.0, 2.0)
	v2 := NewVec(3.0, 4.0)
	result := v1.Add(v2)
	fmt.Println(result)
	// Output: (4,6)
}

// ExampleVec_Sub demonstrates how to use the Sub method.
func ExampleVec_Sub() {
	v1 := NewVec(5.0, 7.0)
	v2 := NewVec(2.0, 3.0)
	result := v1.Sub(v2)
	fmt.Println(result)
	// Output: (3,4)
}

// ExampleVec_AddMutable demonstrates how to use the AddMutable method.
func ExampleVec_AddMutable() {
	v1 := NewVec(1.0, 2.0)
	v2 := NewVec(3.0, 4.0)
	v1.AddMutable(v2)
	fmt.Println(v1)
	// Output: (4,6)
}

// ExampleVec_SubMutable demonstrates how to use the SubMutable method.
func ExampleVec_SubMutable() {
	v1 := NewVec(5.0, 7.0)
	v2 := NewVec(2.0, 3.0)
	v1.SubMutable(v2)
	fmt.Println(v1)
	// Output: (3,4)
}

// ExampleVec_Equals demonstrates how to use the Equals method.
func ExampleVec_Equals() {
	v1 := NewVec(1.0, 2.0)
	v2 := NewVec(1.0, 2.0)
	v3 := NewVec(3.0, 4.0)
	fmt.Println(v1.Equals(v2))
	fmt.Println(v1.Equals(v3))
	// Output:
	// true
	// false
}

// ExampleVec_String demonstrates how to use the String method.
func ExampleVec_String() {
	v := NewVec(1.0, 2.0)
	fmt.Println(v.String())
	// Output: (1,2)
}

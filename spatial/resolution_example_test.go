package spatial

import (
	"fmt"
)

// --- EXAMPLES ---

func ExampleResolution_Side() {
	fmt.Println(Size1x1.Side())
	fmt.Println(Size8x8.Side())
	fmt.Println(Size1024x1024.Side())
	// Output:
	// 1
	// 8
	// 1024
}

func ExampleResolution_MaxCoord() {
	fmt.Println(Size1x1.MaxCoord())
	fmt.Println(Size8x8.MaxCoord())
	fmt.Println(Size1024x1024.MaxCoord())
	// Output:
	// 0
	// 7
	// 1023
}

func ExampleResolution_Cells() {
	fmt.Println(Size1x1.Cells())
	fmt.Println(Size8x8.Cells())
	fmt.Println(Size1024x1024.Cells())
	// Output:
	// 1
	// 64
	// 1048576
}

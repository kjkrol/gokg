package geom_test

import (
	"fmt"

	"github.com/kjkrol/gokg/geom"
)

func ExampleLength() {
	fmt.Println(geom.Length(geom.Vec{X: 3, Y: 4}))
	// Output: 5
}

// Clamp keeps a vector inside the box, in both directions.
func ExampleClamp() {
	size := geom.Vec{X: 5, Y: 5}
	fmt.Println(geom.Clamp(geom.Vec{X: 5.5, Y: 7.5}, size))
	fmt.Println(geom.Clamp(geom.Vec{X: -2, Y: 3}, size))
	// Output:
	// (5,5)
	// (0,3)
}

// Wrap folds a vector back inside, which is what a toroidal world does to
// anything crossing an edge.
func ExampleWrap() {
	size := geom.Vec{X: 5, Y: 5}
	fmt.Println(geom.Wrap(geom.Vec{X: 5.5, Y: 9.5}, size))
	fmt.Println(geom.Wrap(geom.Vec{X: -0.5, Y: 2}, size))
	// Output:
	// (0.5,4.5)
	// (4.5,2)
}

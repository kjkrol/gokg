package spatial

import "fmt"

func ExampleNewMortonCode() {
	for _, mCode := range []MortonCode{
		NewMortonCode(4, 1),
		NewMortonCode(5, 5),
		NewMortonCode(0b1001, 0b10101),
	} {
		fmt.Printf("%b\n", mCode)
	}
	// Output:
	// 10010
	// 110011
	// 1001100011
}

func ExampleMortonCode_Decode() {
	for _, mCode := range []MortonCode{
		NewMortonCode(4, 1),
		NewMortonCode(5, 5),
		NewMortonCode(0b1001, 0b10101),
	} {
		x, y := mCode.Decode()
		fmt.Printf("(%d,%d)\n", x, y)
	}
	// Output:
	//(4,1)
	//(5,5)
	//(9,21)
}

package intset

import "fmt"

func ExampleIntSet_Elems() {
	x := IntSet{words: []uint{123, 234, 345, 345}} // "{0 1 3 4 5 6 65 67 69 70 71 128 131 132 134 136 192 195 196 198 200}"
	fmt.Println(&x)
	fmt.Println(x.Elems())
	fmt.Printf("%T\n", x.Elems())
	// output:
	// {0 1 3 4 5 6 65 67 69 70 71 128 131 132 134 136 192 195 196 198 200}
	// [0 1 3 4 5 6 65 67 69 70 71 128 131 132 134 136 192 195 196 198 200]
	// []int
}

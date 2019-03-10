package intset

import "fmt"

func ExampleIntSet_UnionWith() {
	var x, y IntSet
	x.Add(3)
	x.Add(30)
	y.AddAll(3, 300, 333)
	x.UnionWith(&y)
	fmt.Println(x)
	// outbut:
	// {3 30 300 333}
}

func ExampleIntSet_IntersectionWith() {
	var x, y IntSet
	x.Add(1)
	x.Add(456)
	x.Add(10989)
	x.Add(29374)
	y.AddAll(1, 2, 1021, 10989, 284892, 66)
	x.IntersectionWith(&y)
	fmt.Println(&x)
	// output: {1 10989}
}

func ExampleIntSet_DifferenceWith() {
	var x, y IntSet
	x.AddAll(1, 2, 3, 41, 56, 125, 233, 567, 864, 980, 999)
	y.AddAll(5, 7, 8, 9, 864, 980, 999)
	x.DifferenceWith(&y)
	fmt.Println(&x)
	// output: {1 2 3 41 56 125 233 567}
}

func ExampleIntSet_SymmetricDifferenceWith() {
	var x, y IntSet
	x.AddAll(1, 2, 3, 41, 56, 125, 233, 567, 864, 980, 999)
	y.AddAll(5, 7, 8, 9, 864, 980, 999)
	x.SymmetricDifferenceWith(&y)
	fmt.Println(&x)
	// output: {1 2 3 5 7 8 9 41 56 125 233 567}
}
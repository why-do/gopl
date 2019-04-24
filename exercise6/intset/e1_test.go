package intset

import "fmt"

func ExampleIntSet_Len() {
	var x IntSet
	x.AddAll(100, 200, 300, 888, 1981, 18901821, 1912)
	fmt.Println(x.Len())
	// output: 7
}

func ExampleIntSet_Len2() {
	var x IntSet
	x.AddAll(100, 200, 300, 888, 1981, 18901821, 1912)
	fmt.Println(x.Len2())
	// output: 7
}

func ExampleIntSet_Len3() {
	var x IntSet
	x.AddAll(100, 200, 300, 888, 1981, 18901821, 1912)
	fmt.Println(x.Len3())
	// output: 7
}

func ExampleIntSet_Remove() {
	var x IntSet
	x.Remove(10)
	fmt.Println(&x)
	x.AddAll(1, 2, 3, 4, 5, 6)
	fmt.Println(&x)
	x.Remove(3)
	fmt.Println(&x)
	x.Add(55555)
	x.Add(666666)
	x.Add(7777777)
	fmt.Println(&x)
	x.Remove(55555)
	x.Remove(666666)
	fmt.Println(&x)
	x.Remove(7777777)
	fmt.Println(&x)
	// output:
	// {}
	// {1 2 3 4 5 6}
	// {1 2 4 5 6}
	// {1 2 4 5 6 55555 666666 7777777}
	// {1 2 4 5 6 7777777}
	// {1 2 4 5 6}
}

func ExampleIntSet_Clear() {
	var x IntSet
	x.Add(11)
	fmt.Println(&x)
	x.Clear()
	fmt.Println(&x)
	x.Add(22)
	fmt.Println(&x)
	// output:
	// {11}
	// {}
	// {22}
}

func ExampleIntSet_Copy() {
	var x IntSet
	x.AddAll(333, 4444, 88888888, 999999999)
	y := *x.Copy()
	fmt.Println(&x, &y)
	x.Remove(333)
	y.Remove(999999999)
	fmt.Println(&x, &y)
	// output:
	// {333 4444 88888888 999999999} {333 4444 88888888 999999999}
	// {4444 88888888 999999999} {333 4444 88888888}
}
package intset

import "fmt"

func Example_one() {
	var s1, s2 IntSet
	s1.Add(1)
	s1.Add(144)
	s1.Add(9)
	fmt.Println(s1.String())

	s2.Add(9)
	s2.Add(42)
	fmt.Println(&s2)

	s1.UnionWith(&s2)
	fmt.Println(&s1)

	fmt.Println(s1.Has(9), s1.Has(123))

	// Output:
	// {1 9 144}
	// {9 42}
	// {1 9 42 144}
	// true false
}

func Example_two() {
	var x IntSet
	x.Add(1)
	x.Add(144)
	x.Add(9)
	x.Add(42)

	fmt.Println(&x)
	fmt.Println(x.String())
	fmt.Println(x)
	// Output:
	// {1 9 42 144}
	// {1 9 42 144}
	// {[4398046511618 0 65536]}
}
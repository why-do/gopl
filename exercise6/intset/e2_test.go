package intset

import "fmt"

func ExampleIntSet_AddAll() {
	s1 := IntSet{}
	s1.Add(10)
	fmt.Println(&s1)
	s1.AddAll(100, 1000, 10000, 100, 1000)
	fmt.Println(&s1)
	// output:
	// {10}
	// {10 100 1000 10000}
}

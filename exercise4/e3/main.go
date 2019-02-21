package main

import "fmt"

func reverse(s *[5]int) {
	for i, j := 0, len(*s)-1; i < j; i, j = i+1, j-1 {
		(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
	}
}

func main() {
	l := [...]int{1, 2, 3, 4, 5}
	fmt.Println(l)
	reverse(&l) // 传入数组指针作为参
	fmt.Println(l)
}
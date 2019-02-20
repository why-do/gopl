package main

import "fmt"

func reverse(s *[]int) {
	for i, j := 0, len(*s)-1; i < j; i, j = i+1, j-1 {
		(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
	}
}

func main() {
	l := [...]int{1, 2, 3, 4, 5} // 这个是数组
	fmt.Println(l)
	s := l[:]
	reverse(&s) // 传入切片
	fmt.Println(l)
}
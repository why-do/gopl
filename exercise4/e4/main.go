package main

import "fmt"

func rotate(n int, s []int) {
	l := len(s)
	n = n % l
	if n == 0 {
		return
	}

	tmp := make([]int, l)
	copy(tmp, s)
	for i := 0; i < l; i++ {
		s[(l+i-n)%l] = tmp[i]
	}
}

func main() {
	l := [6]int{1, 2, 3, 4, 5, 6}
	n := 3
	fmt.Println(l, n)
	rotate(n, l[:])
	fmt.Println(l, n)
}

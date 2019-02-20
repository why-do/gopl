package main

import "fmt"

func reverse(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func moveLeft(n int, s []int) {
	reverse(s[:n])
	reverse(s[n:])
	reverse(s)
}

func moveRight(n int, s []int) {
	reverse(s[n:])
	reverse(s[:n])
	reverse(s)
}

func main() {
	l := [...]int{1, 2, 3, 4, 5} // 这个是数组
	fmt.Println(l)
	reverse(l[:]) // 传入切片
	fmt.Println(l)
	moveLeft(2, l[:])
	fmt.Println(l)
	moveRight(1, l[:])
	fmt.Println(l)
}

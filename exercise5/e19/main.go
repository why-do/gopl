package main

import "fmt"

func main() {
	s := noRet()
	fmt.Println(s)
}

func noRet() (s string) {
	defer func() {
		p := recover()
		s = fmt.Sprint(p)
	}()
	panic("Hello")
}


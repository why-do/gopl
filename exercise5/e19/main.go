package main

import "fmt"

func main() {
	defer func() {
		p := recover()
		fmt.Println(p)
	}()
	noRet()
}

func noRet() {
	panic("Hello")
}


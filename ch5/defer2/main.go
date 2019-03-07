package main

import (
	"fmt"
	"os"
	"runtime"
)

func f(x int) {
	fmt.Printf("f(%d)\n", x+0/x)
	defer fmt.Printf("defer %d\n", x)
	f(x - 1)
}

func printStack() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	os.Stdout.WriteString("Stack 中的内容:\n")
	os.Stdout.Write(buf[:n])
	os.Stdout.WriteString("Stack 结束...\n")
}

func main() {
	defer printStack()
	f(3)
}

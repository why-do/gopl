package main

import (
	"flag"
	"fmt"
	"runtime"
)

func main() {
	var n int
	flag.IntVar(&n, "n", 1, "GOMAXPROCS")
	flag.Parse()
	runtime.GOMAXPROCS(n)
	for {
		go fmt.Print(0)
		fmt.Print(1)
	}
}

package main

import (
	"flag"
	"fmt"
	"gopl/ch7/tempconv"
)

var temp = tempconv.CelsiusFlag("temp", 20.0, "温度")

func main() {
	flag.Parse()
	fmt.Println(*temp)
}

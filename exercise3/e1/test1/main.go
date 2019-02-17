package main

import (
	"fmt"
	"math"
)

func main() {
	f(1, 0)
	f(0, 1)
	x := f(1e300, 0)
	fmt.Println(math.IsInf(x, 0))
	fmt.Println(math.IsInf(x, 1))
	fmt.Println(math.IsInf(x, -1))
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y) // 到(0,0)的距离
	fmt.Println(r, math.Sin(r))
	return math.Sin(r)
}

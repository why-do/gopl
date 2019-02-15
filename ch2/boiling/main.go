// 输出水的沸点
package main

import "fmt"

const boilingF = 212.0

func main() {
	var f = boilingF
	var c = (f -32) * 5 /9
	fmt.Printf("沸点 = %g℉  or %g℃\n", f, c)
}
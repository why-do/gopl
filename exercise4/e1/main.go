package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	c1 := sha256.Sum256([]byte("x"))
	fmt.Println(PopCount(c1))
	fmt.Printf("%08b\n", c1)
}

// pc[i] 是 i 的 population count
var pc [256]byte

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

func PopCount(c [32]byte) (n int) {
	for _, v := range c {
		n += int(pc[v])
	}
	return n
}

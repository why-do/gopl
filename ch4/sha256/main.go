package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	c1 := sha256.Sum256([]byte("x"))
	c2 := sha256.Sum256([]byte("X"))
	fmt.Printf("%T\n", c1)
	fmt.Printf("% x\n", c1)
	fmt.Printf("% x\n", c2)
	fmt.Println(c1==c2)
}

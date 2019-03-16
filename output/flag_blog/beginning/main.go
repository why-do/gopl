package main

import (
	"flag"
	"fmt"
)

var name string

func init() {
	flag.StringVar(&name, "name", "Adam", "名字")
}

var ageP = flag.Int("age", 18, "年龄")

func main() {
	flag.Parse()
	fmt.Printf("%T %[1]v\n", name)
	fmt.Printf("%T %[1]v\n", ageP)
	fmt.Printf("%T %[1]v\n", *ageP)
}

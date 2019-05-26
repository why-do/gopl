package main

// #include <float.h>
import "C"
import "fmt"

func main() {
	fmt.Println("Max float value of float is", C.FLT_MAX)
}

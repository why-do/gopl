package main

import (
	"fmt"
	"os"
)

// 把命令行参数，依次打印，每行一个
func main() {
	for _, s := range os.Args[1:] {
		fmt.Println(s)
	}
}

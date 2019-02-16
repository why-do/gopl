package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"gopl/exercise2/e3/popcount"
)

func main() {
	input := bufio.NewScanner(os.Stdin)
	fmt.Println("输入 q 退出")
	for input.Scan() {
		s := input.Text()
		if s == "q" {
			os.Exit(1)
		}
		x, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "需要输入数字: %v\n", err)
			continue
		}
		fmt.Println(popcount.PopCount(x))
	}
}

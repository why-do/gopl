package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counts := make(map[string]int)
	// 从标准输入接收数据，程序执行后可以在命令行继续输入，回车输入换行
	// 要想输入完成，windows系统下，在输入回车换行后的空行位置，按 ctrl+z，再回车确认
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		counts[input.Text()]++
	}
	// 注意：上面忽略了 input.Err() 中可能的错误
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var counts map[string]int  // Unicode 字符数量

func main() {
	files := os.Args[1:]
	if len(files) == 0 {
		fmt.Println("未提供文件路径")
		os.Exit(1)
	}
	counts = make(map[string]int)
	for _, arg := range files {
		f, err := os.Open(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Open file error: %v\n", err)
			continue
		}
		wordfreq(f)
		f.Close()
	}
	// 报告
	for k, v := range counts {
		if v > 1{
			fmt.Printf("%d %q\n", v, k)
		}
	}
}

func wordfreq(f io.Reader) {
	input := bufio.NewScanner(f)
	input.Split(bufio.ScanWords)  // 按单词分割，以空白符进行划分
	for input.Scan() {
		counts[input.Text()]++
	}
}

package main

import (
	"fmt"
	"log"
	"os"

	"gopl/ch5/links"
)

// 对每一个worklist中的元素调用f
// 并将返回的内容添加到worklist中，对每一个元素，最多调用一次f
func breadthFirst(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item)...)
			}
		}
	}
}

func crwal(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

func main() {
	// 开始广度遍历
	// 从命令行参数开始
	breadthFirst(crwal, os.Args[1:])
}

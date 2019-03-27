package main

import (
	"fmt"
	"log"
	"os"

	"gopl/ch5/links"
)

// 令牌 tokens 是一个计数信号量
// 确保并发请求限制在 20 个以内
var tokens = make(chan struct{}, 20)

func crawl(url string) []string {
	fmt.Println(url)
	tokens <- struct{}{} // 获取令牌
	list, err := links.Extract(url)
	<- tokens // 释放令牌
	if err != nil {
		log.Print(err)
	}
	return list
}

func main() {
	worklist := make(chan []string)

	// 从命令行参数开始
	go func() { worklist <- os.Args[1:] }()

	// 并发爬取 Web
	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				go func(link string) {
					worklist <- crawl(link)
				}(link)
			}
		}
	}
}

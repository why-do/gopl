package main

import (
	"fmt"
	"log"
	"os"

	"gopl/ch5/links"
)

func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

func main() {
	worklist := make(chan []string) // 可能有重复的URL列表
	unseenLinks := make(chan string) // 去重后的eURL列表

	// 向任务列表中添加命令行参数
	go func() { worklist <- os.Args[1:] }()

	// 创建20个爬虫 goroutine 来获取每个不可见链接
	for i := 0; i < 20; i++ {
		go func() {
			for link := range unseenLinks {
				foundLinks := crawl(link)
				go func() { worklist <- foundLinks }()
			}
		}()
	}
	
	// 主 goroutine 对 URL 列表进行去重
	// 并把没有爬取过的条目发送给爬虫程序
	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				unseenLinks <- link
			}
		}
	}
}

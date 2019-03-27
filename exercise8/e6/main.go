package main

import (
	"flag"
	"fmt"
	"log"

	"gopl/ch5/links"
)

// 令牌 tokens 是一个计数信号量
// 确保并发请求限制在 20 个以内
var tokens = make(chan struct{}, 20)

func crawl(url string, depth int) urllist {
	fmt.Println(depth, url)
	tokens <- struct{}{} // 获取令牌
	list, err := links.Extract(url)
	<-tokens // 释放令牌
	if err != nil {
		log.Print(err)
	}
	return urllist{list, depth+1}
}

var depth int
var starturl string

func init() {
	flag.IntVar(&depth, "depth", -1, "深度限制") // 小于0就是不限制递归深度，0就是只爬取当前页面
	flag.StringVar(&starturl, "url", "http://lab.scrapyd.cn/", "起始URL")
}

type urllist struct {
	urls  []string
	depth int
}

func main() {
	flag.Parse()
	worklist := make(chan urllist)
	var n int // 等待发送到任务列表的数量

	// 从命令行参数开始
	n++
	go func() { worklist <- urllist{[]string{starturl}, 0} }()

	// 并发爬取 Web
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		list := <-worklist
		if depth >=0 && list.depth > depth {
			fmt.Println("continue") // TODO: 这里还是有问题的
			continue
		}
		for _, link := range list.urls {
			if !seen[link] {
				seen[link] = true
				n++
				go func(url string, depth int) {
					worklist <- crawl(url, depth)
				}(link, list.depth)
			}
		}
	}
}

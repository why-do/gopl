package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

	"gopl/ch5/links"
)

var count = make(chan int) // 统计一共爬取了多个页面

// 令牌 tokens 是一个计数信号量
// 确保并发请求限制在 20 个以内
var tokens = make(chan struct{}, 20)

func crawl(url string, depth int) urllist {
	fmt.Println(depth, <-count, url)
	tokens <- struct{}{} // 获取令牌
	list, err := links.Extract(url)
	<-tokens // 释放令牌
	if err != nil {
		log.Print(err)
	}
	return urllist{list, depth + 1}
}

var depth int

func init() {
	flag.IntVar(&depth, "depth", -1, "深度限制") // 小于0就是不限制递归深度，0就是只爬取当前页面
}

type urllist struct {
	urls  []string
	depth int
}

func main() {
	// 负责 count 值自增的 goroutine
	go func() {
		var i int
		for {
			i++
			count <- i
		}
	}()

	flag.Parse()
	worklist := make(chan urllist)
	// 等待发送到任务列表的数量
	// 因为需要在 goroutine 里修改，需要换成并发安全的计数器
	var n sync.WaitGroup
	starturls := flag.Args()
	if len(flag.Args()) == 0 {
		starturls = []string{"http://lab.scrapyd.cn/"}
	}

	// 从命令行参数开始
	n.Add(1)
	go func() { worklist <- urllist{starturls, 0} }()
	// 等待全部worklist处理完，就关闭worklist
	go func() {
		n.Wait()
		close(worklist)
	}()

	// 并发爬取 Web
	seen := make(map[string]bool)
	for list := range worklist {
		// 处理完一个worklist后才能让 n 计数器减1
		// 而处理 worklist 又是很多个 goroutine，所以需要再用一个计数器
		var n2 sync.WaitGroup
		for _, link := range list.urls {
			if !seen[link] {
				seen[link] = true
				n2.Add(1)
				go func(url string, listDepth int) {
					nextList := crawl(url, listDepth)
					// 如果 depth>0 说明有深度限制
					// 如果当前的深度已经达到（或超过）深度限制，则爬取完这个连接后，不需要再继续爬取，直接返回
					if depth >= 0 && listDepth >= depth {
						// 超出递归深度的页面，在爬取完之后，也输出 URL
						// for _, nextUrl := range nextList.urls {
						// 	fmt.Println(nextList.depth, "stop", nextUrl)
						// }
						n2.Done() // 所有退出的情况都要减计数器n2的值，但是一定要在向通道发送之前
						return
					}
					n.Add(1)             // 添加任务前，计数加1
					n2.Done()            // 先确保计数器n加1了，再减计数器n2的值
					worklist <- nextList // 新的任务加入管道必须在最后，之后再一次for循环迭代的时候，才会接收这个值
				}(link, list.depth)
			}
		}
		n2.Wait()
		n.Done()
		// 把计数器的操作也放到 goroutine 中，这样可以继续下一次 for 循环的迭代
		// go func() {
		// 	n2.Wait()
		// 	n.Done()
		// }()
	}
}

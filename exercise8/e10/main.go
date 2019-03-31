package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"gopl/exercise8/e10/links"
)

// 检测取消状态的工具函数
var done = make(chan struct{})

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

var count = make(chan int) // 统计一共爬取了多个页面

// 令牌 tokens 是一个计数信号量
// 确保并发请求限制在 20 个以内
var tokens = make(chan struct{}, 20)

func crawl(url string, depth int) urllist {
	select {
	case <-done:
		return urllist{nil, depth + 1}
	case tokens <- struct{}{}: // 获取令牌
		fmt.Println(depth, <-count, url)
	}
	list, err := links.Extract(url, done)
	<-tokens // 释放令牌
	if err != nil && !strings.Contains(err.Error(), "net/http: request canceled") {
		log.Print(err)
	}
	return urllist{list, depth + 1}
}

var depth int
var showOverDepthURL bool

func init() {
	flag.IntVar(&depth, "depth", -1, "深度限制") // 小于0就是不限制递归深度，0就是只爬取当前页面
	flag.BoolVar(&showOverDepthURL, "showover", false, "超过深度限制的URL是否打印")
}

type urllist struct {
	urls  []string
	depth int
}

func main() {
	// 当检测到输入时，广播取消
	go func() {
		os.Stdin.Read(make([]byte, 1)) // 读一个字节
		close(done)
	}()

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
loop:
	for {
		var list urllist
		var worklistok bool
		select {
		case <-done:
			// 耗尽 worklist，让已经创建的 goroutine 结束
			for range worklist {
				n.Done()
			}
			// 执行到这里的前提是迭代完 worklist，就是需要 worklist 关闭
			// 关闭 worklist 则需要 n 计数器归0。而 worklist 每一次减1，需要一个 n2 计数器归零
			// 所以，下面的 return 应该不会在其他 goroutine 运行完毕之前执行
			return
		case list, worklistok = <-worklist:
			if !worklistok {
				break loop
			}
		}

		// 处理完一个worklist后才能让 n 计数器减1
		// 而处理 worklist 又是很多个 goroutine，所以需要再用一个计数器
		var n2 sync.WaitGroup
		for _, link := range list.urls {
			if cancelled() {
				break
			}
			if !seen[link] {
				seen[link] = true
				n2.Add(1)
				go func(url string, listDepth int) {
					nextList := crawl(url, listDepth)
					// 如果 depth>0 说明有深度限制
					// 如果当前的深度已经达到（或超过）深度限制，则爬取完这个连接后，不需要再继续爬取，直接返回
					if depth >= 0 && listDepth >= depth {
						// 超出递归深度的页面，在爬取完之后，也输出 URL
						if showOverDepthURL {
							for _, nextUrl := range nextList.urls {
								fmt.Println(nextList.depth, "stop", nextUrl)
							}
						}
						nextList = urllist{nil, listDepth + 1} // 丢弃返回的数据
					}
					// 超出深度限制，或本身就没有新的URL返回，也包括取消操作之后的返回，
					if len(nextList.urls) == 0 {
						n2.Done()
						return
					}
					n.Add(1)             // 添加任务前，计数加1
					n2.Done()            // 先确保计数器n加1了，再减计数器n2的值
					worklist <- nextList // 新的任务加入管道必须在最后，之后再一次for循环迭代的时候，才会接收这个值
				}(link, list.depth)
			}
		}
		// n2.Wait()
		// n.Done()
		// 把计数器的操作也放到 goroutine 中，这样可以继续下一次 for 循环的迭代
		go func() {
			n2.Wait()
			n.Done()
		}()
	}
}

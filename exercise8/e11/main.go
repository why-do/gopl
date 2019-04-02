package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

// 递归解析文档树获取 title
func forEachNode(n *html.Node, titleP *string, pre, post func(n *html.Node)) {
	if *titleP != "" {
		return
	}
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, titleP, pre, post)
	}
	if post != nil {
		post(n)
	}
}

// 使用上面的 forEachNode 函数，递归文档树。返回找到的第一个 title 或者全部遍历返回空字符串
func soleTitle(doc *html.Node) string {
	var title string // 被下面的闭包引用了
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			title = n.FirstChild.Data
		}
	}
	forEachNode(doc, &title, visitNode, nil)
	return title
}

// 解析返回 title 的入口函数
// 把响应体解析为文档树，然后交给 soleTitle 处理，获取 title
func title(url string, body io.Reader) (string, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return "", fmt.Errorf("parseing %s as HTML: %v", url, err)
	}
	title := soleTitle(doc)
	if title == "" {
		return "", errors.New("no title element")
	}
	return title, nil
}

// 上面是解析返回结果的逻辑，不是这里的重点

// 请求镜像资源
func mirroredQuery(urls ...string) string {
	type respData struct { // 返回的数据类型
		resp string
		err  error
	}
	count := len(urls)                      // 总共发起的请求数
	responses := make(chan respData, count) // 有几个镜像，就要多大的容量，不能少
	done := make(chan struct{})             // 取消广播的通道
	var wg sync.WaitGroup                   // 计数器，等所有请求返回后再结束。帮助判断其他连接是否可以取消
	wg.Add(count)
	for _, url := range urls {
		go func(url string) {
			defer wg.Done()
			resp, err := request(url, done)
			responses <- respData{resp, err}
		}(url)
	}
	// 等待结果返回并处理
	var response string
	for i := 0; i < count; i++ {
		data := <-responses
		if data.err == nil { // 只接收第一个无错误的返回
			response = data.resp
			close(done)
			break
		}
		fmt.Fprintf(os.Stderr, "mirror get: %v\n", data.err)
	}
	wg.Wait()
	return response
}

// 负责发起请求并返回结果和可能的错误
func request(url string, done <-chan struct{}) (response string, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Cancel = done
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 检查返回的页面是HTML通过判断Content-Type，比如：Content-Type: text/html; charset=utf-8
	ct := resp.Header.Get("Content-Type")
	if ct != "text/html" && !strings.HasPrefix(ct, "text/html;") {
		return "", fmt.Errorf("%s has type %s, not text/html", url, ct)
	}
	// 如果上面检查响应头没问题，把响应体交给 title 函数解析获取结果
	return title(url, resp.Body)
}

func main() {
	response := mirroredQuery(os.Args[1:]...)
	fmt.Println(response)
}

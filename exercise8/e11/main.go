package main

import (
	"fmt"
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

func soleTitle(doc *html.Node) string {
	var title string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			title = n.FirstChild.Data
		}
	}
	forEachNode(doc, &title, visitNode, nil)
	return title
}

// 请求镜像资源
func mirroredQuery(urls ...string) string {
	responses := make(chan string, len(urls)) // 有几个镜像，就要多大的容量，不能少
	done := make(chan struct{})               // 取消广播的通道
	var wg sync.WaitGroup
	wg.Add(len(urls))
	for _, url := range urls {
		go func(url string) {
			defer wg.Done()
			resp, err := request(url, done)
			if err != nil {
				fmt.Fprintf(os.Stderr, "mirror get: %v\n", err)
				return
			}
			responses <- resp
		}(url)
	}
	resp := <-responses // 返回最快一个获取到的请求结果
	close(done)
	wg.Wait()
	return resp
}

// 这个函数签名不对，还要再修改
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

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", fmt.Errorf("parseing %s as HTML: %v", url, err)
	}
	title := soleTitle(doc)
	if title == "" {
		return "", fmt.Errorf("no title element")
	}
	return title, nil
}

func main() {
	title := mirroredQuery(os.Args[1:]...)
	fmt.Println(title)
}

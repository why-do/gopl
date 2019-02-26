package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// 输出从 URL 获取的内容
func fetch(url string) io.ReadCloser {
	if !strings.HasPrefix(url, "http") {
		url = fmt.Sprint("http://", url)
	}
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch get error: %v\n", err)
		os.Exit(1)
	}
	return resp.Body
}

// 输出 HTML 文档中的所有连接
func findlinks(rc io.ReadCloser) {
	doc, err := html.Parse(rc)
	rc.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks1: %v\n", err)
		os.Exit(1)
	}
	for _, link := range visit(nil, doc) {
		fmt.Println(link)
	}
}

// 将节点 n 中的每个链接添加到结果中
func visit(links []string, n *html.Node) []string {
	if n == nil {
		return links
	}
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	// 可怕的递归，非常不好理解。
	return visit(visit(links, n.FirstChild), n.NextSibling)
}

func main() {
	for _, arg := range os.Args[1:] {
		body := fetch(arg)
		findlinks(body)
	}
}

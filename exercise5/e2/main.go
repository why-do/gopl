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

// 统计元素个数
func parseElements(rc io.ReadCloser) map[string]int {
	doc, err := html.Parse(rc)
	rc.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks1: %v\n", err)
		os.Exit(1)
	}
	elements := make(map[string]int)
	return countElements(elements, doc)
}

// 将节点 n 中的每个链接添加到结果中
func countElements(elements map[string]int, n *html.Node) map[string]int {
	if n == nil {
		return elements
	}
	if n.Type == html.ElementNode {
		elements[n.Data]++
	}
	// 可怕的递归，非常不好理解。
	return countElements(countElements(elements, n.FirstChild), n.NextSibling)
}

func main() {
	elem := make(map[string]int)
	for _, arg := range os.Args[1:] {
		body := fetch(arg)
		e := parseElements(body)
		for k, v := range e {
			elem[k] += v
		}
	}
	fmt.Println(elem)
}

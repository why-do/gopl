package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

// soleTitle 返回文档中一个非空标题元素
// 如果没有标题则返回错误
func soleTitle(doc *html.Node) (title string, err error) {
	type bailout struct{}
	
	defer func() {
		switch p := recover(); p {
		case nil:
			// 没有宕机
		case bailout{}:
			// 预期的宕机
			err = fmt.Errorf("multiple title elements")
		default:
			panic(p) // 未预期的宕机，继续宕机过程
		}
	}()
	// 如果发现多余一个非空标题，退出递归
	forEachNode(doc, func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			if title != "" {
				panic(bailout{}) // 多个标题元素
			}
			title = n.FirstChild.Data
		}
	}, nil)
	if title == "" {
		return "", fmt.Errorf("no title element")
	}
	return title, nil
}

func title(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查返回的页面是HTML通过判断Content-Type，比如：Content-Type: text/html; charset=utf-8
	ct := resp.Header.Get("Content-Type")
	if ct != "text/html" && !strings.HasPrefix(ct, "text/html;") {
		return fmt.Errorf("%s has type %s, not text/html", url, ct)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("parseing %s as HTML: %v", url, err)
	}

	title, err := soleTitle(doc)
	if err != nil {
		return err
	}
	fmt.Println(title)
	return nil
}

func main() {
	for _, arg := range os.Args[1:] {
		if err := title(arg); err != nil {
			fmt.Fprintf(os.Stderr, "title: %v\n", err)
		}
	}
}

// 做不来
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func main() {
	urlP := flag.String("url", "http://baidu.com", "完整的url")
	idP := flag.String("id", "head", "看查找的ID")
	flag.Parse()

	doc, err := htmlDoc(*urlP)
	if err != nil {
		fmt.Fprintf(os.Stderr, "get node error: %v", err)
		os.Exit(1)
	}
	node := ElementByID(doc, *idP)
	node2 := ElementByID2(doc, *idP)
	showNode(node)
	showNode(node2)
}

// 输出 node 的标签名称和属性
func showNode(node *html.Node) {
	if node != nil {
		fmt.Printf("<%s", node.Data)
		for _, a := range node.Attr {
			fmt.Printf(` %s="%s"`, a.Key, a.Val)
		}
		fmt.Println(">")
	} else {
		fmt.Println(node)
	}
}

// Get请求然后解析，返回解析后的文档树
func htmlDoc(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return html.Parse(resp.Body)
}

// 上面的这些不是重点
// 下面开始是重点
var found bool // 控制是否继续遍历

// 这个函数的功能保持不变，依旧提供前序调用和后序调用，虽然这里不需要用到后序调用
func forEachNode(n *html.Node, id string, pre, post func(n *html.Node, id string) bool) {
	if found {
		return
	}
	if pre != nil {
		found = pre(n, id)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, id, pre, post)
	}
	if post != nil {
		found = post(n, id)
	}
}

var foundNode *html.Node // 保存找到的Node

func startElement(n *html.Node, id string) bool {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == "id" && attr.Val == id {
				foundNode = n
				return true
			}
		}
	}
	return false
}

// 使用一些全局变量的话，其实貌似也不难，只是感觉代码很不规范
func ElementByID(doc *html.Node, id string) *html.Node {
	forEachNode(doc, id, startElement, nil) // 不需要后序调用，传入参数nil
	return foundNode
}

// 后面一节的闭包可以解决上面的问题，而且代码方面并没有太多的差别，
// 只要理解闭包，稍加修改就能实现
func ElementByID2(doc *html.Node, id string) *html.Node {
	var found2 bool
	var forEachNode2 func(n *html.Node, id string, pre, post func(n *html.Node, id string) bool)
	forEachNode2 = func(n *html.Node, id string, pre, post func(n *html.Node, id string) bool) {
		if found2 {
			return
		}
		if pre != nil {
			found2 = pre(n, id)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			forEachNode2(c, id, pre, post)
		}
		if post != nil {
			found2 = post(n, id)
		}
	}

	var foundNode2 *html.Node
	startElement2 := func(n *html.Node, id string) bool {
		if n.Type == html.ElementNode {
			for _, attr := range n.Attr {
				if attr.Key == "id" && attr.Val == id {
					foundNode2 = n
					return true
				}
			}
		}
		return false
	}

	forEachNode2(doc, id, startElement2, nil)
	return foundNode2
}

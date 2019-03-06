package main

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "参数不足，至少需要 2，实际 %d", len(os.Args)-1)
		return
	}
	doc, err := getDoc(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "getdoc: %v\n", err)
	}
	nodes := ElementsByTagname(doc, os.Args[2:]...)
	for _, n := range nodes {
		showNode(n)
	}
}

func getDoc(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return html.Parse(resp.Body)
}

func ElementsByTagname(doc *html.Node, name ...string) []*html.Node {
	var nodes []*html.Node
	var forEachNode func(n *html.Node, name ...string)
	forEachNode = func(n *html.Node, name ...string) {
		for _, s := range name {
			if n.Data == s {
				nodes = append(nodes, n)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			forEachNode(c, name...)
		}
	}
	
	forEachNode(doc, name...)
	return nodes
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
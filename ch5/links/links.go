// 提供解析连接的函数
package links

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

// 向给定的URL发起HTTP GET 请求
// 解析HTML并返回HTML文档中存在的链接
func Extract(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("get %s: %s", url, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parse %s: %s", url, err)
	}

	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue  // 忽略不合法的URL
				}
				links = append(links, link.String())
			}
		}
	}
	forEachNode(doc, visitNode, nil)  // 只要前序遍历，后续不执行，传nil
	return links, nil
}

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

/* 使用时写的函数
func main() {
	url := "https://baidu.com"
	urls, err := Extract(url)
	if err != nil {  // 错误处理随便写写，不引入新的包
		fmt.Printf("extract: %v\n", err)
		return
	}
	for n, u := range urls {
		fmt.Printf("%2d: %s\n", n, u)
	}
}
*/
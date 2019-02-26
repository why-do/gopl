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

func findElements(rc io.ReadCloser, w io.Writer) {
	doc, err := html.Parse(rc)
	rc.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks1: %v\n", err)
		os.Exit(1)
	}
	out(w, doc)
}

func out(w io.Writer, n *html.Node) {
	if n == nil {
		return
	}
	if n.Type == html.ElementNode && n.Data != "script" && n.Data != "style" {
		if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
			s := n.FirstChild.Data
			if s := strings.TrimSpace(s); len(s) > 0 {
				fmt.Fprintf(w, "%s: %s\n", n.Data, s)
			}
		}
	}
	out(w, n.FirstChild)
	out(w, n.NextSibling)
}

func main() {
	for _, arg := range os.Args[1:] {
		body := fetch(arg)
		findElements(body, os.Stdout)
	}
}

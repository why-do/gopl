// 发送一个 HTTP GET 请求，并且获取文档的字数与图片数量
package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func CountWordsAndImage(url string) (words, images int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		err = fmt.Errorf("parsing HTML: %s", err)
		return
	}
	words, images = countWordsAndImage(doc)
	return
}

func countWordsAndImage(n *html.Node) (words, images int) {
	if n.Type == html.TextNode {
		s := n.Data
		f := strings.NewReader(s)
		input := bufio.NewScanner(f)
		input.Split(bufio.ScanWords)
		for input.Scan() {
			words++
		}
	}
	if n.Type == html.ElementNode && n.Data == "img" {
		images++
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		words_c, images_c := countWordsAndImage(c)
		words += words_c
		images += images_c
	}
	return 
}

func main() {
	for _, url := range os.Args[1:] {
		words, images, err := CountWordsAndImage(url)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Printf("%s: words: %d images: %d", url, words, images)
	}
}
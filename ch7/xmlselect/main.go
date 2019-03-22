package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func xmlselect(r io.Reader) {
	dec := xml.NewDecoder(r)
	var stack []string // 元素名的栈
	var checkStack = []string{"div", "div", "h2"} // 进行对比的元素名的顺序
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "xmlselect: %v\n", err)
			os.Exit(1)
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			stack = append(stack, tok.Name.Local) // push 入栈
		case xml.EndElement:
			stack = stack[:len(stack)-1] // pop 出栈
		case xml.CharData:
			if containsAll(stack, checkStack) {
				fmt.Printf("%s: %s\n", strings.Join(stack, " "), tok)
			}
		}
	}
}

// containsAll 判断 x 是否包含 y 中的所有元素，且顺序一致
func containsAll(x, y []string) bool {
	for len(y) <= len(x) {
		if len(y) == 0 {
			return true
		}
		if x[0] == y[0] {
			y = y[1:]
		}
		x = x[1:]
	}
	return false
}

// ch1/fetch
func fetch(w io.Writer) {
	url := "http://www.w3.org/TR/2006/REC-xml11-20060816"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch : %v\n", err)
		os.Exit(1) // 进程退出时，返回状态码1
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch reading %s: %v\n", url, err)
		os.Exit(1)
	}
	fmt.Fprintf(w, "%s\n", b)
}

func main() {
	buf := bytes.NewBuffer(nil)
	fetch(buf)
	xmlselect(buf)
}

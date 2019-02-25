// 输出从 URL 获取的内容
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	for _, url := range os.Args[1:] {
		url = checkUrl(url)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR fetch request %s: %v\n", url, err)
			os.Exit(1) // 进程退出时，返回状态码1
		}
		_, err = io.Copy(os.Stdout, resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR fetch reading %s: %v\n", url, err)
			os.Exit(1)
		}
	}
}

func checkUrl(s string) string {
	if strings.HasPrefix(s, "http") {
		return s
	}
	return fmt.Sprint("http://", s)
}
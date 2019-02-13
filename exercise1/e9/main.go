package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	for _, url := range os.Args[1:] {
		url = checkUrl(url)
		resp, err := http.Head(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR fetch request %s: %v\n", url, err)
			os.Exit(1) // 进程退出时，返回状态码1
		}
		status := resp.Status
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR fetch reading %s: %v\n", url, err)
			os.Exit(1)
		}
		fmt.Printf("%s: %s\n", url, status)
	}
}

func checkUrl(s string) string {
	if strings.HasPrefix(s, "http") {
		return s
	}
	return fmt.Sprint("http://", s)
}

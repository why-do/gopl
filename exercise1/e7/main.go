package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	for _, url := range os.Args[1:] {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR fetch request %s: %v\n", url, err)
			os.Exit(1) // 进程退出时，返回状态码1
		}
		// b, err := ioutil.ReadAll(resp.Body)
		_, err = io.Copy(os.Stdout, resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR fetch reading %s: %v\n", url, err)
			os.Exit(1)
		}
	}
}

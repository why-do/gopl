// 迷你回声和计数器服务器
package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var mu sync.Mutex
var count int

func main() {
	fmt.Println("http://localhost:8000/hello")
	http.HandleFunc("/", handler)
	fmt.Println("http://localhost:8000/count")
	http.HandleFunc("/count", counter)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// 处理程序回显请求的 URL 的路径部分
func handler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count++
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	mu.Unlock()
}

// 回显目前为止调用的次数
func counter(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	fmt.Fprintf(w, "Count %d\n", count)
	mu.Unlock()
}

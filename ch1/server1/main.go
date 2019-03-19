// 迷你回声服务器
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("http://localhost:8000/hello")
	http.HandleFunc("/", handler) // 回声请求调用处理程序
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// 处理非持续回显请求 URL r 的路径部分
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

package main

import (
	"fmt"
	"log"
	"net/http"
)

import "gopl/ch12/params"

// search 用于处理 /search URL endpoint.
func search(resp http.ResponseWriter, req *http.Request) {
	var data struct {
		Labels     []string `http:"l"`
		MaxResults int      `http:"max"`
		Exact      bool     `http:"x"`
	}
	data.MaxResults = 10 // 设置默认值
	if err := params.Unpack(req, &data); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// ...其他处理代码...
	fmt.Fprintf(resp, "Search: %+v\n", data)
}

func main() {
	fmt.Println("http://localhost:8000/search")                                 // Search: {Labels:[] MaxResults:10 Exact:false}
	fmt.Println("http://localhost:8000/search?l=golang&l=gopl")                 // Search: {Labels:[golang gopl] MaxResults:10 Exact:false}
	fmt.Println("http://localhost:8000/search?l=gopl&x=1")                      // Search: {Labels:[gopl] MaxResults:10 Exact:true}
	fmt.Println("http://localhost:8000/search?x=true&max=100&max=200&l=golang") // Search: {Labels:[golang] MaxResults:200 Exact:true}
	fmt.Println("http://localhost:8000/search?q=hello")                         // Search: {Labels:[] MaxResults:10 Exact:false}  # 不存在的参数会忽略
	fmt.Println("http://localhost:8000/search?x=123")                           // x: strconv.ParseBool: parsing "123": invalid syntax  # x 提供的参数解析错误
	fmt.Println("http://localhost:8000/search?max=lots")                        // max: strconv.ParseInt: parsing "lots": invalid syntax  # max 提供的参数解析错误
	http.HandleFunc("/search", search)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

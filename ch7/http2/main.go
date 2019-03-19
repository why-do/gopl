package main

import (
	"fmt"
	"log"
	"net/http"
)

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database map[string]dollars

func (db database) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/list":
		for item, price := range db {
			fmt.Fprintf(w, "%s: %s\n", item, price)
		}
	case "/price":
		item := req.URL.Query().Get("item")
		price, ok := db[item]
		if !ok {
			w.WriteHeader(http.StatusNotFound) // 404
			fmt.Fprintf(w, "no such item: %q\n", item)
			// 也可以用 http.Error 实现上面2行的效果
			// http.Error(w, fmt.Sprintf("no such item: %q\n", item), http.StatusNotFound)
			return
		}
		fmt.Fprintf(w, "%s\n", price)
	default:
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such page: %s\n", req.URL)
		// http.Error(w, fmt.Sprintf("no such page: %s\n", req.URL), http.StatusNotFound)
	}
}

func main() {
	db := database{"shoes": 50, "socks": 5}
	fmt.Println("http://localhost:8000/list")
	fmt.Println("http://localhost:8000/price?item=shoes")
	log.Fatal(http.ListenAndServe("localhost:8000", db))
}

package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("http://localhost:8000")
	handler := func(w http.ResponseWriter, r *http.Request) {
		if order := r.URL.Query().Get("o"); order != "" {
			orderBy(tracks, order)
		}
		showTracks(w, tracks)
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

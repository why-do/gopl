package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	const templ = `<p>A: {{.A}}</p><p>B: {{.B}}</p>`
	t := template.Must(template.New("escape").Parse(templ))
	var data struct {
		A string        // 不受信任的纯文本
		B template.HTML // 受信任的HTML
	}
	data.A = "<b>Hello!</b>"
	data.B = "<b>Hello!</b>"

	fmt.Println("http://localhost:8000")
	handler := func(w http.ResponseWriter, r *http.Request) {
		if err := t.Execute(w, data); err != nil {
			log.Fatal(err)
		}
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

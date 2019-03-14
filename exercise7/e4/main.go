package main

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/net/html"
)

type MyReader string

func (r *MyReader) Read(b []byte) (n int, err error) {
	n = copy(b, *r)
	err = io.EOF
	return
}

func NewMyReader(s string) *MyReader {
	myStr := MyReader(s)
	return &myStr
}

func main() {
	var s1 MyReader = `<h1>Hello</h1>`
	fmt.Printf("%T %[1]s\n", s1)

	doc, err := html.Parse(&s1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks1: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(doc)
}

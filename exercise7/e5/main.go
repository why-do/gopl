package main

import (
	"fmt"
	"io"
	"strings"
)

type MyReader struct {
	rd io.Reader
	i  int64
}

func (r *MyReader) Read(b []byte) (n int, err error) {
	var l int
	if len(b) > int(r.i) {
		l = int(r.i)
	} else {
		l = len(b)
	}
	n, err = r.rd.Read(b[:l])
	if err != nil {
		return n, err
	}
	r.i -= int64(n)
	if r.i <= 0 {
		err = io.EOF
	}
	return
}

func LimitReader(r io.Reader, n int64) io.Reader {
	return &MyReader{r, n}
}

func main() {
	s1 := "abcdefghijklmnopqrstuvwxyz"
	r1 := strings.NewReader(s1)
	r2 := LimitReader(r1, 10)
	fmt.Println(r1)
	fmt.Println(r2)
	var b1 = make([]byte, 26)
	var b2 = make([]byte, 6)
	// r1.Read(b1)
	r2.Read(b2)
	fmt.Println(string(b1))
	fmt.Println(string(b2))
}

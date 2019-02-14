package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	fs := []func() string {echo1, echo2, echo3}
	times := 100000  // 重复执行多次
	for _, f := range fs {
		var s string
		start := time.Now()
		for i := 0; i < times; i++ {
			s = f()
		}
		secs := time.Since(start).Seconds()
		fmt.Printf("%fs\t", secs)
		fmt.Println(len(s))
	}
}

func echo1() string {
	var s, sep string
	for i := 1; i < len(os.Args); i++ {
		s += sep + os.Args[i]
		sep = " "
	}
	return s
}

func echo2() string {
	s, sep := "", ""
	for _, arg := range os.Args[1:] {
		s += sep + arg
		sep = " "
	}
	return s
}

func echo3() string {
	return fmt.Sprintf(strings.Join(os.Args[1:], " "))
}

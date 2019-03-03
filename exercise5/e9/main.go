package main

import (
	"fmt"
	"strings"
)

func main() {
	s1 := "This is a $foo test for $foo1 $foo2"
	s2 := expand(s1, foo)
	fmt.Println(s2)
}

// 该函数替换参数 s 中每一个子字符串 "$foo" 为 `f("foo")` 的返回值
func expand(s string, f func(string) string) string {
	// strings.Replace(str string, old string, new string, n int) string
	return strings.Replace(s, "$foo", f("foo"), -1)
}

func foo(s string) string {
	return strings.Join([]string{"{{.", s, "}}"}, "")
}

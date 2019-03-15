package main

import (
	"flag"
	"fmt"
	"strings"
)

type urls []string

func (v *urls) String() string {
	// *v = []string{"baidu.com"} // 通过指针改变初始值
	r := []string{}
	for _, s := range *v {
		r = append(r, fmt.Sprintf("%q", s))
	}
	return strings.Join(r, ", ")
}

var isNew bool
func (v *urls) Set(s string) error {
	if !isNew {
		*v = nil
		isNew = true
	}
	*v = append(*v, s)
	return nil
}

func main() {
	var value urls
	// value = append(value, "baidu.com") // 传递给Var函数前就设定好初始值
	flag.Var(&value, "url", "域名")
	flag.Parse()
	fmt.Printf("%T % [1]q\n", value)
	s := []string(value)
	fmt.Printf("%T % [1]q\n", s)
}
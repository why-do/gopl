package main

import (
	"flag"
	"fmt"
	"strings"
)

type urls []string

func (v *urls) String() string {
	r := []string{}
	for _, s := range *v {
		r = append(r, fmt.Sprintf("%q", s))
	}
	return strings.Join(r, ", ")
}

func (v *urls) Set(s string) error {
	// *v = nil // 不能再清空原有的记录了
	// strings.Fields 可以区分连续的空格
	*v = append(*v, s)
	return nil
}

func Urls(name string, value []string, usage string) *[]string {
	p := new([]string) // value 是传值进来的，取不到地址，new一个内存空间，存放value的值
	*p = value
	flag.CommandLine.Var((*urls)(p), name, usage)
	return p
}


func main() {
	s := Urls("url", []string{"baidu.com"}, "域名")
	flag.Parse()
	fmt.Printf("% q\n", *s)
}
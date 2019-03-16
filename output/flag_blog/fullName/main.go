package main

import (
	"flag"
	"fmt"
	"strings"
)

type fullName []string

func (v *fullName) String() string {
	r := []string{}
	for _, s := range *v {
		r = append(r, fmt.Sprintf("%q", s))
	}
	return strings.Join(r, " ")
}

func (v *fullName) Set(s string) error {
	*v = nil
	// strings.Fields 可以区分连续的空格
	for _, str := range strings.Fields(s) {
		*v = append(*v, str)
	}
	return nil
}

func FullName(name string, value []string, usage string) *[]string {
	p := new([]string) // value 是传值进来的，取不到地址，new一个内存空间，存放value的值
	*p = value
	flag.CommandLine.Var((*fullName)(p), name, usage)
	return p
}

func main() {
	s := FullName("name", []string{"Karl", "Lichter", "Von", "Randoll"}, "全名")
	flag.Parse()
	fmt.Printf("% q\n", *s)
}

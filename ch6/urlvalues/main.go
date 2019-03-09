package main

import "fmt"

// 映射字符串到字符串列表
type Values map[string][]string

func (v Values) Get(key string) string {
	if vs := v[key]; len(vs) > 0 {
		return vs[0]
	}
	return ""
}

func (v Values) Add(key, value string) {
	v[key] = append(v[key], value)
}

func main() {
	m := Values{"lang": {"en"}}
	m.Add("item", "1")
	m.Add("item", "2")

	fmt.Printf("%q\n", m.Get("lang"))
	fmt.Printf("%q\n", m.Get("q"))
	fmt.Printf("%q\n", m.Get("item"))
	fmt.Printf("%q\n", m["item"])

	m = nil
	fmt.Printf("%q\n", m.Get("item"))
	m.Add("item", "3")  // panic，赋值给一个值为nil的map
}

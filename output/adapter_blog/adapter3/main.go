package main

import "fmt"

type People struct {
	name string
}

func (p *People) Greet() {
	fmt.Printf("Hello, I am %s.\n", p.name)
}

var p1 = People{"Adam"}

func init() {
	fmt.Print("init: ")
	p1.Greet()
}

// 接口
type Hello interface {
	SayHello()
}

func SayHello(s Hello) {
	s.SayHello()
}

func main() {
	SayHello(PeopleSayHello(p1.Greet))
}

// 适配器
type PeopleSayHello func()

func (f PeopleSayHello) SayHello() {
	f()
}


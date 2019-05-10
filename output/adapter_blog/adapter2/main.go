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
	SayHello(&PeopleSayHello{&p1})
}

// 适配器
type PeopleSayHello struct {
	*People
}

func (p *PeopleSayHello) SayHello() {
	p.Greet()
}



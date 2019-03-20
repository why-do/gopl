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
	p2 := PeopleSayHello(p1)
	SayHello(&p2)	
}

// 适配器
type PeopleSayHello People

func (p *PeopleSayHello) SayHello() {
	p1 := People(*p)
	p1.Greet()
}



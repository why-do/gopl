package main

import "fmt"

// 这可以表示一个坐标
type Point struct {
	X, Y int
}

// 坐标加上半径就是一个圆
type Circle struct {
	Point
	Radius int
}

// 圆加上辐条数，这表示一个轮子
type Wheel struct {
	Circle
	Spokes int
}

var w Wheel

func main() {
	w = Wheel{Circle{Point{8, 8}, 5}, 20}
	w = Wheel{
		Circle: Circle{
			Point:  Point{X: 8, Y: 8},
			Radius: 5,
		},
		Spokes: 20,
	}

	fmt.Printf("%v\n", w)
	fmt.Printf("%#v\n", w)
}

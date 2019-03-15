package main

import (
	"flag"
	"fmt"
)

type Celsius float64
type Fahrenheit float64

func CToF(c Celsius) Fahrenheit { return Fahrenheit(c*9.0/5.0 + 32.0) }
func FToC(f Fahrenheit) Celsius { return Celsius((f - 32.0) * 5.0 / 9.0) }

func (c Celsius) String() string { return fmt.Sprintf("%g°C", c) }

// 上面这些都是之前定义过的内容，是可以作为包引出过来的
// 为了说明清楚，就单独把需要用到的部分复制过来

// *celsiusValue 满足 flag.Vulue 接口
// 同一个包不必这么麻烦，直接定义 Celsius 类型即可。这里假设是从别的包引入的类型
type celsiusValue Celsius

func (c *celsiusValue) String() string { return fmt.Sprintf("%.2f°C", *c) }
// func (c *celsiusValue) String() string { return (*Celsius)(c).String() }

func (c *celsiusValue) Set(s string) error {
	var unit string
	var value float64
	fmt.Sscanf(s, "%f%s", &value, &unit) // 无须检查错误
	switch unit {
	case "C", "°C":
		*c = celsiusValue(value)
		return nil
	case "F", "°F":
		*c = celsiusValue(FToC(Fahrenheit(value)))
		return nil
	}
	return fmt.Errorf("invalid temperature %q", s)
}

func CelsiusFlag(name string, value Celsius, usage string) *Celsius {
	p := new(Celsius) // value 是传值进来的，取不到地址，new一个内存空间，存放value的值
	*p = value
	flag.CommandLine.Var((*celsiusValue)(p), name, usage)
	return p
}

func main() {
	tempP := CelsiusFlag("temp", 36.7, "温度")
	flag.Parse()
	fmt.Printf("%T, %[1]v\n", tempP)
}

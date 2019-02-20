// 演示切片的就地修改算法
package main

import "fmt"

// nonempty 返回一个新的切片，切片中的元素都是非空字符串
// 在函数的调用过程中，底层数组的元素发生了改变
func nonempty(strings []string) []string {
	i := 0
	for _, s := range strings {
		if s != "" {
			strings[i] = s
			i++
		}
	}
	return strings[:i]
}

func nonempty2(strings []string) []string {
	out := strings[:0]  // 直接引用原始切片的零长度切片
	for _, s := range strings {
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

func main() {
	f := nonempty
	// f := nonempty2
	data := []string{"one", "", "three"}
	fmt.Printf("%q\n", f(data))
	fmt.Printf("%q\n", data)
}
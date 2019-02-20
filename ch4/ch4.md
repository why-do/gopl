# 4.2 slice

## 反转和平移
就地反转slice中的元素：
```go
package main

import "fmt"

func reverse(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func main() {
	l := [...]int{1, 2, 3, 4, 5} // 这个是数组
	fmt.Println(l)
	reverse(l[:]) // 传入切片
	fmt.Println(l)
}
```

将一个切片向左平移n个元素的简单方法是连续调用三次反转函数。第一次反转前n个元素，第二次返回剩下的元素，最后整体做一次反转：
```go
func moveLeft(n int, s []int) {
	reverse(s[:n])
	reverse(s[n:])
	reverse(s)
}

func moveRight(n int, s []int) {
	reverse(s[n:])
	reverse(s[:n])
	reverse(s)
}
```

## 切片的比较
与数组不同，切片无法做比较。标准库中提供了高度优化的函数 bytes.Equal 来比较两个字节切片（[]byte）。但是对其他类型的切片，Go不支持比较。当然自己写一个比较的函数也不难：
```go
func equal(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}
```
上面的方法也只是返回执行函数当时的结果，但是切片的底层数组可以能发生改变，在不同的时间切片所拥有的元素可能不同，不能保证整个生命周期都保持不变。总之，Go不允许直接比较切片。  

**和nil比较**  
切片唯一允许的比较操作是和nil做比较。值为nil的切片长度和容量都是零，但是也有非nil的切片长度和容量也都是零的：
```go
func main() {
	var s []int
	fmt.Println(s == nil)  // true
	s = nil
	fmt.Println(s == nil)  // true
	s = []int(nil)
	fmt.Println(s == nil)  // true
	s = []int{}
	fmt.Println(s == nil)  // flase
}
```
所以要检查一个切片是否为空，应该使用 len(s) == 0，而不是和nil做比较。  
另外，值为nil的切片其表现和其它长度为零的切片是一样的。无论值是否为nil，GO的函数都应该以相同的方式对待所有长度为零的切片。  


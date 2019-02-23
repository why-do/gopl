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

## 初始化
像切片和map这类引用类型，使用前是需要初始化的。仅仅进行声明，是不分配内存的，此时值为nil。  
完成初始化后（大括号或者make函数），此时就是已经完成了初始化，分配内存空间，值不为nil。  

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


# 4.3 map

## 集合
Go 没有提供集合类型，但是利用key唯一的特点，可以用map来实现这个功能。比如说字符串的集合：
```go
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	seen := make(map[string]bool) // 字符串集合
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		if !seen[line] {
			seen[line] = true
			fmt.Println("Set:", line)
		}
	}
	if err := input.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "dedup: %v\n", err)
		os.Exit(1)
	}
}
```
从标准输出获取字符串，用map来存储已经出现过的行，只有首次出现的字符串才会打印出来。  

**使用空结构体作value**  
这里使用bool来作为map的value，而bool也有true和false两种值，而实际只使用了1种值。  
这里还可以使用空结构体（类型：struct{}、值：struct{}{}）。空结构体，没有长度，也不携带任何信息，用它可能是最合适的。但由于这种方式节约的内存很少并且语法复杂，所以一般尽量避免这样使用。  

## 使用切片做key
切片是不能作为key的，并且切片是不可比较的，不过可以有一个间接的方法来实现切片作key。定义一个帮助函数k，将每一个key都映射到字符串：
```go
var m = make(map[string]int)

func k(list []string) string { fmt.Sprint("%q", list) }

func Add(list []string) { m[k(list)]++ }
func Count(list []string) int { return m[k(list)] }
```
这里使用%q来格式化切片，就是包含双引号的字符串，所以（\["ab", "cd"\] 和 \["abcd"\]）是不一样的。就是，当且仅当 x 和 y 相等的时候，才认为 k(x)==k(y)。  
同样的方法适用于任何不可直接比较的key类型，不仅仅局限于切片。同样，k(x) 的类型不一定是字符串类型，任何能够得到想要的比较结果的可比较类型都可以。  


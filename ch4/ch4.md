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

# 4.5 JSON
JSON是一种发送和接收格式化信息的标准。JSON不是唯一的标准，XML、ASN.1 和 Google 的 Protocol Buffer 都是相似的标准。Go通过标准库 encoding/json、encoding/xml、encoding/asn1 和其他的库对这些格式的编码和解码提供了非常好的支持，这些库都拥有相同的API。  

## 序列化输出
首先定义一组数据：
```go
type Movie struct {
	Title  string
	Year   int  `json:"released"`
	Color  bool `json:"color,omitempty"`
	Actors []string
}

var movies = []Movie{
	{Title: "Casablanca", Year: 1942, Color: false,
		Actors: []string{"Humphrey Bogart", "Ingrid Bergman"}},
	{Title: "Cool Hand Luke", Year: 1967, Color: true,
		Actors: []string{"Paul Newman"}},
	{Title: "Bullitt", Year: 1968, Color: true,
		Actors: []string{"Steve McQueen", "Jacqueline Bisset"}},
}
```
然后通过 json.Marshal 进行编码：
```go
data, err := json.Marshal(movies)
if err != nil {
	log.Fatalf("JSON Marshal failed: %s", err)
}
fmt.Printf("%s\n", data)

/* 执行结果
[{"Title":"Casablanca","released":1942,"Actors":["Humphrey Bogart","Ingrid Bergman"]},{"Title":"Cool Hand Luke","released":1967,"color":true,"Actors":["Paul Newman"]},{"Title":"Bullitt","released":1968,"color":true,"Actors":["Steve McQueen","Jacqueline Bisset"]}]
*/
```
这种紧凑的表示方法适合传输，但是不方便阅读。有一个 json.MarshalIndent 的变体可以输出整齐格式化过的结果。多传2个参数，第一个是定义每行输出的前缀字符串，第二个是定义缩进的字符串：
```go
data, err := json.MarshalIndent(movies, "", "    ")
if err != nil {
	log.Fatalf("JSON Marshal failed: %s", err)
}
fmt.Printf("%s\n", data)

/* 执行结果
[
    {
        "Title": "Casablanca",
        "released": 1942,
        "Actors": [
            "Humphrey Bogart",
            "Ingrid Bergman"
        ]
    },
    {
        "Title": "Cool Hand Luke",
        "released": 1967,
        "color": true,
        "Actors": [
            "Paul Newman"
        ]
    },
    {
        "Title": "Bullitt",
        "released": 1968,
        "color": true,
        "Actors": [
            "Steve McQueen",
            "Jacqueline Bisset"
        ]
    }
]
*/
```

只有可导出的成员可以转换为JSON字段，上面的例子中用的都是大写。  
**成员标签定义**（field tag），是结构体成员的编译期间关联的一些元素信息。标签值的第一部分指定了Go结构体成员对应的JSON中字段的名字。  
另外，Color标签还有一个额外的选项 omitempty，它表示如果这个成员的值是零值或者为空，则不输出这个成员到JSON中。所以Title为"Casablanca"的JSON里没有color。  

## 反序列化
反序列化操作将JSON字符串解码为Go数据结构。这个是由 json.Unmarshal 实现的。
```go
var titles []struct{ Title string }
if err := json.Unmarshal(data, &titles); err != nil {
	log.Fatalf("JSON unmarshaling failed: %s", err)
}
fmt.Println(titles)

/* 执行结果
[{Casablanca} {Cool Hand Luke} {Bullitt}]
*/
```
这里接收数据时定义的结构体只有一个Title字段，这样当函数 Unmarshal 调用完成后，将填充结构体切片中的 Title 值，而JSON中其他的字段就丢弃了。  
反序列化的时候，成员标签定义一样有效。Go里可导出的字段总是首字母大写的，而一般定义JSON的字段都是小写，所以总是需要为每个字段定义标签。不过这里还是可以偷个懒，如果只有反序列化，而不需要做序列化。在 Unmarshal 阶段，JSON字段的名称关联到Go结构体成员的名称是忽略大小写的。*序列化是不会有这个便捷，因为序列化的时候没有JSON需要关联。*另外，小写的变量在需要分词的时候，可能会使用下划线分割，这种情况下，还是要用一下成员标签定义。  

## 流式解码器
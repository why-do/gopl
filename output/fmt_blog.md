# fmt 包
功能：fmt包实现了类似C语言printf和scanf的格式化I/O。格式化动作（'verb'）源自C语言但更简单。  

## fmt格式化输出（ch1.md）
Printf 函数有超过10个各种转义字符，Go 程序员称为 verb。下表不完整，但是它说明了很多可用的功能：
| verb | 描述 |
|-----|-----|
| %d | 十进制数 |
| %x, %o, %b | 十六进制，八进制，二进制数 |
| %f, %g, %e | 浮点数 |
| %t | 布尔型 |
| %c | 字符（Unicode码点） |
| %s | 字符串 |
| %q | 带引号字符串或者字符 |
| %v | 内置格式的任何值 |
| %T | 任何值的类型 |
| %% | 百分号本身 |

以ln结尾的，比如fmt.Println，使用%v的方式来格式化参数，并且在最后追加换行符。

## 更详细的占位符说明
上面的表格比较概括，而且也不全，对于不同的类型，还有不同的细节。所有的说明都在fmt包的doc.go文档里有详细的说明：http://docscn.studygolang.com/src/fmt/doc.go

**普通占位符**：  
| 占位符 | 说明 | 举例 | 输出 |
|-----|-----|-----|-----|
| %v | 相应值的默认格式。 | Printf("%v", people) | {zhangsan} |
| %+v | 打印结构体时，会添加字段名 | Printf("%+v", people) | {Name:zhangsan} |
| %#v | 相应值的Go语法表示 | Printf("%#v", people)  | main.Human{Name:"zhangsan"} |
| %T | 相应值的类型的Go语法表示 | Printf("%T", people) | main.Human |
| %% | 字面上的百分号，并非值的占位符 | Printf("%%") | % |

**布尔占位符**：  
| 占位符 | 说明 | 举例 | 输出 |
|-----|-----|-----|-----|
| %t | true 或 false。 | Printf("%t", true) | true |

**整数占位符**：  
| 占位符 | 说明 | 举例 | 输出 |
|-----|-----|-----|-----|
| %b | 二进制表示 | Printf("%b", 5) | 101
| %c | 相应Unicode码点所表示的字符 | Printf("%c", 0x4E2D) | 中 |
| %d | 十进制表示 | Printf("%d", 0x12) | 18 |
| %o | 八进制表示 | Printf("%d", 10) | 12 |
| %q | 单引号围绕的字符字面值，由Go语法安全地转义 | Printf("%q", 0x4E2D) | '中' |
| %x | 十六进制表示，字母形式为小写 a-f | Printf("%x", 13) | d |
| %X | 十六进制表示，字母形式为大写 A-F | Printf("%x", 13) | D |
| %U | Unicode格式，相当于 "%04X" 加上前导 "U+" | Printf("%U", 0x4E2D) | U+4E2D |

**浮点数和复数**：
| 占位符 | 说明 | 举例 | 输出 |
|-----|-----|-----|-----|
| %b | 无小数部分的，指数为二的幂的科学计数法，与 strconv.FormatFloat 的 'b' 转换格式一致。例如 -123456p-78 |
| %e | 科学计数法，例如 -1234.456e+78 | Printf("%e", 10.2) | 1.020000e+01 |
| %E | 科学计数法，例如 -1234.456E+78 | Printf("%e", 10.2) | 1.020000E+01 |
| %f | 有小数点而无指数，例如 123.456 | Printf("%f", 10.2) | 10.200000 |
| %g | 根据情况选择 %e 或 %f 以产生更紧凑的（无末尾的0）输出 | Printf("%g", 10.20) | 10.2 |
| %G | 根据情况选择 %E 或 %f 以产生更紧凑的（无末尾的0）输出 | Printf("%G", 10.20+2i) | (10.2+2i) |

**指针**：  
| 占位符 | 说明 | 举例 | 输出 |
|-----|-----|-----|-----|
| %p | 十六进制表示，前缀 0x  | Printf("%p", &people) | 0x4f57f0 |

**其他标记（副词）**：
+ "\+" ： 总打印数值的正负号；对于%q（%+q）保证只输出字符的编码。`fmt.Printf("%q %+[1]q\n", "中文") // "中文" "\u4e2d\u6587"`
+ "-" ： 在右侧而非左侧填充空格（左对齐该区域）
+ " " ： (空格)为数值中省略的正负号留出空白（% d）； 以十六进制（% x, % X）打印字符串或切片时，在字节之间用空格隔开
+ "0" ： 填充前导的0而非空格；对于数字，这会将填充移到正负号之后
+ "\#" ： 备用格式，不同的类型还不一样。
+ "\*" ： 使用一个变量来控制输出的宽度，实现可变宽度。

副词#的备用格式：
八进制、十六进制，默认没有前导，使用#后会添加前导符号"0"、"0x"、"0X"，防止产生歧义：
```go
fmt.Printf("%o %#[1]o\n", 123)        // 173 0173
fmt.Printf("%x %#[1]x %#[1]X\n", 123) // 7b 0x7b 0X7B
```
指针默认有前导，备用格式就是就掉前导：
```go
var s string
fmt.Printf("%p %#[1]p\n", &s) // 0xc00004c240 c00004c240
```
对于字符串，%#q有些情况下会输出反引号围绕的字符串，不过测试下来不总是这样：
```go
fmt.Printf("%q, %#[1]q\n", "ab\tcd") // "ab\tcd", `ab   cd`
fmt.Printf("%q, %#[1]q\n", "ab\ncd") // "ab\ncd", "ab\ncd"
```
对于Unicode，打印出字符的编码后还会打印该字符：
```go
fmt.Printf("%U, %#[1]U", '中') // U+4E2D, U+4E2D '中'
```

# 使用示例
下面是一些可以通过合理的构造格式化字符串来打到最佳的显示效果的示例。  
可能有些示例之间会有点重复。还有一些示例因为太简单，感觉上一节已经讲过了。不过上一节的介绍偏理论，而这里的内容更注重实际使用。  

## fmt的两个技巧（ch3.md）
一、%后的副词[1]告知Printf重复使用第一个操作数。  
二、%o、%x、%X之前的副词#告知Printf输出相应的前缀 0、0x、0X。  
```go
func main() {
	o := 0666
	fmt.Printf("%d %[1]o %#[1]o\n", o)  // 438 666 0666
	x := int64(0xdeadbeef)
	fmt.Printf("%d %[1]x %#[1]x %#[1]X\n", x)  // 3735928559 deadbeef 0xdeadbeef 0XDEADBEEF
}
```

## 输出字节序列
使用%x可以输出字节序列的UTF-8编码，还可以加上空格也就是% x，这样每个字符还能隔开，在对字节序列进行输出的时候特别有用：
```go
fmt.Printf("%x\n", "abcdefg")   // "61626364656667"
fmt.Printf("% x\n", "abcdefg")  // "61 62 63 64 65 66 67"
fmt.Printf("% #x\n", "abcdefg") // "0x61 0x62 0x63 0x64 0x65 0x66 0x67"
```
对于Unicode字符，还是按照字节处理的，如若要输出Unicode码点，需要先转成[]rune类型：
```go
fmt.Printf("% x\n", "世界")        // "e4 b8 96 e7 95 8c"
fmt.Printf("%x\n", []rune("世界")) // "[4e16 754c]"
```

## 指定宽度
通过指定相同的宽度，可以做到右对齐的效果：
```go
func main() {
	fmt.Printf("%4d\n", 1)
	fmt.Printf("%4d\n", 10)
	fmt.Printf("%4d\n", 100)
	fmt.Printf("%4d\n", 1000)
	fmt.Printf("%4d\n", 10000)  // 这个会超出宽度
}

/* 输出结果
$ go run main.go
   1
  10
 100
1000
10000
*/
```
如果宽度不够，输出时也不会丢失信息，而是把信息全部输出，不受宽度的限制。  
默认使用空格填充，也可以指定填充的内容，比如使用0填充，在输出二进制数的时候非常有用：
```go
func main() {
	fmt.Printf("%04b\n", 1)
	fmt.Printf("%04b\n", 2)
	fmt.Printf("%04b\n", 3)
	fmt.Printf("%04b\n", 4)
	fmt.Printf("%04b\n", 666)  // 这个会超出宽度
}

/* 输出结果
$ go run main.go
0001
0010
0011
0100
1010011010
*/
```

## 宽度和精度
操作数字的时候，宽度为该数值占用区域的最小宽度；精度为小数点之后的位数。  
对于 %g 和 %G 精度是所有数字的总和，而用 %f 打印出来同样是小数，精度是小数点后面的位数。比如：123.45，%.4g 是 "123.5" 而 %.2f 是 "123.45"。
```go
const n float64 = 123.45
fmt.Printf("%.4g %.2[1]f", n) // g:123.5 f:123.45
```

## 打印结构体
打印结构体的时候，使用副词#或者+可以使结构化符号%v以类似Go语法的方式输出对象，这个方法里面包含了成员变量的名字：
```go
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

/* 执行结果
PS H:\Go\src\gopl\ch4\embed> go run main.go
{{{8 8} 5} 20}
main.Wheel{Circle:main.Circle{Point:main.Point{X:8, Y:8}, Radius:5}, Spokes:20}
PS H:\Go\src\gopl\ch4\embed>
*/
```

## 字符串对齐
宽度和精度对于字符串输出同样有效：
+ 宽度为输出的最小字符数，如果必要的话会为已格式化的形式填充空格。
+ 精度为输出的最大字符数，如果必要的话会直接截断。

这样就可以在输出字符串的时候做到左对齐和右对齐，输出类似表格的样式：
```go
package main

import "fmt"

type message struct {
	Title string
	text  string
}

func main() {
	list := []message{
		{"fmt", "Package fmt implements formatted I/O with functions analogous to C's printf and scanf. The format 'verbs' are derived from C's but are simpler. "},
		{"bytes", "Package bytes implements functions for the manipulation of byte slices. It is analogous to the facilities of the strings package. "},
		{"time", "Package time provides functionality for measuring and displaying time. "},
		{"net/http", "Package http provides HTTP client and server implementations. "},
	}
	_ = list
	for _, msg := range list {
		fmt.Printf("%9.9s %.99s\n", msg.Title, msg.text)
	}
}

/* 输出效果
PS H:\Go\src\gopl\output> go run main.go
      fmt Package fmt implements formatted I/O with functions analogous to C's printf and scanf. The format '
    bytes Package bytes implements functions for the manipulation of byte slices. It is analogous to the faci
     time Package time provides functionality for measuring and displaying time.
 net/http Package http provides HTTP client and server implementations.
PS H:\Go\src\gopl\output>
*/
```
不过对于中文输出就不会那么漂亮了，因为这里填充的空格是半角的空格。  

**左对齐**  
其实只要用上制表符\t就能对齐了，不过如果字符串的长度相差比较大，也无法对齐，并且间隔是自动的。下面是另一种精确的实现方法。  
默认是在前面进行填充，似乎也没有在后面填充的格式。不过可以用下面的方法输出左对齐的效果：
```go
fmt.Printf("%s%*s %.99s\n", msg.Title, 10-len(msg.Title), "", msg.text)
```
这里的`%*s`是在s对应的内容输出前，动态的填充由变量指定的前导空格。由于这里s对应的是空字符串，所以就相当于是动态添加空格，而空格的数量是可以根据前一个变量的长度计算后动态变化的。这是下一节的内容，不过这里正好先展示下效果。  

## 动态控制宽度（缩进）
还有一种控制宽度的方法，使用\*号，这种方式输出的宽度是由之后的变量决定的所以是可变的：
```go
func main() {
	for i := 0; i < 3; i++ {
		fmt.Printf("%*d\n", i*4, i)
	}
}

/* 执行结果
PS H:\Go\src\gopl\output> go run main.go
0
   1
       2
PS H:\Go\src\gopl\output> go run main.go
*/
```
输出的宽度由第一个变量控制，而输出的内容是第二个变量。  
这里加了个0，这样就用0来代替原来的空格来填补宽度了，下面显示的效果更加直观：
```go
func main() {
	for i := 0; i < 3; i++ {
		fmt.Printf("%0*d\n", i*4, i)
	}
}

/* 执行结果
PS H:\Go\src\gopl\output> go run main.go
0
0001
00000002
PS H:\Go\src\gopl\output>
*/
```
由于宽度是包含字符串本身的，而缩进效果是不包括字符串的。所以可是使用空字符串作为第二个参数，把要输出的内容放在后面。  
推荐这么用，`fmt.Printf("%*s<%s>\n", n*2, "", "div")` 第一个变量控制宽度，第二个是输出的字符串，这里是空字符串，这里就是两格缩进的效果。把要输出的内容放在后面，这里的第三个变量。  
这种显示效果在处理代码和html标签的时候特别好，下面的程序先通过Get请求获取一个页面，然后解析页面中的文档树，并输出有缩进效果的树结构：
```go
// ch5/outline2

/* 执行结果
PS H:\Go\src\gopl\ch5\outline2> go run main.go http://baidu.com
<html>
  <head>
    <meta>
    </meta>
  </head>
  <body>
  </body>
</html>
PS H:\Go\src\gopl\ch5\outline2>
*/
```
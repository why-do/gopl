# fmt

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
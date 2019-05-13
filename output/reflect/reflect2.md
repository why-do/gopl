# 12.5 使用 reflect.Value 来设置值
到目前为止，反射只是用来**解析**变量值。本节的重点是**改变**值。  

## 可寻址的值（canAddr）
reflect\.Value 的值，有些是可寻址的，有些是不可寻址的。通过 reflect\.ValueOf(x) 返回的 reflect\.Value 都是不可寻址的。但是通过指针提领得来的 reflect\.Value 是可寻址的。可以通过调用 reflect\.ValueOf(&x).Elem() 来获得任意变量 x 可寻址的 reflect\.Value 值。  
可以通过变量的 CanAddr 方法来询问 reflect\.Value 变量是否可寻址：
```go
x := 2                   // value   type    variable?
a := reflect.ValueOf(2)  // 2       int     no
b := reflect.ValueOf(x)  // 2       int     no
c := reflect.ValueOf(&x) // &x      *int    no
d := c.Elem()            // 2       int     yes (x)

fmt.Println(a.CanAddr()) // "false"
fmt.Println(b.CanAddr()) // "false"
fmt.Println(c.CanAddr()) // "false"
fmt.Println(d.CanAddr()) // "true"
```

## 更新变量（Set）
从一个可寻址的 reflect\.Value() 获取变量需要三步：
1. 调用 Addr()，返回一个 Value，其中包含一个指向变量的指针
2. 在这个 Value 上调用 interface()，返回一个包含这个指针的 interface{} 值
3. 如果知道变量的类型，使用类型断言把空接口转换为一个普通指针

之后，就可以通过这个指针来更新变量了：
```go
x := 2
d := reflect.ValueOf(&x).Elem()   // d代表变量x
px := d.Addr().Interface().(*int) // px := &x
*px = 3                           // x = 3
fmt.Println(x)                    // "3"
```

还有一个方法，可以直接通过可寻址的 reflect\.Value 来更新变量，不用通过指针，而是直接调用 reflect\.Value\.Set 方法：
```go
d.Set(reflect.ValueOf(4))
fmt.Println(x) // "4"
```

## 注意事项
如果类型不匹配会导致程序崩溃：
```go
d.Set(reflect.ValueOf(int64(5))) // panic: int64 不可赋值给 int
```
在一个不可寻址的 reflect\.Value 上调用 Set 方法也会使程序崩溃：
```go
x := 2
b := reflect.ValueOf(x)
b.Set(reflect.ValueOf(3)) // panic: 在不可寻址的值上使用 Set 方法
```

另外还提供了一些为基本类型特化的 Set 变种：SetInt、SetUint、SetString、SetFloat等：
```go
d := reflect.ValueOf(&x).Elem()
d.SetInt(3)
fmt.Println(x) // "3"
```
这些方法还有一定的容错性。比如 SetInt 方法，任意有符号整型，甚至是底层类型是有符号整型的命名类型，都可以执行成功。如果值太大了，会无提示地截断它。但是在指向 interface{} 变量的 reflect\.Value 上调用 SetInt 会崩溃（尽管使用 Set 是没有问题的）：
```go
x := 1
rx := reflect.ValueOf(&x).Elem()
rx.SetInt(2)                     // OK, x = 2
rx.Set(reflect.ValueOf(3))       // OK, x = 3
rx.SetString("hello")            // panic: string 不能赋值给 int
rx.Set(reflect.ValueOf("hello")) // panic: string 不能赋值给 int

var y interface{}
ry := reflect.ValueOf(&y).Elem()
ry.SetInt(2)                     // panic: 在指向空接口的 Value 上调用 SetInt
ry.Set(reflect.ValueOf(3))       // OK, y = int(3)
ry.SetString("hello")            // panic: 在指向空接口的 Value 上调用 SetString
ry.Set(reflect.ValueOf("hello")) // OK, y = "hello"
```

## 可修改的值（CanSet）
另外，反射是可以越过 Go 言语的导出规则，读取到未导出的成员。但是利用反射不能修改未导出的成员：
```go
stdout := reflect.ValueOf(os.Stdout).Elem() // *os.Stdout, 一个 os.File 变量
fmt.Println(stdout.Type())                  // "os.File"
fd := stdout.FieldByName("fd")
fmt.Println(fd.Int()) // "1" ，获取到了未导出的成员的值
fd.SetInt(2)          // panic: unexported field ，尝试修改则会崩溃
```
一个可寻址的 reflect\.Value 会记录它是否是通过遍历一个未导出的字段来获得的，如果是这样则不允许修改。  
所以在更新变量前用 CanAddr 来检查不能保证正确。CanSet 方法才能正确地报告一个 reflect\.Value 是否可寻址且可更改：
```go
fmt.Println(fd.CanAddr(), fd.CanSet()) // "true false"
```

# 12.6 示例：解码 S 表达式
本节要为 S 表达式编码实现一个简单的 Unmarshal 函数（解码器）。一个健壮的和通用的实现比这里的例子需要更多的代码，这里精简了很多，只支持 S 表达式有限的子集，并且没有优雅地处理错误。代码的目的是阐释反射，而不是语法分析。  

## 词法分析器
词法分析器 lexer 使用 text\/scanner 包提供的扫描器 Scanner 类型来把输入流分解成一系列的标记（token），包括注释、标识符、字符串字面量和数字字面量。扫描器的 Scan 方法将提前扫描并返回下一个标记（类型为 rune）。大部分标记（比如'('）都只包含单个rune，但 text\/scanner 包也可以支持由多个字符组成的记号。调用 Scan 会返回标记的类型，调用 TokenText 则会返回标记的文本。  
因为每个解析器可能需要多次使用当前的记号，但是 Scan 会一直向前扫描，所以把扫描器封装到一个 lexer 辅助类型中，其中保存了 Scan 最近返回的标记：
```go
type lexer struct {
	scan  scanner.Scanner
	token rune // 当前标记
}

func (lex *lexer) next()        { lex.token = lex.scan.Scan() }
func (lex *lexer) text() string { return lex.scan.TokenText() }

func (lex *lexer) consume(want rune) {
	if lex.token != want { // 注意: 错误处理不是这篇的重点，简单粗暴的处理了
		panic(fmt.Sprintf("got %q, want %q", lex.text(), want))
	}
	lex.next()
}
```

## 函数实现
分析器有两个主要的函数。  
一个是read，它读取从当前标记开始的 S 表达式，并更新由可寻址的 reflect\.Value 类型的变量 v 指向的变量：
```go
func read(lex *lexer, v reflect.Value) {
	switch lex.token {
	case scanner.Ident:
		// 仅有的有标识符是 “nil” 和结构体的字段名
		if lex.text() == "nil" {
			v.Set(reflect.Zero(v.Type()))
			lex.next()
			return
		}
	case scanner.String:
		s, _ := strconv.Unquote(lex.text()) // 注意：错误被忽略
		v.SetString(s)
		lex.next()
		return
	case scanner.Int:
		i, _ := strconv.Atoi(lex.text()) // 注意：错误被忽略
		v.SetInt(int64(i))
		lex.next()
		return
	case '(':
		lex.next()
		readList(lex, v)
		lex.next() // consume ')'
		return
	}
	panic(fmt.Sprintf("unexpected token %q", lex.text()))
}
```
S 表达式为两个不同的目的使用标识符：结构体的字段名和指针的 nil 值。read 函数只处理后一种情况。当它遇到 scanner\.Ident 的值为 “nil” 时，通过 reflect\.Zero 函数把 v 设置为其类型的零值。对于其他标识符，则应该产生一个错误（这里则是采用简单粗暴的方法，直接忽略了）。  

还有一个是 readList 函数。一个 '(' 标记代表一个列表的开始，readList 函数可把列表解码为多种类型：map、结构体、切片或者数组，具体类型根据传入待填充变量的类型决定。对于每种类型都会循环解析内容直到遇到匹配的右括号 ')'，这个是用 endList 函数来检测的。  
比较有趣的地方是递归。最简单的例子是处理数组，在遇到 ')' 之前，使用 Index 方法来获得数组的一个元素，再递归调用 read 来填充数据。切片的流程与数组类似，不同之处是先创建每一个元素变量，再填充，最后追加到切片中。  
结构体和map在循环的每一轮中都必须解析一个关于(key value)的子列表。对于结构体，key 是用来定位字段的符号。与数组类似，通过 FieldByName 函数来获得结构体对应字段的变量，再递归调用 read 来填充。对于 map，key 可以是任何类型。与切片类似，先创建新变量，再递归地填充，最后再把新的键值对添加到 map中：
```go
func readList(lex *lexer, v reflect.Value) {
	switch v.Kind() {
	case reflect.Array: // (item ...)
		for i := 0; !endList(lex); i++ {
			read(lex, v.Index(i))
		}

	case reflect.Slice: // (item ...)
		for !endList(lex) {
			item := reflect.New(v.Type().Elem()).Elem()
			read(lex, item)
			v.Set(reflect.Append(v, item))
		}

	case reflect.Struct: // ((name value) ...)
		for !endList(lex) {
			lex.consume('(')
			if lex.token != scanner.Ident {
				panic(fmt.Sprintf("got token %q, want field name", lex.text()))
			}
			name := lex.text()
			lex.next()
			read(lex, v.FieldByName(name))
			lex.consume(')')
		}

	case reflect.Map: // ((key value) ...)
		v.Set(reflect.MakeMap(v.Type()))
		for !endList(lex) {
			lex.consume('(')
			key := reflect.New(v.Type().Key()).Elem()
			read(lex, key)
			value := reflect.New(v.Type().Elem()).Elem()
			read(lex, value)
			v.SetMapIndex(key, value)
			lex.consume(')')
		}

	default:
		panic(fmt.Sprintf("cannot decode list into %v", v.Type()))
	}
}

func endList(lex *lexer) bool {
	switch lex.token {
	case scanner.EOF:
		panic("end of file")
	case ')':
		return true
	}
	return false
}
```

## 封装解析器
最后，把解析器封装成如下所示的一个导出的函数 Unmarshal，隐藏了实现中多个不完美的地方，比如解析过程中遇到错误会崩溃，因此使用了一个延迟调用来从崩溃中恢复，并且返回错误消息：
```go
// Unmarshal 解析 S 表达式数据并且填充到非 nil 指针 out 指向的变量
func Unmarshal(data []byte, out interface{}) (err error) {
	lex := &lexer{scan: scanner.Scanner{Mode: scanner.GoTokens}}
	lex.scan.Init(bytes.NewReader(data))
	lex.next() // 获取第一个标记
	defer func() {
		// 注意: 错误处理不是这篇的重点，简单粗暴的处理了
		if x := recover(); x != nil {
			err = fmt.Errorf("error at %s: %v", lex.scan.Position, x)
		}
	}()
	read(lex, reflect.ValueOf(out).Elem())
	return nil
}
```

一个具备用于生产环境的质量的实现对任何的输入都不应当崩溃，而且应当对每次错误详细报告信息，可能的话，应当包含行号或者偏移量。通过这个示例有助于了解 encoding\/json 这类包的底层机制，以及如何使用反射来填充数据结构。  
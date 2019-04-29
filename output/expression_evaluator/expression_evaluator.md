# 7.9 示例：表达式求值器
本篇将创建简单算术表达式的一个求值器。  

## 定义接口和类型
开始，先确定要使用一个接口 Expr 来代表这种语言的任意一个表达式。暂时没有任何方法，稍后再逐个添加：
```go
// Expr: 算术表达式
type Expr interface{}
```

我们的表达式语言将包括以下符号：
+ 浮点数字面量
+ 二元操作符：加减乘除（\+、\-、\*、\/）
+ 一元操作符：表示正数和负数的 \-x 和 \+x
+ 函数调用：pow(x,y)、sin(x) 和 sqrt(x)
+ 变量：比如 x、pi，自己定义一个变量名称，每次可以提供不用的值

还要有标准的操作符优先级，以及小括号。所有的值都是 float64 类型。  

下面是几个示例表达式：
```
sqrt(A / pi)
pow(x, 3) + pow(y, 3)
(F - 32) * 5 / 9
```

下面5种具体类型代表特定类型的表达式：
+ Var ： 代表变量引用。这个类型是可导出的，至于为什么，后面会讲明
+ literal ： 代表浮点数常量
+ unary ： 代表有一个操作数的操作符表达式，操作数可以是任意的 Expr
+ binary ： 代表有两个操作数的操作符表达式，操作数可以是任意的 Expr
+ call ： 代表函数调用，这里限制它的 fn 字段只能是 pow、sin、sqrt

为了要计算包含变量的表达式，还需要一个上下文（environment）来把变量映射到数值。所有接口和类型的定义如下：
```go
// output/expression_evaluator/eval/ast.go
```
在定义好各种类型后，发现每个类型都需要提供一个 Eval 方法，于是加把这个方法加到接口中，已经添加到上面的代码中了。  
这个包只导出了 Expr、Var、Env。客户端可以在不接触其他表达式类型的情况下使用这个求值器。  

## 定义方法
接下来实现每个类型的 Eval 方法来满足接口。  

Var 的 Eval 方法从上下文中查询结果，如果变量不存在，则会返回0。  
literal 的 Eval 方法直接返回本身的值。  
unbary 的 Eval 方法首先对操作数递归求值，然后应用 op 操作符。  
binary 的 Eval 方法的处理逻辑和 unbary 一样。  
call 方法先对 pow、sin、sqrt 函数的参数求值，再调用 math 包中的对应函数。  
```go
// output/expression_evaluator/eval/eval.go
```

某些方法可能会失败，有些错误会导致 Eval 崩溃，还有些会导致返回不正确的结果。所有这些错误可以在求值之前做检查来发现，所以还需要一个Check方法。不过暂时可以先不管Check方法，而是把 Eval 方法用起来，并通过测试进行验证。  

## Parse函数
要验证 Eval 方法，首先需要得到对象，然后调用对像的 Eval 方法。而对象需要通过解析字符串来获取，这就需要一个 Parse 函数。  

**text\/scanner 包的使用**  
词法分析器 lexer 使用 text\/scanner 包提供的扫描器 Scanner 类型来把输入流分解成一系列的标记（token），包括注释、标识符、字符串字面量和数字字面量。扫描器的 Scan 方法将提前扫描并返回下一个标记（类型为 rune）。大部分标记（比如'('）都只包含单个rune，但 text\/scanner 包也可以支持由多个字符组成的记号。调用 Scan 会返回标记的类型，调用 TokenText 则会返回标记的文本。  
因为每个解析器可能需要多次使用当前的记号，但是 Scan 会一直向前扫描，所以把扫描器封装到一个 lexer 辅助类型中，其中保存了 Scan 最近返回的标记。下面是一个简单的用法示例：
```go
package main

import (
	"fmt"
	"os"
	"strings"
	"text/scanner"
)

type lexer struct {
	scan  scanner.Scanner
	token rune // 当前标记
}

func (lex *lexer) next()        { lex.token = lex.scan.Scan() }
func (lex *lexer) text() string { return lex.scan.TokenText() }

// consume 方法并没有被使用到，包括后面的Pause函数
// 不过这是一个可复用的处理逻辑
func (lex *lexer) consume(want rune) {
	if lex.token != want { // 注意: 错误处理不是这篇的重点，简单粗暴的处理了
		panic(fmt.Sprintf("got %q, want %q", lex.text(), want))
	}
	lex.next()
}

func main() {
	for _, input := range os.Args[1:] {
		lex := new(lexer)
		lex.scan.Init(strings.NewReader(input))
		lex.scan.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats
	
		fmt.Println(input, ":")
		lex.next()
		for lex.token != scanner.EOF {
			fmt.Println("\t", scanner.TokenString(lex.token), lex.text())
			lex.next()
		}
	}
}
```

执行效果如下：
```
PS G:\Steed\Documents\Go\src\localdemo\parse> go run main.go "sqrt(A / pi)" "pow(x, 3) + pow(y, 3)" "(F - 32) * 5 / 9"
sqrt(A / pi) :
         Ident sqrt
         "(" (
         Ident A
         "/" /
         Ident pi
         ")" )
pow(x, 3) + pow(y, 3) :
         Ident pow
         "(" (
         Ident x
         "," ,
         Int 3
         ")" )
         "+" +
         Ident pow
         "(" (
         Ident y
         "," ,
         Int 3
         ")" )
(F - 32) * 5 / 9 :
         "(" (
         Ident F
         "-" -
         Int 32
         ")" )
         "*" *
         Int 5
         "/" /
         Int 9
PS G:\Steed\Documents\Go\src\localdemo\parse>
```

**Parse 函数**  
Parse 函数，递归地将字符串解析为表达式，下面是完整的代码：
```go
// output/expression_evaluator/eval/parse.go
```
整体的逻辑都比较难理解。parseBinary 函数是负责解析二元表达式的，其中包括了对运算符优先级的处理。

## 测试函数
下面的 TestEval 函数用于测试 evaluator 
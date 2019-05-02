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
整体的逻辑都比较难理解。parseBinary 函数是负责解析二元表达式的，其中包括了对运算符优先级的处理（*逻辑比较难懂，自己想不出来，看也没完全看懂，以后有类似的实现或许可以借鉴*）。  

## 测试函数
下面的 TestEval 函数用于测试求值器，它使用 testing 包，使用**基于表**的测试方式。表格中定义了三个表达式并为每个表达式准备了不同的上下文。第一个表达式用于根据圆面积A求半径，第二个用于计算两个变量x和y的立方和，第三个把华氏温度F转为摄氏温度：
```go
// output/expression_evaluator/eval/eval_test.go
```
对于表格中的每一行记录，先解析表达式，在上下文中求值，再输出表达式。启用 \-v 选项查看测试的输出：
```
PS G:\Steed\Documents\Go\src\gopl\output\expression_evaluator\eval> go test -v
=== RUN   TestEval

sqrt(A / pi)
        map[A:87616 pi:3.141592653589793] => 167

pow(x, 3) + pow(y, 3)
        map[x:12 y:1] => 1729
        map[x:9 y:10] => 1729

5 / 9 * (F - 32)
        map[F:-40] => -40
        map[F:32] => 0
        map[F:212] => 100
--- PASS: TestEval (0.00s)
PASS
ok      gopl/output/expression_evaluator/eval   0.329s
PS G:\Steed\Documents\Go\src\gopl\output\expression_evaluator\eval>
```

## check 方法
到目前为止，所有的输入都是合法的，但是并不是总能如此。即使在解释性语言中，通过语法检查来发现**静态**错误（即不用运行程序也能检测出来的错误）也是很常见的。通过分离静态检查和动态检查，可以更快发现错误，也可以只在运行前检查一次，而不用在表达式求值时每次都检查。  
现在就给 Expr 加上一个 Check 方法，用于在表达式语法树上检查静态错误。这个 Check 方法有一个 vars 参数，并不是因为需要传参，而是为了让递归调用的实现起来更方便，具体看后面的代码和说明：
```go
// Expr: 算术表达式
type Expr interface {
	// 返回表达式在 env 上下文下的值
	Eval(env Env) float64
	// Check 方法报告表达式中的错误，并把表达式中的变量加入 Vars 中
	Check(vars map[Var]bool) error
}
```

具体的 Check 方法如下所示。literal 和 Var 的求值不可能出错，所以直接返回 nil。unary 和 binary 的方法首先检查操作符是否合法，再递归地检查操作数。类似地，call 的方法首先检查函数是否是已知的，然后检查参数个数，最后递归地检查每个参数：
```go
// output/expression_evaluator/eval/check.go
```
关于递归的实现。Check 的输入参数是一个 Var 集合，这个集合是表达式中的变量名。要让表达式能成功求值，上下文必须包含所有的变量。从逻辑上来讲，这个集合应当是 Check 的输出结果而不是输入参数，但因为这个方法是递归调用的，在这种情况下使用参数更为方便。调用方最初调用时需要提供一个空的集合。  

## Web 应用
这篇里已经有一个绘制函数 z=f(x,y) 的 SVG 图形的实现了：https://blog.51cto.com/steed/2356431  
不过当时的函数 f 是在编译的时候指定的。既然这里可以对字符串形式的表达式进行解析、检查和求值，那么就可以构建一个 Web 应用，在运行时从客户端接收一个表达式，并绘制函数的曲面图。可以使用 vars 集合来检查表达式是否是一个只有两个变量x、y的函数（为了简单起见，还提供了半径r，所以实际上是3个变量）。使用 Check 方法来拒绝掉不规范的表达式，这样就不会在下面函数的40000个计算过程中（100x100的格子，每一个有4个角）重复这些检查。  
表达式求值器已经完成了，把它作为一个包引入。然后把绘制函数图形加上 Web 应用的代码重新实现一遍，完整的代码如下：
```go
// ch7/suface/main.go
```
重点看 parseAndCheck 函数，组合了解析和检查的步骤。  
还有 plot 函数，函数的签名与 http\.HandlerFunc 类似。解析并检查 HTTP 请求中的表达式，并用它来创建一个有两个变量的匿名函数。这个匿名函数与原始曲面绘制程序中的 f 有同样的签名，并且能对用户提供的表达式进行求值。上下文定义了x、y和半径r。最后，plot 调用了 surface 函数，这里略做了修改，原本直接使用函数 f，现在把函数 f 作为参数传入。  
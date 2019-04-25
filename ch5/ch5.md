# 5.3 多返回值

## 裸返回
一个函数如果有命名的返回值，可以省略 return 语句的操作数，这称为**裸返回**。  
在一个函数中如果存在许多返回语句且有多个返回结果，裸返回可以消除重复代码，但是并不能使代码更加易于理解。比如，对于这种方式，在第一眼看来，不能直观地看出返回的值具体是什么。如果之前一直没有使用过返回值的变量名，返回变量的零值，如果赋过值了，则返回新的值，这就有可能会看漏。鉴于这个原因，应该保守使用裸返回。  

# 5.6 匿名函数

## 图的遍历
在下面的例子中，变量 prereqs 的 map 提供了很多课程（key），以及学习该课程的前置条件（value）：
```go
var prereqs = map[string][]string{
	"algorithems": {"data structures"},
	"calculus":    {"linear algebra"},
	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},
	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}
```
**图**  
这样的问题是一种拓扑排序。概念上，先决条件的内容构成了一张有向图，每一个节点代表一门课程。每一条边代表一门课程所依赖的另一门课程的关系。  
图是无环的：没有节点可以通过图上的路径回到它自己。  

可以使用深度优先的搜索计算得到合法的学习路径，代码入下所示：
```go
// ch5/toposort
```
当一个匿名函数需要进行递归，必须先声明一个变量然后将匿名函数赋给这个变量。如果将两个步骤合并成一个声明，函数字面量将不会存在于该匿名函数的作用域中，这样就不能递归地调用自己了。  
下面是拓扑排序的程序输出，它是确定的结果，就是每次执行都一样。这里输出时调用的是切片而不是 map，所以迭代的顺序是确定的并且在调用最初的 map 之前是对它的 key 进行了排序的。  
```
PS H:\Go\src\gopl\ch5\toposort> go run main.go
1:      intro to programming
2:      discrete math
3:      data structures
4:      algorithems
5:      linear algebra
6:      calculus
7:      formal languages
8:      computer organization
9:      compilers
10:     databases
11:     operating systems
12:     networks
13:     programming languages
PS H:\Go\src\gopl\ch5\toposort>
```

# 警告：捕获迭代变量
首先，看下面的代码：
```go
package main

import "fmt"

func main() {
	var shows []func()
	for _, v := range []int{1, 2, 3, 4, 5} {
		shows = append(shows, func() { fmt.Println(v) })
	}

	for _, f := range shows {
		f()
	}
}
```
这里的期望是依次打印每个数。但实际打印出来的全部都是5。  
在for循环引进的一个块作用域内声明了变量v，然后到了循环里使用的这类变量共享相同的变量，即一个可访问的存储位置，而不是固定的值。v的值在不断地迭代中更新，因此当之后调用打印的时候，v变量已经被每一次的for循环更新多次。所以打印出来的是最后一次迭代时的值。  
这里可以通过引入一个内部变量来解决这个问题，可以换个名字，也可以使用一样的变量名：
```go
func main() {
	var shows []func()
	for _, v := range []int{1, 2, 3, 4, 5} {
		v := v // 这句是关键
		shows = append(shows, func() { fmt.Println(v) })
	}

	for _, f := range shows {
		f()
	}
}
```
看起来奇怪，但却是一个关键性的声明。for循环内也可以随意定义一个不一样的变量名，这样看着更好理解一些。  
也可以用匿名函数（闭包）来理解，这里确实是一个闭包，匿名函数内引用了外部变量。第一个示例中，变量v会在for循环的每次迭戈中更新。第二个示例，匿名函数引用的变量v是在for循环内部声明的，不会随着迭代而更新，并且在for循环内部也没有变化过。  
这样的隐患不仅仅存在于使用range的for循环里。在 `for i := 0; i < 10; i++ {}` 这样的循环里作用域也是同样的，这里的变量i也是会有同样的问题，需要避免。  
另外在go语句和derfer语句的使用当中，迭代变量捕获的问题是最频繁的，这是因为这两个逻辑都会推迟函数的执行时机，直到循环结束。但是这个问题并不是有go或者defer语句造成的。  

**goroutine 中同样的问题**  
下面的用法是错误的：
```go
for _, f := range names {
	go func() {
		call(f) // 注意：不正确
	}
}
```
需要作为一个字面量函数的显式参数传递 f，而不是在 for 循环中声明 f。正确的做法如下：
```go
for _, f := range names {
	go func(f string) {
		call(f)
	}(f) // 显式的传递 f 给函数
}
```
像上面这样，通过添加显式参数，可以确保当 go 语句执行的时候，使用 f 的当前值。  

# 5.8 延迟函数调用
defer 语句也可以用来调试一个复杂的函数，即在函数的“入口”和“出口”处设置调试行为。下面的 bigSlowOperation 函数在开头调用 trace 函数，在函数刚进入的时候执行输出，然后返回一个函数变量，当其被调用的时候执行退出函数的操作。以这种方式推迟返回函数的调用，就可以使一个语句在函数入口和所有出口添加处理，甚至可以传递一些有用的值，比如每个操作的开始时间：
```go
package main

import (
	"log"
	"time"
)

func bigSlowOperation() {
	defer trace("bigSlowOperation")()  // 这个小括号很重要
	// ...这里假设有一些操作...
	time.Sleep(3 * time.Second) // 模拟慢操作
}

func trace(msg string) func() {
	start := time.Now()
	log.Printf("enter %s", msg)
	return func() { log.Printf("exit %s (%s)", msg, time.Since(start)) }
}

func main() {
	bigSlowOperation()
}
```
通常的defer语句提供一个函数，会在函数退出时再调用。  
上面的defer语句，最后面有两个小括号。trace函数调用后会返回一个匿名函数，加上后面的小括号才是延迟调用执行的部分。而trace函数本身则会在当前位置就执行，并且返回匿名函数给defer语句。在trace函数获取返回值的过程中，也就是trace函数里，会先执行两行语句，获取start变量的值以及输出一行信息，这个是在函数开头就执行的。最后函数返回的匿名函数是提供给defer语句在退出的时候进行延迟调用的。  

# 5.9 宕机

## 主动调用 panic
可以直接调用内置的 panic 函数。如果碰到“不可能发生”的状况，panic 是最好的处理方式，比如语句执行到逻辑上不可能到达的地方时。

## 转储栈信息
runtime 包提供了转储栈的方法是程序员可以诊断错误，下面的代码在 main 函数中延迟 printStack 的执行：
```go
// ch5/defer2
```
Panic之后，在退出前会调用 defer 的内容，输出 buf 中的栈信息。最后还会输出宕机消息到标准输出流。  
runtime.Stack 能够输出函数栈信息，在其他语言中，此时函数栈的信息应该已经不存在了。但是 Go 语言的宕机机制让延迟执行的函数在栈清理之前调用。  

# 5.10 恢复
退出程序通常是正常的处理panic异常的方式。但有时需要从异常中恢复，至少可以在程序崩溃前做一些操作。  

## recover函数
将内置的 recover 函数在延迟函数的内部调用，当定义了该 defer 语句的函数发生了 panic 异常，recover 就会终止当前的 panic 状态并且返回 panic value。函数不会从之前 panic 的地方继续运行而是正常返回。在未发生 panic 时调用 recover 则没有任何效果并且返回 nil。  

## 举例说明
假设有一个语言解析器。即使看起来运行正常，但考虑到工作的复杂性，还是会存在只在特殊情况下发生的 bug。此时我们更希望返回一个错误 error 而不是导致程序崩溃 panic。所以 panic 发生后，不要立即终止运行，而是将一些有用的附加消息提供给用户来报告这个bug。下面是使用 recover 部分的代码：
```go
func Parse(input string) (s *Syntax, err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("internal error: %v", p)
		}
	}()
	// ...parser...
}
```

## 恢复的原则
对于 panic 采用无差别的恢复措施是不可靠的。  
从同一个包内发生的 panic 进行恢复有助于简化处理复杂和未知的错误，但一般的原则是，不应该尝试去恢复从另一个包内发生的 panic。公共的 API 应该直接报告错误。同样，也不应该恢复一个 panic，而这段代码却不是由你来维护的，比如调用这提供的回调函数，因为你不清楚这样做是否安全。  
有时也很难完全遵循规范，举个例子，net\/http包中提供了一个web服务器，将收到的请求分发给用户提供的处理函数。很显然，我们不能因为某个处理函数引发的panic异常，影响整个进程导致退出。web服务器遇到处理函数导致的panic时会调用recover，输出堆栈信息，继续运行。这样的做法在实践中很便捷，但也会有一定的风险，比如导致资源泄漏或是因为recover操作，导致其他问题。  
所以，最安全的做法就是选择性地使用 recover。当 panic 之后需要进行恢复的情况本来就不多。为了标识某个 panic 是否应该被恢复，我们可以将 panic value 设置成特殊类型。在 recover 时对 panic value 进行检查，如果发现 panic value 是特殊类型，就将这个 panic 作为 errror 处理。如果不是，则按照正常的 panic 进行处理。
下面示例代码中的 soleTitle 函数就是一个这样的例子：
```go
// ch5/title3
```
defer 调用 recover，检查 panic value，如果该值是 bailout{} 则返回一个普通的错误。所有其他非空的值都是预料外的 panic，这时继续使用 panic value 的值作为参数调用 panic。  

这个示例里，违反了 panic 不处理"预期"错误的建议，但是这里是为了展示这种处理 panic 的机制：
```go
if title != "" {
	panic(bailout{}) // 多个标题元素
}
```
对于一个预期的错误，比如这里标题为空的情况。正常编写程序的时候，不应该调用panic，而是进行处理，比如返回 error。  

有些情况下是没有恢复动作的。比如，内存耗尽会使 Go 运行时发生严重错误而直接终止进程。

## 练习
使用 panic 和 recover 写一个函数，它没有 return 语句，但是能够返回一个非零的值。
```go
// exercise5/e19
```
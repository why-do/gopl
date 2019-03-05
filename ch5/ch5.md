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

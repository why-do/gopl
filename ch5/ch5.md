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

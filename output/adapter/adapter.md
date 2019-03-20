# 适配器
适配器模式是设计模式中的一种：将一个类的接口转换成客户希望的另一个接口。适配器模式使得原本由于接口不兼容而不能一起工作的那些类可以一起工作。  

## OOP的实现
在一般的向对象的语言里，可能是定义一个适配器类来实现的。Go 自然也是可以的。  
比如已经有了一个自己的类，并且实现了一些方法：
```go
type People struct {
	name string
}

func (p *People) Greet() {
	fmt.Printf("Hello, I am %s.\n", p.name)
}
```

现在有一个接口，下面是接口的定义和一个调用该接口的函数：
```go
type Hello interface {
	SayHello()
}

func SayHello(s Hello) {
	s.SayHello()
}
```

问题来了，Peopel并没有实现接口的方法，所以不能直接传递给 SayHello 函数。  
解决方案一，直接为People再写一个方法就解决了。但是如果这个People是别的包里的类型，就不能这么做了。  
解决方案二，为原来的People定义一个别名类型，然后定义新的别名类型的方法来满足接口。实现后大概是下面这样的：
```go
func main() {
	p2 := PeopleSayHello(p1)
	SayHello(&p2)
}

// 适配器
type PeopleSayHello People

func (p *PeopleSayHello) SayHello() {
	p1 := People(*p)
	p1.Greet()
}
```
由于类型转换之后的新的类型是无法取地址的，所以必须要引入一个临时变量。在main函数中多出这么一行也不是很好。

解决方案三，定义一个新的结构体。这个是适配器模式的标准做法，网上也有很多类似的文章：
```go
func main() {
	SayHello(&PeopleSayHello{&p1})
}

// 适配器
type PeopleSayHello struct {
	*People
}

func (p *PeopleSayHello) SayHello() {
	p.Greet()
}
```
这里的方法实现也可以完全自定义。这里只是简单的调用了一个原结构体里的方法，对于这种简单的情况，有更轻量级的方式来适配接口。

## 为函数类型定义方法
Go 语言的接口机制有一些不常见的特性。就是作为一个函数类型，还可以拥有自己的方法。先看下面的实现：
```go
func main() {
	SayHello(PeopleSayHello(p1.Greet))
}

// 适配器
type PeopleSayHello func()

func (f PeopleSayHello) SayHello() {
	f()
}
```
表达式`PeopleSayHello(p1.Greet)`并不是一个并不是一个函数调用，而是一个类型转换。`p1.Greet`是一个方法值，就是如下类型的一个值：
```go
func SayHello()
```
*可以简单的在这里就把接口值就当做是一个函数值。*  
新的类型中定义了 SayHello 方法，现在满足接口了，而这个方法就是调用函数本身，所以这里的 PeopleSayHello 就是一个让函数值满足接口的一个适配器。  
这里用到的是方法值，其实本质上就是为一个函数值定义方法来适配接口。所以，对于方法表达式或函数，也是同样适用的。  
在标准库的 net\/http 中也有类似的用法，就是下面的 HandlerFunc 类型：
```go
package http
type HandlerFunc func(ResponseWriter, *Request)
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r)
}
```

## 完整的示例代码
最后一个例子的完整的示例代码。其他几个例子也只是 main 函数中的调用，和适配器实现部分的代码不同：
```go
// output/adapter/adapter3
```
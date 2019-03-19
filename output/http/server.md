# 1.7 一个 Web 服务器
使用 Go 的库非常容易实现一个 Web 服务器。

## 请求的 URL 路径
这是一个迷你服务器，返回访问服务器的 URL 的路径部分。例如，如果请求的 URL 是 `http://localhost:8000/hello`，响应将是 `URL.Path= "/hello"`。  
下面是完整程序的程序：
```go
// ch1/server1
```
请求的 URL 的路径就是 `r.URL.Path`。  

## 多个处理函数
为服务器添加功能很容易。一个有用的扩展是一个特定的 URL，下面的版本对 \/count 请求会有特殊的响应：
```go
// ch1/server2
```
这个服务器有两个处理函数，通过请求的 URL 来决定哪一个被调用。

## 请求头和表单信息
下面这个示例中的处理函数，报告它接收到的请求头和表单数据，这样还方便服务器审查和调试请求：
```go
// ch1/server3
```
这里汇报了很多的内容：
+ 请求方法 ： r.Method
+ 请求路径 ： r.URL，这里就是 r.URL.Path。r.URL是个结构体，这里应该只有 Path 字段有内容。然后 %s 是调用它的 String 方法输出
+ 请求协议 ： r.Proto
+ 请求头 ： r.Header，这是个 map，这里一项一项输出了
+ 服务端地址 ： r.Host，包括主机名和端口号
+ 客户端地址 ： r.RemoteAddr，包括主机名和端口号
+ 表单信息 ： r.Form，这个先要用 r.ParseForm() 进行解析后才会有内容。包括 Get 请求和 Post 请求的信息都会在 r.Form 这个 map 里。

# 7.7 http.Handler 接口
进一步了解基于 http.Handler 接口的服务器API。

## 接口
下面是源码中接口的定义：
```go
package http

type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}

func ListenAndServe(addr string, handler Handler) error {
	server := &Server{Addr: addr, Handler: handler}
	return server.ListenAndServe()
}
```
ListenAndServe 函数，这里关注接口，只看函数的签名，忽略函数体的内容。函数的第二个参数接收一个 Handler 接口的实例（用来接受所有的请求）。这个函数会一直执行，直到服务出错时返回一个非空的错误值。  

## 简单的示例
下面的程序展示一个简单的例子。使用map类型的database变量记录商品和价格的映射。再加上一个 ServeHTTP 方法来满足 http.Handler 接口。这个函数遍历整个 map 并且输出其中的元素：
```go
// ch7/http1
```

## 添加功能
上面的示例中，服务器只能列出所有的商品，并且完全不管 URL，对每个请求都是同样的功能。一般的 Web 服务会定义过个不同的 URL，每个触发不同的行为。把现有的功能的 URL 设置为 \/list，再加上另一个 \/price 用来显示单个商品的价格，商品可以在请求参数中指定，比如：`/price?item=socks`：
```go
// ch7/http2
```
现在，出差函数基于 URL 的路径部分（req.URL.Path）来决定执行哪部分逻辑。  

**返回错误页面 404**  
如果处理函数不能识别这个路径，那么它通过调用`w.WriteHeader(http.StatusNotFound)`来返回一个 HTTP 错误。这个调用必须在网 w 中写入内容之前执行。这里还可以使用 http.Error 这个工具函数了达到同样的目的：
```go
msg := fmt.Sprintf("no such item: %q\n", item)
http.Error(w, msg, http.StatusNotFound) // 404
```

**Get请求参数**  
对应 \/price 的场景，它调用了 URL 的 Query 方法，把 HTTP 的请求参数解析为一个map，或者更精确来讲，解析为一个 multimap，由 net\/url 包的 url\.Values 类型实现。这里的 url\.Values 是一个 map 映射：
```go
type Values map[string][]string
```
它的 value 是一个 字符串切片，这里用了 Get 方法，只会提取切片的第一个值。如果是要提取某个 key 所有的值，简单的通过 map 的 key 提取 value 应该就好了。  

## 优化功能添加
如果要继续给 ServeHTTP 方法添加功能，应当把每部分逻辑分到独立的函数或方法。net\/http 包提供了一个**请求多工转发器 ServeMux**，用来简化 URL 和处理程序之间的关联。一个 ServeMux 把多个 http\.Handler 组合成单个 http.Handler。在这里，可以看到满足同一个接口的多个类型是可以互相替代的，Web 服务器可以把请求分发到任意一个 http\.Handlr，而不用管后面具体的类型。  
对于更加复杂的应用，多个 ServeMux 会组合起来，用来处理更复杂的分发需求。Go 语言并不需要一个类似于 Python 的 Django 那样的权威 Web 框架。因为 Go 语言的标准库提供的基础单元足够灵活，以至于那样的框架通常不是必须的。进一步来了讲，尽管框架在项目初期带来很多便利，但框架带来了额外复杂性，增加长时间维护的难度。*不过这样的Web框架也是有的，比如：beego。*
将程序修改为使用 ServeMux，用于将 \/list、\/prics 这样的 URL 和对应的处理程序关联起来，这些处理程序也已经拆分到不同的方法中。最后作为主处理程序在 ListenAndServe 调用中使用这个 ServeMux：
```go
// ch7/http3
```

**注册处理程序**
先关注一下用于注册程序的两次 mux.Handle 调用。在第一个调用中，db.list是一个方法值，即如下类型的一个值：
```go
func(w http.ResponseWriter, req *http.Request)
```
当调用 db.list 时，等价于以 db 为接收者调用 database.list 方法。所以 db.list 是一个实现了处理功能的函数。然而他没有接口所需的方法，所以它不满足 http.Handler 接口，也不能直接传给 mux.Handle。  
表达式`http.HandlerFunc(db.list)`其实是一个类型转换，而不是函数调用。注意，http.HandlerFunc 是一个类型，它有如下定义：
```go
package http
type HandlerFunc func(ResponseWriter, *Request)
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r)
}
```
http.HandlerFunc 这个函数类型它有自己的 ServeHTTP 方法，因此它满足接口。而 http.HandlerFunc 的函数签名和 db.list 这个方法值的函数签名是一样的，因此也能够进行类型转换。  
这个是 Go 语言接口机制的一个不常见的特性。它不仅是一个函数类型，还可以拥有自己的方法，它的 ServeHTTP 方法就是调用函数本身，所以 HandlerFunc 是一个让函数值满足接口的一个适配器。在这个例子里，函数和接口的唯一方法拥有同样的签名。这个小技巧让 database 类型可以用不同的方式来满足 http.Handler 接口，一次通过 list 方法，一次通过 price 方法。  

## 简化注册处理
因为这种注册处理程序的方法太常见了，所以 ServeMux 引入了一个 HandleFunc 便捷方法来简化调用，处理程序注册部分的代码可以简化为如下的形式：
```go
// mux.Handle("/list", http.HandlerFunc(db.list))
mux.HandleFunc("/list", db.list)
// mux.Handle("/prics", http.HandlerFunc(db.price))
mux.HandleFunc("/price", db.price)
```

**全局 ServeMux 实例**  
通过 ServeMux，如果需要有两个不同的 Web 服务，在不同的端口监听。那么就定义不同的 URL，分发到不同的处理程序。只须简单地构造两个 ServeMux，再调用一次 ListenAndServe 即可（*建议并发调用*）。不过很多时候一个 Web 服务足够了，另外也不需要多个 ServeMux 实例。对于这种简单的应用场景，建议用下面的简化的调用方法。  
net\/http 包还提供了一个全局的 ServeMux 实例 DefaultServeMux，以及包级别的注册函数 http.Handle 和 http.HandleFunc。要让 DefaultServeMux 作为服务器的主处理程序，无须把它传给 ListenAndServe，直接传nil即可。文章开头的例子里就是这么用的。
服务器的主函数可以进一步简化：
```go
func main() {
	db := database{"shoes": 50, "socks": 5}
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
```

## 并发安全问题
Web 服务器每次都用一个新的 goroutine 来调用处理程序，所以处理程序必须要注意并发问题。比如在访问变量时的锁问题，这个变量可能会被其他 goroutine 访问，包括由同一个处理程序出厂的其他请求。文章开头的第二个例子就要类似的处理。  
并发安全是另外一块内容，需要单独研究和解决，这里去简单提一下。如果要添加创建、更新商品的功能，就需要注意并发安全。  
TODO: 完成练习 7.11 后 把新功能添加进来。

# 部署

## 反向代理
Go 语言原生支持 http，所有 Go 的http服务性能和nginx比较接近。如果用 Go 写的 Web 程序上线，程序前面不需要再部署nginx的Web服务器，这样就省掉的是Web服务器。这是单应用的部署。  
对于多应用部署，服务器需要部署多个Web应用，这时就需要反向代理了，一般这也是nginx或apache。  
**反向代理**，有个很棒的说法是流量转发。我获取到客户端来的请求，将它发往另一个服务器，从服务器获取到响应再回给原先的客户端。**反向**的意义简单来说在于这个代理自身决定了何时将流量发往何处。  
Go 的反向代理，可以参考下这篇：[1 行 Go 代码实现反向代理 ](https://studygolang.com/articles/14246)  

## Panic 处理
https://blog.51cto.com/steed/2321827
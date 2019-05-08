第9章之后

# 9.6 竞态检测器
即使再仔细的检查，仍然可能在并发上犯错。Go 的 runtime 提供了动态分析工具：**竞态检测器**（race detectotr）。  
在下一节的示例中会用到竞态检测器，所以在用之前，先了解一下这个工具。  

**开启竞态检测器**  
简单地把 -race 命令行参数加到 go build、go run、go test 命令里即可使用该功能。它会让编译器为你的应用或测试构建一个修改后的版本，这个版本有额外的手法可以高效记录在执行时对共享变量的所有访问，以及读写这些变量的 goroutine 标识。除此之外，还会记录所有的同步事件、包括 go 语句、通道操作、锁的调用等。（完整的同步事件集合可以在语言规范中的 “The Go Memory Model” 文档中找到。）  

**如何检查到竞态**  
竞态检测器会研究事件流，找到那些有问题的案例，即一个 goroutine 写入一个变量后，中间没有任何同步的操作，就有另外一个 goroutine 读写了该变量。这种情况表明有对共享变量的并发访问，即数据竞态。工具会输出一份报告，包括变量的标识以及读写 goroutine 当时的调用栈。通常情况下这些信息足以定位问题了，下一章的示例会应用到实战中。  

**哪些竞态可能查不到**  
竞态检测器报告所有实际运行了的数据竞态。但只能检测到那些在运行时发生的竞态，无法用来保证肯定不发生竞态。所以为了保证效果，需要全包测试包含了并发调用的场景。  

**可以在生产环境开启竞态检测器**  
由于存在额外的记录工作，带竞态检测功能的程序在执行时需要更长的时间和更多的内存，但即使对于生成环境的任务，这种额外开支也是可以接受的。对于那些偶发的竞态条件，使用竞态检测器可以节省很多调试的时间。  

# 9.7 示例：并发非阻塞缓存
创建一个**并发非阻塞的缓存**系统，它能解决**函数记忆**（memoizing）的问题，即缓存函数的结果，达到多次调用但只须计算一次结果。这个问题在并发实战中很常见但已有的库不能很好地解决这个问题。这里的解决方案将会是并发安全的，并且要避免简单地对整个缓存使用单个锁而带来的锁争夺问题。  

## 被缓存结果的函数
在做系统之前，先准备一个将要被测试的函数。这里将使用下面的 httpGetBody 函数作为示例来演示函数记忆。调用 HTTP 请求相当昂贵，所以我希望只在第一次请求的时候去发起请求，而之后都可以在缓存中找到结果直接返回：
```go
func httpGetBody(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
```
先保证能缓存这个函数的执行结果，之后再使用更多个函数来测试和验证功能。  

## 串行的版本
这是一个并发不安全的版本，不过把基本功能先实现，并发安全的问题之后再进行优化：
```go
// out/memo/memo1/memo.go
```
Memo 实例包含了被记忆的函数 f （类型为Func），以及缓存，类型为一个 key 为字符串，value 为 result 的 map。每个 result 都是调用 f 产生的结果：一个值和一个错误，在设计的推进过程中会展示 Memo 的几种变体，但所有变体都会遵守这些基本概念。  
下面的例子展示了如何使用 Memo。下面是完整的测试源码文件，包括上一小节写的被测试的函数，以及一串 URL。每个 URL 会发起两次请求。对于每个 URL，首先调用 Get，打印延时和返回的数据长度：
```go
// out/memo/memo1/memo.go
```
这里使用 testing 包系统的测试效果。上面有两个测试函数，先只用 TestSequential 进行测试，串行的发起请求。从下面的测试结果看，每一个 URL 第一次调用都会消耗一定的时间，但对 URL 第二次的请求会立刻返回结果：
```
PS H:\Go\src\gopl\output\memo\memo1> go test -run=TestSequential -v
=== RUN   TestSequential
http://docscn.studygolang.com/, 87.1978ms, 6612 bytes
https://studygolang.com/, 203.3312ms, 81819 bytes
https://studygolang.com/pkgdoc, 33.0053ms, 1261 bytes
https://github.com/adonovan/gopl.io/tree/master/ch9, 1.4428937s, 61185 bytes
http://docscn.studygolang.com/, 0s, 6612 bytes
https://studygolang.com/, 0s, 81819 bytes
https://studygolang.com/pkgdoc, 0s, 1261 bytes
https://github.com/adonovan/gopl.io/tree/master/ch9, 0s, 61185 bytes
--- PASS: TestSequential (1.81s)
PASS
ok      gopl/output/memo/memo1  2.063s
PS H:\Go\src\gopl\output\memo\memo1>
```
默认在测试成功的时候不打印这类日志，不过可以加上 -v 参数在成功时也打印测试日志。  
这次测试中所有的 Get 都是串行的。因为 HTTP 请求通过并发来改善的空间很大，所以这次使用 TestConcurrent 进行测试，让所有的请求并发进行。这个测试要使用 sync\.WaitGroup 等待所有的请求完成后再返回结果。  
这次的测试结果基本上都是缓存无效的情况，不过偶尔还会出现无法正常运行的情况。除了缓存无效，可能还会有缓存命中后返回错误结果，甚至崩溃：
```
PS H:\Go\src\gopl\output\memo\memo1> go test -run=TestConcurrent -v
=== RUN   TestConcurrent
http://docscn.studygolang.com/, 92.9972ms, 6612 bytes
http://docscn.studygolang.com/, 98.9889ms, 6612 bytes
https://studygolang.com/pkgdoc, 204.8383ms, 1261 bytes
https://studygolang.com/pkgdoc, 205.8387ms, 1261 bytes
https://studygolang.com/, 234.1566ms, 81819 bytes
https://studygolang.com/, 235.1749ms, 81819 bytes
https://github.com/adonovan/gopl.io/tree/master/ch9, 1.5041445s, 61184 bytes
https://github.com/adonovan/gopl.io/tree/master/ch9, 2.1051443s, 61184 bytes
--- PASS: TestConcurrent (2.11s)
PASS
ok      gopl/output/memo/memo1  2.346s
PS H:\Go\src\gopl\output\memo\memo1>
```
更糟糕的是，多数时候这样都能正常运行，所以甚至很难注意到这样并发调用是有问题的。但是如果加上 -race 标志后再运行，那么竞态检测器就会输出如下的报告：
```
PS H:\Go\src\gopl\output\memo\memo1> go test -run=TestConcurrent -v -race
=== RUN   TestConcurrent
==================
WARNING: DATA RACE
Write at 0x00c000062cf0 by goroutine 11:
  runtime.mapassign_faststr()
      D:/Go/src/runtime/map_faststr.go:190 +0x0
  gopl/output/memo/memo1.(*Memo).Get()
      H:/Go/src/gopl/output/memo/memo1/memo.go:27 +0x1d8
  gopl/output/memo/memo1.TestConcurrent.func1()
      H:/Go/src/gopl/output/memo/memo1/memo_test.go:57 +0xc0

Previous write at 0x00c000062cf0 by goroutine 7:
  runtime.mapassign_faststr()
      D:/Go/src/runtime/map_faststr.go:190 +0x0
  gopl/output/memo/memo1.(*Memo).Get()
      H:/Go/src/gopl/output/memo/memo1/memo.go:27 +0x1d8
  gopl/output/memo/memo1.TestConcurrent.func1()
      H:/Go/src/gopl/output/memo/memo1/memo_test.go:57 +0xc0
...
FAIL    gopl/output/memo/memo1  2.883s
```
这里就是因为两个 goroutine 在没使用同步的情况下更新了 Memo.cache 这个 map。因为整个 Get 并不是并发安全的，它存在数据竞态：
```go
// 注意：并发不安全
func (memo *Memo) Get(key string) (interface{}, error) {
	res, ok := memo.cache[key]
	if !ok {
		res.value, res.err = memo.f(key)
		memo.cache[key] = res
	}
	return res.value, res.err
}
```
所以，接下来就是要进行，实现并发安全。

## 使用互斥锁
让缓存并发安全最简单的方法就是用一个基于监控的同步机制。需要给 Memo 加一个互斥量，并在 Get 开始就获取互斥锁，在返回前释放互斥锁，这样就可以让 cache 相关的操作发生在临界区域内了：
```go
// Memo 缓存了调用 Func 的结果
type Memo struct {
	f     Func
	mu    sync.Mutex // 保护 cache
	cache map[string]result
}

// Get 是并发安全的
func (memo *Memo) Get(key string) (interface{}, error) {
	memo.mu.Lock()
	res, ok := memo.cache[key]
	if !ok {
		res.value, res.err = memo.f(key)
		memo.cache[key] = res
	}
	memo.mu.Unlock()
	return res.value, res.err
}
```
加上锁之后，再运行并发测试函数，竞态检测器不报警了。但是这次的修改后，之前对性能的优化就失效了。由于每次调用 Memo\.f 时都上锁，所以现在的 Get 方法运行的使用实际又是串行的了。这里需要一个**非阻塞的**缓存，一个不会把他需要记忆的函数串行运行的缓存。  
调用 Get 是不需要锁保护的。调用 Get 的判断依据是之前的获取 map 的 key，这个操作需要加锁。调用 Get 返回后，需要把返回结果更新到 map 中去，这个操作也需要加锁。在 map 查询结束后，先释放锁。不加锁的情况下调用 Get。等到结果返回需要更新 map 的时候，再加锁更新 map。具体修改如下：
```go
func (memo *Memo) Get(key string) (interface{}, error) {
	memo.mu.Lock()
	res, ok := memo.cache[key]
	memo.mu.Unlock()
	if !ok {
		res.value, res.err = memo.f(key)
		memo.mu.Lock()
		memo.cache[key] = res
		memo.mu.Unlock()
	}
	return res.value, res.err
}
```
现在，可以安全的并行运行了，但是缓存又失效了。某些URL被获取了两次。修改一下测试源码文件的被测试函数 httpGetBody，在开头输出一行日志，可以观察到每个URL被调用的次数：
```go
func httpGetBody(url string) (interface{}, error) {
	log.Printf("httpGetBody: %s", url) // 输出哪些 url 被函数调用了，从缓存获取结果时不会有这个输出
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
```
修改之后，可以用最初的串行版本再测试一下。那个版本是一定用到缓存的效果的。而现在的版本，在并发的情况下无法用上缓存。  
在几个 goroutine 几乎同时调用的 Get 来获取同一个 URL 时，每个 goroutine 都首先查询缓存，发现缓存中没有需要的数据，然后就都去执行 Get 来获取结果，最后又都用获得的结果来更新 map，其中一个结果会被另外一个覆盖。  
在理想的情况下，应该要避免这种额外的处理。这个功能有时称为**重复抑制**（duplicate suppression）。  

## 重复抑制
下面这个版本，map 的每个元素是一个指向 entry 结构的指针。除了与之前一样包含一个已经记住的函数 f 调用结果之外，每个 entry 还新加一个通道 ready。在设置了 entry 和 result 字段后，通道会关闭，正在等待的 goroutine 会收到广播，然后就可以从 entry 字段读取结果了：
```go
// out/memo/memo4/memo.go
```
关于这里的 map 是否包含某个元素的判断，之前都是返回两个值，通过ok来判断。之前的示例中，map的元素是结构体，由于结构体类型的零值不是nil，通过ok来判断比较好。这里的元素类型是结构体指针，当然可以继续使用ok来判断。不同现在是指针类型了，零值是nil也不会和非零值的情况搞混，所以也可以直接通过nil来判断。  
现在调用 Get 会获取锁，然后去 map 中查询，如果没有找到，就直接分配并插入一个新的值，然后释放锁。之后其他 goroutine 来查询的时候，会发现值存在，那么就直接获取到 map 的值，然后释放锁。  
map 里的值并不是 Get 返回的数据，而是数据是否准备好的通道，和存放数据的字段。此时数据可能还没准备好，数据是否准备好，可以从 ready 通道进行判断。对 ready 通道的读取操作，会在数据没有准备好的时候一直阻塞。一旦数据准备好了，就会关闭 ready 通道，所有从 ready 通道的读取操作就会立刻返回。这是利用通道进行广播的方式。所以查询 map 后获取值的步骤就是先读取 ready 通道等待，一旦通道的读取返回，就表示数据已经准备好了，此时就可以去读取字段 res 里的内容并返回。  
注意，entry 中的变量 e\.res\.value 和 e\.res\.err 被多个 goroutine 共享。创建 entry 的 goroutine 会对这两个变量的值进行设置，其他 goroutine 在收到数据准备完毕的广播后才会开始读取这两个变量。尽管被多个 goroutine 访问，但是此处不需要加锁。ready 通道的关闭先于其他 goroutine 收到广播事件，所以第一个 goroutine 对变量的写入也先于后续多个 goroutine 的读取事件。这种情况下数据竞态不存在。  
到此，并发、重复抑制、非阻塞缓存就完成了。  

## 使用监控goroutine
上面的示例是使用一个互斥量来保护 map 变量的并发安全。下面是另一种设计，让 map 变量限制在一个**监控** goroutine 中。  
首先是类型声明，New 函数在创建实例并返回的同时，还会启动一个 server 方法。该方法会集中处理所有的 Get 调用。我们在获取实例后，依然是调用 Get 来获取结果：
```go
// memo包提供了一个对类型 Func 并发安全的函数记忆功能
// 并发、重复抑制、非阻塞的缓存
// 通过监控 goroutine 来实现并发安全
package memo

// Func 是用于记忆的函数类型
type Func func(key string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

type entry struct {
	res   result
	ready chan struct{} // res 准备好之后会被关闭
}

// Func、result、entry 的声明和之前一致

// request 是一条请求消息
type request struct {
	key      string        // 需要 Func 运行的参数
	response chan<- result // 每个客户端接收结果的通道
}

type Memo struct{ requests chan request }

func New(f Func) *Memo {
	memo := &Memo{requests: make(chan request)} // 创建实例
	go memo.server(f)                           // 启动服务端 goroutine
	return memo                                 // 返回实例，供客户端调用
}
```
可以先往后看客户端和服务端的处理逻辑，在回过来看这里声明的数据类型已经通道的作用。  

**客户端**  
现在 Get 就需要要给监控 goroutine 的通道发送请求和一个接收返回结果的通道。服务端会在收到处理请求后进行处理，之后再通过客户端发来的通道返回结果。而客户端发送请求之后，只需要从自己创建的这个通道中接收，直到接收到数据后，再返回即可：
```go
func (memo *Memo) Get(key string) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{key, response}
	res := <- response
	return res.value, res.err
}
func (memo *Memo) Close() { close(memo.requests) }
```
客户端使用完之后，可以调用 Close 方法关闭发送请求的通道。  

**服务端**  
上面的 Get 相当于一个客户端，还需要一个服务端来处理 Get 发来的请求：
```go
func (memo *Memo) server(f Func) {
	cache := make(map[string]*entry)
	for req := range memo.requests { // 一次处理收到的请求
		e := cache[req.key]
		if e == nil {
			// 对这个 key 的第一次请求
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			go e.call(f, req.key) // 调用 f(key)
		}
		// 无论是否第一次请求，最后要回复结果，都有等待 ready 通道返回后，再去读取结果
		go e.deliver(req.response)
	}
}

func (e *entry) call(f Func, key string) {
	// 执行函数
	e.res.value, e.res.err = f(key)
	// 发送广播通知，数据已经准备好了
	close(e.ready)
}

func (e *entry) deliver(response chan<- result) {
	// 等待数据准备完毕
	<-e.ready
	// 向客户端发送结果
	response <- e.res
}
```
变量 cache 被限制在监控 goroutine 中，就是上面的 server 方法。监控 goroutine 从 requests 的通道中读取请求，直到这个通道被关闭。对于每个请求，先查询缓存，如果没有找到就插入一个新的 entry。  
这里 call 和 deliver 方法需要在独立的 goroutine 中运行，以确保监控 goroutine 内持续处理新请求。  

**完整示例代码**  
下面贴上这个实现方式的完整代码以及测试源码：
```go
// out/memo/memo5/memo.go
// out/memo/memo5/memo_test.go
```

## 小结
这里的例子展示了可以使用两种方案来构建并发结构：
+ 共享变量并上锁
+ 通信顺序进程（communicating sequential process）

第一种是大家普遍认知的，也是Java或者C++等语言中的多线程开发。  
第二种是 Go 语言特有的，也是 Go 语言推荐的。下面是一句推荐的原话：
>Do not communicate by sharing memory; instead, share memory by communicating.  
Go 箴言：“不要通过共享内存来通信，而应该通过通信来共享内存”。

在给定的情况下也许很难判定哪种方案更好，不过了解他们还是有价值的。有时候从一种方案切换到另外一种方案能让代码更简单。  

**CSP并发模型**  
CSP 是 Communicating Sequential Process 的简称，中文可以叫做通信顺序进程，是一种并发编程模型。  
CSP 模型由并发执行的实体（线程或者进程）所组成，实体之间通过发送消息进行通信，这里发送消息时使用的就是通道（channel）。CSP 模型的关键是关注 channel，而不关注发送消息的实体。Go 语言就是借用 CSP 模型的一些概念为之实现并发进行理论支持。Go 语言并没有完全实现 CSP 模型的所有理论，仅仅是借用了 process 和 channel 这两个概念。process 在 Go 语言上的表现就是 goroutine 是实际并发执行的实体，每个实体之间通过 channel 通讯来实现数据共享。  
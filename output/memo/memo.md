# 9.6 竞态检测器
即使在自己的检查，仍然可能在并发上犯错。Go 的 runtime 提供了动态分析工具：**竞态检测器**（race detectotr）。  
在下一节的示例中会用到竞态检测器，所以在用之前，先了解一下这个工具。  

**开启竞态检测器**  
简单地把 -race 命令行参数加到 go build、go run、go test 命令里即可使用该功能。它会让编译器为你的应用或测试构建一个修改后的版本，这个版本有额外的手法可以高效记录在执行时对共享变量的所有访问，以及读写这些变量的 goroutine 标识。除此之外，还会记录所有的同步事件、包括 go 语句、通道操作、锁的调用等。（完整的同步事件集合可以在语言规范中的 “The Go Memory Model” 文档中找到。）  

**如何检查到竞态**  
竞态检测器会研究事件流，找到那些有问题的案例，即一个 goroutine 写入一个变量后，中间没有任何同步的操作，就有另外一个 goroutine 读写了该变量。这种情况表明有对共享变量的并发访问，即数据竞态。工具会输出一份报告，包括变量的标识以及读写 goroutine 当时的调用栈。通常情况下这些信息足以定位问题了。  

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

## 初始版本
这是一个并发不安全的版本，不过基本功能先写出来：
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
// momo2
```
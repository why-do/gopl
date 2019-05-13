第8章之后

# 8.2 示例：并发时钟服务器
本节介绍 net 包，它提供构建客户端和服务器程序的组件，这些程序通过 TCP、UDP 或者 UNIX 套接字进行通信。网络服务 net\/http 包是在 net 包的基础上构建的。  

## 时钟服务器
这个示例是一个时钟服务器，它以每秒一次的频率向客户端发送当前时间：
```go
// ch8/clock1
```
Listen 函数创建一个 net\.Listener 对象，它在一个网络端口上监听进来的连接，这里是 TCP 端口 localhost:8000。监听器的 Accept 方法被阻塞，知道有连接请求进来，然后返回 net\.Conn 对象来代表一个连接。  
handleConn 函数处理一个完整的客户端连接。在循环里，它将 time.Now() 获取的当前时间发送给客户端。因为 net\.Conn 满足 io\.Writer 接口，所以可以直接向它进行写入。当写入失败时循环结束，很多时候是客户端断开连接，这是 handleConn 函数使用延迟（defer）的 Close 调用关闭自己这边的连接，然后继续等待下一个连接请求。  
为了连接到服务器，还需要一个 socket 客户端，这里可以先使用系统的 telnet 来进行验证：
```
$ telnet localhost 8000
```
这里可以开两个 telnet 尝试进行连接，只有第一个可以连接上，而其他的连接会阻塞。当把第一个客户端的连接断开后，服务端会重新返回到 main 函数的 for 循环中等待新的连接。此时之前阻塞的一个连接就能连接进来，继续显示时间。服务端程序暂时先这样，先来实现一个 socket 客户端程序。  

## 客户端 netcat
下面的客户端使用 net\.Dial 实现了 Go 版本的 netcat 程序，用来连接 TCP服务器：
```go
// ch8/netcat1
```
这个程序从网络连接中读取，然后写到标准输出，直到到达 EOF 或者出错。  

## 支持并发的服务器
如果打开多个客户端，同时只有一个客户端能正常工作。第二个客户端必须等到第一个结束才能正常工作，这是因为服务器是**顺序**的，一次只能处理一个客户请求。让服务器支持并发只需要一个很小的改变：在调用 handleConn 的地方添加一个 **go** 关键字，使它在自己的 goroutine 内执行：
```go
for {
	conn, err := listener.Accept()
	if err != nil {
		log.Print(err) // 例如，连接终止
		continue
	}
	go handleConn(conn) // 并发处理连接
}
```
现在的版本，多个客户端可以同时接入并正常工作了。  

# 8.3 示例：并发回声服务器
上面的时钟服务器每个连接使用一个 goroutine。下面要实现的这个回声服务器，每个连接使用多个 goroutine 来处理。大多数的回声服务器仅仅将读到的内容写回去，所以可以使用下面简单的 handleConn 版本：
```go
func handleConn(c net.Conn) {
	io.Copy(c, c)
	c.Close()
}
```

## 有趣的回声服务端
下面的这个版本可以重复3次，第一个全大写，第二次正常，第三次全消息：
```go
// reverb1
// ch8/reverb1
```
在上一个示例中，已经知道需要使用 go 关键字调用 handleConn 函数。不过在这个例子中，重点不是处理多个客户端的连接，所以这里不是重点。


## 升级客户端
现在来升级一下客户端，使它可以在终端上向服务器输入，还可以将服务器的回复复制到输出，这里提供了另一个使用并发的机会：
```go
// ch8/netcat2
```

## 优化服务端
使用上面的服务端版本，如果有多个连续的输入，新输入的内容不会马上返回，而是要等待之前输入的内容全部返回后才会处理之后的内容。要想做的更好，需要更多的 goroutine。再一次，在调用 echo 时需要加入 go 关键字：
```go
// reverb2
func handleConn(c net.Conn) {
	input := bufio.NewScanner(c)
	for input.Scan() {
		go echo(c, input.Text(), 1*time.Second)
	}
	// 注意：忽略 input.Err() 中可能的错误
	c.Close()
}
```
这个改进的版本，回声也是并发的，在时间上互相重合。  

## 小结
这就是使服务器变成并发所要做的，不仅处理来自多个客户端的链接，还包括在一个连接处理中，使用多个 go 关键字。在这个例子里，单个客户端连接也可以同时发起多个请求。在最初的版本里，没有使用 go 调用 echo，所以处理单个客户端的请求不是并发的，只有前一个处理完才会继续处理下一个。之后改进的版本，使用 go 调用 echo，这里对每一个请求的处理都是并发的了。  
然而，在添加这些 go 关键字的同时，必须要仔细考虑方法 net\.Conn 的并发调用是不是安全的，对大多数类型来讲，这都是不安全的。  

# 接收完回声再结束
之前的客户端在主 goroutine 中将输入复制到服务器中，这样的客户端在输入接收后立即退出，即使后台的 goroutine 还在继续。为了让程序等待后台的 goroutine 在完成后再退出，使用一个通道来同步两个 goroutine：
```go
func main() {
	conn, err := net.Dial("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn) // 注意：忽略错误
		log.Println("done")
		done <- struct{}{} // 通知主 goroutine 的信号
	}()
	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done // 等待后台 goroutine 完成
}
```
当用户关闭标准输入流（Windows系统使用Ctrl+Z）时，mustCopy 返回，主 goroutine 调用 conn.Close() 来关闭两端网络连接。关闭写半边的连接会导致服务器看到 EOF。关闭读半边的连接导致后台 goroutine 调用 io.Copy 返回 “read from closed connection” 错误，所以这个版本里去掉了打印错误日志。  

## 客户端优化
上面这个版本使用起来的效果和之前的版本并没有太大的差别，几乎看不到差别。虽然多了等待连接关闭，但是依然不会等待接收完毕所有服务器的返回。不过这步解决了等待 goroutine 运行完毕后，主 goroutine 才会结束。使用下面的 TCP 链接，就可以实现接收完毕所有信息后，goroutine 才会结束。在 net 包中，conn 接口有一个具体的类型 \*net\.TCPConn，它代表一个 TCP 连接：
```go
tcpAddr, err := net.ResolveTCPAddr("tcp", ":8000")
if err != nil {
	log.Fatal(err)
}
conn, err := net.DialTCP("tcp", nil, tcpAddr)
if err != nil {
	log.Fatal(err)
}
```
TCP 链接由两半边组成，可以通过 CloseRead 和 CloseWrite 方法分别关闭。修改主 goroutine，仅仅关闭连接的写半边，这样程序可以继续执行输出来自 reverb1 服务器的回声，即使标准输入已经关闭：
```go
// exercise8/e3
```
现在只对第一个回声服务器版本 reverb1 有效，对于之后改进的可以并发处理同一个客户端多个请求的 reverb2 服务器，服务端还需要做一些修改。  

## 服务端优化
在 reverb2 服务器的版本中，因为对于每一个连接，每一次回声的请求都会生成一个新的 goroutine 进行处理。为了知道什么时候最后一个 goroutine 结束（有时候不一定是最后启动的那个），需要在每一个 goroutine 启动千递增计数，在每一个 goroutine 结束时递减计数。这需要一个特殊设计的计数器，它可以被多个 goroutine 安全地操作，然后又一个方法一直等到他变为 0。这个计数器类型是 sync\.WaitGroup。下面是完整的服务器代码：
```go
// exercise8/e4
```
注意 Add 和 Done 方法的不对称性。Add 递增计数器，它必须工作在 goroutine 开始之前执行，而不是在中间。另外，Add 有一个参数，但 Done 没有，它等价于 Add(-1)。使用 defer 来确保计数器在任何情况下都可以递减。在不知道迭代次数的情况下，上面的代码结构是通用的，符合习惯的并行循环模式。  

## 超时断开
下面的版本增加了超时断开的功能。这样服务端和客户端就各有两个断开连接的情况了，原本只有一种。  
服务端原本只要被动等待客户端断开就可以了，这个逻辑原本原本放在主 goroutine 中。现在服务端超时需要主动断开，客户端断开了，需要被动断开，这2个逻辑都需要一个单独的 goroutine，而主 goroutine 则阻塞接收这两个情况的通道，任意一个通道有数据，就断开并退出。  
客户端原本只需要响应接收标准输入的 Ctrl+Z 然后断开写半边的连接，这个逻辑也需要从主 goroutine 放到一个新的 goroutine 中。另外一种断开的连接是被动响应服务端的断开连接然后客户端也退出。这里还要稍微在复杂一点，如果是服务端的超时断开，则直接断开。如果是客户端的主动断开，则还需要继续等待服务端的断开，然后再退出。  
这里用到了大量的 select 多路复用：
```go
// exercise8/e6
```
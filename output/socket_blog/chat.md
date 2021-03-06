# 8.10 示例：聊天服务器
实现一个聊天服务器，它可以在几个用户之间相互广播文本消息。  
这个程序中有四种 goroutine：
+ 主 goroutine，就是 main 函数
+ 广播（broadcaster）goroutine。非常好的展示了 select 用法，因为它需要处理三种不同类型的消息
+ 每一个连接里有一个连接处理（handleConn）goroutine 
+ 每一个连接里还有一个客户写入（clientWriter）goroutine

## 主函数
主函数的工作是监听端口，接受连接请求。对每一个连接，它创建一个新的 handleConn。就像之前的并发回声服务器中那样：
```go
// main 函数
```

## 广播器
广播器，它的变量 clients 会记录当前连接的客户集合。其记录的内容是每一个客户端对外发送消息的通道：
```go
// 广播器
```
广播器监听两个全局的通道 entering 和 leaving。通过它们通知有客户进入和离开，如果从一个通道中接收到事件，它将更新 clients 集合。如果是客户离开，还会关闭对应客户对外发送消息的通道。  
广播器还监听 messages 通道，所有的客户都会将要广播的消息发送到这个通道。当收到一个消息后，就会把消息广播给所有客户。  

## 客户端处理函数
handleConn 函数创建一个对外发送消息的新通道，然后通过 entering 通道通知广播器新客户进入。接着，要读取客户发来的每一条消息，通过 messages 通道将每一条消息发送给广播器，发送时再每条消息前面加上发送者的ID作为前缀。一旦客户端将消息读取完毕，handleConn 通过 leaving 通道通知客户离开，然后关闭连接：
```go
// 客户端处理函数: handleConn
```
另外，handleConn 函数还为每一个客户创建了写入（clientWriter）goroutine，每个客户都从自己的通道中接收消息发送给客户端的网络连接。在广播器收到 leaving 通知并关闭这个接收消息的通道后，clientWriter 会结束通道的遍历后运行结束：
```go
// 客户端处理函数: clientWriter
```
给客户端发送的消息字符串需要用"\\n"结尾。如果换成"\\r\\n"结尾，平台的兼容性应该会更好。至少windows上的telnet客户端可以直接使用了。  

## 使用客户端进行聊天
完整的源码就是上面的四段代码，拼在一起就能运行了。  
和之前使用回声服务器一样，可以用 telnet 或者也可以用之前写的 netcat 作为客户端来聊天。  
当有 n 个客户 session 在连接的时候，程序并发运行着 2n+2 个相互通信的 goroutine，它不需要隐式的加锁操作也能做到并发安全。clients map 被限制在广播器这一个 goroutine 中，所以不会被并发的访问。唯一被多个 goroutine 共享的变量是通道以及 net\.Conn 的实例，它们也都是并发安全的。  

# 聊天服务器功能扩展
上面的聊天服务器提供了一个很好的架构，现在再在其之上扩展功能就很方便了。  

## 通知当前的用户列表
在新用户到来之后，告知该新用户当前在聊天室的所有的用户列表。每个用户加入后，系统都会自动生成一个用户名（基于用户的网络连接，之后会添加设置用户名的功能），就是要把这些存在的用户名打印出来。  
所有的用户列表只在广播器的 clients map 中，但是这个 map 又不包括用户名。所以先要修改数据类型，把每个连接的数据结构加上一个新的用户名字段：
```go
type client chan<- string // 对外发送消息的通道
type clientInfo struct {
	name string
	ch   client
}
```
原本使用 client 作为元素的通道和 map，现在全部也都要换成 clientInfo 作为元素。像新用户发送当前用户列表的任务也在广播器中完成：
```go
// e12 广播器
```
客户端处理函数还需要做少量的修改，主要是因为数据结构变了。原本给 entering 和 leaving 通道发送的是 ch。现在要发送封装好 who 的结构体。客户端处理函数的代码略，之后的扩展中会贴出来：
```go
cli := clientInfo{who, ch}
entering <- cli
```

## 断掉长时间空闲的客户端
如果在一段时间里，客户端没有任何输入，服务器就将客户端断开。之前的逻辑是，客户端处理函数会一直在阻塞在 input.Scan() 这里等待客户端输入。只要在另外一个 goroutine 中调用 conn.Close()，就可以让当前阻塞的读操作变成非阻塞，就像 input.Scan() 输入完成的读操作一样。不过这么做的话会有一点小问题，原本在主 goroutine 的结尾有一个`conn.Close()`操作，现在在定时的 goroutine 中还需要有一个关闭的操作。如果因为定时而结束的，就会有两次关闭操作。  
这里关闭的是 socket 连接，本质上就是文件句柄。尝试多次关闭貌似不会有什么问题，不过要解决这个问题也不难。一种是把响应用户输入的操作也放到 goroutine 中。现有有两个 goroutine 在运行，主 goroutine 则只要一直阻塞，通过一个通道等待其中任何一个 goroutine 完成后发送的信号即可。这样关闭的操作只在主 goroutine 中操作。下面的是客户端处理函数，包括上一个功能里修改的部分：
```go
// e13 客户端处理函数: handleConn
```
这里还简单加了一个限制客户端发送空消息的功能，在 input.Scan() 循环中。空消息不会发送广播，但是可以重置定时器的时间。  

## 客户端可以输入名字
在客户端连接后，不立刻进入聊天室，而是先输入一个名字。考虑到名字不能和已有的名字重复，而现有的名字都保存在广播器里的 clients 这个 map 中。所以客户端输入的名字需要在 clients 中查找一下是否已经有人用了。现在有了按名字进行查找的需求，clients 类型更适合使用一个以名字为 key 的 map 而不是原本的集合。这个 map 的 value 就是向该客户发送消息的通道，也就是最初这个集合的 key 的值：
```go
clients := make(map[string]client) // 所有连接的客户端集合
```

**客户端处理函数**  
在客户端处理函数的开头，需要增加注册用户名的过程。用户名注册的处理过程比较复杂，所以单独封装到了一个函数 clientRegiste 中：
```go
// 客户端处理函数
func handleConn(conn net.Conn) {
	who := clientRegiste(conn) // 新增这一行，注册获取用户名

	ch := make(chan string) // 对外发送客户消息的通道
	go clientWriter(conn, ch)

	// who := conn.RemoteAddr().String() // 去掉这一行
	// 之后的代码不变
}
```
这里使用一个交互的方式来获取用户名，代替原本通过连接的信息自动生成。这个函数是串行的，只有在返回用户名后，才会继续执行下去。之后的代码和之前是一样的。  
在 clientRegiste 函数中，不停的和终端进行交互，处理收到的消息，如果用户名可用，继续执行之后的流程。如果用户名不可用，则提示用户继续处理：
```go
// e14 客户端处理函数 clientRegiste
```
这里只有最简单的功能，还可以增加输入超时，以及尝试次数的限制。所以把这个函数独立出来完成功能，更方便之后对注册函数进行扩展。  
函数的主要逻辑就是 input.Scan() 的循环，这和 handleConn 中的循环十分相似。如果之后再加上输入超时，这两段的处理逻辑只有极小部分的差别，所以这部分代码也可以单独写一个函数。这里避免过早的优化，暂时就先这样，看着也比较清晰。之后要添加超时功能的时候，再把这部分重复的代码独立出来。*这部分优化最后完整的代码里会有。*  

**广播器**  
在广播器的 select 里要加一个分支，用来处理用户名的请求。收到请求后，判断是否已经存在，把结果返回给 clientRegiste。因为 clients 是只有广播器可见的，所以这里要使用通道传递过来，判断后再用通道把结果传回去。这样可以保证 clients 变量只在这一个 goroutine 里被使用（包括修改）。另外，每个客户端的注册都使用一个通道将注册信息发送给广播器，但是广播器返回的内容，需要对每个客户端使用不同的通道。所以这里，广播器新创建了专门用于注册交互的数据结构：
```go
type registeInfo struct {
	name string
	ch   chan<- bool
}

var register = make(chan registeInfo) // 注册用户名的通道
```
客户注册的函数创建一个布尔型的通道，加上用户的名字封装到 registeInfo 结构体中。然后广播器判断后，把结果通道 registeInfo 里的 ch 字段这个通道，把结果返回给对应的客户注册函数。  
下面是广播器 broadcaster 的代码，主要是 select 新增了一个分支，处理注册用户名：
```go
// e14 广播器
```

## 预防客户端延迟影响
最后还有一个问题，就是客户端可能会卡或者延迟，但是客户端的问题不能影响到服务器的正常运行。不过我没法实现一个这样的有延迟的客户端，默认操作系统应该就已经非常友好的帮我们处理掉了，把从网络上接收到的数据暂存在缓冲区里（*对于TCP连接还有乱序重组和超时重传，这些我们都不需要关心了*），等待程序去读取。代码里接收的操作应该是直接从缓冲区读取，这时服务的已经发送完毕了。所以现在只能照着下面的思路写了：
>任何客户程序读取数据的时间很长最终会造成所有的客户卡住。修改广播器，使它满足如果一个向客户写入的通道没有准备好接受它，那么跳过这条消息。还可以给每一个向客户发送消息的通道增加缓冲，这样大多数的消息不会丢弃；广播器在这个通道上应该使用非阻塞的发送方式。

客户端处理函数中创建的发送消息的通道改用有缓冲区的通道：
```go
// 客户端处理函数
func handleConn(conn net.Conn) {
	defer conn.Close() // 退出时关闭客户端连接，现在有分支了，并且可能会提前退出

	who, ok := clientRegiste(conn) // 注册获取用户名
	if !ok { // 用户名未注册成功
		fmt.Fprintln(conn, "\r\nName registe failed...")
		return
	}

	ch := make(chan string, 10) // 有缓冲区，对外发送客户消息的通道
	go clientWriter(conn, ch)

	// 省略后面的代码
}
```
然后广播器的 select 对应的 messages 通道的分支，改成非阻塞的方式：
```go
select {
case msg := <-messages:
	// 把所有接收的消息广播给所有的客户
	// 发送消息通道
	for name, cli := range clients {
		select {
		case cli <- msg:
		default:
			fmt.Fprintf(os.Stderr, "send message failed: %s: %s\n", name, msg)
		}
	}
// 其他分支略过
}
```

下面是聊天服务器最后完整的代码。这里的改变还包括了上一节最后提到的注册用户名时的输入的超时。已经两次用到了输入超时，分别在 handleConn 和 clientRegiste 中，这里也就把这部分代码单独写了一个函数 inputWithTimeout。完整代码如下：
```go
// e15/main.go
```

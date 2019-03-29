# goroutine 和通道
Go 有两种并发编程风格：
+ goroutine 和通道（chennle），支持**通信顺序进程**（Communicating Sequential Process, CSP），CSP 是一个并发的模式，在不同的执行体（goroutine）之间传递值，但是变量本身局限于单一的执行体。
+ **共享内存多线程**的传统模型，和在其他主流语言中使用线程类似。

# 8.4 通道
如果说 goroutine 是 Go 程序并发的执行体，**通道**就是它们之间的连接。每一个通道是一个具体类型的导管，叫作通道的**元素类型**。  
像 map 一样，通道是一个使用 make 创建的数据结构的应用。当复制或者作为参数传递到一个函数时，复制的是引用，这样调用者和被调用者都引用了同一份数据结构。和其他引用类型一样，通道的零值是 nil。  
通道有两个主要操作：**发送**（send）和**接收**（receive），两者统称为**通信**。通道还支持第三个操作：**关闭**（close），它设置一个标志位来指示值当前已经发送完毕。  
使用简单的 make 调用创建的通道叫**无缓冲（unbuffered）通道**，但 make 还可以接受第二个可选参数，一个表示通道**容量**的整数。如果容量是0，创建的也是无缓冲通道。  

## 8.4.1 无缓冲通道
使用无缓冲通道进行的通信导致发送和接收 goroutine **同步化**。因此无缓冲通道也称为**同步通道**。  
通过通道发送消息有两个重要的方面需要考虑：
+ 每条消息有一个值
+ 通信本身以及通信发生的时间。当我们强调这方面的时候，把消息叫做**事件**（event）

当事件没有携带额外的信息时，它单纯的目的是进行同步。和 map 实现的集合一样，可以使用一个 struct{} 元素类型的通道来强调它，尽管通常使用 bool 或 int 类型的通道来做相同的事情。因为`done <- 1`更简短。*书上讲集合的时候，使用的是 bool类型，这里讲事件同步，使用的是空结构体。*

## 8.4.2 管道
通道可以用来连接 goroutine，这样一个的输出是另一个的输入，这个叫**管道**（pipline）。  

**关闭通道**  
如果发送方知道没有更多的数据要发送，告诉接收者所在的 goroutine 可以停止等待是很有用的。这可以通过调用内置的 Close 函数来关闭通道：
```go
ch1 := make(chan bool) // 创建通道 ch1
// 下面是关闭通道
close(ch1)
```
在通道关闭后，任何后续的发送操作将会导致应用崩溃。当关闭的通道被**读完**（就是最后一个发送的值被接收）后，所有后续的接收操作都会立即返回，返回值是对应类型的零值。  
关闭通道还可以作为一个广播机制，后面的章节会具体讲。  

**检查通道的关闭**  
没有一个直接的方式来判断是否通道已经关闭，不过可以接收返回两个参数：接收到的元素，以及一个布尔值（通常是ok），返回 true 表示接收成功，返回 false 表示当前的接收操作在一个关闭的并且读完的通道上。*这个方法检查的也不是通道是否关闭了，而是通道里的值是否已经取完了。只有关闭的通道，才能保证不会有新值进入，把里面的值都取完后，会返回 false 表示这次取到的是通道关闭后的零值，而不是原本就是一个值为零的数据。*  

另外，还提供了一个 range 循环语法可以在通道上迭代。这个语法更为方便接收在通道上所有发送的值，接收完最后一个值后结束循环。  

**垃圾回收**  
结束时，关闭没一个通道不是必需的。只有在通知接收方 goroutine 所有的数据都发送完毕的时候才需要关闭通道。通道也可以通过垃圾回收器根据它是否可以访问来决定是否回收它，而不是根据它是否关闭。  

## 8.4.3 单向通道类型
Go 还提供了**单向通道**类型，仅仅导出发送或接收操作。类型`chan<- int`是一个**只能发送**的通道，允许发送单不允许接收。反之，类型`<-chan int`是一个**只能接收**的通道，允许接收但是不能发送。这里像箭头一样的操作符相对于 chan 关键字的位置是一个帮助记忆的点。如果违反这里的接收或发送的原则，在编译时会被检查出来。  
在函数定义时，指定了单向通道的类型。在函数调用时，依然是把正常定义的双向通道类型传值给函数的参数。函数的调用会隐式地将普通的通道类型转化为要求的单向通道的类型。在任何赋值操作中将双向通道转换为单向通道都是允许的，但是反过来是不行的。一旦有一个单向通道，是没有办法通过它获取到引用同一个数据结构的双向通道的类型的。  

## 8.4.4 缓冲通道
缓冲通道有一个元素队列，队列的最大长度在创建的时候通过 make 的容量参数来设置：
```go
ch1 := make(chan string, 3)
```
通过调用内置的 cap 函数，可以获取通道缓冲区的容量。这种需求不太常见。  
通过调用内置的 len 函数，可以获取通道内的元素个数。不过在并发程序中这个信息会随着检索操作很快过时，所以它的价值很低，但是它在错误诊断和性能优化的时候很有用。  

**这不是队列**  
发送和接收操作可以在同一个 goroutine 中，但在真实的程序中通常由不同的 goroutine 执行。因为语法简单，新手有时候粗暴地将缓冲通道作为队列在单个 goroutine 中使用，但是这是个错误的用法。通道和 goroutine 的调度深度关联，如果没有另一个 goroutine 从通道进行接收，发送者（也许是整个程序）有被永久阻塞的风险。如果仅仅需要一个简单的队列，使用切片创建一个就好了。  

**示例：并发请求最快的镜像资源**  
下面的例子展示一个使用缓冲通道的应用。它并发地向三个**镜像地址**发请求，镜像指相同但分布在不同地理区域的服务器。它将它们的响应通过一个缓冲通道进行发送，然后只接收第一个返回的响应，因为它是最早到达的。所以 mirroredQuery 函数甚至在两个比较慢的服务器还没有响应之前返回了一个结果。（偶然情况下，会出现像这个例子中的几个 goroutine 同时在一个通道上并发发送，或者同时从一个通道接收的情况。）：
```go
func mirroredQuery() string {
	responses := make(chan string, 3) // 有几个镜像，就要多大的容量，不能少
	go func () { responses <- request("asia.gopl.io") }()
	go func () { responses <- request("europe.gopl.io") }()
	go func () { responses <- request("americas.gopl.io") }()
	return <- responses // 返回最快一个获取到的请求结果
}

func request(hostname string) (response string) { return "省略获取返回的代码" }
```

**goroutine 泄露**  
在上面的示例中，如果使用的是无缓冲通道，两个比较慢的 goroutine 将被卡住，因为在它们发送响应结果到通道的时候没有 goroutine 来接收。这个情况叫做 **goroutine 泄漏**。它属于一个 bug。不像回收变量，泄漏的 goroutine 不会自动回收，所以要确保 goroutine 在不再需要的时候可以自动结束。  

## 通道缓冲的选择
无缓冲和缓冲通道的选择，缓冲通道容量大小的选择，都会对程序的正确性产生影响。无缓冲通道提供强同步保障，因为每一次发送都需要和一次对应的接收同步；对于缓冲通道，这些操作则是解耦的。如果知道要发送的值数量的上限，通常会创建一个容量是使用上限的缓冲通道，在接收第一个值前就完成所有的发送。在内存无法提供缓冲容量的情况下，可能导致程序死锁。  

TODO： 通道缓冲对程序性能的影响。蛋糕店示例：
TODO: gopl.io/ch8/cake ，性能基准（参考 11.4 节）

# 8.5 并行循环
重申了 5.6.1 里提到的**捕获迭代变量**的问题，这次是在 goroutine 中的情况，和避免方式。  
TODO：生成图像的缩略图，gopl.io/ch8/thumbnail 包提供的 ImageFile 函数，书上不展示这个函数的细节。  

# 8.7 使用 select 多路复用
有时候需要在多个通道上接收，不能只从一个通道上接收，因为任何一个操作都会在完成前阻塞。所以需要**多路复用**那些操作过程，为了实现这个目的，需要一个 select 语句：
```go
select {
case <-ch1:
	// ...
case x := <-ch2:
	// ...use x...
case ch3 <- y:
	// ...
default:
	// ...
}
```
上面展示的是 select 语句的通用形式。像 switch 语句一样，它有一系列的情况和一个可选的默认分支。每一个情况指定一次**通信**（在一些通道上进行发送或接收操作）和关联的一段代码块。接收表达式操作可能出现在它本身上，像第一个情况，或者在一个短变量声明中，像第二个情况。第二种形式可以让你引用所接收的值。  
select 一直等待，直到一次通信来告知有一些情况可以执行。然后，它进行这次通信，执行此情况所对应的语句，其他的通信将不会发生。  

## 使用示例
下面是一个微妙的例子。通道 ch 的缓冲区大小为 1，它要么是空的，要么是满的，因此只有在其中一个状况下可以执行，要么在 i 是偶数时发送，要么在 i 是奇数时接收。它总是输出 0 2 4 6 8：
```go
func main() {
	ch := make(chan int, 1)
	for i := 0; i < 10; i++ {
		select {
		case x := <-ch:
			fmt.Println(x)
		case ch <- i:
		}
	}
}
```
如果多个情况同时满足，select 随机选择一个，这样保证每一个通道有相同的机会被选中。在前一个例子中增加缓冲区的容量，会使输出变得不可确定，因为当缓冲既不空也不满的情况，相当于 select 语句在随机做选择。  
## 非阻塞模式
有时候我们试图在一个通道上发送或接收，但是不想在通道没有准备好的情况下被阻塞，**非阻塞通信**。这使用 select 语句也可以做到。select 可以有一个默认情况，它用来指定在没有其他的通信发生时可以立即执行的动作。  
下面的 select 语句尝试从 abort 通道中接收一个值，如果没有值，它什么也不做。这是一个非阻塞的接收操作。重复这个动作称为对通道**轮询**：
```go
select {
case <-abort:
	fmt.Println("Launch aborted!")
	return
default:
	// 不执行任何操作
}
```

## 通道的零值
通道的零值是 nil。令人惊讶的是，nil 通道有时候很有用。因为在 nil 通道上发送和接收将永远阻塞。对于 select 语句中的情况，如果其通道是 nil，它将永远不会被选择。可以用 nil 来开启或禁用特性所对应的情况，比如超时处理或者取消操作，响应其他的输入事件或者发送事件。  

# 8.8 示例：并发目录遍历
这里要构建一个程序，根据命令行指定的输入，报告一个或多个目录的磁盘使用情况，类似 UNIX 的 du 命令。  

## 递归遍历目录
大多数的工作由下面的 walkDir 函数完成，它使用 dirents 辅助函数来枚举目录中的条目：
```go
// walkDir 递归地遍历以 dir 为根目录的整个文件树
// 并在 fileSizes 上发送每个已找到的文件的大小
func walkDir(dir string, fileSizes chan<- int64) {
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			walkDir(subdir, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}

// dirents 返回 dir 目录中的条目
func dirents(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	}
	return entries
}
```
ioutil\.ReadDir 函数返回一个 os\.FileInfo 类型的切片，针对单个文件同样的信息可以通过调用 os\.Stat 函数来返回。对每一个子目录，walkDir 递归调用它自己，对于每一个文件，walkDir 发送一条消息到 fileSizes 通道。消息是文件所占用的字节数。  

## 计算大小并输出
下面的 main 函数使用两个 goroutine。后台 goroutine 调用 walkDir 遍历命令行上指定的每一个目录，最后关闭 fileSizes 通道。主 goroutine 计算从通道中接收的文件的大小的和，最后输出总数：
```go
func main() {
	// 确定初始目录
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}
	// 遍历文件树
	fileSizes := make(chan int64)
	go func() {
		for _, root := range roots {
			walkDir(root, fileSizes)
		}
		close(fileSizes)
	}()
	// 输出结果
	var nfiles, nbytes int64
	for size := range fileSizes {
		nfiles++
		nbytes += size
	}
	printDiskUsage(nfiles, nbytes)
}

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files   %.2f GB\n", nfiles, float64(nbytes)/(1<<30)) // 1<<30 就是 2**30 就是 1024*1024*1024
}
```
现在程序可以正常的工作。

## 汇报进度
如果程序可以汇报进度的话，会更加友好。如果仅仅只是把 printDiskUsage 调用移动到循环内部，会有非常多的输出。  
下面的示例，修改了主 goroutine 中记录结果的部分。不是在每次迭代中输出，而是加了一个定时器，通过 select 定期输出一次结果。另外还加上了 -v 参数来控制，可以选择性的开启这个功能。如果不开启功能，那么 tick 通道的值就是 nil，它对应的分支在select 中就永远是阻塞的。相当于没有开启这个选项，很直观的理解：
```go
var verbose = flag.Bool("v", false, "周期性的输出进度")

func main() {
	// 确定初始目录，没变化
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}
	// 遍历文件树，没变化
	fileSizes := make(chan int64)
	go func() {
		for _, root := range roots {
			walkDir(root, fileSizes)
		}
		close(fileSizes)
	}()

	// 定期输出结果
	var tick <-chan time.Time
	if *verbose {
		tick = time.Tick(500 * time.Millisecond)
	}
	var nfiles, nbytes int64
loop:
	for {
		select {
		case size, ok := <-fileSizes:
			if !ok {
				break loop // fileSizes 关闭，则退出，相当于原来的遍历结束
			}
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		}
	}
	printDiskUsage(nfiles, nbytes)
}
```
因为这个版本有两个通道需要接收 size、tick，所以无法使用 range 循环了。所以第一个 select 的分支需要通过第二个参数 ok  来判断通道是否已经关闭。这个的 break 退出使用了标签，因为没有标签的 break 只能跳出当前的 select 这层，而这里是需要跳出外层的 for 循环。  
这里的 flag 的解析也值得借鉴，非常简单。首先是解析指定的参数，这里是 -v 参数。多余的参数会通过 flag.Args() 返回一个字符串切片。调用的时候，必须把解析的参数放在前面：
```
PS H:\Go\src\gopl\ch8\du2> go run main.go -v E:\BaiduNetdiskDownload E:\XMPCache E:\Downloads
4 files   0.02 GB
41 files   2.16 GB
177 files   6.99 GB
567 files   46.66 GB
605 files   50.26 GB
PS H:\Go\src\gopl\ch8\du2>
```

## 提高并发效率
还可以进一步提高效率，这里的 walkDir 也是可以并发调用从而充分利用磁盘系统的并行机制。这个版本使用了 sycn.WaitGroup 来为并发调用的 walkDir 计数。当计数器为减为 0 的时候，关闭 fileSizes 通道：
```go
var verbose = flag.Bool("v", false, "周期性的输出进度")

func main() {
	// 确定初始目录，没变化
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	// 并行遍历每一个文档树
	fileSizes := make(chan int64)
	var n sync.WaitGroup
	for _, root := range roots {
		n.Add(1)
		go walkDir(root, &n, fileSizes) // 注意，多传了一个参数
	}
	go func() {
		n.Wait()
		close(fileSizes)
	}()

	// 定期输出结果，没变化
	var tick <-chan time.Time
	if *verbose {
		tick = time.Tick(500 * time.Millisecond)
	}
	var nfiles, nbytes int64
loop:
	for {
		select {
		case size, ok := <-fileSizes:
			if !ok {
				break loop // fileSizes 关闭，则退出，相当于原来的遍历结束
			}
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		}
	}
	printDiskUsage(nfiles, nbytes)
}

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files   %.2f GB\n", nfiles, float64(nbytes)/(1<<30)) // 1<<30 就是 2**30 就是 1024*1024*1024
}

func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) { // 注意，多了个参数
	defer n.Done() // 记得退出时计数器要减1
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			walkDir(subdir, n, fileSizes) // 注意，多了个参数
		} else {
			fileSizes <- entry.Size()
		}
	}
}
```

## 限制并发
还需要限制一下并发数，这里要修改一下 dirents 函数来使用计数信号量进行限制，防止同时打开太多的文件：
```go
// 用于限制目录并发数的计数信号量
var sema = make(chan struct{}, 20)

// dirents 返回 dir 目录中的条目
func dirents(dir string) []os.FileInfo {
	sema <- struct{}{}        // 获取令牌
	defer func() { <-sema }() // 释放令牌
	entries, err := ioutil.ReadDir(dir) // 这个打开文件的操作需要限制并发，在这句之前加上计数信号量，非常合适
	if err != nil {
		fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	}
	return entries
}
```
现在这个版本的是最好的了。不过下面还会再增加一个取消的操作，这里的取消会用到广播的机制。  

# 8.9 取消（广播）
一个 goroutine 无法直接终止另一个，因为这样会让所有的共享变量状态处于不确定状态。正确的做法是使用通道来传递一个信号，当 goroutine 接收到信号时，就终止自己。这里要讨论的是如何同时取消多个 goroutine。  
一个可选的做法是，给通道发送你要取消的 goroutine 同样多的信号。但是如果一些 goroutine 已经自己终止了，这样计数就多了，就会在发送过程中卡住。如果某些 goroutine 还会自我繁殖，那么信号的数量又会太少。通常，任何时刻都很难知道有多少个 goroutine 正在工作。对于取消操作，这里需要一个可靠的机制在一个通道上**广播**一个事件，这样所以的 goroutine 就都能收到信号，而不用关心具体有多少个 goroutine。  
当一个通道关闭且已经取完所有发送的值后，接下来的接收操作都会立刻返回，得到零值。就可以利用这个特性来创建一个广播机制。第一步，创建一个取消通道，在它上面不发送任何的值，但是它的关闭表明程序需要停止它正在做的事前。

## 查询状态
还要定义一个工具函数 cancelled，在它被调用的时候检测或**轮询**取消状态：
```go
var done = make(chan struct{})

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}
```

## 发送取消广播
接下来，创建一个读取标准输入的 goroutine，它通常连接到终端，当用户按回车后，这个 goroutine 通过关闭 done 通道来广播取消事件：
```go
// 当检测到输入时，广播取消
go func() {
	os.Stdin.Read(make([]byte, 1)) // 读一个字节
	close(done)
}()
```

## 响应取消操作
现在要让所有的 goroutine 来响应这个取消操作。在主 goroutine 中的 select 中，尝试从 done 接收。如果接收到了，就需要进行取消操作，但是在结束之前，它必须耗尽 fileSizes 通道，丢弃它所有的值，知道通道关闭。这么做是为了保证所有的 walkDir 调用可以执行完，不会卡在向 fileSizes 通道发送消息上：
```go
for {
	select {
	case <-done:
		// 耗尽 fileSizes，让已经创建的 goroutine 结束
		for range fileSizes {
			// 什么也不做
		}
		return
	case siez, ok := <-fileSizes:
		if !ok {
			break loop
		}
		nfiles++
		nbytes += siez
	case <-tick:
		printDiskUsage(nfiles, nbytes)
	}
}
```
walkDir 的 goroutine 在开始的时候轮询取消状态。如果是取消的状态，就什么都不做立即返回。这样在取消后创建的 goroutine 就会什么都不做而是立刻返回：
```go
func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
	defer n.Done()
	if cancelled() {
		return
	}
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			walkDir(subdir, n, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}
```
现在基本就避免了在取消后创建新的 goroutine。但是其他已经创建的 goroutine 则会等待他们执行完毕。要想更快的响应，就需要更多的程序逻辑变更入侵。要确保在取消事件之后没有更多昂贵的操作发生。这就需要更新更多的代码，但是通常可以通过在少量重要的地方检察取消装来来打到目的。在 dirents 中获取信号量令牌的操作也可需要快速结束：
```go
func dirents(dir string) []os.FileInfo {
	select {
	case sema <- struct{}{}: // 获取令牌
	case <-done:
		return nil // 取消
	}
	defer func() { <-sema }() // 释放令牌
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	}
	return entries
}
```
现在，当取消事件发生时，已经进入 dirents 函数的调用，如果已经获取到了令牌，则会执行完毕，但是返回后，在地递归调用 walkDir 的时候就会快速退出。那些还没获取令牌的调用，此时在 select 中会因为从 done 通道中接收到取消的广播而直接返回 nil。  

## 测试的技巧
期望的情况是，当然是当取消事件到来时 main 函数可以返回，然后程序随之退出。如果发现在取消事件到来的时候 main 函数没有返回，可以执行一个 panic 调用。从崩溃的转存储信息中通常含有足够的信息来帮助我们分析，发现哪些 goroutine 还没有合适的取消。也可能是已经取消了，但是需要的时间比较长。总之，使用 panic 可以帮助查找原因。  
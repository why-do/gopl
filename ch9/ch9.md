# 9.1 竞态
**并发**，如果无法确定一个事件先于另外一个事件，那么这两个事件就是并发的。  
**并发安全**（concurrency-safe），如果一个函数在并发调用时仍然能正确工作，那么这个函数就是并发安全的。如果一个类型的所有可访问方法和操作都是并发安全的，则它可称为并发安全的类型。  
并发安全的类型是特例而不是普遍存在的，对于绝大部分变量，如要回避并发访问，只有下面几种办法：
+ **限制**变量只存在于一个 goroutine 内。
+ 维护一个更高层的**互斥不变量**。

竞态是指在多个 goroutine 按某些交错顺序执行时程序无法给出正确的结果。  
**数据竞态**（data race）是竞态的一种。数据竞态发生于两个 goroutine 并发读写同一个变量并且至少其中一个是写入时。有三种方法来避免数据竞态：
+ 不要修改变量
+ 避免从多个 goroutine 访问同一个变量。就是**限制**
+ 允许多个 goroutine 访问同一个变量，但在同一时间只有一个 goroutine 可以访问。这种方法称为**互斥机制**

Go 箴言：“不要通过共享内存来通信，而应该通过通信来共享内存”。  

# 9.2 互斥锁：sync.Mutex
使用缓冲通道可以实现一个计数信号量，可以用于同时发起的 goroutine 的数量。一个计数上限为 1 的信号量称为**二进制信号量**（binary semaphore）。  
使用二进制信号量就可以实现互斥锁：
```go
var (
	sema    = make(chan struct{}, 1) // 用来保护 balance 的二进制信号量
	balance int
)

func Deposit(amount int) {
	sema <- struct{}{} // 获取令牌
	balance = balance + amount
	<-sema // 释放令牌
}

func Balance() int {
	sema <- struct{}{} // 获取令牌
	b := balance
	<-sema // 释放令牌
	return b
}
```
**互斥锁**模式应用非常广泛，所以 sync 包有一个单独的 Mutex 类型来支持这种模式：
```go
import "sync"

var (
	mu      sync.Mutex // 保护 balance
	balance int
)

func Deposit(amount int) {
	mu.Lock()
	balance = balance + amount
	mu.Unlock()
}

func Balance() int {
	mu.Lock()
	b := balance
	mu.Unlock()
	return b
}
```

互斥量**保护**共享变量。按照惯例，被互斥量保护的变量声明应当紧接在互斥量的声明之后。如果实际情况不是如此，请加注释说明。  

**临界区域**，在 Lock 和 Unlock 之间的代码，可以自由地读取和修改共享变量，这一部分称为临界区域。

**封装**，即通过在程序中减少对数据结构的非预期交互，来帮助我们保证数据结构中的不变量。类似的原因，封装也可以用来保持并发中的不变性。所以无论是为了保护包级别的变量，还是结构中的字段，当使用一个互斥量时，都请确保互斥量本身以及被保护的变量都没有导出。  

# 9.3 读写互斥锁：sync.RWMutex
**多读单写锁**，允许只读操作可以并发执行，但写操作需要获得完全独享的访问权限。Go 语言中的 sync\.RWMutex 提供了这种功能：
```go
var mu sync.RWMutex
var balance int

func Balance() int {
	mu.RLock() // 读取
	defer mu.RUnlock()
	return balance
}
```
Balance 函数尅调用 mu\.RLock 和 mu\.RUnlock 方法来分别获取和释放一个**读锁**（也称为**共享锁**）。而之前的 mu\.Lock 和 mu\.Unlock 方法则是分别获取和释放一个**写锁**（也称为**互斥锁**）。  
一般情况下，不应该假定那些**逻辑上**只读的函数和方法不会更新一些变量。比如，一个看起来只是简单访问的方法，可能会递增内部使用的计数器，或者更新一个缓存来让重复的调用更快。如果不确定，就应该使用互斥锁。  

**读锁的应用场景**  
仅在绝大部分 goroutine 都在获取读锁并且锁竞争比较激烈时，RWMutex 才有优势。因为 RWMutex 需要更复杂的内部实现，所以在竞争不激烈时它比普通的互斥锁慢。  

# 9.4 内存同步
现代的计算机一般会有多个处理器，每个处理器都有内存的本地缓存。为了提高效率，对内存的写入是缓存在每个处理器中的，只在必要时才刷回内存。甚至刷会内存的顺序都可能与 goroutine 的写入顺序不一致。像通道通信或者互斥锁操作这样的同步源语都会导致处理器把累积的写操作刷回内存并提交。但这个时刻之前 goroutine 的执行结果就无法保证能被运行在其他处理器的 goroutine 观察到。  
考虑如下的代码片段可能的输出：
```go
var x, y int
go func() {
	x = 1
	fmt.Print("y:", y, " ")
}
go func() {
	y = 1
	fmt.Print("x:", x, " ")
}
```
下面4个是显而易见的可能的输出结果：
```
y:0 x:1
x:0 y:1
x:1 y:1
y:1 x:1
```
但是下面的输出也是可能出现的：
```
x:0 y:0
y:0 x:0
```
在某些特定的编译器、CPU 或者其他情况下，这些确实可能发生。  

单个 goroutine 内，每个语句的效果保证按照执行的顺序发生，也就是说，goroutine 是**串行一致**的（sequentially consistent）。但在缺乏使用通道或者互斥量来显式同步的情况下，并不能保证所有的 goroutine 看到的事件顺序都是一致的。  
上面的两个 goroutine 尽管打印语句是在赋值另外一个变量之后，但是一个 goroutine 并不一定能观察到另一个 goroutine 对变量的效果。所以可能输出的是一个变量的**过期值**。  
尽管很容易把并发简单理解为多个 goroutine 中语句的某种交错执行方式。如果两个 goroutine 在不同的 CPU 上执行，每个 CPU 都有自己的缓存，那么一个 goroutine 的写入操作在同步到内存之前对另外一个 goroutine 的打印变量的语句是不可见的。  
这些并发的问题都可以通过采用简单、成熟的模式来避免，即在可能的情况下，把变量限制到单个 goroutine 中，对于其他变量，使用互斥锁。  

# 9.5 延迟初始化：sync.Once
延迟一个昂贵的初始化步骤到有实际需求的时刻是一个很好的实践。预先初始化一个变量会增加程序的启动延迟，并且如果实际执行时有可能根本用不上这个变量，那么初始化也不是必需的。  
sync 包提供了针对一次性初始化问题的特化解决方案：sync.Once。从概念上来讲，Once 包含一个布尔变量和一个互斥量，布尔变量记录初始化是否已经完成，互斥量则负责保护这个布尔变量和客户端的数据结构。Once 唯一的方法 Do 以初始化函数作为它的参数：
```go
var loadIconsOnce sync.Once
var icons map[string]image.Image

// 这是个昂贵的初始化步骤
func loadIcons() {
	icons = map[string]image.Image{
		"spades.png":   loadIcon("spades.png"),
		"hearts.png":   loadIcon("hearts.png"),
		"diamonds.png": loadIcon("diamonds.png"),
		"clubs.png":    loadIcon("clubs.png"),
	}
}

// 并发安全
func Icon(name string) image.Image {
	loadIconsOnce.Do(loadIcons)
	return icons[name]
}
```
每次调用 Do 方法时，会先锁定互斥量并检查里边的布尔变量。在第一次调用时，这个布尔变量为 false，Do 会调用它参数的方法，然后把布尔变量设置为 true。之后 DO 方法的调用相当于空操作，只是通过互斥量的同步来保证初始化操作对内存产生的效果对所有的 goroutine 可见。以这种方式来使用 sync.Once，可以避免变量在构造完成之前就被其他 goroutine 访问。  

# 9.8 goroutine 与线程
**这章节都是概念，基本就是在抄书**  
goroutine 与操作系统（OS）线程之间的差异本质上属于量变。但是足够大的量变会变成质变，所以还是要区分一下两者的差异。  

## 9.8.1 可增长的栈
每个 OS 线程都有一个固定大小的栈内存（通常为 2MB），栈内存区域用户保存在其他函数调用期间那些正在执行或临时暂停的函数中的局部变量。这个固定的大小对小的 goroutine 来说太大了，对于要创建数量巨大的 goroutine 来说，就会有巨大的浪费。另外，对于更复杂或者深度递归的函数，固定大小的栈又会不够大。改变这个固定大小，调小了可以允许创建更多的线程，改大了则可以容许更深的递归，但两者无法同时兼容。  
gotouine 也用于存放那些正在执行或临时暂停的函数中的局部变量。但栈的大小不是固定的，它可与按需增大或缩小。goroutine 的栈大小限制可以达到 1GB。当然，只有极少的 goroutine 会使用这么大的栈。  

## 9.8.2 goroutine 调度
这节完全是照书抄的。  

**OS线程调度器**  
OS线程由OS内核来调度。每隔几毫秒，一个硬件时钟中断发送到CPU、CPU调用一个叫**调度器**的内核函数。这个函数暂停当前正在运行的线程，把它的寄存器信息保存到内存，查看线程列表并决定接下来运行哪一个线程，再从内存恢复线程的注册表信息，最后继续执行选中的线程。因为OS线程由内核来调度，所以控制权限从一个线程到另外一个线程需要一个完整的**上下文切换**（context switch）：即保存一个线程的状态到内存，再恢复另外一个线程的状态、最后更新调度器的数据结构。考虑这个操作涉及的内存局域性以及涉及的内存访问数量，还有访问内存所需的CPU周期数量的增加，这个操作其实是很慢的。  

**Go调度器**  
Go 运行时包含一个自己的调度器，这个调度器使用一个称为**m:n 调度**的技术（因为它可以复用/调度 m 个 goroutine 到 n 个OS线程）。Go 调度器与内核调度器的工作类似，但 Go 调度器值需关心单个 Go 程序的 goroutine 调度问题。  

**差别**  
与操作系统的线程调度器不同的是，Go 调度器不是由硬件时钟来定期触发的，而是由特定的 Go 语言结构来触发的。比如当一个 goroutine 调用 time\.Sleep 或被通道阻塞或对互斥量操作时，调度器就会将这个 goroutine 设为休眠模式，并运行其他 goroutine 直到前一个可重新唤醒为止。因为它不需要切换到内核语境，所以调用一个 goroutine 比调度一个线程成本低很多。  

## 9.8.3 GOMAXPROCS
Go 调度器使用 GOMAXPROCS 参数来确定需要使用多少个OS线程来同时执行 Go 代码，默认值是机器上的CPU数量（GOMAXPROCS 是 m:n 调度中的 n）。正在休眠或者正被通道通信阻塞的 goroutine 不需要占用线程。阻塞在 I\/O 和其他系统调用中或调用非 Go 语言写的函数的 goroutine 需要一个独立的OS线程，但这个线程不计算在 GOMAXPROCS 内。  
可以用 GOMAXPROCS 环境变量或者 runtime\.GOMAXPROCS 函数来显式控制这个参数。可以用一个小程序来看看 GOMAXPROCS 的效果，这个程序无止境地输出0和1：
```go
func main() {
	var n int
	flag.IntVar(&n, "n", 1, "GOMAXPROCS")
	flag.Parse()
	runtime.GOMAXPROCS(n)
	for {
		go fmt.Print(0)
		fmt.Print(1)
	}
}
```
这里使用命令行参数来控制线程数量。  
Linux 中应该可以直接设置 GOMAXPROCS 环境变量来运行程序：
```
$ GOMAXPROCS=1 go run main.go
$ GOMAXPROCS=2 go run main.go
```
GOMAXPROCS 为1时，每次最多只能由一个 goroutine 运行。最开始是主 goroutine，它会连续输出很多1。在运行了一段时间之后，Go 调度器让主 goroutine 休眠，并唤醒另一个输出0的 goroutine，让它有机会执行。所以执行结果能看到大段的连续的0或1。  
GOMAXPROCS 为2时，就有两个可用的OS线程，所以两个 goroutine 可以同时运行，输出的0和1就会交替出现（我看到的是小段小段的交替）。  

## 9.8.4 goroutine 没有标识
在大部分支持多线程的操作系统和编程语言里，当前线程都有一个独特的标识，它通常可以取一个整数或者指针。这个特性让我们可以轻松构建一个**线程的局部存储**，它本质上就是一个全局的 map，以线程的标识为 key，这样各个线程都可以独立地用这个 map 存储和获取值，而不受其他线程的干扰。  
goroutine 没有可供程序员访问的表示。这个是有设计来决定的，因为线程局部存储有一个被滥用的的倾向。  
Go 语言鼓励一种更简单的编程风格。其中，能影响一个函数行为的参数应当是显式指定的。  
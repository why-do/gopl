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

# 9.6 竞态检测器
即使在自己的检查，仍然可能在并发上犯错。Go 的 runtime 提供了动态分析工具：**竞态检测器**（race detectotr）。  

**开启竞态检测器**  
简单地把 -race 命令行参数加到 go build、go run、go test 命令里即可使用该功能。它会让编译器为你的应用或测试构建一个修改后的版本，这个版本有额外的手法可以高效记录在执行时对共享变量的所有访问，以及读写这些变量的 goroutine 标识。除此之外，还会记录所有的同步事件、包括 go 语句、通道操作、锁的调用等。（完整的同步事件集合可以在语言规范中的 “The Go Memory Model” 文档中找到。）  

**如何检查到竞态**  
竞态检测器会研究事件流，找到那些有问题的案例，即一个 goroutine 写入一个变量后，中间没有任何同步的操作，就有另外一个 goroutine 读写了该变量。这种情况表明有对共享变量的并发访问，即数据竞态。工具会输出一份报告，包括变量的标识以及读写 goroutine 当时的调用栈。通常情况下这些信息足以定位问题了。  

**哪些竞态可能查不到**  
竞态检测器报告所有实际运行了的数据竞态。但只能检测到那些在运行时发生的竞态，无法用来保证肯定不发生竞态。所以为了保证效果，需要全包测试包含了并发调用的场景。  

**可以在生产环境开启竞态检测器**  
由于存在额外的记录工作，带竞态检测功能的程序在执行时需要更长的时间和更多的内存，但即使对于生成环境的任务，这种额外开支也是可以接受的。对于那些偶发的竞态条件，使用竞态检测器可以节省很多调试的时间。  
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
# goroutine 和通道
Go 有两种并发编程风格：
+ goroutine 和通道（chennle），支持**通信顺序进程**（Communicating Sequential Process, CSP），CSP 是一个并发的模式，在不同的执行体（goroutine）之间传递值，但是变量本身局限于单一的执行体。
+ **共享内存多线程**的传统模型，和在其他主流语言中使用线程类似。


# 
TODO：ch5.md 捕获迭代变量
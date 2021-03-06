# 9.1 竞态
+ 练习9.1：向 gopl.io/ch9/bank1 程序添加一个函数 Withdraw(amount int)bool。结果应当反映交易成功还是由于余额不足而失败。函数发送到监控 goroutine 的消息应当包含取款金额和一个新的通道，这个通道用于监控 goroutine 把布尔型的结果发送回 Withdraw 函数。

# 9.5 延迟初始化：sync.Once
+ 练习9.2：重新 2.6.2 节的 PopCount 示例，使用 sync.Once 来把查找表的初始化延迟到第一次使用时。（从实际效果看，像 PopCount 这种既小又经高优化的函数无法承担同步的成本。）

# 9.7 示例：并发非阻塞缓存
+ 练习9.3：扩展 Func 类型和 (*Memo)\.Get 方法，让调用者可选择性地提供一个 done 通道，方便取消操作（参考 8.9 节）。不要缓存被取消的 Func 调用结果。

# 9.8 goroutine 与线程

## 9.8.1 可增长的栈
+ 练习9.4：使用通道构造一个把任意多个 goroutine 串联在一起的流水线程序。在内存耗尽之前你能创建的最大流水线级数是多少？一个值穿过整个流水线需要多久？

## 9.8.2 goroutine 调度
+ 练习9.5：写一个程序，两个 goroutine 通过两个无缓冲通道来互相转发消息。这个程序能维持每秒多少次通信？

## 9.8.3 GOMAXPROCS
+ 练习9.6：测量计算密集型并行程序（见练习8.5）在 GOMAXPROCS 参数变化时的性能变化。在你的计算机上最优值是多少？你的计算机有多少个CPU？
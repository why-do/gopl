# 9.1 竞态
+ 练习9.1：向 gopl.io/ch9/bank1 程序添加一个函数 Withdraw(amount int)bool。结果应当反映交易成功还是由于余额不足而失败。函数发送到监控 goroutine 的消息应当包含取款金额和一个新的通道，这个通道用于监控 goroutine 把布尔型的结果发送回 Withdraw 函数。
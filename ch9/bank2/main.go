// 这是一个只有一个账户的并发安全银行
package bank

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


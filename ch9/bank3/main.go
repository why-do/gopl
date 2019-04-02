// 这是一个只有一个账户的并发安全银行
package bank

import "sync"

var (
	mu      sync.Mutex // 保护 balance
	balance int
)

func Withdraw(amount int) bool {
	mu.Lock()
	defer mu.Unlock()
	deposit(-amount)
	if balance < 0 {
		deposit(amount)
		return false // 金额不足
	}
	return true
}

func Deposit(amount int) {
	mu.Lock()
	defer mu.Unlock()
	deposit(amount)
}

func Balance() int {
	mu.Lock()
	b := balance
	mu.Unlock()
	return b
}

// 这个函数要求已经获取互斥锁
func deposit(amount int) { balance += amount }

// 这是一个只有一个账户的并发安全银行
package bank

var deposits = make(chan int) // 发送存款额
var balances = make(chan int) // 接收余额

func Deposit(amount int) { deposits <- amount }
func Balance() int       { return <-balances }

var withdraw = make(chan int)    // 扣款通道
var withdrawOK = make(chan bool) // 扣款结果反馈通道
func Withdraw(amount int) bool {
	withdraw <- amount
	return <-withdrawOK
}

func teller() {
	var balance int // balance 被限制在 teller 这个 goroutine 中
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		case amount := <-withdraw:
			if balance >= amount {
				balance -= amount
				withdrawOK <- true
			} else {
				withdrawOK <- false
			}
		}
	}
}

func init() {
	go teller() // 启动监控 goroutine
}

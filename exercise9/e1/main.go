package main

import (
	"fmt"
	bank "gopl/exercise9/e1/bank1"
)

func main() {
	fmt.Println(bank.Balance())
	bank.Deposit(100)
	fmt.Println(bank.Balance())
	fmt.Println(bank.Withdraw(80))
	fmt.Println(bank.Balance())
	fmt.Println(bank.Withdraw(80))
	fmt.Println(bank.Balance())
}

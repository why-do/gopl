package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// pc[i] 是 i 的 population count
var pc [256]byte
var loadPcOnce sync.Once

// 这种既小又经高优化的函数，进行同步的成本反而高。就是多此一举，不过这是个练习。
func loadPc() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

// PopCount 返回 x 的 population count (number of set bits: 置位的个数)
func PopCount(x uint64) int {
	loadPcOnce.Do(loadPc)
	return int(pc[byte(x>>(0*8))] +
		pc[byte(x>>(1*8))] +
		pc[byte(x>>(2*8))] +
		pc[byte(x>>(3*8))] +
		pc[byte(x>>(4*8))] +
		pc[byte(x>>(5*8))] +
		pc[byte(x>>(6*8))] +
		pc[byte(x>>(7*8))])
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	x := rand.Uint64()
	fmt.Printf("%d %[1]b\n", x)
	var count int
	for _, i := range fmt.Sprintf("%b", x) {
		if i == '1' {
			count++
		}
	}
	fmt.Println(PopCount(x), count)
}

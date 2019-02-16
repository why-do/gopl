package popcount

// pc[i] 是 i 的 population count
var pc [256]byte

// 空间换时间：先生成一张表，之后这部分就可以不做计算，直接查表获取值。
func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

// PopCount 返回 x 的 population count (number of set bits: 置位的个数)
func PopCount(x uint64) int {
	var count int
	for i := uint(0); i < 8; i++ {
		count += int(pc[byte(x>>(i*8))])
	}
	return count
}
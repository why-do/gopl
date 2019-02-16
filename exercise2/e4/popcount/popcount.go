package popcount

// PopCount 返回 x 的 population count
// 数字小的情况下，效率高
func PopCount(x uint64) int {
	var count int
	for x != 0 {
		if x & 1 == 1 {
			count++
		}
		x >>= 1
	}
	return count
}

package basename

// 移除路径和扩展名
func basename(s string) string {
	// 从后往前找，把最后一个/以及之前的内容截掉
	for i := len(s) - 1; i > 0; i-- {
		if s[i] == '/' {
			s = s[i+1:]
			break
		}
	}
	// 从后往前找，把最后一个.指点的内容保留
	for i := len(s) -1; i > 0; i-- {
		if s[i] == '.' {
			s = s[:i]
			break
		}
	}
	return s
}
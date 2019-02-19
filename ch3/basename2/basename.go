package basename

import "strings"

func basename(s string) string {
	slash := strings.LastIndex(s, "/") // 如果没找到，则结果是-1
	s = s[slash+1:]  // 没找到返回-1,这样加1后就是全部字符串。如果找到了，加1后就是去掉查找字符及前面的字符串
	if dot := strings.LastIndex(s, "."); dot >=0 {
		s = s[:dot]
	}
	return s
}

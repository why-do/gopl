package main

import (
	"bytes"
	"fmt"
)

// 函数 intsToString 与 fmt.Sprintf(values) 类似，但插入了逗号
func intsToString(values []int) string {
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i, v := range values {
		if i > 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "%d", v)
	}
	buf.WriteByte(']')
	return buf.String()
}

func main() {
	fmt.Println(intsToString([]int{1, 2, 3}))  // "[1, 2, 3]"
	fmt.Println([]int{1, 2, 3})  // "[1 2 3]"
}
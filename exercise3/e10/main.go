package main

import (
	"bytes"
	"fmt"
	"unicode/utf8"
)

// 向十进制非负整数的字符串中插入逗号
func comma(s string) string {
	buff := bytes.NewBuffer(nil)
	n := utf8.RuneCountInString(s)  // 即使是非ASCII也能操作
	for _, v := range s {
		buff.WriteRune(v)
		n--
		if n%3 == 0 && n != 0 {
			buff.WriteByte(',')
		}
	}
	return buff.String()
}

var s1 string = "1234567"
var s2 string = "12345678"
var s3 string = "123456789"
var s4 string = "壹贰叁肆伍陆柒"
var s5 string = "壹贰叁肆伍陆柒捌"
var s6 string = "壹贰叁肆伍陆柒捌玖"

func main() {
	fmt.Println(comma(s1))
	fmt.Println(comma(s2))
	fmt.Println(comma(s3))
	fmt.Println(comma(s4))
	fmt.Println(comma(s5))
	fmt.Println(comma(s6))
}

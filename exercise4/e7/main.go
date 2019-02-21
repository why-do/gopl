package main

import (
	"fmt"
	"unicode/utf8"
)

// 逻辑简单，但是算法复杂度应该太高了，不过正好可以用来做验证
func reverse_byte(slice []byte) {
	for l := len(slice); l > 0; {
		r, size := utf8.DecodeRuneInString(string(slice[0:]))
		copy(slice[0:l], slice[0+size:l])
		copy(slice[l-size:l], []byte(string(r)))
		l -= size
	}
}

func reverse(s []byte) {
	var (
		lRd, rRd           int  // 读指针
		lWr, rWr           int  // 写指针
		lHasRune, rHasRune bool // 是否有字符
		lr, rr             rune // 读取到的字符
		lsize, rsize       int  // 读取到字符的宽度
	)
	rRd, rWr = len(s), len(s)
	for lRd < rRd {
		if !lHasRune {
			lr, lsize = utf8.DecodeRune(s[lRd:])
			lRd += lsize
			lHasRune = true
		}
		if !rHasRune {
			rr, rsize = utf8.DecodeLastRune(s[:rRd])
			rRd -= rsize
			rHasRune = true
		}

		if lsize <= rWr-rRd {
			utf8.EncodeRune(s[rWr-lsize:], lr)
			rWr -= lsize
			lHasRune = false
		}
		if rsize <= lRd-lWr {
			utf8.EncodeRune(s[lWr:], rr)
			lWr += rsize
			rHasRune = false
		}
	}

	// 最后还可能会剩个字符没写
	if lHasRune {
		utf8.EncodeRune(s[rWr-lsize:], lr)
	}
	if rHasRune {
		utf8.EncodeRune(s[lWr:], rr)
	}
}

func main() {
	str := "abc你好大家好1"
	s := []byte(str)
	fmt.Printf("%d %q\n", len(s), s)
	reverse(s)
	fmt.Printf("%d %q\n", len(s), s)

	s1 := []byte(str)
	reverse_byte(s1)
	fmt.Printf("%d %q\n", len(s), s1)
}

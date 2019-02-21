package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func repSpace(b []byte) []byte {
	const space = ' '
	var size int
	var r rune
	var n int
	var lastIsSpace bool
	for i := 0; i < len(b); i += size {
		r, size = utf8.DecodeRune(b[i:])
		if unicode.IsSpace(r) {
			if lastIsSpace {
				continue
			}
			b[n] = space
			n++
			lastIsSpace = true
		} else {
			utf8.EncodeRune(b[n:], r)
			n += size
			lastIsSpace = false
		}
	}
	return b[:n]
}

func main() {
	s := "abc  大家好   defg \u0085 \u00a0 \u00a1\u00a2ABC  \n DEF\r\t\r\tGH"
	s1 := repSpace([]byte(s))
	fmt.Printf("%q\n", s)
	fmt.Printf("%q\n", s1)
}

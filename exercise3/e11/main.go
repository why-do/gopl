package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// 向十进制非负整数的字符串中插入逗号
func comma(s string) (string, error) {
	if _, err := strconv.ParseFloat(s, 64); err != nil {
		return s, err
	}
	sl := strings.Split(s, ".")
	var s1, s2 string
	switch len(sl) {
	case 1:
		s1 = sl[0]
	case 2:
		s1, s2 = sl[0], sl[1]
	default:
		return s, fmt.Errorf("split error, len = %v want(2)\n", len(sl))
	}

	buff1 := bytes.NewBuffer(nil)
	n1 := len(s1)
	for _, v := range s1 {
		buff1.WriteRune(v)
		n1--
		if v == '+' || v == '-' {
			continue
		}
		if n1%3 == 0 && n1 > 0 {
			buff1.WriteByte(',')
		}
	}

	buff2 := bytes.NewBuffer(nil)
	var n2 int
	for _, v := range s2 {
		buff2.WriteRune(v)
		n2++
		if n2%3 == 0 && n2 < len(s2) {
			buff2.WriteByte(',')
		}
	}

	return strings.Join([]string{buff1.String(), buff2.String()}, "."), nil
}

// var s string = "1234567.987654"
var s1 string = "123.4567.987.654"
var s2 string = "1234567.987654"
var s3 string = "+123456.987654"
var s4 string = "-123456.987654"

func main() {
	var s string
	var err error

	s, err = comma(s1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
	} else {
		fmt.Println(s)
	}

	s, err = comma(s2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
	} else {
		fmt.Println(s)
	}

	s, err = comma(s3)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
	} else {
		fmt.Println(s)
	}

	s, err = comma(s4)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
	} else {
		fmt.Println(s)
	}
}

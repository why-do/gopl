package main

import (
	"fmt"
	"sort"
)

type PalindromeString string

func (s PalindromeString) Len() int           { return len(s) }
func (s PalindromeString) Less(i, j int) bool { return s[i] < s[j] }
func (s PalindromeString) Swap(i, j int)      {} // 字符串不能修改的

func IsPalindrome(s sort.Interface) bool {
	if s.Len() == 0 {
		return false
	}
	i, j := 0, s.Len()-1
	for i < j {
		if !s.Less(i, j) && !s.Less(j, i) {
			i++
			j--
		} else {
			return false
		}
	}
	return true
}

func main() {
	fmt.Println(IsPalindrome(PalindromeString("")))
	fmt.Println(IsPalindrome(PalindromeString("12321")))
	fmt.Println(IsPalindrome(PalindromeString("abcdcba")))
	fmt.Println(IsPalindrome(sort.IntSlice([]int{1, 2, 3, 4, 5})))
	fmt.Println(IsPalindrome(sort.IntSlice([]int{1, 2, 3, 2, 1})))
	fmt.Println(IsPalindrome(sort.StringSlice([]string{"", "", ""})))
	fmt.Println(IsPalindrome(sort.StringSlice([]string{"你好", "世界", "你好"})))
}

// word 包提供了文字游戏相关的工具函数
package word

import "unicode"

// IsPalindrome 判断一个字符串是否是回文
func IsPalindrome(s string) bool {
	var letters []rune
	for _, r := range s {
		if unicode.IsLetter(r) {
			letters = append(letters, unicode.ToLower(r))
		}
	}
	for i := range letters {
		if letters[i] != letters[len(letters)-1-i] {
			return false
		}
	}
	return true
}

// IsPalindrome 的优化版本
func IsPalindrome2(s string) bool {
	// var letters []rune
	letters := make([]rune, 0, len(s))
	for _, r := range s {
		if unicode.IsLetter(r) {
			letters = append(letters, unicode.ToLower(r))
		}
	}
	n := len(letters)
	for i := 0; i < n; i++ {
		if letters[i] != letters[len(letters)-1-i] {
			return false
		}
	}
	return true
}
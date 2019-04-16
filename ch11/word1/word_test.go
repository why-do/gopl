package word

import "testing"

func TestPalindrome(t *testing.T) {
	if !IsPalindrome("civic") {
		t.Error(`IsPalindrome("civic") = false`)
	}
	if !IsPalindrome("madam") {
		t.Error(`IsPalindrome("madam") = false`)
	}
}

func TestNonPalindrome(t *testing.T) {
	if IsPalindrome("palindrome") {
		t.Error(`IsPalindrome("palindrome") = true`)
	}
}

func TestChinesePalindrome(t *testing.T) {
	input := "上海自来水来自海上"
	if !IsPalindrome(input) {
		t.Errorf(`IsPalindrome(%q) = false`, input)
	}
}

func TestSentencePalindrome(t *testing.T) {
	input := "Madam, I'm Adam"
	if !IsPalindrome(input) {
		t.Errorf(`IsPalindrome(%q) = false`, input)
	}
}
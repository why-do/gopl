package main

import (
	"math/rand"
	"testing"
	"time"
)

var tests = []struct {
	input string
	want  string
}{
	{"abc", "cba"},
	{"123", "321"},
	{"你好，世界!", "!界世，好你"},
	{"a一二三,四五.六,z", "z,六.五四,三二一a"},
}

func TestReverse_rune(t *testing.T) {
	for _, test := range tests {
		s := []byte(test.input)
		reverse_rune(s)
		if string(s) != test.want {
			t.Errorf("reverse(%q) = %q, want %q\n", test.input, string(s), test.want)
		}
	}
}

func TestReverse_byte(t *testing.T) {
	for _, test := range tests {
		s := []byte(test.input)
		reverse_byte(s)
		if string(s) != test.want {
			t.Errorf("reverse(%q) = %q, want %q\n", test.input, string(s), test.want)
		}
	}
}

func TestReverse(t *testing.T) {
	for _, test := range tests {
		s := []byte(test.input)
		reverse(s)
		if string(s) != test.want {
			t.Errorf("reverse(%q) = %q, want %q\n", test.input, string(s), test.want)
		}
	}
}

// randomSte 返回一个随机字符串，它的长度和内容都是随机生成的
func randomStr(rng *rand.Rand) string {
	n := rng.Intn(25) // 随机字符串最大长度24
	runes := make([]rune, n)
	for i := 0; i < n; i++ {
		var r rune
		switch rune(rng.Intn(6)) {
		case 0: // ASCII 字母，1个字节
			r = rune(rng.Intn(0x4B) + 0x30)
		case 1: // 希腊字母，2个字节
			r = rune(rng.Intn(57) + 0x391)
		case 2: // 日文
			r = rune(rng.Intn(0xBF) + 0x3041)
		case 3: // 韩文
			r = rune(rng.Intn(0x2BA4) + 0xAC00)
		case 4, 5, 6: // 中文
			r = rune(rng.Intn(0x4E00) + 0x51D6)
		}
		runes[i] = r
	}
	return string(runes)
}

func TestRandomReverse(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	for i := 0; i < 1000; i++ {
		test := randomStr(rng)
		s1 := []byte(test)
		reverse_rune(s1)
		t.Logf("%s => %s\n", test, string(s1))

		s2 := []byte(test)
		reverse_byte(s2)
		if string(s1) != string(s2) {
			t.Errorf("reverse_byte(%q) = %q, want %q\n", test, string(s2), string(s1))
		}

		s3 := []byte(test)
		reverse(s3)
		if string(s1) != string(s3) {
			t.Errorf("reverse_byte(%q) = %q, want %q\n", test, string(s3), string(s1))
		}
	}
}

// 效率测试
func BenchmarkReverse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, test := range tests {
			reverse([]byte(test.input))
		}
	}
}

func BenchmarkReverse_rune(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, test := range tests {
			reverse_rune([]byte(test.input))
		}	}
}

func BenchmarkReverse_byte(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, test := range tests {
			reverse_byte([]byte(test.input))
		}	}
}

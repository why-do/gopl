# 反转 UTF-8 编码的字符串
练习4.7  
实现一个 reverse 函数，反转 UTF-8 编码的字符串中的字符元素。传入的参数是字符串对应的字节切片类型（[]byte）。

## 简单的实现
首先，不考虑效率，先用一个简单的逻辑来实现。切片的反转方法如下：
```go
func reverse(s []int) {
    for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
        s[i], s[j] = s[j], s[i]
    }
}
```

只要将数据转成字符切片，然后套用切片反转的代码就可以了：
```go
// 先转成字符切片，然后再用切片的反转的方法，最后更新参数指向的底层数组
func reverse_rune(slice []byte) {
	r := []rune(string(slice))
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
        r[i], r[j] = r[j], r[i]
	}
	for i := range slice {
		slice[i] = []byte(string(r))[i]
	}
}
```
这个方法逻辑清晰，适合之后拿来做随机测试

## 原地反转
下面的函数实现了原地反转。每次读取第一个字符，把尾部标志位之前的一串字符移到开头，再把之前读的字符放到标志位的位置，然后向前移动标志位：
```go
func reverse_byte(slice []byte) {
	for l := len(slice); l > 0; {
		r, size := utf8.DecodeRuneInString(string(slice[0:]))
		copy(slice[0:l], slice[0+size:l])
		copy(slice[l-size:l], []byte(string(r)))
		l -= size
	}
}
```
这个实现效率还是稍差，之后会做测试比较。  

## 高效的原地反转
下面的函数，用了很多的标志位，实现了一个高效的原地反转：
```go
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
```

# 测试验证
下面是测试代码，来验证上面的函数的正确性以及效率。  

## 功能测试
基于表的测试很直观也很简单，可以方便的添加更多测试用例：
```go
var tests = []struct {
	input string
	want string
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
```

测试结果：
```
PS H:\Go\src\gopl\exercise4\e7> go test -run TestReverse -v
=== RUN   TestReverse_rune
--- PASS: TestReverse_rune (0.00s)
=== RUN   TestReverse_byte
--- PASS: TestReverse_byte (0.00s)
=== RUN   TestReverse
--- PASS: TestReverse (0.00s)
PASS
ok      gopl/exercise4/e7       0.263s
PS H:\Go\src\gopl\exercise4\e7>
```

## 随机测试
随机测试也是功能测试的一种，通过构建随机输入来扩展测试的覆盖范围。有两种策略：
+ 额外写一个函数，这个函数使用低效但是清晰的算法，然后检查两种实现的输出是否一致
+ 构建符合某种模式的输入，这样就可以知道对应的输出。

下面就是用第一种策略写的随机测试。为了让输出的内容有更好的可读性，选择了一些熟悉的字符生成随机字符：
```go
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
```

测试结果：
```
PS H:\Go\src\gopl\exercise4\e7> go test -run Random
PASS
ok      gopl/exercise4/e7       0.298s
PS H:\Go\src\gopl\exercise4\e7>
```
还可以加上 \-v 参数，查看详细的测试日志。  

## 基准测试
基准测试也没什么特别的，把功能测试的测试用例全部跑一遍：
```go
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
```

测试结果：
```
PS H:\Go\src\gopl\exercise4\e7> go test -benchmem -bench .
goos: windows
goarch: amd64
pkg: gopl/exercise4/e7
BenchmarkReverse-8               5000000               286 ns/op               0 B/op          0 allocs/op
BenchmarkReverse_rune-8           500000              3610 ns/op               0 B/op          0 allocs/op
BenchmarkReverse_byte-8          3000000               583 ns/op               0 B/op          0 allocs/op
PASS
ok      gopl/exercise4/e7       6.226s
PS H:\Go\src\gopl\exercise4\e7>
```
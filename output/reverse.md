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
下面的函数实现了原地反转。每次将一个字符移动到末尾，效率比较差
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

## 高效的原地反转
下面的函数，用了很多的标志位，实现了一个高效的原地反转：
```go

```

# 测试验证

## 功能测试

## 随机测试

## 基准测试
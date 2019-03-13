# 7.1 接口即约定
+ 练习7.1：使用类似 ByteCounter 的想法，实现单词和行的计数器。实现时考虑使用 bufio.ScanWords。
+ 练习7.2：实现一个满足如下签名的 CountingWriter 函数，输入一个 io.Writer，输出一个封装了输入值的新 Writer，以及一个指向 int64 的指针，该指针对应的值是新的 Writer 写入的字节数。
```go
func CountingWriter(w io.Writer) (io.Writer, *int64)
```
+ 练习7.3：为 gopl.io/ch4/treesort 中的 *tree 类型（见4.4节）写一个String方法，用于展示其中的值序列。
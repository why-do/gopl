# 7.1 接口即约定
+ 练习7.1：使用类似 ByteCounter 的想法，实现单词和行的计数器。实现时考虑使用 bufio.ScanWords。
+ 练习7.2：实现一个满足如下签名的 CountingWriter 函数，输入一个 io.Writer，输出一个封装了输入值的新 Writer，以及一个指向 int64 的指针，该指针对应的值是新的 Writer 写入的字节数。
```go
func CountingWriter(w io.Writer) (io.Writer, *int64)
```
+ 练习7.3：为 gopl.io/ch4/treesort 中的 *tree 类型（见4.4节）写一个String方法，用于展示其中的值序列。

# 7.2 接口类型
+ 练习7.4：strings.NewReader 函数输入一个字符串，返回一个从字符串读取数据且满足 io.Reader 接口（也满足其他接口）的值。请自己实现该函数，并且通过它来让 HTML 分析器（参考5.2节）支持以字符串作为输入。
+ 练习7.5：io 包中的 LimitReader 函数接受 io.Reader r和字节数n，返回一个 Reader，该返回值从 r 读取数据，但在读取 n 字节后报告文件结束。请实现该函数。
```go
func LimitReader(r io.Reader, n int64) io.Reader
```

# 7.4 使用 flag.Value 来解析参数
+ 练习7.6：在 tempflag 中支持热力学温度。
+ 练习7.7：请解释为什么默认值 20。0 没写 °C，而帮助消息中却包含 °C。

# 7.6 使用 sort.Interface 来排序
+ 练习7.8：很多图形界面提供了一个表格控件，它支持有状态的多层排序：先按照最近单击的列来排序，接着是上一次使用单击的列，依次类推。请定义 sort.Interface 接口实现来满足如上需求。试比较这个方法与多次使用 sort.Stable 排序的异同。
+ 练习7.9：利用 html/template 包（见4.6节）来替换 printTracks 函数，使用 HTML 表格来显示音乐列表。结合上一个练习，来实现通过单击列头来发送 HTTP 请求，进而对表格排序。
+ 练习7.10：sort.Interface 也可以用于其他用途。试写一个函数 `IsPalindrome(s sort.Interface)bool` 来判断一个序列是否是回文，即序列反转后是否保持不变。可以假定对于下标分别为 i、j 的元素，如果 `!s.Less(i,j) && !s.Less(j,i)`，那么两个元素相等。 
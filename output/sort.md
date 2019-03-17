# 7.6 使用 sort.Interface 来排序
排序是一个在很多程序中广泛使用的操作。sort 包提供了针对任意序列根据任意排序函数原地排序的功能。  
这样的设计号称并不常见。在很多语言中，排序算法跟序列数据类型绑定，排序函数跟序列元素类型绑定。但 Go 语言的 sort.Sort 函数对序列和其中元素的布局无任何要求，它使用 sort.Interface 接口来实现。

## 接口实现
一个原地排序算法需要知道三个信息：
1. 序列长度
2. 比较两个元素的含义
3. 如何交换两个元素

所以 sort.Interface 接口就有三个方法：
```go
package sort

type Interface interface {
	Len() int
	Less(i, j int) bool
	Swap(i, j int)
}
```

## 字符串排序
要对序列排序，需要先确定一个实现了上面三个方法的类型，接着把 sort.Sort 函数应用到上面这类方法的示例上。以字符串切片 []string 为例，定义一个新的类型 StringSlice 以及它的三个方法如下：
```go
type StringSlice []string

func (p StringSlice) Len() int           { return len(p) }
func (p StringSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p StringSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
```
现在要对一个字符串切片进行排序，先把类型转换为 StringSlice 类型，然后调用 sort.Sort 函数即可：
```go
sort.Sort(StringSlice(names))
```
字符串切片的排序太常见了，所以 sort 包提供了 StringSLice 类型，以及一个直接排序的 Strings 函数。上面对 StringSlice 类型及其方法的定义就是源码里实现的代码。所以要对字符串切片进行排序，直接调用 srot.Strings 函数即可：
```go
sort.Strings(names)
```
为了简便，sort 包专门对 \[\]int、\[\]string、\[\]float64 这三个类型提供了排序的函数和相关类型。对于其他类型，就需要自己写了，不过写起来也不复杂。

## 反转 Reverse
sort.Reverse 函数值得仔细看一下，它使用了一个重要的概念**组合**。sort 包定义了一个 reverse 类型，这个类型是一个嵌入了 sort.Interface 的结构。reverse 的 Less 方法，直接调用内嵌的 sort.Interface 值的 Less 方法，但是调换了下标，这样就实现了颠倒排序的结果了：
```go
package sort

type reverse struct { Interface } // 这个是在sort包里的，所以就是 sort.Interface
func (r reverse) Less(i, j int) bool { return r.Interface.Less(j, i) }  // 这里调换了函数体中 i 和 j 的位置
func Reverse(data Interface) Interface { return &reverse{data} }
```
reverse 的另外两个方法 Len 和 Swap，没有定义，就会由内嵌的 sort.Interface 隐式提供。导出函数 Reverse 则返回一个包含原始的 sort.Interface 值的 reverse 实例。最终反转排序的调用如下，先做类型转换，然后加一步通过 Reverse 函数生成 reverse 实例，最后就是调用 sort.Sort 函数完成排序：
```go
sort.Sort(sort.Reverse(StringSlice(names)))
```

# 复杂类型的排序
对于一个复杂的类型，比如结构体，它会有多个字段，也就是可能需要对多个维度进行排序。

## 数据和结构
下面是一个复杂的结构体类型，并且数据也已经准备好了：
```go
type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

var tracks = []*Track{
	{"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
	{"Go", "Moby", "Moby", 1992, length("3m37s")},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, length("4m24s")},
}

func length(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(s)
	}
	return d
}
```

## 用表格输出结构体
这里使用 text/tabwriter 包可以生成一个干净整洁的表格。这里注意，\*tabwriter.Writer 满足 io.Writer 接口，先用它收集所有写入的数据，在 Flush 方法调用时才会格式化整个表格并且输出：
```go
func printTracks(tracks []*Track) {
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Title", "Artist", "Album", "Year", "Length")
	fmt.Fprintf(tw, format, "-----", "------", "-----", "----", "------")
	for _, t := range tracks {
		fmt.Fprintf(tw, format, t.Title, t.Artist, t.Album, t.Year, t.Length)
	}
	tw.Flush()
}
```

## 排序
按照 Artist 字段进行排序：
```go
type byArtist []*Track

func (x byArtist) Len() int           { return len(x) }
func (x byArtist) Less(i, j int) bool { return x[i].Artist < x[j].Artist }
func (x byArtist) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
```

按照 Year 字段进行排序：
```go
type byYear []*Track

func (x byYear) Len() int           { return len(x) }
func (x byYear) Less(i, j int) bool { return x[i].Year < x[j].Year }
func (x byYear) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
```

## 优化
上面分别对2个字段定义了排序。像这样，对于每一个排序都需要实现一个新的 sort.Interface。但是其实只有 Less 方法不同，而 Len 和 Swap 方法都是一样的。  
在下面的例子中，具体类型 customSort 组合了一个要排序的切片类型和一个函数，这样每次只要写一个比较函数就可以定义一个新的排序。从这个例子也可以看到，实现 sort.Interface 的具体类型并不一定都是切片，比如这里的 customSort 就是一个结构体：
```go
type customSort struct {
	t    []*Track
	less func(x, y *Track) bool
}

func (x customSort) Len() int           { return len(x.t) }
func (x customSort) Less(i, j int) bool { return x.less(x.t[i], x.t[j]) }
func (x customSort) Swap(i, j int)      { x.t[i], x.t[j] = x.t[j], x.t[i] }
```
接下来就定义一个多层的比较函数，先对 Title 排序，然后对 Year，最后是时长：
```go
sort.Sort(customSort{tracks, func(x, y *Track) bool {
	if x.Title != y.Title {
		return x.Title < y.Title
	}
	if x.Year != y.Year {
		return x.Year < y.Year
	}
	if x.Length != y.Length {
		return x.Length < y.Length
	}
	return false
}})
```

## 检查排序
一个长度为 n 的序列进行排序需要 O(n logn) 次比较操作，而判断一个序列是否已经排好序值需要最多 (n-1) 次比较。sort 包提供的 IsSorted 函数可以判断序列是否是排好序的。
```go
values := []int{3, 1, 4, 1}
fmt.Println(sort.IntsAreSorted(values))                         // "false"
sort.Ints(values)                                               // 排序
fmt.Println(values)                                             // "[1 1 3 4]"
fmt.Println(sort.IntsAreSorted(values))                         // "true"
sort.Sort(sort.Reverse(sort.IntSlice(values)))                  // 反转
fmt.Println(values)                                             // "[4 3 1 1]"
fmt.Println(sort.IntsAreSorted(values))                         // "false"
fmt.Println(sort.IsSorted(sort.Reverse(sort.IntSlice(values)))) // "true"
```

## 完整示例代码
上面的代码的源码文件，可以运行：
```go
// ch7/sorting
```
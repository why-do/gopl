# 读写操作

## bufio包，逐行读取（ch1.md）
关于bufio包，使用它可以简便和高效地处理输入和输出。其中一个最有用的特性是称为扫描器（Scanner）的类型，它可以读取输入，以行或者单词为单位断开，这是处理以行为单位输入内容的最简单方式。  
先声明一个 bufio.Scanner 类型的变量：
```go
input := bufio.NewScanner(os.Stdin)
```
扫描器从程序的标准输入进行读取。每一次调用 inout.Scan() 读取下一行，并且将结尾的换行符去掉；通过调用 input.Text() 来获取读到的内容。Scan 方法在读到新行的时候返回true，在没有更多内容的时候返回flase。  
下面的例子，从标准输入获取到内容后，转成全大写再打印出来。Windows下使用 Ctrl+Z 后回车可以退出：
```go
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		fmt.Println(strings.ToUpper(input.Text()))
	}
}
```
上面的示例是每从标准输入获取到一行的内容，就进行处理。还有一种做法是，从标准输入获取全部的内容后（Ctrl+Z 后回车表示输入完成），最后再一次全部输出。下面的示例是从标准输入获取到全部内容，然后全部输出，每一行的内容之间插入一个空格。最后也会有一个空格，输出是把最后一个空格截断：
```go
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

func main() {
	input := bufio.NewScanner(os.Stdin)
	buff := bytes.NewBuffer(nil)
	for input.Scan() {
		buff.Write(input.Bytes()) // 忽略错误
		buff.WriteByte(' ')
	}
	fmt.Printf("%q\n", buff.String()[:buff.Len()-1])  // 截掉最后一个空格
}
```

## 读取内容并丢弃（ch1.md）
使用io.Copy函数读取响应的内容，比如直接复制内容到标准输出，这样就不需要把数据流装到缓冲区：
```go
n, err := io.Copy(os.Stdout, resp.Body)
```
还可以通过写入 ioutil.Discard 输出流进行丢弃，这样做应该是为了要有一个读取的过程：
```go
n, err := io.Copy(ioutil.Discard, resp.Body)
```

# 向 io\.Writer 写入字符串(7.12)
这里要讲的是通过类型断言来查询特性，下面是一个标准库中使用的示例。并且，**这个是向 io\.Writer 写入字符串的推荐方法。**  
下面定义一个方法，往 io\.Writer 接口接入字符串信息：
```go
func writeMsg(w io.Writer, msg string) error {
	if _, err := w.Write([]byte(msg)); err != nil {
		return err
	}
	return nil
}
```
因为 Write 方法需要一个字节切片（[]byte），而需要写入的是一个字符串，所以要做类型转换。这种转换需要进行内存分配和内存复制，但复制后内存又会被马上抛弃。这里就会有性能问题，这个内存分配会导致性能下降，需要避开这个内存分配。  
在很多包中，实现了 io\.Writer 的重要类，都会提供一个对应的高效写入字符串的 WriteString 方法，比如：
+ \*http\.response ： 这个类型不可导出，我们一般使用的是 http\.ResponseWriter 接口
+ \*bytes\.Buffer
+ \*os\.File
+ \*bufio\.Write

由于 Writer 接口并不包括 WriteString 方法，不能直接调用。这里可以先定义一个新的接口，这个接口只包含 WriteString 方法，然后使用类型断言来判断 w 的动态类型是否满足这个新接口：
```go
// 将s写入w，如果w有WriteString方法，就直接调用该方法
func writeString(w io.Writer, s string) (n int, err error) {
	type stringWriter interface {
		WriteString(string) (n int, err error)
	}
	if sw, ok := w.(stringWriter); ok {
		return sw.WriteString(s) // 避免了内存复制
	}
	return w.Write([]byte(s)) // 分配了临时内存
}

func writeMsg(w io.Writer, msg string) error {
	if _, err := writeString(w, msg); err != nil {
		return err
	}
	return nil
}
```

由于上面的操作太常见了，io 包已经提供了一个函数 io\.WriteString 可以直接使用。上面的代码可以简化为如下的方式：
```go
func writeMsg2(w io.Writer, msg string) error {
	if _, err := io.WriteString(w, msg); err != nil {
		return err
	}
	return nil
}
```
这里本质上还是一样的，接口的定义、类型检查都封装到了io包的内部，并且内部实现的逻辑和上面是一样的。  
*之前想要向某个具体的类型写入字符串的时候，会查看该类型是否有 WriteString 方法。但是如果要操作的类型是 io\.Writer 接口，虽然实际背后的动态类型还是那个类型，但是就无法调用 WriteString 方法了。这时可以直接使用 io 包提供的工具函数完成同样的操作。*  
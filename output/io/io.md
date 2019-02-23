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

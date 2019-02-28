# 错误
简单的错误处理是使用 Fprintf 和 %v 在标准错误流上输出一条消息，%v 可以使用默认格式显示任意类型的值。  
为了保持示例代码简短，有时会对错误处理有意进行一定程度的忽略。明显的错误还是要处理的。但是有些出现概率很小的错误，就忽略了，不过要标记所跳过的错误检查，就是加上注释。  

根据情形，将有许多可能的处理场景，接下来是5个例子。  

## 将错误传递下去
最常见的情形是将错误传递下去，使得在子例程中发生的错误变为主调例程的错误。  
一种是不做任何操作立即向调用者返回错误：
```go
resp, err := http.Get(url)
if err != nil {
    return nil, err
}
```
还有一种，不会直接返回，因为错误信息中缺失一些关键信息：
```go
doc, err := html.Parse(resp.Body)
resp.Body.Close()
if err != nil {
    return nil, fmt.Errorf("parsing %s as HTML: %v\n", url, err)
}
```
这里格式化了一条错误消息并且返回一个新的错误值。可以为原始的错误消息不断地添加上下文信息来建立一个可读的错误描述。当错误最终被程序的 main 函数处理时，它应该能够提供一个从最根本问题到总体故障的清晰因果链、这里有一个 NASA 的事故调查的例子：
```
genesis: crashed: no parachute: G-switch failed: bad relay orientation
```
因为错误频繁地串联起来，所以消息字符串首字母不应该大写而且应该避免换行。错误结果可能会很长，但能能够使用 grep 这样的工具找到需要的信息。  

**需要添加的关键信息**  
有时候可以不用添加信息直接返回，有时候需要添加一些关键信息，因为错误信息里没有。比如 os.Open 打开文件时，返回的错误不仅仅包括错误的信息，还包含文件的名字，因此调用者构造错误消息的时候不需要包含文件的名字这类信息。具体哪些信息是缺少的关键信息需要在原始的错误消息的基础上添加？  
一般地，f(x) 调用只负责报告函数的行为 f 和参数值 x，因为它们和错误的上下文相关。调用者则负责添加进一步的信息，但是 f(x) 本身并不会，并且在函数内部也没有这些信息。  
比如上面的 html.Parse 返回的错误信息里不可能有 url 的信息，但是，是关键信息需要添加。而 os.Open 中，文件名字也是关键信息，但是这个正是函数的参数值，所以函数本身会返回这个信息，不需要另外添加。  

## 尝试重试
对于不固定或者不可预测的错误，在短暂的间隔后对操作进行重试是合乎情理的。超出一定的重试次数和限定的时间后再报错退出。  
下面给出了完整的代码，暂时只看 WaitForServer 函数：
```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// 尝试连接 url 对应的服务器
// 在一分钟内使用指数退避策略进行重试
// 所有的尝试失败后返回错误
func WaitForServer(url string) error {
	const timeout = 1 * time.Minute
	deadline := time.Now().Add(timeout)
	for tries := 0; time.Now().Before(deadline); tries++ {
		_, err := http.Head(url)
		if err == nil {
			return nil // 成功
		}
		log.Printf("server not responding (%s); retrying...", err)
		time.Sleep(time.Second << uint(tries)) // 指数退避策略
	}
	return fmt.Errorf("server %s failed to respond after %s", url, timeout)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "需要提供 url 参数\n")
		os.Exit(1)
	}
	url := os.Args[1]
	if err := WaitForServer(url); err != nil {
		fmt.Fprintf(os.Stderr, "Site is down: %v\n", err)
		os.Exit(1)
	}
}
```
这里的**指数退避策略**，已经尝试多次简单的超时退出的实现也很有意思。  

## 输出日志并退出
接着看上面的代码，如果多次重试后依然不能成功，调用者能够输出错误然后优雅地停止程序，但一般这样的处理应该留给主程序部分：
```go
if err := WaitForServer(url); err != nil {
    fmt.Fprintf(os.Stderr, "Site is down: %v\n", err)
    os.Exit(1)
}
```
通常，如果是库函数，应该将错误传递给调用者，除非这个错误表示一个内部的一致性错误，这意味着库内部存在 bug。  
这里还有一个更加方便的方法是通过调用 log.Fatalf 实现上面相同的效果。和所有的日志函数一样，它默认会将时间和日期作为前缀添加到错误消息前：
```go
if err := WaitForServer(url); err != nil {
    log.Fatalf(os.Stderr, "Site is down: %v\n", err)
}
```
这种带日期时间的默认格式有助于长期运行的服务器，而对于交互式的命令行工具则意义不大。  
还可以自定义命令的名称作为 log 包的前缀，并且将日期和时间略去：
```go
log.SetPrefix("waid: ")
log.SetFlags(0)
```

## 记录log日志
在一些错误情况下，只记录下错误信息然后程序继续运行。同样地，可以选择使用 log 包来增加日志的常用前缀：
```go
if err := Ping(): err != nil {
    log.Printf("Ping failed: %v; networking disabled", err)
}
```
所有 log 函数都会为缺少换行符的日志补充一个换行符。  
或者是，直接输出到标准错误流：
```go
if err := Ping(): err != nil {
    fmt.Fprintf(os.Stderr, "Ping failed: %v; networking disabled\n", err)
}
```
没有用 log 函数，所以没有时间日期，当然也不需要。上面说了，对于交互式的命令工具意义不大。  

## 忽略错误
在某些罕见的情况下，还可以直接安全地忽略掉整个日志：
```go
dir, err := ioutil.TempDir("", "scratch")
if err != nil {
    return fmt.Errorf("failed to create temp dir: %v", err)
}
// 使用临时的目录
os.RemoveAll(dir)  // 忽略错误，$TMPDIR 会被周期性删除
```
调用 os.RemoveAll 可能会失败，但程序忽略了这个错误，原因是操作系统会周期性地清理临时目录。在这个例子中，有意的抛弃了错误，但程序的逻辑看上去就和忘记去处理一样了。要习惯考虑到每一个函数调用可能发生的出错情况，当有意忽略一个错误的时候，要清楚地注释一下你的意图。  


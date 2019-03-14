# 7.4 使用 flag.Value 来解析参数

## 解析时间
首先看示例，这里实现了暂停指定时间的功能：
```go
var period = flag.Duration("period", 1*time.Second, "sleep period")

func main() {
	flag.Parse()
	fmt.Printf("Sleeping for %v...", *period)
	time.Sleep(*period)
	fmt.Println()
}
```
默认是1秒，但是可以通过参数来控制。flag.Duration函数创建了一个 \*time.Duration 类型的标志变量，并且允许用户用一种友好的方式来指定时长。就是用 String 方法对应的记录方法。这种对称的设计提供了一个良好的用户接口。
```
PS H:\Go\src\gopl\ch7\sleep> go run main.go -period 3s
Sleeping for 3s...
PS H:\Go\src\gopl\ch7\sleep> go run main.go -period 1m
Sleeping for 1m0s...
PS H:\Go\src\gopl\ch7\sleep> go run main.go -period 1.5h
Sleeping for 1h30m0s...
```
因为时间长度类的命令行标志广泛应用，所以这个功能内置到了 flag 包中。


## 备忘录
这么做，是为了可以使用 Celsius 的 String 方法
```go
type celsiusFlag struct{ Celsius }
```
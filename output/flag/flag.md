# 命令行参数
命令行参数可以直接通过 os.Args 获取，另外标准库的 flag 包专门用于接收和解除命令行参数

## os.Args
简单的只是从命令行获取一个或一组参数，可以直接使用 os.Args。下面的这种写法，无需进行判断，无论是否提供了命令行参数，或者提供了多个，都可以处理：
```go
// 把命令行参数，依次打印，每行一个
func main() {
	for _, s := range os.Args[1:] {
		fmt.Println(s)
	}
}
```

## flag 基本使用
下面的例子使用了两种形式的调用方法：
```go
package main

import (
	"flag"
	"fmt"
)

var name string

func init() {
	flag.StringVar(&name, "name", "Adam", "名字")
}

var ageP = flag.Int("age", 18, "年龄")

func main() {
	flag.Parse()
	fmt.Printf("%T %[1]v\n", name)
	fmt.Printf("%T %[1]v\n", ageP)
	fmt.Printf("%T %[1]v\n", *ageP)
}
```
第一种是直接把变量的指针传递给函数作为第一个参数，函数内部会对该变量进行赋值。这种形式必须写在一个函数体的内部。  
第二种是函数会把数据的指针作为函数的返回值返回，这种形式就是给变量赋值，不需要现在函数体内，不过拿到的返回值是指针。  

## 解析时间（7.4 使用 flag.Value 来解析参数）
时间长度类的命令行标志应用广泛，这个功能内置到了 flag 包中。  
先看看源码中的示例，之后在自定义命令行标志的时候也能有个参考。下面的示例，实现了暂停指定时间的功能：
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

# 自定义类型
更多的情况下，是需要自己实现接口来进行自定义的。

## 接口说明
支持自定义类型，需要定义一个满足 flag.Value 接口的类型：
```go
package flag

// Value 接口代表了存储在标志内的值
type Value interface {
	String() string
	Set(string) error
}
```
String 方法用于格式化标志对应的值，可用于输出命令行帮助消息。  
Set 方法解析了传入的字符串参数并更新标志值。可以认为 Set 方法是 String 方法的逆操作，这两个方法使用同样的记法规格是一个很好的实践。  

## 自定义温度解析
下面定义 celsiusFlag 类型来允许在参数中使用摄氏温度或华氏温度。因为 Celsius 类型原本就已经实现了 String 方法，这里把 Celsius 内嵌到了 celsiusFlag 结构体中，这样结构体有就有了 String 方法（*外围结构体类型不仅获取了匿名成员的内部变量，还有相关方法*）。所以为了满足接口，只须再定一个 Set 方法：
```go
type Celsius float64
type Fahrenheit float64

func CToF(c Celsius) Fahrenheit { return Fahrenheit(c*9.0/5.0 + 32.0) }
func FToC(f Fahrenheit) Celsius { return Celsius((f - 32.0) * 5.0 / 9.0) }

func (c Celsius) String() string { return fmt.Sprintf("%g°C", c) }
// 上面这些都是之前在别处定义过的内容，是可以作为包引出过来的
// 为了说明清楚，就单独把需要用到的部分复制过来

// *celsiusFlag 满足 flag.Vulue 接口
type celsiusFlag struct{ Celsius }

func (f *celsiusFlag) Set(s string) error {
	var unit string
	var value float64
	fmt.Sscanf(s, "%f%s", &value, &unit) // 无须检查错误
	switch unit {
	case "C", "°C":
		f.Celsius = Celsius(value)
		return nil
	case "F", "°F":
		f.Celsius = FToC(Fahrenheit(value))
		return nil
	}
	return fmt.Errorf("invalid temperature %q", s)
}
```
fmt.Sscanf 函数用于从输入 s 解析一个浮点值和一个字符串。通常是需要检查错误的，但是这里如果出错，后面的 switch 里的条件也是无法满足的，是可以通过switch之后的错误处理来一并进行处理的。  
这里还需要写一个 CelsiusFlag 函数来封装上面的逻辑。这个函数返回了一个 Celsius 的指针，它指向嵌入在 celsiusFlag 变量 f 中的一个字段。Celsius 字段在标志处理过程中会发生变化（经由Set
方法）。调用 Var 方法可以把这个标志加入到程序的命令行标记集合中，即全局变量 flag.CommandLine。*如果一个程序有非常复杂的命令行接口，那么单个全局变量就不够用了，需要多个类似的变量来支撑。关于 flag.CommandLine 的自定制，后面会单独展开。*  
调用 Var 方法是会把 \*celsiusFlag 实参赋给 flag.Value 形参，编译器会在此时检查 \*celsiusFlag 类型是否有 flag.Value 所必需的方法：
```go
func CelsiusFlag(name string, value Celsius, usage string) *Celsius {
	f := celsiusFlag{value}
	flag.CommandLine.Var(&f, name, usage)
	return &f.Celsius
}
```
现在就可以在程序中使用这个标志了，使用代码如下：
```go
var temp = CelsiusFlag("temp", 20.0, "温度")

func main() {
	flag.Parse()
	fmt.Println(*temp)
}
```

接下来还可以把上面的例子简单改一下，不用结构体了，而是换成变量的别名，这样就需要额外再实现一个String方法，完整的代码如下：
```go
// output/flag/tempconv2
```

## 打印默认值
使用上面最后一个例子，打印帮助，查看默认值的提示：
```
PS G:\Steed\Documents\Go\src\gopl\output\flag\tempconv2> go run main.go -h
Usage of C:\Users\Steed\AppData\Local\Temp\go-build840446178\b001\exe\main.exe:
  -temp value
        温度 (default 36.70°C)
exit status 2
PS G:\Steed\Documents\Go\src\gopl\output\flag\tempconv2> go run main.go -temp 36.7C
*main.Celsius, 36.7°C
PS G:\Steed\Documents\Go\src\gopl\output\flag\tempconv2>
```
默认值打印的格式和打印只的格式是有区别的，这是因为类型不同，调用了不同的 String 方法。  
这里默认值显示的格式是根据接口类型的String方法定义的，在这里就是 \*celsiusValue 类型的String方法。而后面打印的是 Celsius 类型，使用的是 Celsius 类型的 String 方法。这里定义了两个String方法，但是打印的效果又不同，显示不统一，这样的做法不够好。这里可以看出两个问题：
1. 最初，使用结构体匿名封装的形式，避免了重复定义 String 方法。这样就保证了自定义的结构体类型 celsiusFlag 的String方法就是原本的 Celsius 类型的String方法。
2. 帮助消息中打印的默认值，实际是打印自定义类型的值。而自定义类型只在flag包中有用，解析完成后使用的都是原本的类型，这里就是 Celsius 类型。这两个类型的String方法最好能保持一致。

所以，使用结构体封装应该是一种不错的实现方式。不过flag包中的 time.Duration 类型用的就是类型别名来实现的：
```go
type durationValue time.Duration

func (d *durationValue) Set(s string) error {
	v, err := time.ParseDuration(s)
	*d = durationValue(v)
	return err
}

func (d *durationValue) String() string { return (*time.Duration)(d).String() }
```
上面是源码中的部分代码，可以看出这里保持一致的方法是进行类型转换后，调用原来类型的String方法。可能原本定义的是值类型的String方法，也可能直接就是定义了指针类型的String方法，不过指针类型的方法包括了所有值类型的方法，所以这里不必关系原本类型的方法具体是指针方法还是值方法。  
所以最后一个示例中的String方法也可以做同样的修改：
```go
func (c *celsiusValue) String() string { return (*Celsius)(f).String() }
```

## 自定义切片
以字符串切片为例，这里有两种实现的思路。一种是直接提供一个字符串，然后做分隔得到切片：
```go
// output/flag/nameValue
```
这里就不管 String 方法和 Set 方法展示规格的一致了，String方法才用 %q 的输出形式可以更好的把每一个元素清楚的展示出来。  

还有一种方式是，可以多次调用同一个参数，每一次调用，就添加一个元素：
```go
// output/flag/urls
```
由于每出现一个参数，都会调用一次 Set 方法，所以只要在 Set 里对切片进行append就可以了。不过这也带来一个问题，就是默认值无法被覆盖掉：
```
PS G:\Steed\Documents\Go\src\gopl\output\flag\urls> go run main.go -h
Usage of C:\Users\Steed\AppData\Local\Temp\go-build727433198\b001\exe\main.exe:
  -url value
        域名 (default "baidu.com")
exit status 2
PS G:\Steed\Documents\Go\src\gopl\output\flag\urls> go run main.go
["baidu.com"]
PS G:\Steed\Documents\Go\src\gopl\output\flag\urls> go run main.go -url shuxun.net -url 51cto.com
["baidu.com" "shuxun.net" "51cto.com"]
PS G:\Steed\Documents\Go\src\gopl\output\flag\urls>
```
下面这个版本的Set方法引入了一个全局变量，可以改进上面的问题：
```go
var newUrls urls
func (v *urls) Set(s string) error {
	newUrls = append(newUrls, s)
	*v = newUrls
	return nil
}
```
这里是一个方法，无法改成闭包。如果不想使用全局变量，可以把自定义类型改成结构体，添加一个newUrls字段。

## 简易的自定义版本
要实现自定义类型，只需要实现接口就可以了。不过上面的例子中都额外写了一个函数，用于返回自定义类型的指针，并且还设置了默认值。这个方法内部也是调用 Var 方法。这里可以直接使用 flag 包里的 Var 函数调用全局的Var方法：
```go
// output/flag/urls3
```
这里提供了两个设置初始值的方法，示例中都注释掉了。  
String方法由于内部是获得指针的，所以可以对变量进行修改。并且该方法调用的时机是在解析开始时只调用一次。所以在String方法里设置默认值是可行的。不过无法在打印帮助的时候把默认值打印出来。*黑科技？*  
另外，由于Var函数需要接收一个变量，所以在定义变量的时候，就可以赋一个初始值。并且在打印帮助的时候是可以把这个初始值打印出来的。  
不过简易版本最大的问题就是 Var 函数接收和返回的值都是 Value 接口类型。所以在使用之前，需要对返回值做一次类型转换。而设置初始值也是对 Value 接口类型的值进行设置。主要问题就是对外暴露了 Value 类型。现在调用者必须知道并且使用 Value 类型，对 Value 类型进行处理，这样就不是很友好。而之前的示例中，调用方（就是main函数中的那些代码）是完全可以忽略 Value 的存在的。  

# 自定义命令参数容器

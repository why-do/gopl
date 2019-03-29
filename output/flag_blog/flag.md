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

## 切片加参数
这是上面两个实现的结合，可以提供一组 os\.Args 的参数，另外还可以使用 flag 来进行参数设置。  
首先是不需要参数设置的情况。仅仅就是使用 flag 包提供的方法来代替使用 os\.Args 的实现：
```go
func main() {
	flag.Parse()
	for _, arg := range flag.Args() {
		fmt.Println(arg)
	}
}
```
基本上没什么差别，不过引入 flag 包之后，就可以使用参数了，比如加上一个 -upper 参数，让输出全大写：
```go
func main() {
	var upper bool
	flag.BoolVar(&upper, "upper", false, "是否大写")
	flag.Parse()
	for _, arg := range flag.Args() {
		if upper {
			fmt.Println(strings.ToUpper(arg))
		} else {
			fmt.Println(arg)
		}
	}
}
```
自定义切片类型的实现下面会讲。不过像这样简单的使用，只有一个切片类型，也不需要使用自定义类型就可以方便的实现了：
```
PS G:\Steed\Documents\Go\src\localdemo\flag> go run main.go -upper hello hi bye
HELLO
HI
BYE
PS G:\Steed\Documents\Go\src\localdemo\flag>
```
命令行参数必须放在前面，把不需要解析的参数全部放在最后。  

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
方法）。调用 Var 方法可以把这个标志加入到程序的命令行标记集合中，即全局变量 flag.CommandLine。*如果一个程序有非常复杂的命令行接口，那么单个全局变量就不够用了，需要多个类似的变量来支撑。最后一节“创建私有命令参数容器”会做简单的展开，不过也没有实现到这个程度。*  
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
// output/flag/fullName
```
这里就不管 String 方法和 Set 方法展示规格的一致了，String方法采用 %q 的输出形式可以更好的把每一个元素清楚的展示出来。  

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
var isNew bool
func (v *urls) Set(s string) error {
	if !isNew {
		*v = nil
		isNew = true
	}
	*v = append(*v, s)
	return nil
}
```
这里是一个方法，无法改成闭包。最好的做法就是将这个变量和原本的字符串切片封装为一个结构体：
```go
type urls struct {
	data  []string
	isNew bool
}
```
剩下的修改，参考之前自定义温度解析的实现就差不多了。  

## 简易的自定义版本
要实现自定义类型，只需要实现接口就可以了。不过上面的例子中都额外写了一个函数，用于返回自定义类型的指针，并且还设置了默认值。这个方法内部也是调用 Var 方法。这里可以直接使用 flag 包里的 Var 函数调用全局的Var方法：
```go
// output/flag/urls3
```
这里提供了两个设置初始值的方法，示例中都注释掉了。  
String 方法由于内部是获得指针的，所以可以对变量进行修改。并且该方法调用的时机是在解析开始时只调用一次。所以在 String 方法里设置默认值是可行的。不过无法在打印帮助的时候把默认值打印出来。*不需要这么做，但是正好可以对String方法有进一步的了解，还有就是这里利用指针修改参数原值的思路。*  
另外，由于 Var 函数需要接收一个变量，所以在定义变量的时候，就可以赋一个初始值。并且在打印帮助的时候是可以把这个初始值打印出来的。  
不过简易版本最大的问题就是 Var 函数接收和返回的值都是 Value 接口类型。所以在使用之前，需要对返回值做一次类型转换。而设置初始值也是对 Value 接口类型的值进行设置。主要问题就是对外暴露了 Value 类型。现在调用者必须知道并且使用 Value 类型，对 Value 类型进行处理，这样就不是很友好。而之前的示例中，调用方（就是main函数中的那些代码）是完全可以忽略 Value 的存在的。  
**小结：**这一小段主要是为了说明，之前示例中额外定义的函数是非常好的做法，封装了 flag 内部接口的细节。经过这个函数封装后再提供给用户使用，用户就可以完全忽略 flag.Value 这个接口而直接操作真正需要的类型了。这个函数的作用就是封装接口的所有细节，调用者只需要关注真正需要的操作的类型。  

# 自定义命令参数容器
接下来就是通过包提供的方法行进一步的自定制。以下3小节是一层一层更加接近底层的调用，做更加深入的定制。   

## 定制 Usage
回到最基本的使用，打印一下帮助消息可以得到以下的内容：
```
PS H:\Go\src\gopl\output\flag\beginning> go run main.go -h
Usage of C:\Users\Steed\AppData\Local\Temp\go-build926710106\b001\exe\main.exe:
  -age int
        年龄 (default 18)
  -name string
        名字 (default "Adam")
exit status 2
PS H:\Go\src\gopl\output\flag\beginning>
```
这里关注第一行，在 Usage of 后面是一长串的路径，这个是go run命令在构建上述命令源码文件时临时生成的可执行文件的完整路径。如果是编译之后再执行，就是可执行文件的相对路径，就没那么难看了。  
这一行的内容也是可以自定制的，但是首先来看看源码里的实现：
```go
var Usage = func() {
	fmt.Fprintf(CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	PrintDefaults()
}

func (f *FlagSet) Output() io.Writer {
	if f.output == nil {
		return os.Stderr
	}
	return f.output
}
```
看上面的代码就清楚了，输出的内容就是执行的命令本身 os.Args[0]。就会输出的位置默认就是标准错误 os.Stderr。  
这个 Usage 是可导出的变量，值是一个匿名函数，只要重新为 Usage 赋一个新值就可以完成内容的自定制：
```go
var name string

func init() {
	flag.StringVar(&name, "name", "Adam", "名字")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "请指定名字和年龄：")
		flag.PrintDefaults()
	}
}

var ageP = flag.Int("age", 18, "年龄")

func main() {
	flag.Parse()
	fmt.Printf("%T %[1]v\n", name)
	fmt.Printf("%T %[1]v\n", ageP)
	fmt.Printf("%T %[1]v\n", *ageP)
}
```
只要在 flag.Parse() 执行前覆盖掉 flag.Usage 即可。  
下面那行 flag.PrintDefaults() 则是打印帮助信息中其他的内容。完全可以把这行去掉，这里完全可以自定义打印更多其他内容，甚至是执行其他操作。  

## 定制 CommandLine
在调用flag包中的一些函数（比如StringVar、Parse等等）的时候，实际上是在调用flag.CommandLine变量的对应方法。  
flag.CommandLine相当于默认情况下的命令参数容器。通过对flag.CommandLine重新赋值，就可以更深层次地定制当前命令源码文件的参数使用说明。  
flag包提供了NewFlagSet函数用于创建自定制的 CommandLine 。在上一个简单例子的基础上，修改一下其中的init函数的内容：
```go
var name string
var age int

func init() {
	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)
	flag.StringVar(&name, "name", "Adam", "名字")
	age = *flag.Int("age", 18, "年龄")
	// 和上面两句效果一样
	// flag.CommandLine.StringVar(&name, "name", "Adam", "名字")
	// var ageP = flag.CommandLine.Int("age", 18, "年龄")
	flag.CommandLine.Usage = func() {
		fmt.Fprintln(os.Stderr, "请指定名字和年龄：")
		flag.PrintDefaults()
	}
}
```
其实这里只加了一行语句。所有flag包的操作都要在flag.NewFlagSet执行之后，否则之前执行的内容会被覆盖掉。所以这里把flag.Int的调用移到了包内，否则在全局中的赋值语句会在这之前就运行了，然后被flag.NewFlagSet方法覆盖掉。  
这里无论是 `flag.StringVar` 或者是 `flag.CommandLine.StringVar`，最终都是使用flag.NewFlagSet创建的 \*FlagSet 对象的方法来调用的。不过本质上是有差别的：
+ `flag.StringVar` ： 使用默认 CommandLine 的对象调用，但是第一行语句 `flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)` 则是把它的值覆盖为新创建的对象。
+ `flag.CommandLine.StringVar` ： 使用 flag.NewFlagSet 函数创建的对象来调用，所以和上面是一个东西。

第一个方式是专门为默认的容器提供的便捷调用方式。第二个是则是通用的方法，之后创建私有命令参数容器的时候就需要用通用的方式来调用了。  
Usage 必须用flag.CommandLine调用。另外不定制的话，包里也准备了默认的方法可以使用：
```go
func (f *FlagSet) defaultUsage() {
	if f.name == "" {
		fmt.Fprintf(f.Output(), "Usage:\n")
	} else {
		fmt.Fprintf(f.Output(), "Usage of %s:\n", f.name)
	}
	f.PrintDefaults()
}
```
第一个参数的作用基本就是显示一个名称，也可以用空字符串，向上面这样。而第二个参数可以是下面三种常量：
```go
const (
	ContinueOnError ErrorHandling = iota // Return a descriptive error.
	ExitOnError                          // Call os.Exit(2).
	PanicOnError                         // Call panic with a descriptive error.
)
```
效果一看就明白了。定义在解析遇到问题后，是执行何种操作。默认的就是ExitOnError，所以在--help执行打印说明后，最后一行会出现“exit status 2”，以状态码2退出。这里可以根据需要定制为抛出Panic。  
使用-h参数打印帮助信息也算是解析出错，如果是Panic则会在打印帮助信息后Panic，如果是Continue则先打印帮助信息然后按照默认值执行。所以如果要使用另外两种模式，最好修改一下-h参数的行为，就是上面讲的定制Usage。使用-h参数之后程序将执行的就是Usage指定的函数。  

## 创建私有命令参数容器
上一个例子依然是使用flag包提供的命令参数容器，只是重新进行了创建和赋值。这里依然是调用flag.NewFlagSet()函数创建命令参数容器，不过这次赋值给自定义的变量：
```go
package main

import (
	"flag"
	"fmt"
	"os"
)

var cmdLine = flag.NewFlagSet("", flag.ExitOnError)
var name string
var age int

func init() {
	cmdLine.StringVar(&name, "name", "Adam", "名字")
	age = *cmdLine.Int("age", 18, "年龄")
}

func main() {
	cmdLine.Parse(os.Args[1:])
	fmt.Printf("%T %[1]v\n", name)
	fmt.Printf("%T %[1]v\n", age)
}

```
首先通过 flag.NewFlagSet 函数创建了私有的命令参数容器。然后调用其他方法的接收者都使用这个容器。另外还有很多方法可以调用，可以继续探索。
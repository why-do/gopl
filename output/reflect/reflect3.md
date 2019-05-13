# 12.7 访问结构体字段标签
*这里的“成员”和“字段”两个词有点混用，但都是同一个意思。*  
可以使用结构体**成员标签**（field tag）在进行JSON反序列化的时候对应JSON中字段的名字。json 成员标签让我们可以选择其他的字段名以及忽略输出的空字段。这小节将通过反射机制获取结构体字段的标签，然后填充字段的值，就和JSON反序列化一样，目标和结果是一样的，只是获取的数据源不同。  
有一个 Web 服务应用的场景，在 Web 服务器中，绝大部分 HTTP 处理函数的第一件事就是提取请求参数到局部变量中。这里将定义一个工具函数 params\.Unpack，使用结构体成员标签直接将参数填充到结构体对应的字段中。因为 URL 的长度有限，所以参数的名称一般比较短，含义也比较模糊。这需要通过成员标签将结构体的字段和参数名称对应上。  

## 在HTTP处理函数中使用
首先，展示这个工具函数的用法。就是假设已经实现了这个 params\.Unpack 函数，下面的 search 函数就是一个 HTTP 处理函数，它定义了一个结构体变量 data，data 也定义了成员标签来对应请求参数的名字。Unpack 函数从请求中提取数据来填充这个结构体，这样不仅可以更方便的访问，还避免了手动转换类型：
```go
package main

import (
	"fmt"
	"net/http"
)

import "gopl/ch12/params"

// search 用于处理 /search URL endpoint.
func search(resp http.ResponseWriter, req *http.Request) {
	var data struct {
		Labels     []string `http:"l"`
		MaxResults int      `http:"max"`
		Exact      bool     `http:"x"`
	}
	data.MaxResults = 10 // 设置默认值
	if err := params.Unpack(req, &data); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// ...其他处理代码...
	fmt.Fprintf(resp, "Search: %+v\n", data)
}

// 这里还缺少一个 main 函数，最后会补上
```
 
## 工具函数 Unpack 的实现
下面的 Unpack 函数做了三件事情：  
一、调用 req\.ParseForm() 来解析请求。在这之后，req\.Form 就有了所有的请求参数，这个方法对 HTTP 的 GET 和 POST 请求都适用。  
二、Unpack 函数构造了一个从每个**有效**字段名到对应字段变量的映射。在字段有标签时，有效字段名与实际字段名可以不同。reflect\.Type 的 Field 方法会返回一个 reflect\.StructField 类型，这个类型提供了每个字段的名称、类型以及一个可选的标签。它的 Tag 字段类型为 reflect\.StructTag，底层类型为字符串，提供了一个 Get 方法用于解析和提取对于一个特定 key 的子串，比如下面例子中会用到的 http:"..."。  
三、Unpack 遍历 HTTP 参数中的所有 key\/value 对，并且更新对应的结构体字段。同一个参数可以出现多次。如果对应的字段是切片，则参数所有的值都会追加到切片里。否则，这个字段会被多次覆盖，只有最后一次的值才有效。  

Unpack 函数的代码如下：
```go
// Unpack 从 HTTP 请求 req 的参数中提取数据填充到 ptr 指向的结构体的各个字段
func Unpack(req *http.Request, ptr interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	// 创建字段映射表，key 为有效名称
	fields := make(map[string]reflect.Value)
	v := reflect.ValueOf(ptr).Elem() // reflect.ValueOf(&x).Elem() 获得任意变量 x 可寻址的值，用于设置值。
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField，提供了每个字段的名称、类型以及一个可选的标签
		tag := fieldInfo.Tag           // a reflect.Structtag，底层类型为字符串，提供了一个 Get 方法，下一行就用到了
		name := tag.Get("http")        // Get 方法用于解析和提取对于一个特定 key 的子串
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		fields[name] = v.Field(i)
	}

	// 对请求中的每个参数更新结构体中对应的字段
	for name, values := range req.Form {
		f := fields[name]
		if !f.IsValid() {
			continue // 忽略不能识别的 HTTP 参数
		}
		for _, value := range values {
			if f.Kind() == reflect.Slice {
				elem := reflect.New(f.Type().Elem()).Elem()
				if err := populate(elem, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
				f.Set(reflect.Append(f, elem))
			} else {
				if err := populate(f, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
			}
		}
	}
	return nil
}
```

这里还调用了一个 populate 函数，负责从单个 HTTP 请求参数值填充单个字段 v （或者切片字段中的单个元素）。目前，它仅支持字符串、有符号整数和布尔值。要支持其他类型可以再添加：
```go
func populate(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(value)

	case reflect.Int:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)

	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		v.SetBool(b)

	default:
		return fmt.Errorf("unsupported kind %s", v.Type())
	}
	return nil
}
```

## 执行效果
接着把 search 处理程序添加到一个 Web 服务器中，直接在 search 所在的 main 包的命令源码文件中添加下面的 main 函数：
```go
func main() {
	fmt.Println("http://localhost:8000/search")                                 // Search: {Labels:[] MaxResults:10 Exact:false}
	fmt.Println("http://localhost:8000/search?l=golang&l=gopl")                 // Search: {Labels:[golang gopl] MaxResults:10 Exact:false}
	fmt.Println("http://localhost:8000/search?l=gopl&x=1")                      // Search: {Labels:[gopl] MaxResults:10 Exact:true}
	fmt.Println("http://localhost:8000/search?x=true&max=100&max=200&l=golang") // Search: {Labels:[golang] MaxResults:200 Exact:true}
	fmt.Println("http://localhost:8000/search?q=hello")                         // Search: {Labels:[] MaxResults:10 Exact:false}  # 不存在的参数会忽略
	fmt.Println("http://localhost:8000/search?x=123")                           // x: strconv.ParseBool: parsing "123": invalid syntax  # x 提供的参数解析错误
	fmt.Println("http://localhost:8000/search?max=lots")                        // max: strconv.ParseInt: parsing "lots": invalid syntax  # max 提供的参数解析错误
	http.HandleFunc("/search", search)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
```
这里提供了几个示例以及输出的结果，直接使用浏览器，输入URL就能返回对应的结果。

# 12.8 显示类型的方法
通过反射的 reflect\.Type 来获取一个任意值的类型并枚举它的方法。下面的例子是把类型和方法都打印出来：
```go
// ch12/methods/methods.go
```
reflect\.Type 和 reflect\.Value 都有一个叫作 Method 的方法：
+ 每个 t.Method(i) 都会返回一个 reflect\.Method 类型的实例，这个结构类型描述了这个方法的名称和类型。
+ 每个 v.Method(i) 都会返回一个 reflect\.Value，代表一个方法值，即一个已经绑定接收者的方法。

下面是两个示例测试，展示以及验证上面的函数：
```go
// ch12/methods/methods_test.go
```

另外还有一个 reflect\.Value\.Call 方法，可以调用 Func 类型的 Value，这里没有演示。  

# 12.9 注意事项
还有很多反射API，这里的示例展示了反射能做哪些事情。  
反射是一个功能和表达能力都很强大的工具，但是要慎用，主要有三个原因。

## 代码脆弱
基于反射的代码是很脆弱的。一般编译器在编译时就能报告错误，但是反射错误则要等到执行时才会以崩溃的方式来报告。这可能是等待程序运行很久以后才会发生。  
比如，尝试读取一个字符串然后填充一个 Int 类型的变量，那么调用 reflect\.Value\.SetString 就会崩溃。很多使用反射的程序都会有类似的风险。所以对每一个 reflect\.Value 都需要仔细检查它的类型、是否可寻址、是否可设置。  
要回避这种缺陷的最好的办法就是确保反射的使用完整的封装在包里，并且如果可能，在包的 API 中避免使用 reflect\.Value，尽量使用特定的类型来确保输入是合法的值。如果做不到，那就需要在每个危险的操作前都做额外的动态检查。比如标准库的 fmt\.Printf 可以作为一个示例，当遇到操作数类型不合适时，它不会崩溃，而是输出一条描述性的错误消息。这尽管仍然会有 bug，但定位起来就简单多了：
```go
fmt.Printf("%d %s\n", "hello", 123) // %!d(string=hello) %!s(int=123)
```

反射还降低了自动重构和分析工具的安全性与准确度，因为它们无法检测到类型的信息。  

## 难理解、难维护
类型也算是某种形式的文档，而反射的相关操作则无法做静态类型检查，所以大量使用反射的代码是很难理解的。对应接收 interface{} 或者reflect\.Value 的函数，一定要写清楚期望的参数类型和其他限制条件。  

## 运行慢
基于反射的函数会比为特定类型优化的函数慢一到两个数量级。在一个典型的程序中，大部分函数与整体性能无关，所以为了让程序更清晰可以使用反射。比如测试就和适合使用反射，因为大部分测试都使用小数据集。但对性能关键路径上的函数，最好避免使用反射。  
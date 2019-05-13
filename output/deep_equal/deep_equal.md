这篇解决反射章节第一个例子 dispaly 中没有处理的循环引用的问题
拼接在unsafe包的内容之后，主要是需要unsafe.Pointer的理论为基础。

# 13.3 示例：深度相等
reflect 包中的 DeepEqual 函数用来报告两个变量的值是否深度相等。DeepEqual 函数的基本类型使用内置的 == 操作符进行比较。对于组合类型，它逐层深入比较相应的元素。因为这个函数适合于任意的一对变量值的比较，甚至是那些无法通过 == 来比较的值，所以在一些测试代码中广泛地使用这个函数。下面的代码就是用 DeepEqual 来比较两个 []string 类型的值：
```go
func TestSplit(t *testing.T) {
	got := strings.Split("a:b:c", ":")
	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) { /* ... */ }
}
```

## DeepEqual 的不足
虽然 DeepEqual 很方便，可以支持任意的数据类型，但是它的不足是判断过于武断。例如，一个值为 nil 的 map 和一个值不为 nil 的空 map 会判断为不相等，一个值为 nil 的切片和不为 nil 的空切片同样也会判断为不相等：
```go
var c, d map[string]int = nil, make(map[string]int)
fmt.Println(reflect.DeepEqual(c, d)) // "false"

var a, b []string = nil, []string{}
fmt.Println(reflect.DeepEqual(a, b)) // "false"
```

## 自定义比较函数
所以，接下来要自己定义一个 Equal 函数。和 DeepEqual 类似，但是可以把一个值为 nil 的切片或 map 和一个值不为 nil 的空切片或 map 判断为相等。对参数的基本递归检查可以通过反射来实现。需要定义一个未导出的函数 equal 用来进行递归检查，隐藏反射的细节。参数 seen 是为了检查循环引用，并且因为要递归所以作为参数进行传递。对于每对要进行比较的值 x 和 y，equal 函数检查两者是否合法（IsValid）以及它们是否具有相同的类型（Type）。函数的结果通过 switch 的 case 语句返回，在 case 中比较两个相同类型的值：
```go
// output/deep_equal/equal/equal.go
```
在 API 中不暴露反射的细节，所以最后的可导出的 Equel 函数对参数显式调用 reflect\.ValueOf 函数。  

## 支持循环引用
为了确保算法终止设置可以对循环数据结果进行比较，它必须记录哪两对变量已经比较过了，并且避免再次进行比较。Equal 函数定义了一个叫做 comparison 的结构体集合，每个元素都包含两个变量的地址（unsafe\.Pointer 表示）以及比较的类型。比如切片的比较，x 和 x[0] 的地址是一样的，这时候就要分开是两个切片的比较 x 和 y，还是切片的两个元素的比较 x[0] 和 y[0]。  
当 equal 确认了两个参数都是合法的并且类型也一样，在执行 switch 语句进行比较之前，先检查这两个变量是否已经比较过了，如果已经比较过了，则直接返回结果并终止这次递归比较。  

**unsafe.Pointer**  
就是上一节讲的问题，reflect\.UnsafeAddr 返回的是一个 uintptr 类型（字母意思就是不安全的地址），这里需要直接转成 unsafe.Pointer 类型来保证地址可以始终指向最初的那个变量。  

## 测试验证
下面输出完整的测试代码：
```go
// output/deep_equal/equal/equal_test.go
```

在最后的示例测试函数 Example_equalCycle 中，验证了一个循环链表也能完成比较，而不会卡住：
```go
type link struct {
	value string
	tail  *link
}
a, b, c := &link{value: "a"}, &link{value: "b"}, &link{value: "c"}
a.tail, b.tail, c.tail = b, a, c
```
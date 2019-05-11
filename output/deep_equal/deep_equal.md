反射之后就可以发表

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
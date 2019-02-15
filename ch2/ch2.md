# 2.1 名称

## 关键字
共25个关键字，只能用在语法允许的地方，不能作为名称：
```
break       //退出循环
default     //选择结构默认项（switch、select）
func        //定义函数
interface   //定义接口
select      //channel
case        //选择结构标签
chan        //定义channel
const       //常量
continue    //跳过本次循环
defer       //延迟执行内容（收尾工作）
go          //并发执行
map         //map类型
struct      //定义结构体
else        //选择结构
goto        //跳转语句
package     //包
switch      //选择结构
fallthrough //switch里继续检查后面的分支
if          //选择结构
range       //从slice、map等结构中取元素
type        //定义类型
for         //循环
import      //导入包
return      //返回
var         //定义变量
```

## 内置预声明
内置的预声明的常量、类型和函数：
+ 常量
  + true、false
  + iota
  + nil
+ 类型
  + int、int8、int16、int32、int64
  + uint、uint8、uint16、uint32、uint64、uintptr
  + float32、float64、complex128、complex64
  + bool、byte、rune、string、error
+ 函数
  + make、len、cap、new、append、copy、close、delete
  + complex、real、imag
  + panic、recover

这些名称不是预留的，可以在声明中使用它们。也可能会看到对其中的名称进行重声明，但是要知道这会有冲突的风险。

## 命名规则
单词组合时，使用驼峰式。如果是缩写，比如：ASCII或HTML，要么全大写，要么全小写。比如组合 html 和 escape，可以是下面几种写法：
+ htmlEscape
+ HTMLEscape
+ EscapeHTML

但是不推荐这样的写法：
+ Escapehtml : 这样完全区分不了html是一个词，所以这样HTML要全大写
+ EscapeHtml : 这样虽然能区分，但是违反了全大写或全小写的建议


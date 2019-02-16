# 基本数据

Go的数据类型分四大类：
1. 基础类型（basic type）
    + 数字（number）
    + 字符串（string）
    + 布尔型（boolean）
2. 聚合类型（aggregate type）
    + 数组（array）
    + 结构体（struct）
3. 引用类型（reference type）
    + 指针（pointer）
    + 切片（slice）
    + 哈希表（map）
    + 函数（function）
    + 通道（channel）
4. 接口类型（interface type）

# 3.1 整数

## 二元操作符
二元操作符分五大优先级，按优先级降序排列：
```
*    /   %   <<  >>  &   &^
+    -   |   ^
==    !=  <   <=  >   >=
&&
||
```

位运算符：
| 符号 | 说明 | 集合 |
| --- | --- | --- |
| & | AND | 交集 |
| \| | OR | 并集 |
| ^ | XOR | 对称差 |
| &^ | 位清空（AND NOT） | 差集 |
| << | 左移 | N/A |
| >> | 右移 | N/A |

位清空，表达式 z=x&^y ，把y中是1的位在x里对应的那个位，置0。  
差集，就是集合x去掉集合y中的元素之后的集合。对称差则是再加上集合y去掉集合x中的元素的集合，就是前后两个集合互相求差集，之后再并集。

## fmt的两个技巧
一、%后的副词[1]告知Printf重复使用第一个操作数。  
二、%o、%x、%X之前的副词#告知Printf输出相应的前缀 0、0x、0X。  
```go
func main() {
	o := 0666
	fmt.Printf("%d %[1]o %#[1]o\n", o)  // 438 666 0666
	x := int64(0xdeadbeef)
	fmt.Printf("%d %[1]x %#[1]x %#[1]X\n", x)  // 3735928559 deadbeef 0xdeadbeef 0XDEADBEEF
}
```
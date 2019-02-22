package main

import "fmt"

type tree struct {
	value       int
	left, right *tree
}

// 就地排序
func sort(values []int) {
	var root *tree
	// 先添加到一个二叉树中。add函数中的逻辑：小的放左边，大的放右边
	for _, v := range values {
		root = add(root, v)
	}
	// 原来的切片不要了，从原来的切片的第一个元素的位置，从小到大追加元素
	appendValues(values[:0], root)
}

// 递归调用，从根开始一层一层往下找，比节点小就往左找下一层，比节点大就往右找下一层
// 直到找到nil，在那个位置创建一个新节点，value就是自己的值，left和right都不管，默认是指针的零值就是nil
func add(t *tree, value int) *tree {
	// 递归，先写退出条件
	if t == nil {
		t = new(tree)
		t.value = value
		return t
	}
	// 然后是递归的调用
	if value < t.value {
		t.left = add(t.left, value)
	} else {
		t.right = add(t.right, value)
	}
	return t
}

// 将元素按顺序追加到 values 里，然后返回切片
func appendValues(values []int, t *tree) []int {
	if t != nil {
		values = appendValues(values, t.left)
		values = append(values, t.value)
		values = appendValues(values, t.right)
	}
	return values
}

func main() {
	l := []int{3, 1, 6, 3, 7, 2, 0, 2, 3}
	fmt.Println(l)
	sort(l)
	fmt.Println(l)
}

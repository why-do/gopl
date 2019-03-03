package main

import (
	"fmt"
	"sort"
)

var prereqs = map[string][]string{
	"algorithems":    {"data structures"},
	"calculus":       {"linear algebra"},
	"linear algebra": {"calculus"},  // 加了这个后，现在这个图有环存在了，一个节点可以通过路径回到它自己
	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},
	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

func main() {
	for i, course := range topoSort(prereqs) {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
}

// 概念：有向图、拓扑排序
// 如果是有向无环图，就可以输出序列
// 如果图有环，则要发现环
// 1、选一个没有前驱的顶点，并输出
// 2、删除此顶点。
// 重复这两边直到图空，或者图不空但找不到无前驱的顶点
func topoSort(m map[string][]string) []string {
	// 闭包的部分
	var order []string
	seen := make(map[string]bool)
	var visitAll func(items []string)
	visitAll = func(items []string) {
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				visitAll(m[item])
				order = append(order, item)
			}
		}
	}
	// 主体
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	visitAll(keys)
	return order
}

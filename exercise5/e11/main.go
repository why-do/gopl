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
	order, ok := topoSort(prereqs)
	if !ok {
		fmt.Printf("发现成环: \n")
	}
	for i, course := range order {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
}

// 概念：有向图、拓扑排序
// 如果是有向无环图，就可以输出序列
// 如果图有环，则要发现环
func topoSort(m map[string][]string) ([]string, bool) {
	// 闭包的部分
	var order []string
	seen := make(map[string]bool)
	var visitAll func(items []string)
	var checkLoop, loopOrder []string
	var foundLoop bool
	visitAll = func(items []string) {
		if foundLoop {
			return
		}
		for _, item := range items {
			for _, i := range checkLoop {
				if i == item {
					for _, j := range checkLoop {
						loopOrder = append(loopOrder, j)
					}
					loopOrder = append(loopOrder, item)
					foundLoop = true
					return
				}
			}
			if !seen[item] {
				seen[item] = true
				checkLoop = append(checkLoop, item)
				visitAll(m[item])
				checkLoop = checkLoop[:len(checkLoop)-1]
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
	if foundLoop {
		return loopOrder, !foundLoop
	}
	return order, !foundLoop
}

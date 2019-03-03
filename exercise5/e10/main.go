package main

import (
	"fmt"
)

var prereqs = map[string][]string{
	"algorithems": {"data structures"},
	"calculus":    {"linear algebra"},
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
	topo := topoSort(prereqs)
	for i := 1; i <= len(topo); i++ {
		fmt.Printf("%d\t%s\n", i, topo[i])
	}
	check(topo)
}

func check(m map[int]string) {
	// 闭包
	learned := make(map[string]bool)
	checkLearned := func(item string) bool {
		fmt.Printf("checking: %q%*s", item, 25-len(item), "")
		for i, need := range prereqs[item] {
			if i == 0 {
				fmt.Print("need:")
			}
			fmt.Printf(" %q", need)
			if !learned[need] {
				return false
			}
		}
		fmt.Println()
		learned[item] = true
		return true
	}
	// 主体
	for i := 1; i <= len(m); i++ {
		if ok := checkLearned(m[i]); !ok {
			fmt.Println("验证到不合法:", m[i])
		}
	}
}

func topoSort(m map[string][]string) map[int]string {
	// 闭包的部分
	var order = make(map[int]string)
	index := 1
	seen := make(map[string]bool)
	var visitAll func(items []string)
	visitAll = func(items []string) {
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				visitAll(m[item])
				order[index] = item
				index++
			}
		}
	}
	// 主体
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	visitAll(keys)
	return order
}

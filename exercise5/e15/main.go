package main

import "fmt"

func sum(vals ...int) int {
	total := 0
	for _, val := range vals {
		total += val
	}
	return total
}

func max(vals ...int) int {
	if len(vals) == 0 {
		return 0
	}
	var res int
	for i, val := range vals {
		if i == 0 {
			res = val
			continue
		}
		if res < val {
			res = val
		}
	}
	return res
}

func min(vals ...int) int {
	if len(vals) == 0 {
		return 0
	}
	var res int
	for i, val := range vals {
		if i == 0 {
			res = val
			continue
		}
		if res > val {
			res = val
		}
	}
	return res
}

func main() {
	fmt.Println(max())
	fmt.Println(max(3))
	fmt.Println(max(1, 2, 3, 4))
}
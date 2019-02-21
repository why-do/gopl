package main

import "fmt"

func nonrepeat(s []string) []string {
	var tmp string
	var n int
	for i, v := range s {
		// fmt.Printf("%d %q %q %q\n", i, tmp, v, s)
		if tmp == v && i != 0 {
			tmp = v
			continue
		}
		tmp = v
		s[n] = v
		n++
	}
	return s[:n]
}

func main() {
	s := []string{"", "", "a", "a", "b", "b", "a"}
	s = nonrepeat(s)
	fmt.Printf("%q\n", s)
}

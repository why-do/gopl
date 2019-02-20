package main

import (
	"fmt"
	"strings"
)

func isomerism(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}

	for _, r := range s1 {
		rs := string(r)
		if strings.Count(s1, rs) != strings.Count(s2, rs) {
			return false
		}
	}
	return true
}

var s string = "122333444455555"
var s1 string = "543215432543545"
var s2 string = "543214321321211"
var s3 string = "123451234123121"
var s4 string = "12345"

func main() {
	fmt.Println(isomerism(s, s1))
	fmt.Println(isomerism(s, s2))
	fmt.Println(isomerism(s, s3))
	fmt.Println(isomerism(s, s4))
	fmt.Println(isomerism(s2, s3))
}

package main

import (
	"bufio"
	"bytes"
	"fmt"
)

type WordCounter int

func (c *WordCounter) Write(p []byte) (int, error) {
	input := bufio.NewScanner(bytes.NewReader(p))
	input.Split(bufio.ScanWords)
	var words int
	for input.Scan() {
		*c++
		words++
	}
	return words, nil
}

func main() {
	var c WordCounter
	s1 := "It is a good day to die!"
	n, _ := c.Write([]byte(s1))
	fmt.Println(c, n)
	n, _ = c.Write([]byte(s1))
	fmt.Println(c, n)
}

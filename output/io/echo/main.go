package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

func main() {
	input := bufio.NewScanner(os.Stdin)
	buff := bytes.NewBuffer(nil)
	for input.Scan() {
		buff.Write(input.Bytes()) // 忽略错误
		buff.WriteByte(' ')
	}
	fmt.Printf("%q\n", buff.String()[:buff.Len()-1])  // 截掉最后一个空格
}

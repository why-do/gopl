package main

import (
	"bufio"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"os"
)

var shaP = flag.Int("sha", 256, "SHA256 or SHA384 or SHA512")

func main() {
	flag.Parse()
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		SumN(*shaP, input.Bytes())
	}
}

func SumN(n int, data []byte) {
	switch n {
	case 384:
		fmt.Printf("% x\n", sha512.Sum384(data))
	case 512:
		fmt.Printf("% x\n", sha512.Sum512(data))
	default:
		fmt.Printf("% x\n", sha256.Sum256(data))
	}
}

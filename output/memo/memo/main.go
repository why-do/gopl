package main

import (
	"fmt"
	"os"

	"gopl/output/memo/memotest"
)

func main() {
	resp, err := memotest.HTTPGetBody(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Printf("%s\n", resp)
}

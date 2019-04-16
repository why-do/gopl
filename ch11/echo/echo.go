// 输出命令行参数
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	n   = flag.Bool("n", false, "omit trailing newline")
	sep = flag.String("s", " ", "separator")
)

var out io.Writer = os.Stdout // 测试过程中将会被更改

func main() {
	flag.Parse()
	if err := echo(!*n, *sep, flag.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "echo: %v\n", err)
		os.Exit(1)
	}
}

func echo(newline bool, sep string, args []string) error {
	fmt.Fprintf(out, strings.Join(args, sep))
	if newline {
		fmt.Fprintln(out)
	}
	return nil
}

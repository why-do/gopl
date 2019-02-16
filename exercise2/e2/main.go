// cf 把它的数值参数转换为摄氏度和华氏度
package main

import (
	"bufio"
	"fmt"
	"gopl/exercise2/e2/unitconv"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			out(arg)
		}
	} else {
		fmt.Println("Input from os.Stdin. Use Ctrl+Z to exit.")
		input := bufio.NewScanner(os.Stdin)
		for input.Scan() {
			out(input.Text())
		}
	}
}

func out(s string) {
	u, err := strconv.ParseFloat(s, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cf: %v\n", err)
		os.Exit(1)
	}

	f := unitconv.Fahrenheit(u)
	c := unitconv.Celsius(u)
	fmt.Printf("%s = %s, %s = %s\n",
		f, unitconv.FToC(f), c, unitconv.CTOF(c))

	i := unitconv.Inch(u)
	m := unitconv.Meter(u)
	fmt.Printf("%s = %s, %s = %s\n",
		i, unitconv.IToM(i), m, unitconv.MToI(m))

	p := unitconv.Pound(u)
	k := unitconv.Kilogram(u)
	fmt.Printf("%s = %s, %s = %s\n",
		p, unitconv.PToK(p), k, unitconv.KToP(k))
}

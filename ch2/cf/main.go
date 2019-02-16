// cf 把它的数值参数转换为摄氏度和华氏度
package main

import (
	"fmt"
	"gopl/exercise2/e1/tempconv"
	"os"
	"strconv"
)

func main() {
	for _, arg := range os.Args[1:] {
		t, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cf: %v\n", err)
			os.Exit(1)
		}
		f := tempconv.Fahrenheit(t)
		c := tempconv.Celsius(t)
		k := tempconv.Kelvin(t)
		fmt.Printf("%s = %s, %s = %s\n",
			f, tempconv.FToC(f), c, tempconv.CTOF(c))
		fmt.Printf("%s = %s, %s = %s\n",
			k, tempconv.KToC(k), c, tempconv.CToK(c))
	}
}

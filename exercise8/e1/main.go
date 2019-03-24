package main

import (
	"fmt"
	"time"
)

func main() {
	loc := new(time.Location)
	loc, _ = time.LoadLocation("Asia/Shanghai")
	fmt.Printf("%s:\t%s\n", loc, time.Now().In(loc).Format("2006/01/02 15:04:05"))
	loc, _ = time.LoadLocation("US/Eastern")
	fmt.Printf("%s:\t%s\n", loc, time.Now().In(loc).Format("2006/01/02 15:04:05"))
	loc, _ = time.LoadLocation("Asia/Tokyo")
	fmt.Printf("%s:\t%s\n", loc, time.Now().In(loc).Format("2006/01/02 15:04:05"))
	loc, _ = time.LoadLocation("Europe/London")
	fmt.Printf("%s:\t%s\n", loc, time.Now().In(loc).Format("2006/01/02 15:04:05"))
}

// 将符合条件的issue输出为一个表格
package main

import (
	"fmt"
	"gopl/ch4/github"
	"log"
	"os"
	"time"
)

func main() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issue: \n", result.TotalCount)
	for _, item := range result.Items {
		if item.CreateAt.After(time.Now().AddDate(-1, 0, 0)) {
			fmt.Printf("%v ", item.CreateAt.Format("2006-01-02"))
			fmt.Printf("#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title)
		}
	}
}

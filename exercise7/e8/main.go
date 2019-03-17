/* 排序的稳定性
如果要按多个维度进行排序，比如排序第一关键字是姓名，第二关键字是年龄。
那么先做一个年龄的排序（稳定不稳定无所谓），然后再在原来的基础上做一次姓名的排序（必须稳定），才能得到正确的结果。
sort.Sort 是快速排序，是不稳定的排序。所以只能进行一次排序，每次比较的时候把所有字段都进行比较，才能避免不稳定的情况
sort.Stable 能保证排序的稳定性，但是牺牲了效率
*/
package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"text/tabwriter"
	"time"
)

// 随机种子
func init() {
	rand.Seed(time.Now().UnixNano())
}

// 结构体
type Score struct {
	Id                      uint
	Class                   string
	Gender                  string
	Chinese, Maths, English float64
}

// 随机生成一个数据
func NewScore(id uint) *Score {
	return &Score{
		Id:      id,
		Class:   []string{"A", "B", "C", "D"}[rand.Intn(4)],
		Gender:  []string{"Boy", "Girl"}[rand.Intn(2)],
		Chinese: float64(rand.Intn(31) + 70),
		Maths:   float64(rand.Intn(31) + 70),
		English: float64(rand.Intn(31) + 70),
	}
}

// 随机生成一组数据
var scores []*Score

func init() {
	for i := 1; i <= 30; i++ {
		scores = append(scores, NewScore(uint(i)))
	}
}

// 打印的方法
func printScores(scores []*Score) {
	const formatTitle = "%v\t%v\t%v\t%v\t%v\t%v\t\n"
	const formatData = "%v\t%v\t%v\t%.1f\t%.1f\t%.1f\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, formatTitle, "Id", "Class", "Gender", "Chinese", "Maths", "English")
	fmt.Fprintf(tw, formatTitle, "--", "-----", "------", "-------", "-----", "-------")
	for _, t := range scores {
		fmt.Fprintf(tw, formatData, t.Id, t.Class, t.Gender, t.Chinese, t.Maths, t.English)
	}
	tw.Flush()
}

// 实现排序的接口
type customSort struct {
	t      []*Score
	orders []string
}

func (x customSort) Len() int      { return len(x.t) }
func (x customSort) Swap(i, j int) { x.t[i], x.t[j] = x.t[j], x.t[i] }
func (x customSort) Less(i, j int) bool { // TODO 可以考虑反射的实现
	for _, o := range x.orders {
		switch o {
		case "Id":
			if x.t[i].Id != x.t[j].Id {
				return x.t[i].Id < x.t[j].Id
			}
		case "Class":
			if x.t[i].Class != x.t[j].Class {
				return x.t[i].Class < x.t[j].Class
			}
		case "Gender":
			if x.t[i].Gender != x.t[j].Gender {
				return x.t[i].Gender < x.t[j].Gender
			}
		case "Chinese":
			if x.t[i].Chinese != x.t[j].Chinese {
				return x.t[i].Chinese < x.t[j].Chinese
			}
		case "-Chinese":
			if x.t[i].Chinese != x.t[j].Chinese {
				return x.t[j].Chinese < x.t[i].Chinese
			}
		case "Maths":
			if x.t[i].Maths != x.t[j].Maths {
				return x.t[i].Maths < x.t[j].Maths
			}
		case "-Maths":
			if x.t[i].Maths != x.t[j].Maths {
				return x.t[j].Maths < x.t[i].Maths
			}
		case "English":
			if x.t[i].English != x.t[j].English {
				return x.t[i].English < x.t[j].English
			}
		case "-English":
			if x.t[i].English != x.t[j].English {
				return x.t[j].English < x.t[i].English
			}
		}
	}
	return false
}

func main() {
	sort.Sort(customSort{scores, []string{"Class", "-Chinese", "Maths", "-Enlish"}})
	printScores(scores)
}

package main

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
)

type Student struct {
	name string
	age  int
}

var students []*Student

func printStudents(students []*Student) {
	const format = "%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "name", "age")
	fmt.Fprintf(tw, format, "----", "---")
	for _, s := range students {
		fmt.Fprintf(tw, format, s.name, s.age)
	}
	tw.Flush()
}

func init() {
	students = []*Student{
		&Student{"Adam", 20},
		&Student{"Bob", 18},
		&Student{"Clark", 19},
		&Student{"Daisy", 18},
		&Student{"Eva", 20},
		&Student{"Frank", 20},
		&Student{"Gideon", 19},
	}
}

type studentSort struct {
	s    []*Student
	less func(x, y *Student) bool
}

func (x studentSort) Len() int           { return len(x.s) }
func (x studentSort) Less(i, j int) bool { return x.less(x.s[i], x.s[j]) }
func (x studentSort) Swap(i, j int)      { x.s[i], x.s[j] = x.s[j], x.s[i] }

func main() {
	if len(os.Args[1:]) > 0 {
		switch os.Args[1] {
		case "stable":
			stable()
		case "unstable":
			unstable()
		case "best":
			best()
		}
	}
	printStudents(students)
}

func stable() {
	sort.Sort(studentSort{students, func(x, y *Student) bool {
		if x.name != y.name {
			return x.name < y.name
		}
		return false
	}})

	sort.Stable(studentSort{students, func(x, y *Student) bool {
		if x.age != y.age {
			return x.age < y.age
		}
		return false
	}})
}

func unstable() {
	sort.Sort(studentSort{students, func(x, y *Student) bool {
		if x.name != y.name {
			return x.name < y.name
		}
		return false
	}})

	sort.Sort(studentSort{students, func(x, y *Student) bool {
		if x.age != y.age {
			return x.age < y.age
		}
		return false
	}})
}

func best() {
	sort.Sort(studentSort{students, func(x, y *Student) bool {
		if x.age != y.age {
			return x.age < y.age
		}
		if x.name != y.name {
			return x.name < y.name
		}
		return false
	}})
}

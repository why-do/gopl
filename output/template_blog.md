# 4.6 文本和 HTML 模板
进行简单的格式化输出，使用fmt包就足够了。但是要实现更复杂的格式化输出，并且有时候还要求格式和代码彻底分离。这可以通过 text/templat 包和 html/template 包里的方法来实现，通过这两个包，可以将程序变量的值代入到模板中。  

## 模板表达式
模板，是一个字符串或者文件，它包含一个或者多个两边用双大括号包围的单元，这称为**操作**。大多数字符串是直接输出的，但是操作可以引发其他的行为。  
每个操作在模板语言里对应一个表达式，功能包括：
+ 输出值
+ 选择结构体成员
+ 调用函数和方法
+ 描述控逻辑
+ 实例化其他的模板

继续使用 GitHub 的 issue 接口返回的数据，这次使用模板来输出。一个简单的字符串模板如下所示：
```go
const templ = `{{.TotalCount}} issues:
{{range .Items}}----------------------------------------
Number: {{.Number}}
User:   {{.User.Login}}
Title:  {{.Title | printf "%.64s"}}
Age:    {{.CreatedAt | daysAgo}} days
{{end}}`
```
点号（\.）表示当前值的标记。最开始的时候表示模板里的参数，也就是 github.IssuesSearchResult。  
操作 {{.TotalCount}} 就是 TotalCount 字段的值。  
{{range .Items}} 和 {{end}} 操作创建一个循环，这个循环内部的点号（\.）表示Items里的每一个元素。  
在操作中，管道符（|）会将前一个操作的结果当做下一个操作的输入，这个和UNIX里的管道类似。  
`{{.Title | printf "%.64s"}}`，这里的第二个操作是printf函数，在包里这个名称对应的就是fmt.Sprintf，所以会按照fmt.Sprintf函数返回的样式输出。  
`{{.CreatedAt | daysAgo}}`，这里的第二个操作数是 daysAgo，这是一个自定义的函数，具体如下：
```go
func daysAgo(t time.Time) int {
	return int(time.Since(t).Hours() / 24)
}
```

## 模板输出的过程
通过模板输出结果需要两个步骤：
1. 解析模板并转换为内部表示的方法
2. 在指定的输入上执行（就是执行并输出）

解析模板只需要执行一次。下面的代码创建并解析上面定义的文本模板：
```go
report, err := template.New("report").
    Funcs(template.FuncMap{"daysAgo": daysAgo}).
    Parse(templ)
if err != nil {
    log.Fatal(err)
}
```
这里使用了方法的链式调用。template.New 函数创建并返回一个新的模板。  
Funcs 方法将自定义的 daysAgo 函数到内部的函数列表中。之前提到的printf实际对应的是fmt.Sprintf，也是在包内默认就已经在这个函数列表里了。如果有更多的自定义函数，就多次调用这个方法添加。  
最后就是调用Parse进行解析。  
上面的代码完成了创建模板，添加内部可调用的 daysAgo 函数，解析（Parse方法），检查（检查err是否为空）。现在就可以调用report的 Execute 方法，传入数据源（github.IssuesSearchResult，这个需要先调用github.SearchIssues函数来获取），并指定输出目标（使用 os.Stdout）：
```go
if err := report.Execute(os.Stdout, result); err != nil {
    log.Fatal(err)
}
```

之前的代码比较凌乱，下面出完整可运行的代码：
```go
package main

import (
	"log"
	"os"
	"text/template"
	"time"

	"gopl/ch4/github"
)

const templ = `{{.TotalCount}} issues:
{{range .Items}}----------------------------------------
Number: {{.Number}}
User:   {{.User.Login}}
Title:  {{.Title | printf "%.64s"}}
Age:    {{.CreatedAt | daysAgo}} days
{{end}}`

// 自定义输出格式的方法
func daysAgo(t time.Time) int {
	return int(time.Since(t).Hours() / 24)
}

func main() {
	// 解析模板
	report, err := template.New("report").
		Funcs(template.FuncMap{"daysAgo": daysAgo}).
		Parse(templ)
	if err != nil {
		log.Fatal(err)
	}
	// 获取数据
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	// 输出
	if err := report.Execute(os.Stdout, result); err != nil {
		log.Fatal(err)
	}
}
```
这个版本还可以改善，下面对解析错误的处理进行了改进

## 帮助函数 Must
由于目标通常是在编译期间就固定下来的，因此无法解析将会是一个严重的bug。上面的版本如果无法解析（去掉个大括号试试），只会以比较温和的方式报告出来。  
这里推荐使用帮助函数 template.Must，模板错误会Panic：
```go
package main

import (
	"log"
	"os"
	"text/template"
	"time"

	"gopl/ch4/github"
)

const templ = `{{.TotalCount}} issues:
{{range .Items}}----------------------------------------
Number: {{.Number}}
User:   {{.User.Login}}
Title:  {{.Title | printf "%.64s"}}
Age:    {{.CreatedAt | daysAgo}} days
{{end}}`

// 自定义输出格式的方法
func daysAgo(t time.Time) int {
	return int(time.Since(t).Hours() / 24)
}

// 使用帮助函数
var report = template.Must(template.New("issuelist").
	Funcs(template.FuncMap{"daysAgo": daysAgo}).
	Parse(templ))

func main() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	if err := report.Execute(os.Stdout, result); err != nil {
		log.Fatal(err)
	}
}
```
和上个版本的区别就是解析的过程外再包了一层 template.Must 函数。而效果就是原本解析错误是调用 `log.Fatal(err)` 来退出，这个调用也是自己的代码里指定的。  
而现在是调用 `panic(err)` 来退出，并且会看到一个更加严重的错误报告（错误信息是一样的），并且这个也是包内部提供的并且推荐的做法。  
最后是输出的结果：
```
PS H:\Go\src\gopl\ch4\issuesreport> go run main.go repo:golang/go is:open json decoder tag
6 issues:
----------------------------------------
Number: 28143
User:   Carpetsmoker
Title:  proposal: encoding/json: add "readonly" tag
Age:    135 days
----------------------------------------
Number: 14750
User:   cyberphone
Title:  encoding/json: parser ignores the case of member names
Age:    1079 days
----------------------------------------
...
```

# HTML 模板
接着看 html/template 包。它使用和 text/template 包里一样的 API 和表达式语法，并且额外地对出现在 HTML、JavaScript、CSS 和 URL 中的字符串进行自动转义。这样可以避免在生成 HTML 是引发一些安全问题。  

## 使用模板输出页面
下面是一个将 issue 输出为 HTML 表格代码。由于两个包里的API是一样的，所以除了模板本身以外，GO代码没有太大的差别：
```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

import (
	"gopl/ch4/github"
	"html/template"
)

var issueList = template.Must(template.New("issuelist").Parse(`
<h1>{{.TotalCount}} issues</h1>
<table>
<tr style='text-align: left'>
  <th>#</th>
  <th>State</th>
  <th>User</th>
  <th>Title</th>
</tr>
{{range .Items}}
<tr>
  <td><a href='{{.HTMLURL}}'>{{.Number}}</a></td>
  <td>{{.State}}</td>
  <td><a href='{{.User.HTMLURL}}'>{{.User.Login}}</a></td>
  <td>{{.Title}}</td>
</tr>
{{end}}
</table>
`))

func main() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("http://localhost:8000")
	handler := func(w http.ResponseWriter, r *http.Request) {
		showIssue(w, result)
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func showIssue(w http.ResponseWriter, result *github.IssuesSearchResult) {
	if err := issueList.Execute(w, result); err != nil {
		log.Fatal(err)
	}
}
```

## template.HTML 类型
通过模板的操作导入的字符串，默认都会按照原样显示出来。就是会把HTML的特殊字符自动进行转义，效果就是无法通过模板导入的内容生成html标签。  
如果就是需要通过模板的操作再导入一些HTML的内容，就需要使用 template.HTML 类型。使用 template.HTML 类型后，可以避免模板自动转义受信任的 HTML 数据。同样的类型还有 template.CSS、template.JS、template.URL 等，具体可以查看源码。  
下面的操作演示了普通的 string 类型和 template.HTML 类型在导入一个 HTML 标签后显示效果的差别：
```go
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	const templ = `<p>A: {{.A}}</p><p>B: {{.B}}</p>`
	t := template.Must(template.New("escape").Parse(templ))
	var data struct {
		A string        // 不受信任的纯文本
		B template.HTML // 受信任的HTML
	}
	data.A = "<b>Hello!</b>"
	data.B = "<b>Hello!</b>"

	fmt.Println("http://localhost:8000")
	handler := func(w http.ResponseWriter, r *http.Request) {
		if err := t.Execute(w, data); err != nil {
			log.Fatal(err)
		}
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
```
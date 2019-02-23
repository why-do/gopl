# 4.5 JSON
JSON是一种发送和接收格式化信息的标准。JSON不是唯一的标准，XML、ASN.1 和 Google 的 Protocol Buffer 都是相似的标准。Go通过标准库 encoding/json、encoding/xml、encoding/asn1 和其他的库对这些格式的编码和解码提供了非常好的支持，这些库都拥有相同的API。  

## 序列化输出
首先定义一组数据：
```go
type Movie struct {
	Title  string
	Year   int  `json:"released"`
	Color  bool `json:"color,omitempty"`
	Actors []string
}

var movies = []Movie{
	{Title: "Casablanca", Year: 1942, Color: false,
		Actors: []string{"Humphrey Bogart", "Ingrid Bergman"}},
	{Title: "Cool Hand Luke", Year: 1967, Color: true,
		Actors: []string{"Paul Newman"}},
	{Title: "Bullitt", Year: 1968, Color: true,
		Actors: []string{"Steve McQueen", "Jacqueline Bisset"}},
}
```

然后通过 json.Marshal 进行编码：
```go
data, err := json.Marshal(movies)
if err != nil {
	log.Fatalf("JSON Marshal failed: %s", err)
}
fmt.Printf("%s\n", data)

/* 执行结果
[{"Title":"Casablanca","released":1942,"Actors":["Humphrey Bogart","Ingrid Bergman"]},{"Title":"Cool Hand Luke","released":1967,"color":true,"Actors":["Paul Newman"]},{"Title":"Bullitt","released":1968,"color":true,"Actors":["Steve McQueen","Jacqueline Bisset"]}]
*/
```

这种紧凑的表示方法适合传输，但是不方便阅读。有一个 json.MarshalIndent 的变体可以输出整齐格式化过的结果。多传2个参数，第一个是定义每行输出的前缀字符串，第二个是定义缩进的字符串：
```go
data, err := json.MarshalIndent(movies, "", "    ")
if err != nil {
	log.Fatalf("JSON Marshal failed: %s", err)
}
fmt.Printf("%s\n", data)

/* 执行结果
[
    {
        "Title": "Casablanca",
        "released": 1942,
        "Actors": [
            "Humphrey Bogart",
            "Ingrid Bergman"
        ]
    },
    {
        "Title": "Cool Hand Luke",
        "released": 1967,
        "color": true,
        "Actors": [
            "Paul Newman"
        ]
    },
    {
        "Title": "Bullitt",
        "released": 1968,
        "color": true,
        "Actors": [
            "Steve McQueen",
            "Jacqueline Bisset"
        ]
    }
]
*/
```

只有可导出的成员可以转换为JSON字段，上面的例子中用的都是大写。  
**成员标签**（field tag），是结构体成员的编译期间关联的一些元素信息。标签值的第一部分指定了Go结构体成员对应的JSON中字段的名字。  
另外，Color标签还有一个额外的选项 omitempty，它表示如果这个成员的值是零值或者为空，则不输出这个成员到JSON中。所以Title为"Casablanca"的JSON里没有color。  

## 反序列化
反序列化操作将JSON字符串解码为Go数据结构。这个是由 json.Unmarshal 实现的。
```go
var titles []struct{ Title string }
if err := json.Unmarshal(data, &titles); err != nil {
	log.Fatalf("JSON unmarshaling failed: %s", err)
}
fmt.Println(titles)

/* 执行结果
[{Casablanca} {Cool Hand Luke} {Bullitt}]
*/
```
这里接收数据时定义的结构体只有一个Title字段，这样当函数 Unmarshal 调用完成后，将填充结构体切片中的 Title 值，而JSON中其他的字段就丢弃了。  

# Web 应用
很多的 Web 服务器都提供 JSON 接口，通过发送HTTP请求来获取想要得到的JSON信息。下面通过查询Github提供的 issue 跟踪接口来演示一下。  

## 定义结构体
首先，定义好类型，顺便还有常量：
```go
// ch4/github/github.go
// https://api.github.com/ 提供了丰富的接口
// 提供查询GitHub的issue接口的API
// GitHub上有详细的API使用说明：https://developer.github.com/v3/search/#search-issues-and-pull-requests
package github

import "time"

const IssuesURL = "https://api.github.com/search/issues"

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}

type Issue struct {
	Number   int
	HTMLURL  string `json:"html_url"`
	Title    string
	State    string
	User     *User
	CreateAt time.Time `json:"created_at"`
	Body     string    // Markdown 格式
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}
```
关于字段名称，即使对应的JSON字段的名称都是小写的，但是结构体中的字段必须首字母大写（*不可导出的字段也无法把JSON数据导入*）。这种情况很普遍，这里可以偷个懒。在 Unmarshal 阶段，JSON字段的名称关联到Go结构体成员的名称是忽略大小写的，这里也不需要考虑序列化的问题，所以很多地方都不需要写成员标签。不过，小写的变量在需要分词的时候，可能会使用下划线分割，这种情况下，还是要用一下成员标签的。  
这里也是选择性地对JSON中的字段进行解码，因为相对于这里演示的内容，GitHub的查询返回的信息是相当多的。  

## 请求获取JSON并解析
函数 SearchIssues 发送HTTP请求并将返回的JSON字符串进行解析。  
关于Get请求的参数，参数中可能会出现URL格式里的特殊字符，比如 ?、&。因此要使用 url.QueryEscape 函数进行转义。  
```go
// ch4/github/search.go
package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// 查询GitHub的issue接口
func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	q := url.QueryEscape(strings.Join(terms, " "))
	resp, err := http.Get(IssuesURL + "?q=" + q)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var result IssuesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
```

**流式解码器**  
之前是使用 json.Unmarshal 进行解码，而这里使用流式解码器。它可以依次从字节流中解码出多个JSON实体，不过这里没有用到该功能。另外还有对应的 json.Encoder 的流式编码器。  
调用 Decode 方法后，就完成了对变量 result 的填充。  

## 调用执行
最后就是将 result 中的内容进行格式化输出，这里用了固定宽度的方法将结果输出为类似表格的形式：
```go
// ch4/issues/main.go
// 将符合条件的issue输出为一个表格
package main

import (
	"fmt"
	"gopl/ch4/github"
	"log"
	"os"
)

func main() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issue: \n", result.TotalCount)
	for _, item := range result.Items {
		fmt.Printf("#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title)
	}
}
```

使用命令行参数指定搜索条件，该命令搜索 Go 项目里的 issue 接口，查找 open 状态的列表。由于返回的还是很多，后面的参数是对内容再进行筛选：
```
PS H:\Go\src\gopl\ch4\issues> go run main.go repo:golang/go is:open json decoder tag
6 issue:
#28143 Carpetsmo proposal: encoding/json: add "readonly" tag
#14750 cyberphon encoding/json: parser ignores the case of member names
#17609 nathanjsw encoding/json: ambiguous fields are marshalled
#22816 ganelon13 encoding/json: include field name in unmarshal error me
#19348 davidlaza cmd/compile: enable mid-stack inlining
#19109  bradfitz proposal: cmd/go: make fuzzing a first class citizen, l
PS H:\Go\src\gopl\ch4\issues>
```

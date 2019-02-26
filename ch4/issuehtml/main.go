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
	// 建议用这些参数调用执行：go run main.go repo:golang/go is:open json decoder tag
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

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
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string    // Markdown 格式
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

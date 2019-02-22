package github

import "time"

const IssuesURL = "https://api.gihub.com/search/issues"

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
	CreateAt time.Time `json:"crete_at"`
	Body     string    // Markdown 格式
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

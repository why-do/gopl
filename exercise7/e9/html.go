package main

import (
	"html/template"
	"log"
	"net/http"
)

var musicList = template.Must(template.New("musiclist").Parse(`
<h1>Track List</h1>
<table border=1>
<tr style='text-align: left'>
  <th><a href='/?o=Title'>Title</a></th>
  <th><a href='/?o=Artist'>Artist</a></th>
  <th><a href='/?o=Album'>Album</a></th>
  <th><a href='/?o=Year'>Year</a></th>
  <th><a href='/?o=Length'>Length</a></th>
</tr>
{{range .}}
<tr>
  <td>{{.Title}}</td>
  <td>{{.Artist}}</td>
  <td>{{.Album}}</td>
  <td>{{.Year}}</td>
  <td>{{.Length}}</td>
</tr>
{{end}}
</table>
`))

func showTracks(w http.ResponseWriter, tracks []*Track) {
	if err := musicList.Execute(w, tracks); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database struct {
	Items map[string]dollars
	sync.RWMutex
}

func (db *database) list(w http.ResponseWriter, req *http.Request) {
	// db.RLock()
	// defer db.RUnlock()
	// for item, price := range db.Items {
	// 	fmt.Fprintf(w, "%s: %s\n", item, price)

	// }
	showItem(w, db)
}

func (db *database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	db.RLock()
	defer db.RUnlock()
	price, ok := db.Items[item]
	if !ok {
		http.Error(w, fmt.Sprintf("no such item: %q\n", item), http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "%s\n", price)
}

// 从 URL 解析获取item和price
func getItemPrice(req *http.Request) (string, dollars, error) {
	item := req.URL.Query().Get("item")
	if item == "" {
		return "", 0, errors.New("item not get")
	}
	priceStr := req.URL.Query().Get("price")
	if priceStr == "" {
		return item, 0, errors.New("price not get")
	}
	price64, err := strconv.ParseFloat(priceStr, 32)
	price := dollars(price64)
	if err != nil {
		return item, price, fmt.Errorf("Parse Price: %v\n", err)
	}
	return item, price, err
}

func (db *database) add(w http.ResponseWriter, req *http.Request) {
	item, price, err := getItemPrice(req)
	if err != nil {
		http.Error(w, fmt.Sprintln(err), http.StatusNotFound)
		return
	}
	db.Lock()
	defer db.Unlock()
	if _, ok := db.Items[item]; ok {
		http.Error(w, fmt.Sprintf("%s is already exist.\n", item), http.StatusNotFound)
		return
	}
	db.Items[item] = dollars(price)
	fmt.Fprintf(w, "success add %s: %s\n", item, dollars(price))
}

func (db *database) update(w http.ResponseWriter, req *http.Request) {
	item, price, err := getItemPrice(req)
	if err != nil {
		http.Error(w, fmt.Sprintln(err), http.StatusNotFound)
		return
	}
	db.Lock()
	defer db.Unlock()
	if _, ok := db.Items[item]; !ok {
		http.Error(w, fmt.Sprintf("%s is not exist.\n", item), http.StatusNotFound)
		return
	}
	db.Items[item] = dollars(price)
	fmt.Fprintf(w, "success udate %s: %s\n", item, dollars(price))
}

func (db *database) delete(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	func() {
		db.Lock()
		defer db.Unlock()
		delete(db.Items, item)
	}()
	db.list(w, req)
}

func main() {
	db := database{
		Items: map[string]dollars{"shoes": 50, "socks": 5},
	}
	fmt.Println("http://localhost:8000/list")
	fmt.Println("http://localhost:8000/price?item=shoes")
	fmt.Println("http://localhost:8000/add?item=football&price=11")
	fmt.Println("http://localhost:8000/update?item=football&price=12.35")
	fmt.Println("http://localhost:8000/delete?item=shoes")
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)
	http.HandleFunc("/add", db.add)
	http.HandleFunc("/update", db.update)
	http.HandleFunc("/delete", db.delete)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

var itemList = template.Must(template.New("itemlist").Parse(`
<body>
<h3>Items</h3>
<table border="1">
<tr style='text-align: left'>
  <th>Item</th>
  <th>Prise</th>
</tr>
{{range $i, $v :=  .Items}}
<tr>
  <td>{{$i}}</td>
  <td>{{$v}}</td>
</tr>
{{end}}
</table>
</body>
`))

func showItem(w http.ResponseWriter, db *database) {
	db.RLock()
	defer db.RUnlock()
	if err := itemList.Execute(w, db); err != nil {
		log.Fatal(err)
	}
}

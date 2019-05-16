package memo

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

func httpGetBody(url string) (interface{}, error) {
	log.Printf("httpGetBody: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

var urls = []string{ // 换一批慢一些的网站，加载时间1s左右的国外资源
	"https://github.com/adonovan/gopl.io/tree/master/ch9",
	"https://www.djangoproject.com/",
	"https://getbootstrap.com/",
	"https://www.python.org/",
}

func TestSequential(t *testing.T) { // 串行
	m := New(httpGetBody)
	urls = append(urls, urls...) // 每个 URL 请求两次
	for _, url := range urls {
		start := time.Now()
		value, err := m.Get(url)
		if err != nil {
			log.Print(err)
		}
		t.Logf("%s, %s, %d bytes\n", url, time.Since(start), len(value.([]byte)))
	}
}

func TestConcurrent(t *testing.T) { // 并行
	m := New(httpGetBody)
	var n sync.WaitGroup
	urls = append(urls, urls...) // 每个 URL 请求两次
	n.Add(len(urls))
	for _, url := range urls {
		go func(url string) {
			defer n.Done()
			start := time.Now()
			value, err := m.Get(url)
			if err != nil {
				log.Print(err)
			}
			t.Logf("%s, %s, %d bytes\n", url, time.Since(start), len(value.([]byte)))
		}(url)
	}
	n.Wait()
}

package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	memo "gopl/exercise9/e3/memo5"
)

var urls = []string{ // 这些国外的网站，加载时间相对会比较长
	"https://github.com/adonovan/gopl.io/tree/master/ch9",
	"https://www.djangoproject.com/",
	"https://getbootstrap.com/",
	"https://www.python.org/",
}

func httpGetBody(url string) (interface{}, error) {
	log.Printf("httpGetBody: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func main() {
	// 闭包的方式计算函数运行时间
	defer func() func() {
		start := time.Now()
		return func() {
			log.Printf("总共运行时间1: %s", time.Since(start))
		}
	}()()
	// 直接获取开始的时间然后传参，更好理解更简洁
	defer func(start time.Time) {
		log.Printf("总共运行时间2: %s", time.Since(start))
	}(time.Now())
	
	// 当检测到输入时，广播取消
	done := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1)) // 读一个字节
		close(done)
	}()

	m := memo.New(httpGetBody, done)
	defer m.Close()
	var n sync.WaitGroup
	urls = append(urls, urls...) // 每个 URL 请求两次
	n.Add(len(urls))
	for _, url := range urls {
		go func(url string) {
			defer n.Done()
			start := time.Now()
			value, err := m.Get(url)
			if err != nil {
				log.Printf("%s: %v", url, err)
				return
			}
			var lenth int
			if v, ok := value.([]byte); ok {
				lenth = len(v)
			}
			log.Printf("%s, %s, %d bytes\n", url, time.Since(start), lenth)
		}(url)
	}
	n.Wait()
}

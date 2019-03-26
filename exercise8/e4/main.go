package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup // 工作 goroutine 的个数

func echo(c net.Conn, shout string, delay time.Duration) {
	defer wg.Done()
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func handleConn(c net.Conn) {
	input := bufio.NewScanner(c)
	for input.Scan() {
		wg.Add(1)
		go echo(c, input.Text(), 2*time.Second)
	}
	// 注意：忽略 input.Err() 中可能的错误
	wg.Wait()
	c.Close()
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // 例如，连接终止
			continue
		}
		go handleConn(conn) // 并发处理连接
	}
}

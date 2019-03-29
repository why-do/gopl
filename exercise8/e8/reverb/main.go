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
	stop1 := make(chan struct{})
	stop2 := make(chan struct{})

	inputSignal := make(chan struct{}) // 有任何输入，就发送一个信号
	go func() { // 接收客户端发送回声的goroutine
		input := bufio.NewScanner(c)
		for input.Scan() { // 注意：忽略 input.Err() 中可能的错误
			inputSignal <- struct{}{}
			wg.Add(1)
			go echo(c, input.Text(), 1*time.Second)
		}
		// 退出上面的for循环，表示客户端断开
		stop1 <- struct{}{}
	}()

	delay := 5 * time.Second
	timer := time.NewTimer(delay)
	go func() { // 计算超时的goroutine
		for {
			select {
			case <-inputSignal:
				timer.Reset(delay)
			case <-timer.C:
				// 超时，断开连接
				stop2 <- struct{}{}
				return
			}
		}
	}()

	select {
	case <-stop1:
	case <-stop2:
	}
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

package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn) // 注意：忽略错误
		log.Println("done")
		done <- struct{}{} // 通知主 goroutine 的信号
	}()
	mustCopy(conn, os.Stdin)
	conn.CloseWrite()
	<-done // 等待后台 goroutine 完成
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

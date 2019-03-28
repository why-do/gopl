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

	done1 := make(chan struct{})
	go func() { // 打印回声的goroutine
		io.Copy(os.Stdout, conn) // 注意：忽略错误
		log.Println("done")
		done1 <- struct{}{} // 通知主 goroutine 的信号
	}()

	done2 := make(chan struct{})
	go func() { // 发送请求的goroutine
		mustCopy(conn, os.Stdin)
		conn.CloseWrite()
		done2 <- struct{}{}
	}()

	select { // 等待后台 goroutine 完成
	case <-done1:
	case <-done2: // 客户端主动断开后，值关闭写半边连接
		<-done1 // 继续等待服务端断开，就是等待全是打印完毕
	}
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

var port int
var tz string

func init() {
	flag.IntVar(&port, "port", 8000, "端口号")
	flag.StringVar(&tz, "tz", "", "时区")
}

func main() {
	flag.Parse()
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
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

func handleConn(c net.Conn) {
	defer c.Close()
	var loc *time.Location
	if tz != "" {
		loc, _ = time.LoadLocation(tz) // 忽略错误
	}
	for {
		var err error
		if loc == nil {
			_, err = io.WriteString(c, time.Now().Format("2006/01/02 15:04:05\r\n"))
		} else {
			_, err = io.WriteString(c, time.Now().In(loc).Format("2006/01/02 15:04:05\r\n"))
		}
		
		if err != nil {
			return // 例如，连接断开
		}
		time.Sleep(1 * time.Second)
	}
}

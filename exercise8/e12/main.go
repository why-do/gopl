package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

// 广播器
type client chan<- string // 对外发送消息的通道
type clientInfo struct {
	name string
	ch   client
}

var (
	entering = make(chan clientInfo)
	leaving  = make(chan clientInfo)
	messages = make(chan string) // 所有接受的客户消息
)

func broadcaster() {
	clients := make(map[clientInfo]bool) // 所有连接的客户端集合
	for {
		select {
		case msg := <-messages:
			// 把所有接收的消息广播给所有的客户
			// 发送消息通道
			for cli := range clients {
				cli.ch <- msg
			}
		case cli := <-entering:
			// 在每一个新用户到来的时候，通知当前存在的用户
			var users []string
			for cli := range clients {
				users = append(users, cli.name)
			}
			if len(users) > 0 {
				cli.ch <- fmt.Sprintf("Other users in room: %s", strings.Join(users, "; "))
			} else {
				cli.ch <- "You are the only user in this room."
			}

			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli.ch)
		}
	}
}

// 客户端处理函数
func handleConn(conn net.Conn) {
	ch := make(chan string) // 对外发送客户消息的通道
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	cli := clientInfo{who, ch}       // 打包好用户名和通道
	ch <- "You are " + who           // 这条单发给自己
	messages <- who + " has arrived" // 这条进行进行广播，但是自己还没加到广播列表中
	entering <- cli                  // 然后把自己加到广播列表中

	input := bufio.NewScanner(conn)
	for input.Scan() {
		if len(strings.TrimSpace(input.Text())) == 0 { // 禁止发送纯空白字符
			continue
		}
		messages <- who + ": " + input.Text()
	}
	// 注意，忽略input.Err()中可能的错误

	leaving <- cli
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		// windows 需要 \r 了正常显示
		fmt.Fprintln(conn, msg+"\r") // 注意，忽略网络层面的错误
	}
}

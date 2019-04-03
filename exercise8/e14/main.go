package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
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

type registeInfo struct {
	name string
	ch   chan<- bool
}

var register = make(chan registeInfo) // 注册用户名的通道

func broadcaster() {
	clients := make(map[string]client) // 所有连接的客户端集合
	for {
		select {
		case msg := <-messages:
			// 把所有接收的消息广播给所有的客户
			// 发送消息通道
			for _, cli := range clients {
				cli <- msg
			}
		case user := <-register:
			// 先判断新用户名是否有重复
			_, ok := clients[user.name]
			user.ch <- !ok
		case cliSt := <-entering:
			// 在每一个新用户到来的时候，通知当前存在的用户
			var users []string
			for user := range clients {
				users = append(users, user)
			}
			if len(users) > 0 {
				cliSt.ch <- fmt.Sprintf("Other users in room: %s", strings.Join(users, "; "))
			} else {
				cliSt.ch <- "You are the only user in this room."
			}

			clients[cliSt.name] = cliSt.ch
		case cliSt := <-leaving:
			delete(clients, cliSt.name)
			close(cliSt.ch)
		}
	}
}

// 客户端处理函数
func handleConn(conn net.Conn) {
	who := clientRegiste(conn) // 注册获取用户名

	ch := make(chan string) // 对外发送客户消息的通道
	go clientWriter(conn, ch)

	// who := conn.RemoteAddr().String() // 去掉这一行
	cli := clientInfo{who, ch}       // 打包好用户名和通道
	ch <- "You are " + who           // 这条单发给自己
	messages <- who + " has arrived" // 现在这条广播自己也能收到
	entering <- cli

	done := make(chan struct{}, 2) // 等待下面两个 goroutine 其中一个执行完成。使用缓冲通道防止 goroutine 泄漏
	// 计算超时的goroutine
	inputSignal := make(chan struct{}) // 有任何输入，就发送一个信号
	timeout := 100 * time.Second       // 客户端空闲的超时时间
	go func() {
		timer := time.NewTimer(timeout)
		for {
			select {
			case <-inputSignal:
				timer.Reset(timeout)
			case <-timer.C:
				// 超时，断开连接
				done <- struct{}{}
				return
			}
		}
	}()

	go func() {
		input := bufio.NewScanner(conn)
		for input.Scan() {
			inputSignal <- struct{}{}
			if len(strings.TrimSpace(input.Text())) == 0 { // 禁止发送纯空白字符
				continue
			}
			messages <- who + ": " + input.Text()
		}
		// 注意，忽略input.Err()中可能的错误
		done <- struct{}{}
	}()

	<-done
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

// 注册用户名
func clientRegiste(conn net.Conn) (who string) {
	ch := make(chan bool)
	fmt.Fprint(conn, "input nickname: ") // 注意，忽略网络层面的错误
	input := bufio.NewScanner(conn)
	for input.Scan() {
		if len(strings.TrimSpace(input.Text())) == 0 { // 禁止发送纯空白字符
			continue
		}
		who = input.Text()
		register <- registeInfo{who, ch}
		if <-ch {
			break
		}
		fmt.Fprintf(conn, "name %q is existed\r\ntry other name: ", who)
	}
	// 注意，忽略input.Err()中可能的错误
	return who
}
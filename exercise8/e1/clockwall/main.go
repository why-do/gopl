package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var repo []string

func main() {
	for i, arg := range os.Args[1:] {
		repo = append(repo, "")
		go wall(arg, i)
	}
	for {
		fmt.Println(strings.Join(repo, "\t"))
		time.Sleep(time.Second)
	}
}

func wall(info string, i int) {
	infos := strings.Split(info, "=")
	if len(infos) != 2 {
		return
	}
	tz, addr := infos[0], infos[1]
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	input := bufio.NewScanner(conn)
	for input.Scan() {
		repo[i] = fmt.Sprintf("%s: %s", tz, input.Text())
	}
}

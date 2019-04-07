package main

import (
	"fmt"
	"runtime"
	"time"
)

type data struct {
	start time.Time
	id    int
	num   int
}

func main() {
	var in, out chan data
	lineInput := make(chan int) // 流水线的入口
	done := make(chan struct{})

	out = make(chan data)
	go func(out chan data) {
		for {
			id := <-lineInput
			x := data{
				start: time.Now(),
				id:    id,
			}
			out <- x
		}
	}(out)

	var id int // 每个 goroutine 有一个唯一 id
	for {
		new := make(chan data)
		in, out = out, new
		go func(id int, in, out chan data) {
			for {
				x := <-in
				switch {
				case x.id < id:
					panic("不应该进入这个分支")
				case x.id == id: // 自己创建的数据打印出来
					fmt.Printf("%s id: %d num: %d\n", time.Since(x.start), x.id, x.num)
					done <- struct{}{}
				case x.id > id: // id更大的 goroutine 创建的数据，则自增后继续往后传
					x.num++
					out <- x
				}
			}
		}(id, in, out)

		// 获取内存使用情况
		mem := runtime.MemStats{}
		runtime.ReadMemStats(&mem)
		if mem.Sys > 1024*1024*1024*4 { // 超过 4G 停止继续创建
			// 这个空间不是系统使用的空间大小，而是 Go 使用的空间大小。就是程序运行后内存的增量
			fmt.Println("id:", id, "mem:", mem.Sys)
			lineInput <- id // 停止创建 goroutine 后，向流水线的入口放入一个值
			<-done
			break
		}

		// 打印进度
		if id%10000 == 0 {
			fmt.Println("id:", id, "mem:", mem.Sys)
		}
		id++
	}
}

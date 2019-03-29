package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var done = make(chan struct{})

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

func main() {
	roots := os.Args[1:]
	if len(roots) == 0 {
		roots = []string{"."}
	}

	// 当检测到输入时，广播取消
	go func() {
		os.Stdin.Read(make([]byte, 1)) // 读一个字节
		close(done)
	}()

	fileSizes := make(chan int64)
	var n sync.WaitGroup
	for _, root := range roots {
		n.Add(1)
		go walkDir(root, &n, fileSizes)
	}
	go func() {
		n.Wait()
		close(fileSizes)
	}()

	tick := time.Tick(500 * time.Millisecond)
	var nfiles, nbytes int64
loop:
	for {
		select {
		case <-done:
			// 耗尽 fileSizes，让已经创建的 goroutine 结束
			for range fileSizes {
				// 什么也不做
			}
			return
		case siez, ok := <-fileSizes:
			if !ok {
				break loop
			}
			nfiles++
			nbytes += siez
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		}
	}
	printDiskUsage(nfiles, nbytes)
}

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files   %.2f GB\n", nfiles, float64(nbytes)/(1<<30)) // 1<<30 就是 2**30 就是 1024*1024*1024
}

// walkDir 递归地遍历以 dir 为根目录的整个文件树
// 并在 fileSizes 上发送每个已找到的文件的大小
func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
	defer n.Done()
	if cancelled() {
		return
	}
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			walkDir(subdir, n, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}

// 用于限制目录并发数的计数信号量
var sema = make(chan struct{}, 20)

// dirents 返回 dir 目录中的条目
func dirents(dir string) []os.FileInfo {
	select {
	case sema <- struct{}{}: // 获取令牌
	case <-done:
		return nil // 取消
	}
	defer func() { <-sema }() // 释放令牌
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	}
	return entries
}

package main

import (
	"fmt"
	"io"
	"os"
)

type WriterCounter struct {
	Writer io.Writer
	count  *int64
}

func (c *WriterCounter) Write(p []byte) (int, error) {
	n, err := c.Writer.Write(p)
	*c.count += int64(n)
	return n, err
}

func CountingWriter(w io.Writer) (io.Writer, *int64) {
	var count int64
	wc := WriterCounter{w, &count}
	return &wc, &count
}

func main() {
	wc, count := CountingWriter(os.Stdout)
	wc.Write([]byte("Hello"))
	fmt.Println()
	fmt.Println(*count)
	wc.Write([]byte("World"))
	fmt.Println()
	fmt.Println(*count)
}

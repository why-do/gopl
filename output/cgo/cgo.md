# 13.4 使用 cgo 调用 C 代码
cgo 是用来为 C 函数创建 Go 绑定的工具。诸如此类的工具都叫作**外部函数接口**（FFI）。  
其他的工具还有，比如SWIG（sig\.org）是另一个工具，它提供了更加复杂的特性用来集成C++的类，这个不讲。

## 使用cgo的场景
如果一个程序已经有现成的C语言的实现，但是还没有Go语言的实现的时候，那没有一下3种选择：
1. 如果是一个比较小的C语言库，可以使用纯 Go 语言来移植它（重新实现一遍）。
2. 如果性能不是很关键，可以用 os/exec 包以辅助子进程的方式来调用C程序。
3. 当需要使用复杂而且性能要求高的底层C接口时，就是使用cgo的场景了

简单说就是，如果是简单的实现，那么就再造一个Go语言的轮子。如果性能要求不高，可以直接通过系统来调用这个程序。只有不想重新造轮子又不想间接的通过系统来调用的时候，就需要用到 cgo 了。  

bzip2 压缩程序正是这样的一个情况。接下来就要使用 cgo 来构建一个简单的数据压缩程序。  
标准库的 compress/... 子包中提供了流行压缩算法的压缩器和解压缩器，包括流行的LZW压缩算法（Unix的compress命令用的算法）和DEFLATE压缩算法（GNU gzip命令用的算法）。这些包中的 API 有些许的不同，但都提供一个对 io\.Writer 的封装用来对写入的数据进行压缩，并且还有一个对 io\.Reader 的封装，在读取数据的同时进行压缩。例如：
```go
package gzip // compress/gzip
func NewWriter(w io.Writer) io.WriteCloser
func NewReader(r io.Reader) (io.ReadCloser, error)
```

bzip2 算法基于优雅的 Burrows-Wheeler 变换，它和 gzip 相比速度要慢但是压缩比更高。标准库的 compress/bzip2 提供了 bzip2 的解压缩器，但是目前还没有提供压缩功能。从头开始实现这个压缩算法比较麻烦，而且 http://bzip.org 已经有现成的libbzip2的开源实现了，这是一个文档完善且高性能的开源 C 语言实现。  

要使用C语言的libbzip2包，需要先构建一个bz_stream结构体，这个结构体包含输入和输出缓冲区，以及三个C函数：
+ BZ2_bzCompressInit: 初始化缓存，分配流的缓冲区
+ BZ2_bzCompress: 将输入缓存的数据压缩到输出缓存
+ BZ2_bzCompressEnd: 释放不需要的缓存

## C代码
可以在Go代码中直接调用BZ2_bzCompressInit和BZ2_bzCompressEnd。  
但是对于BZ2_bzCompress，我们将定义一个C语言的包装函数，用它完成真正的工作。下面是C代码，和其他Go文件放在同一个包下：
```c
// bzip 包中的文件 bzip2.c

// 对 libbzip2 的简单包装，适合 cgo 使用
#include <bzlib.h>

int bz2compress(bz_stream *s, int action,
                char *in, unsigned *inlen, char *out, unsigned *outlen) {
  s->next_in = in;
  s->avail_in = *inlen;
  s->next_out = out;
  s->avail_out = *outlen;
  int r = BZ2_bzCompress(s, action);
  *inlen -= s->avail_in;
  *outlen -= s->avail_out;
  s->next_in = s->next_out = NULL;
  return r;
}
```

**安装gcc**  
可能会出现如下的错误提示：
```
exec: "gcc": executable file not found in %PATH%
```
这个应该是缺少gcc编译器，所以需要安装配置好。在Windows系统上可能要麻烦一点。  

## cgo注释
然后是Go代码，这里只是源码文件开头的部分，第一部分如下所示。  
声明 import "C" 很特别， 并没有这样的一个包，但是这会让编译程序在编译之前先运行cgo工具：
```go
// bzip 包中的文件 bzip2.go 的第一部分

// 包 bzip 封装了一个使用 bzip2 压缩算法的 writer (bzip.org).
package bzip

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -L/usr/lib -lbz2
#include <bzlib.h>
#include <stdlib.h>
bz_stream* bz2alloc() { return calloc(1, sizeof(bz_stream)); }
int bz2compress(bz_stream *s, int action,
                char *in, unsigned *inlen, char *out, unsigned *outlen);
void bz2free(bz_stream* s) { free(s); }
*/
import "C"

import (
	"io"
	"unsafe"
)

type writer struct {
	w      io.Writer // 基本输出流
	stream *C.bz_stream
	outbuf [64 * 1024]byte
}

// NewWriter 对于 bzip2 压缩的流返回一个 writer
func NewWriter(out io.Writer) io.WriteCloser {
	const blockSize = 9
	const verbosity = 0
	const workFactor = 30
	w := &writer{w: out, stream: C.bz2alloc()}
	C.BZ2_bzCompressInit(w.stream, blockSize, verbosity, workFactor)
	return w
}
```
在预处理过程中，cgo 产生一个临时包，这个包里包含了所有C语言的函数和类型对应的Go语言声明。例如 C.bz_stream 和 C.BZ2_bzCompressInit。cgo 工具通过以一种特殊的方式调用C编译器来发现在Go源文件中 `import "C"` 声明之前的注释中包含的C头文件中的内容。  

在cgo注释中还可以包含 #cgo 指令，用来指定C工具链中其他的选项。CFLAGS 和 LDFLAGS 分别对应传给C语言编译器的编译参数和链接器参数，使它们可以从特定目录找到bzlib\.h头文件和libbz2\.a库文件。这个例子假定已经在 `/usr` 目录成功安装了bzip2库。根据个人的安装情况，可以修改或者删除这些标记。（这里还有一个纯C生成的cgo绑定，不依赖bzip2静态库和操作系统的具体环境，具体访问github： https://github.com/chai2010/bzip2 ，这里就顺带提一下现在有更方便的实现方式了）  

NewWriter 调用C函数 BZ2_bzCompressInit 来初始化流的缓冲区。在 writer 结构体中还包含一个额外的缓冲区用来耗尽解压缩器的输出缓冲区。  

## Write方法
下面所示的 Write 方法将未解压的数据写入压缩器中，然后在一个循环中调用 bz2compress 函数，直到所有的数据压缩完毕。Go程序可以访问C的类型（比如 bz_stream、char 和 uint），C的函数（比如 bz2compress），甚至是类似C的预处理宏的对象（比如 BZ_RUN），这些都通过 C.x 的方式来访问。即使类型 C.unit 和 Go 的 uint 长度相同，它们的类型也是不同的：
```go
// bzip 包中的文件 bzip2.go 的第二部分

func (w *writer) Write(data []byte) (int, error) {
	if w.stream == nil {
		panic("closed")
	}
	var total int // 写入的未压缩字节数

	for len(data) > 0 {
		inlen, outlen := C.uint(len(data)), C.uint(cap(w.outbuf))
		C.bz2compress(w.stream, C.BZ_RUN,
			(*C.char)(unsafe.Pointer(&data[0])), &inlen,
			(*C.char)(unsafe.Pointer(&w.outbuf)), &outlen)
		total += int(inlen)
		data = data[inlen:]
		if _, err := w.w.Write(w.outbuf[:outlen]); err != nil {
			return total, err
		}
	}
	return total, nil
}
```
每一次的迭代首先计算传说数据 data 剩余的长度以及输出缓冲 w\.outbuf 的容量。然后把这两个值的地址以及 data 和 w\.outbuf 的地址都传递给 bz2compress 函数。两个长度信息传地址而不传值，这样C函数就可以更新这两个值。这两个值记录的分别是已压缩的数据和压缩后数据的大小。然后把每块压缩后的数据写入底层的 io\.Writer（w\.w\.Write方法）。  

## Close方法
Close方法和Write方法结构类似，通过一个循环将剩余的压缩后的数据从输出缓冲区写入底层：
```go
// bzip 包中的文件 bzip2.go 的第三部分

// Close 方法清空压缩的数据并关闭流
// 它不会关闭底层的 io.Writer
func (w *writer) Close() error {
	if w.stream == nil {
		panic("closed")
	}
	defer func() {
		C.BZ2_bzCompressEnd(w.stream)
		C.bz2free(w.stream)
		w.stream = nil
	}()
	for {
		inlen, outlen := C.uint(0), C.uint(cap(w.outbuf))
		r := C.bz2compress(w.stream, C.BZ_FINISH, nil, &inlen,
			(*C.char)(unsafe.Pointer(&w.outbuf)), &outlen)
		if _, err := w.w.Write(w.outbuf[:outlen]); err != nil {
			return err
		}
		if r == C.BZ_STREAM_END {
			return nil
		}
	}
}
```
压缩完成后，Close 方法最后会调用 C\.BZ2_bzCompressEnd 来释放流缓冲区，这写语句写在 defer 中来确保所有路径返回后都会释放资源。这个时候，w\.stream 指针就不能安全地解引用了，要把它设置为 nil，并且在方法调用的开头添加显式的 nil 检查。这样如果用户在 Close 之后错误地调用方法，程序就会panic。  

## 使用bzip包
下面的程序，使用上面的程序包实现bzip2压缩命令。用起来和很多UNIX系统上面的 bzip2 命令相似：
```go
// bzipper 读取输入、使用 bzip2 压缩然后输出数据
package main

import (
	"io"
	"log"
	"os"

	"gopl/bzip"
)

func main() {
	w := bzip.NewWriter(os.Stdout)
	if _, err := io.Copy(w, os.Stdin); err != nil {
		log.Fatalf("bzipper: %v\n", err)
	}
	if err := w.Close(); err != nil {
		log.Fatalf("bzipper: close: %v\n", err)
	}
}
```

## 总结
这里演示了如何将C库链接进Go程序中。（*反过来，可以将Go程序编译为静态库然后链接进C程序中，也可以编译为动态库通过C程序来加载和共享*）  
另外还有一些别的问题。  

**没有bzip2库**  
这里的例子是假设已经安装了 bzip2 库。如果是安装位置不对，可以修改 #cgo 来解决。另外，也有人提供了不用依赖本机上的 bzip2 库的实现。  
这里有一个从纯C代码生成的cgo绑定，不依赖bzip2静态库和操作系统的具体环境，具体请访问 https://github.com/chai2010/bzip2  

**并发安全问题**  
上面的实现中，结构体 writer 不是并发安全的。并且并发调用 Close 和 Write 也会导致C代码崩溃。这个问题可以用加锁的方式来避免，使用 sync\.Mutex 可以保证 bzip2\.writer 在多个goroutines中被并发调用是安全的。  

**os/exec 包的实现**  
开篇提到了还有一种实现方式：用 os/exec 包以辅助子进程的方式来调用C程序。  
可以使用 os/exec 包将 /bin/bzip2 命令作为一个子进程执行。实现一个纯Go的 bzip\.NewWriter 来替代原来的实现。这样就是一个纯Go的实现，不需要C言语的基础。不过虽然是纯Go的实现，但还是要依赖本机能够运行 /bin/bzaip2 命令的。  

# 扩展内容
「GCTT 出品」Cgo 和 Python：  
https://studygolang.com/articles/13019?fr=sidebar
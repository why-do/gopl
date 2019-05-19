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
// cgo/bzip/bzip2.c
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

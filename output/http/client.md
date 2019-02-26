# 1.5 获取一个 URL
下面的程序展示从互联网获取信息，获取URL的内容，然后不加解析地输出：
```go
// output/http/fetch/
```
这个程序使用里使用了 net/http 包。http.Get 函数产生一个 HTTP 请求，如果没有出错，返回结果存在响应结构 resp 里面。其中 resp 的 Body 域包含服务器端响应的一个可读取数据流。这里可以用 ioutil.ReadAll 读取整个响应的结果。不过这里用的是 io.Copy(dst, src) 函数，这样不需要把整个响应的数据流都装到缓冲区之中。读取完数据后，要关闭 Body 数据流来避免资源泄露。  

# 1.6 并发获取多个 URL
这个程序和上一个一样，获取URL的内容，并且是并发获取的。这个版本丢弃响应的内容，只报告每一个响应的大小和花费的时间：
```go
// output/http/fetchall/
```
io.Copy 函数读取响应的内容，然后通过写入 ioutil.Discard 输出流进行丢弃。Copy 返回字节数以及出现的任何错误。只所以要写入 ioutil.Discard 来丢弃，这样就会有一个读取的过程，可以获取返回的字节数。  

# 5.2 递归
函数可以**递归**调用，这意味着函数可以直接或者间接地调用自己。递归是一种实用的技术，可以处理许多带有递归特性的数据结构。下面就会使用递归处理HTML文件。  

## 解析 HTML
下面的代码示例使用了 golang.org/x/net/html 包。它提供了解析 HTML 的功能。下面会用到 golang.org/x/net/html API 如下的一些代码。函数 html.Parse 读入一段字节序列，解析它们，然后返回 HTML 文档树的根节点 html.Node。这里可以只关注函数签名的部分，函数内部实现细节就先了解到上面文字说明的部分。HTML 有多种节点，比如文本、注释等。这里只关注 a 标签和里面的 href 的值：
```go
// golang.org/x/net/html
package html

// A NodeType is the type of a Node.
type NodeType uint32

const (
	ErrorNode NodeType = iota
	TextNode
	DocumentNode
	ElementNode
	CommentNode
	DoctypeNode
	scopeMarkerNode
)

type Node struct {
	Parent, FirstChild, LastChild, PrevSibling, NextSibling *Node

	Type      NodeType
	DataAtom  atom.Atom
	Data      string
	Namespace string
	Attr      []Attribute
}

type Attribute struct {
	Namespace, Key, Val string
}

func Parse(r io.Reader) (*Node, error) {
	p := &parser{
		tokenizer: NewTokenizer(r),
		doc: &Node{
			Type: DocumentNode,
		},
		scripting:  true,
		framesetOK: true,
		im:         initialIM,
	}
	err := p.parse()
	if err != nil {
		return nil, err
	}
	return p.doc, nil
}
```
主函数从标准输入中读入 HTML，使用递归的 visit 函数获取 HTML 文本的超链接，并且把所有的超链接输出。  
visit 函数遍历 HTML 树上的所有节点，从 HTML 的 a 标签中得到 href 属性的内容，将获取到的链接内容添加到字符串切片中，最后再返回这个切片。
```go
// output/http/findlinks1
```
要对树中的任意节点 n 进行递归，visit 递归地调用自己去访问节点 n 的所有子节点，并且将访问过的节点保存在 FirstChild 链表中。

分别将两个程序编译后，使用管道将 fetch 程序的输出定向到 findlinks1。编译后执行：
```
PS H:\Go\src\gopl\output\http> go build gopl/output/http/fetch
PS H:\Go\src\gopl\output\http> go build gopl/output/http/findlinks1
PS H:\Go\src\gopl\output\http> ./fetch studygolang.com | ./findlinks1
/readings?rtype=1
/dl
#
http://docs.studygolang.com
http://docscn.studygolang.com
/pkgdoc
http://tour.studygolang.com
/account/register
/account/login
/?tab=all
/?tab=hot
https://e.coding.net/?utm_source=studygolang
```

这是另一个版本，把 fetch 和 findlinks1 合并到一起了，直接把 get 请求返回的读接口 resp.Body 传递给 findlinks1 读取返回的内容进行处理。最后还对 visit 进行了修改，现在使用递归调用 visit （而不是循环）遍历 n.FirstChild 链表。
```go
// exercise5/e1
```

## 遍历 HTML 节点树
下面的程序使用递归遍历所有 HTML 文本中的节点数，并输出树的结构。当递归遇到每个元素时，它都会讲元素标签压入栈，然后输出栈：
```go
// output/http/outline
```
注意一个细节，尽管 outline 会将元素压栈但并不会出栈。当 outline 递归调用自己时，被调用的函数会接收到栈的副本。尽管被调用者可能会对栈（切片类型）进行元素的添加、修改甚至创建新数组的操作，但它并不会修改调用者原来传递的元素，所以当被调用函数返回时，调用者的栈依旧保持原样。  
现在可以找一些网页输出
```
PS H:\Go\src\gopl\output\http> ./fetch baidu.com | ./outline
[html]
[html head]
[html head meta]
[html body]
PS H:\Go\src\gopl\output\http>
```
许多编程语言使用固定长度的函数调用栈，大小在 64KB 到 2MB 之间。递归的深度会受限于固定长度的栈大小，所以当进行深度调用时必须谨防栈溢出。固定长度的栈甚至会造成一定的安全隐患。相比固定长度的栈，Go 语言的实现使用了可变长度的栈，栈的大小会随着使用而增长，可达到 1GB 左右的上限。这使得我们可以安全地使用递归而不用担心溢出的问题。  
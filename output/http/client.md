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

## 合并的版本（5.3 多返回值）
这是另一个版本，把 fetch 和 findLinks 合并到一起了，FindLInks函数自己发送 HTTP 请求。最后还对 visit 进行了修改，现在使用递归调用 visit （而不是循环）遍历 n.FirstChild 链表：
```go
// ch5/findlinks2
```
findLinks 函数有4个返回语句：
+ 第一个返回语句中，错误直接返回
+ 后两个返回语句则使用 fmt.Errorf 格式化处理过的附加上下文信息
+ 如果函数调用成功，最后一个返回语句返回字符串切片，且 error 为空

**关闭 resp.Body**  
这里必须保证 resp.Body 正确关闭使得网络资源正常释放。即使在发生错误的情况下也必须释放资源。  
Go 语言的垃圾回收机制将回收未使用的内存，但不能指望它会释放未使用的操作系统资源，比如打开的文件以及网络连接。必须显示地关闭它们。  

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

## 遍历 HTML 节点树2
这里使用函数变量。可以将每个节点的操作逻辑从遍历树形结构的逻辑中分开。这次不重用 fetch 程序了，全部写在一起了：
```go
// ch5/outline2
```
这里的 forEachNode 函数接受两个函数作为参数，一个在本节点访问子节点前调用，另一个在所有子节点都访问后调用。这样的代码组织给调用者提供了很多的灵活性。  
这里还巧妙的利用了 fmt 的缩进输出。%\*s 中的 \* 号输出带有可变数量空格的字符串。输出的宽度和字符串由后面的两个参数确定，这里只需要输出空格，字符串用的是空字符串。  
这次输出的是要缩进效果的结构：
```
PS H:\Go\src\gopl\ch5\outline2> go run main.go http://baidu.com
<html>
  <head>
    <meta>
    </meta>
  </head>
  <body>
  </body>
</html>
PS H:\Go\src\gopl\ch5\outline2>
```

# 5.6 匿名函数
网络爬虫的遍历。

## 解析链接
在之前遍历节点树的基础上，这次来获取页面中所有的链接。将之前的 visit 函数替换为匿名函数（闭包），现在可以直接在匿名函数里把找到的链接添加到 links 切片中，这样的改变之后，逻辑上更加清晰也更好理解了。因为 Extract 函数只需要前序调用，这里就把 post 部分的参数值传nil。这里做成一个包，后面要继续使用：
```go
// ch5/links
```
**解析URL成为绝对路径**  
这里不直接把href原封不动地添加到切片中，而将它解析成基于当前文档的相对路径 resp.Request.URL。结果的链接是绝对路径的形式，这样就可以直接用 http.Get 继续调用。  

## 图的遍历
网页爬虫的核心是解决图的遍历，使用递归的方法可以实现深度优先遍历。对于网络爬虫，需要广度优先遍历。*另外还可以进行并发遍历，这里不讲这个。*  
下面的示例函数展示了广度优先遍历的精髓。调用者提供一个初始列表 worklist，它包含要访问的项和一个函数变量 f 用来处理每一个项。每一个项有字符串来识别。函数 f 将返回一个新的列表，其中包含需要新添加到 worklist 中的项。breadthFirst 函数将在所有节点项都被访问后返回。它需要维护一个字符串集合来保证每个节点只访问一次。  
在爬虫里，每一项节点都是 URL。这里需要提供一个 crawl 函数传给 breadthFirst 函数最为f的值，用来输出URL，然后解析链接并返回：
```go
// ch5/findlinks3
```

## 遍历输出链接
接下来就是找一个网页来测试，下面是一些输出的链接：
```
PS H:\Go\src\gopl\ch5\findlinks3> go run main.go http://lab.scrapyd.cn/
http://lab.scrapyd.cn/
http://lab.scrapyd.cn/archives/57.html
http://lab.scrapyd.cn/tag/%E8%89%BA%E6%9C%AF/
http://lab.scrapyd.cn/tag/%E5%90%8D%E7%94%BB/
http://lab.scrapyd.cn/archives/55.html
http://lab.scrapyd.cn/archives/29.html
http://lab.scrapyd.cn/tag/%E6%9C%A8%E5%BF%83/
http://lab.scrapyd.cn/archives/28.html
http://lab.scrapyd.cn/tag/%E6%B3%B0%E6%88%88%E5%B0%94/
http://lab.scrapyd.cn/tag/%E7%94%9F%E6%B4%BB/
http://lab.scrapyd.cn/archives/27.html
......
```
整个过程将在所有可达的网页被访问到或者内存耗尽时结束。
# 5.2 递归
+ 练习5.1：改变 findlinks 程序，使用递归调用 visit （而不是循环）遍历 n.FirstChild 链表。
+ 练习5.2：写一个函数，用于统计 HTML 文档树内所有的元素个数，如p、div、span等。
+ 练习5.3：写一个函数，用于输出 HTML 文档树中所有文本节点的内容。但不包括 script 或 style 标签，因为这些内容在 Web 浏览器中是不可见的。
+ 练习5.4：扩展 visit 函数，使之能够获得到其他种类的链接地址，比如图片、脚本或样式表的链接。

# 5.3 多返回值
+ 练习5.5：实现函数 countWordsAndImages （参照练习 4.9 中的单词分隔）。
+ 练习5.6：修改 gopl.io/ch3/surface （参考 3.2 节）中的函数 corner，以使用命名的结果以及裸返回语句。

# 5.5 函数变量
+ 练习5.7：开发 startElement 和 endElement 函数并应用的一个通用的 HTML 输出代码中。输出注释节点，文本节点以及每个元素的属性（`<a href="...">`）。当一个元素没有子节点时，使用简短的形式。比如 `<img/>` 而不是 `<img></img>`。写一个测试程序保证输出可以正确解析（参考第11章）。
+ 练习5.8：修改 forEachNode 使得 pre 和 post 函数返回一个布尔型的结果来确定遍历是否继续下去。使用它写一个函数 ElementByID，该函数使用下面的函数签名并且找到第一个符合 id 属性的 HTML 元素。函数在找到符合条件的元素时应该尽快停止遍历。 `func ElementByID(doc *html.Node, id string) *html.Node`
+ 练习5.9：写一个函数 `expand(s string, f func(string) string) string`，该函数替换参数 s 中每一个子字符串 "$foo" 为 `f("foo")` 的返回值。
# 10.5 空导入
+ 练习10.1：扩展 jpeg 程序，使其可以把任意支持的输入格式转换为任意输出格式，使用 image\.Decode 来检测输入格式，并且添加一个标记来选择输出格式。
+ 练习10.2：定义一个通用的归档文件读取函数，它可以读取 ZIP（archive/zip）文件和 POSIX tar（archive/tar）文件。使用一个类似前面描述的注册机制，使用空白导入以插件方式支持各种文件格式。

 # 10.7 go 工具

 ## 10.7.2 包的下载
 + 练习10.3：使用`http://gopl.io/ch1/helloworld?go-get=1`，找到本书的示例代码是由哪个服务托管的（go get 发出的 HTTP 请求包含 go-get 参数，这样服务器可以区分出普通的浏览器请求。）

 ## 10.7.6 包的查询
 + 练习10.4：构建一个工具，它可以汇报工作空间中所有包的过度依赖中，是否含有参数中指定的包。提示：你将需要执行两次 go list，一次是针对初始包，一次是针对所有包。你也许想使用 encoding/json 包（参数 4.5 节）来解析它的JSON格式输出内容。
# 7.2 接口类型

io包定义了很多有用的接口：
+ io.Writer ： 抽象了所有写入字节的类型，下面会列举
+ io.Reader ： 抽象了所有可以读取字节的类型
+ io.Closer ： 抽象了所有可以关闭的类型，比如文件或者网络连接

io.Writer 是一个广泛使用的接口，它负责所有可以写入字节的抽象，包括但不限于下面列举的这些：
+ 文件
+ 内存缓冲区
+ 网络连接
+ HTTP客户端
+ 打包器（archiver）
+ 散列器（hasher）

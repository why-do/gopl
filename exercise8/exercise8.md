# 8.2 示例：并发时钟服务器
+ 练习8.1：修改 clock2 来接收一个端口号，写一个程序 clockwall，作为多个时钟服务器的客户端，读取每一个服务器的时间，类似于不同地区办公室的时钟，然后显示在一个表中。如果可以访问不同地域的计算机，可以远程运行示例程序；否则可以伪装不同的地区，在不同的端口收本地运行：
```
$ TZ=US/Eastern ./clock2 -port 8010 &
$ TZ=Asia/Tokyo ./clock2 -port 8020 &
$ TZ=Europe/London ./clock2 -port 8030 &
$ clockwall NewYork=localhost:8010 London=localhost:8020 Tokyo=localhost:8030
```
+ 练习8.2：实现一个并发的 FTP 服务器。服务器可以解释从客户端发来的命令，例如 cd 用来改变目录，ls 用来列出目录，get 用来发送一个文件的内容，close 用来关闭连接。可以使用标准的 ftp 命令作为客户端，或者自己写一个。

# 8.9 取消
+ TODO：练习1.11
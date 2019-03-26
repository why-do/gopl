# 8.2 示例：并发时钟服务器
+ 练习8.1：修改 clock2 来接收一个端口号，写一个程序 clockwall，作为多个时钟服务器的客户端，读取每一个服务器的时间，类似于不同地区办公室的时钟，然后显示在一个表中。如果可以访问不同地域的计算机，可以远程运行示例程序；否则可以伪装不同的地区，在不同的端口收本地运行：
```
$ TZ=US/Eastern ./clock2 -port 8010 &
$ TZ=Asia/Tokyo ./clock2 -port 8020 &
$ TZ=Europe/London ./clock2 -port 8030 &
$ clockwall NewYork=localhost:8010 London=localhost:8020 Tokyo=localhost:8030
```
+ 练习8.2：实现一个并发的 FTP 服务器。服务器可以解释从客户端发来的命令，例如 cd 用来改变目录，ls 用来列出目录，get 用来发送一个文件的内容，close 用来关闭连接。可以使用标准的 ftp 命令作为客户端，或者自己写一个。

# 8.4.1 无缓冲通道
+ 练习8.3：在 netcat3 中，conn 接口有一个具体的类型 \*net\.TCPConn，它代表一个 TCP 连接。TCP 链接由两半边组成，可以通过 CloseRead 和 CloseWrite 方法分别关闭。修改主 goroutine，仅仅关闭连接的写半边，这样程序可以继续执行输出来自 reverb1 服务器的回声，即使标准输入已经关闭。（对 reverb2 程序来说更难一些，见练习 8.4。）

# 8.5 并行循环
+ 练习8.4：修改 reverb2 程序来使用 sync.WaitGroup  来计算每一个连接上面的活动的回声 goroutine 的个数。当它变成 0 时，关闭练习 8.3 中描述的写半边的 TCP 链接。验证你修改好的 netcat3 客户端，等待最后几个并发的呼喊回声，即使标准输入已经关闭。
+ 练习8.5：使用一个已有的 CPU 绑定的顺序程序，例如 3.3 节的 mandelbrot 程序（复数分形图），或者 3.2 节的 3D 平面计算（浮点数和SVG图），在主循环中并行执行它们，使用通道来通信。在多 CPU 的机器上它的运行速度有多快？goroutine 的最优数量是多少？

# 8.9 取消
+ TODO：练习1.11
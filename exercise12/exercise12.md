# 12.3 Display：一个递归的值显示器
+ 练习12.1：扩展 Display，让它可以处理 map 中键为结构体或者数组的情形。
+ 练习12.2：通过限制递归的层数，让 Display 能安全处理循环引用的数据结构。（在13.3 节中，我们可以看到另外一个检测循环引用的方法。）

# 12.4 示例：编码 S表达式
+ 练习12.3：实现 encode 函数缺失的功能。把布尔值编码为 t 和 nil，浮点数则用 Go 语言的表示法，像 1+2j 这样的复数则编码为 #C(1.0 2.0)。接口编码为成对的类型名和值，比如 ("[]int"(1 2 3))，但要注意这个方法是会有歧义的，因为 reflect\.Type\.String 方法可能会对不同的类型生成同样的字符串。
+ 练习12.4：修改 encode 函数，输出如上所示的美化后的 S 表达式。
+ 练习12.5：改写 encode 函数，从输出 S 表达式改为输出 JSON。使用标准库的解码器 json.Unmarshal 来测试编码器。
+ 练习12.6：改写 encode 函数，优化输出，如果字段值是其类型的零值则不须编码。
+ 练习12.7：参考 json\.encoder （参见 4.5 节）的风格，创建一个 S 表达式编码器的流式 API。

# 12.6 示例：解码 S 表达式
+ 练习12.8：类似于 json\.UnMarshal 函数，sexpr\.Unmarshal 函数在解码之前就需要完善的字节 slice。仿照 json\.Decoder，定义一个 sexpr\.Decoder 类型，允许从一个 io\.Reader 接口解码一系列的值。使用这个新类型来重新实现 sexpr\.Unmarshal。
+ 练习12.9：仿照 xml\.Decoder（参考 7.14 节），写一个基于标记的 S 表达式解码 API。你需要5个类型的标记：Symbol、String、Int、StartList 和 EndList。
+ 练习12.10：扩展 sexpr\.Unmarshal，以处理练习 12.3 中按你的答案编码的布尔值、浮点数和接口。（提示：为了解码接口，你需要一个 map，其中包含每个支持类型从名字到 reflect\.Type 的映射。）

# 12.7 访问结构体字段标签
+ 练习12.11：写一个与 Unpack 对应的 Pack 函数。给定一个结构体的值，Pack 应当返回一个 URL，这个 URL 的参数与输入的结构体对应。
+ 练习12.12：扩展字段标签语法来支持参数有效性检验。比如，一个字符串应当是一个有效的 email 地址或者有效的信用卡号码，一个整数应当是一个有效的美国邮编（美国邮编为5位整数）。修改 Unpack 函数来支持这些功能。
+ 练习12.13：修改 S 表达式编码器（参数 12.4 节）和解码器（参数 12.6 节），支持 sexpr:"..." 形式的字段标签，标签含义同 encoding\/json 包（参数 4.5 节）。
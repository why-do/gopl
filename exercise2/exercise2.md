# 2.6 包和文件
+ 练习2.1：添加类型、常量和函数到 tempconv 包中，处理以开尔文为单位（K）的温度值，0K=-273.15℃，变化1K和变化1℃是等价的。

## 2.6.1 导入
+ 练习2.2：写一个类似于 cf 的通用的单位转换程序，从命令行参数或者标准输入（如果没有参数）获取数字，然后将每一个数字转换为以摄氏度和华氏度表示的温度，以英寸和米表示的长度单位，以磅和千克表示的重量，等等。

## 2.6.2 包初始化
+ 练习2.3：使用循环重写 PopCount 来代替单个表达式。对比两个版本的效率（11.4节会展示如何系统性地对比不同实现的性能。）
+ 练习2.4：写一个用于统计位的 PopCount，它在其实际参数的64位上执行移位操作，每次判断最右边的位，进而实现统计功能。把它与快查表版本的性能进行对比。
+ 练习2.5：使用 x&(x-1) 可以清除x最右边的非零位，利用该特点写一个 PopCount，然后评价它的性能。
// 简单算术表达式的一个求值器
package eval

import (
	"fmt"
	"math"
)

// 使用一个接口Expr来代表这种语言中的任意一个表达式
// 接口的方法后面添加
// type Expr interface{}

// Var 表示一个变量，比如 x。可导出的
type Var string

// literal 是一个数字常量，比如 3.14
type literal float64

// unary 表示一元操作符表达式，比如 -x
type unary struct {
	op rune // '+', '-' 中的一个
	x  Expr
}

// binary 表示二元操作符表达式，比如 x+y
type binary struct {
	op   rune // '+', '-', '*', '/' 中的一个
	x, y Expr
}

// call 表示函数调用表达式，比如 sin(x)
type call struct {
	fn   string // "pow", "sin", "sqrt" 中的一个
	args []Expr
}

// 对包含变量的表达式进行求值，需要一个上下文（environment）来把变量映射到数值
type Env map[Var]float64

// 需要为每种变量表达式定义 Eval 方法来返回表达式在一个给定上下文下的值。
// 每个表达式都必须提供提供这个方法，把它加到Expr接口中去
type Expr interface {
	Eval(env Env) float64
}

// 为各个类型定义 Eval 方法

// 从 Env 中查询结果，如果变量不存在，返回值是0
func (v Var) Eval(env Env) float64 {
	return env[v]
}

// 直接返回本身的值
func (l literal) Eval(env Env) float64 {
	return float64(l)
}

// unary 和 binary 的 Eval 方法首先对它们的操作数递归求值，然后应用 op 操作
func (u unary) Eval(env Env) float64 {
	switch u.op {
	case '+':
		return +u.x.Eval(env)
	case '-':
		return -u.x.Eval(env)
	}
	panic(fmt.Sprintf("unsupported unary operator: %q", u.op))
}
func (b binary) Eval(env Env) float64 {
	switch b.op {
	case '+':
		return b.x.Eval(env) + b.y.Eval(env)
	case '-':
		return b.x.Eval(env) - b.y.Eval(env)
	case '*':
		return b.x.Eval(env) * b.y.Eval(env)
	case '/':
		return b.x.Eval(env) / b.y.Eval(env)
	}
	panic(fmt.Sprintf("unsupported unary operator: %q", b.op))
}

// call 方法，先对 pow、sin、sqrt 函数的参数进行求值，在调用 math 包中的对应函数
func (c call) Eval(env Env) float64 {
	switch c.fn {
	case "pow":
		return math.Pow(c.args[0].Eval(env), c.args[1].Eval(env))
	case "sin":
		return math.Sin(c.args[0].Eval(env))
	case "sqrt":
		return math.Sqrt(c.args[0].Eval(env))
	}
	panic(fmt.Sprintf("unsupported function call: %s", c.fn))
}

package eval

// Expr: 算术表达式
type Expr interface {
	// 返回表达式在 env 上下文下的值
	Eval(env Env) float64
	// Check 方法报告表达式中的错误，并把表达式中的变量加入 Vars 中
	Check(vars map[Var]bool) error
}

// Var 表示一个变量，比如：x.
type Var string

// Env 变量到数值的映射关系
type Env map[Var]float64

// literal 是一个数字常量，比如：3.1415926
type literal float64

// unary 表示一元操作符表达式，比如：-x
type unary struct {
	op rune // one of '+', '-'
	x  Expr
}

// binary 表示二元操作符表达式，比如：x+y.
type binary struct {
	op   rune // one of '+', '-', '*', '/'
	x, y Expr
}

// call 表示函数调用表达式，比如：sin(x).
type call struct {
	fn   string // one of "pow", "sin", "sqrt"
	args []Expr
}

// memo包提供了一个对类型 Func 并发不安全的函数记忆功能
package memo

// Memo 缓存了调用 Func 的结果
type Memo struct {
	f     Func
	cache map[string]result
}

// Func 是用于记忆的函数类型
type Func func(key string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]result)}
}

// 注意：并发不安全
func (memo *Memo) Get(key string) (interface{}, error) {
	res, ok := memo.cache[key]
	if !ok {
		res.value, res.err = memo.f(key)
		memo.cache[key] = res
	}
	return res.value, res.err
}
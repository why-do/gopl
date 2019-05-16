// memo包提供了一个对类型 Func 并发安全的函数记忆功能
// 但是还缺少重复抑制，导致同样的请求在第一个返回之前再次发起的请求，无法等待第一个请求返回后取缓存里取
// 而是依然也会去再次发起一次请求
package memo

import "sync"

// Func 是用于记忆的函数类型
type Func func(key string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]result)}
}

// Memo 缓存了调用 Func 的结果
type Memo struct {
	f     Func
	mu    sync.Mutex // 保护 cache
	cache map[string]result
}

// Get 是并发安全的
func (memo *Memo) Get(key string) (interface{}, error) {
	memo.mu.Lock()
	res, ok := memo.cache[key]
	memo.mu.Unlock()
	if !ok {
		res.value, res.err = memo.f(key)
		memo.mu.Lock()
		memo.cache[key] = res
		memo.mu.Unlock()
	}
	return res.value, res.err
}

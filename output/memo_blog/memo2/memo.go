// memo包提供了一个对类型 Func 并发安全的函数记忆功能
// 使用了互斥锁，来保证并发安全，而实际则是限制了并发导致了整个过程是串行的
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
	if !ok {
		res.value, res.err = memo.f(key)
		memo.cache[key] = res
	}
	memo.mu.Unlock()
	return res.value, res.err
}

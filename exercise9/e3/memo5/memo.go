// memo包提供了一个对类型 Func 并发安全的函数记忆功能
// 并发、重复抑制、非阻塞的缓存
// 通过监控 goroutine 来实现并发安全
package memo

import "errors"

// Func 是用于记忆的函数类型
type Func func(key string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

type entry struct {
	res   result
	ready chan struct{} // res 准备好之后会被关闭
}

// Func、result、entry 的声明和之前一致

// request 是一条请求消息
type request struct {
	key      string        // 需要 Func 运行的参数
	response chan<- result // 每个客户端接收结果的通道
}

type Memo struct {
	requests chan request
	done     <-chan struct{}
}

func New(f Func, done chan struct{}) *Memo {
	memo := &Memo{ // 创建实例
		requests: make(chan request),
		done:     done,
	}
	go memo.server(f) // 启动服务端 goroutine
	return memo       // 返回实例，供客户端调用
}

func (memo *Memo) Close() { close(memo.requests) }

func (memo *Memo) Get(key string) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{key, response}
	select {
	case res := <-response:
		return res.value, res.err
	case <-memo.done: // 客户取消操作后立刻返回
		return nil, errors.New("客户取消操作")
	}
}

func (memo *Memo) server(f Func) {
	cache := make(map[string]*entry)
	for req := range memo.requests { // 一次处理收到的请求
		e := cache[req.key]
		if e == nil {
			// 对这个 key 的第一次请求
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			go e.call(f, req.key, memo.done) // 调用 f(key)
		}
		// 无论是否第一次请求，最后要回复结果，都有等待 ready 通道返回后，再去读取结果
		go e.deliver(req.response)
	}
}

func (e *entry) call(f Func, key string, done <-chan struct{}) {
	// 执行函数
	v, err := f(key)
	select {
	case <-done:
		// 不要缓存被取消的 Func 调用结果
	default:
		e.res.value, e.res.err = v, err
	}
	// 发送广播通知，数据已经准备好了
	close(e.ready)
}

func (e *entry) deliver(response chan<- result) {
	// 等待数据准备完毕
	<-e.ready
	// 向客户端发送结果
	response <- e.res
}

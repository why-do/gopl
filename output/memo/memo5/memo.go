// memo包提供了一个对类型 Func 并发安全的函数记忆功能
// 并发、重复抑制、非阻塞的缓存
// 通过监控 goroutine 来实现并发安全
package memo

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

type Memo struct{ requests chan request }

func New(f Func) *Memo {
	memo := &Memo{requests: make(chan request)} // 创建实例
	go memo.server(f)                           // 启动服务端 goroutine
	return memo                                 // 返回实例，供客户端调用
}

func (memo *Memo) Close() { close(memo.requests) }

func (memo *Memo) Get(key string) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{key, response}
	res := <-response
	return res.value, res.err
}

func (memo *Memo) server(f Func) {
	cache := make(map[string]*entry)
	for req := range memo.requests { // 一次处理收到的请求
		e := cache[req.key]
		if e == nil {
			// 对这个 key 的第一次请求
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			go e.call(f, req.key) // 调用 f(key)
		}
		// 无论是否第一次请求，最后要回复结果，都有等待 ready 通道返回后，再去读取结果
		go e.deliver(req.response)
	}
}

func (e *entry) call(f Func, key string) {
	// 执行函数
	e.res.value, e.res.err = f(key)
	// 发送广播通知，数据已经准备好了
	close(e.ready)
}

func (e *entry) deliver(response chan<- result) {
	// 等待数据准备完毕
	<-e.ready
	// 向客户端发送结果
	response <- e.res
}

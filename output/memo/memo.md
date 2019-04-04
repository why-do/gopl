# 9.7 示例：并发非阻塞缓存
创建一个**并发非阻塞的缓存**系统，它能解决**函数记忆**（memoizing）的问题，即缓存函数的结果，达到多次调用但只须计算一次结果。这个问题在并发实战中很常见但已有的库不能很好地解决这个问题。这里的解决方案将会是并发安全的，并且要避免简单地对整个缓存使用单个锁而带来的锁争夺问题。  

## 被缓存结果的函数
在做系统之前，先准备一个将要被测试的函数。这里将使用下面的 httpGetBody 函数作为示例来演示函数记忆。调用 HTTP 请求相当昂贵，所以我希望只在第一次请求的时候去发起请求，而之后都可以在缓存中找到结果直接返回：
```go
func httpGetBody(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
```
先保证能缓存这个函数的执行结果，之后再使用更多个函数来测试和验证功能。  
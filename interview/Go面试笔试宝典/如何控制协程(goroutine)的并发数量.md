## 面试问题
如果一个接口下面有个异步的操作需要起goroutine去执行，那么如何限制goroutine的数量呢

实际上本质问题为：如何控制协程(goroutine)的并发数量

答案：
1. 利用 channel 的缓存区
2. 线程池/协程池
3. 调整系统资源的上限, 提高某类资源的上限，例如：ulimit -n 999999，将同时打开的文件句柄数量调整为 999999

## 利用信道 channel 的缓冲区大小来实现

```go
func main() {
	var wg sync.WaitGroup
	ch := make(chan struct{}, 3)
	for i := 0; i < 10; i++ {
		ch <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			log.Println(i)
			time.Sleep(time.Second)
			<-ch
		}(i)
	}
	wg.Wait()
}
```
1. make(chan struct{}, 3) 创建缓冲区大小为 3 的 channel，在没有被接收的情况下，至多发送 3 个消息则被阻塞。
2. 开启协程前，调用 ch <- struct{}{}，若缓存区满，则阻塞。
3. 协程任务结束，调用 <-ch 释放缓冲区。
4. sync.WaitGroup 并不是必须的，例如 http 服务，每个请求天然是并发的，此时使用 channel 控制并发处理的任务数量，就不需要 sync.WaitGroup。

## 线程池/协程池
目前有很多第三方库实现了协程池，可以很方便地用来控制协程的并发数量，比较受欢迎的有：

1. [Jeffail/tunny](https://github.com/Jeffail/tunny)
2. [panjf2000/ants](https://github.com/panjf2000/ants)

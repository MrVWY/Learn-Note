1. slice map channel 是线程安全的吗，怎么保证线程安全
  不是线程安全的

- slice 是一个动态数组，并不是线程安全的。当多个 goroutine 并发读写同一个 slice 时，如果没有额外的同步机制（如锁），会导致数据竞争和不确定的行为。
  并发读取：多个 goroutine 并发读取相同的 slice 通常是安全的，
  同时修改：如果有任何 goroutine 同时修改 slice，就会导致数据不一致或其他问题。
  并发写入：多个 goroutine 同时写入或修改 slice 是不安全的，可能导致数据损坏、越界或程序崩溃。
  解决方案：通过 sync.Mutex 或 sync.RWMutex 等锁机制来保护对 slice 的并发访问。
```
package main

import (
    "fmt"
    "sync"
)

func main() {
    var mu sync.Mutex
    slice := []int{}

    wg := sync.WaitGroup{}

    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            mu.Lock()
            slice = append(slice, i)
            mu.Unlock()
        }(i)
    }

    wg.Wait()
    fmt.Println(slice)
}

```

- map 并发读取和写入：如果有一个 goroutine 正在读取 map，而另一个 goroutine 同时对其进行写入，可能会导致严重的崩溃错误。
  解决方案：可以使用 sync.Mutex 或 sync.RWMutex 来显式加锁，确保对 map 的操作是线程安全的。
  或者，使用 Go 提供的 sync.Map，这是一个专门为并发设计的线程安全的 map 实现。
```
package main

import (
    "fmt"
    "sync"
)

func main() {
    var m sync.Map

    wg := sync.WaitGroup{}

    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            m.Store(i, i)
        }(i)
    }

    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            if val, ok := m.Load(i); ok {
                fmt.Println(val)
            }
        }(i)
    }

    wg.Wait()
}

```

- channel 是 Go 语言设计时内置的线程安全数据结构。Go 通过 channel 进行 goroutine 之间的安全通信，保证并发环境下的数据传递。

 并发读写：多个 goroutine 可以安全地同时向 channel 发送数据或从 channel 接收数据，而不需要显式的锁。
 但要注意的是，虽然 channel 本身是线程安全的，但在并发场景下关闭 channel 需要特别小心。
 关闭一个 channel 是一个写操作，必须保证只有一个 goroutine 负责关闭 channel，否则会引发 panic（"close of closed channel" 错误）。

 解决方案：使用 sync.Once 确保 channel 只被关闭一次，或者通过协议让某个固定的 goroutine 负责关闭。

```
package main

import (
    "fmt"
    "sync"
)

func main() {
    ch := make(chan int)
    var once sync.Once

    wg := sync.WaitGroup{}

    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            ch <- i
        }(i)
    }

    go func() {
        wg.Wait()
        once.Do(func() { close(ch) })  // 保证只关闭一次
    }()

    for val := range ch {
        fmt.Println(val)
    }
}
```

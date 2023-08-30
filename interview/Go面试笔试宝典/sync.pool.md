## sync.pool

三个方法：
1. new() 实例化,对象池中没有对象时，将会调用 New 函数创建
  ```go
    var studentPool = sync.Pool{
    New: func() interface{} {
    return new(Student)
    },
    }
   ```
2. get() 用于从对象池中获取对象，因为返回值是 interface{}，因此需要类型转换。
3. put() 则是在对象使用完毕后，返回对象池。
```go
stu := studentPool.Get().(*Student)
json.Unmarshal(buf, stu)
studentPool.Put(stu)
```


###  Pool 池里的元素个数设置
pool中对象的数量不可控，由系统决定

sync.Pool 的 poolLocal 数量受 p 的数量影响，会开辟 runtime.GOMAXPROCS(0) 个 poolLocal

### 注意
1. sync.Pool 本质用途是增加临时对象的重用率，减少 GC 负担
2. 不能对 Pool.Get 出来的对象做预判，有可能是新的（新分配的），有可能是旧的（之前人用过，然后 Put 进去的）
3. 不能对 Pool 池里的元素个数设置和限制
4. sync.Pool 的大小是可伸缩的，高负载时会动态扩容，存放在池中的对象如果不活跃了会被自动清理。
5. 当用完一个从 Pool 取出的实例时候，一定要记得调用 Put，否则 Pool 无法复用这个实例
6. 任何存储在sync.Pool中的变量随时都有可能被销毁，并且销毁的同时没有任何通知
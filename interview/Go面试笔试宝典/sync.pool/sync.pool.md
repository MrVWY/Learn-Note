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

### 池对象Get/Put开销
为每一个绑定协程的P都分配一个子池。每个子池又分为私有池private和共享列表shared。共享列表是分别存放在各个P之上的共享区域，
而不是各个P共享的一块内存。
协程拿自己P里的子池对象不需要加锁，拿共享列表中的就需要加锁了。

这样就可以减少开销，可以尽量避免并发冲突

### sync.Pool的特性
1. 无大小限制。 
2. 自动清理，每次GC前会清掉Pool里的所有对象。所以不适用于做连接池。 
3. 每个P都会有一个本地的poolLocal，Get和Put优先在当前P的本地poolLocal操作。其次再进行跨P操作。 
4. 所以Pool的最大个数是runtime.GOMAXPROCS(0)。
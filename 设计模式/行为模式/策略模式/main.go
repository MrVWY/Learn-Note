package main

//最少最近使用 （LRU）： 移除最近使用最少的一条条目。
//先进先出 （FIFO）： 移除最早创建的条目。
//最少使用 （LFU）： 移除使用频率最低一条条目。

//现在， 我们的主要缓存类将嵌入至 evictionAlgo接口中。 缓存类会将全部类型的移除算法委派给 evictionAlgo接口，
//而不是自行实现。 鉴于 evictionAlgo是一个接口， 我们可在运行时将算法更改为 LRU、 FIFO 或者 LFU， 而不需要对缓存类做出任何更改。

func main() {
	lfu := &Lfu{}
	cache := initCache(lfu)

	cache.add("a", "1")
	cache.add("b", "2")

	cache.add("c", "3")

	lru := &Lru{}
	cache.setEvictionAlgo(lru)

	cache.add("d", "4")

	fifo := &Fifo{}
	cache.setEvictionAlgo(fifo)

	cache.add("e", "5")

}

package main

// 策略接口

// Eviction Algo移除算法

type EvictionAlgo interface {
	evict(c *Cache)
}

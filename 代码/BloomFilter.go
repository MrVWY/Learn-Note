package main

import (
	"fmt"
	"hash"
	"hash/fnv"
	"github.com/spaolacci/murmur3"
	"github.com/bits-and-blooms/bitset"
)

type BloomFilter struct {
	bitset    *bitset.BitSet
	size      uint
	hashFuncs []hash.Hash64
}

// 创建一个新的布隆过滤器
func NewBloomFilter(size uint, hashCount int) *BloomFilter {
	hashFuncs := make([]hash.Hash64, hashCount)
	for i := 0; i < hashCount; i++ {
		hashFuncs[i] = murmur3.New64()
	}

	return &BloomFilter{
		bitset:    bitset.New(size),
		size:      size,
		hashFuncs: hashFuncs,
	}
}

// 添加元素到布隆过滤器
func (bf *BloomFilter) Add(item string) {
	for i, h := range bf.hashFuncs {
		h.Reset()
		h.Write([]byte(item))
		hashValue := h.Sum64() % uint64(bf.size)
		bf.bitset.Set(uint(hashValue))
	}
}

// 检查元素是否在布隆过滤器中
func (bf *BloomFilter) Check(item string) bool {
	for _, h := range bf.hashFuncs {
		h.Reset()
		h.Write([]byte(item))
		hashValue := h.Sum64() % uint64(bf.size)
		if !bf.bitset.Test(uint(hashValue)) {
			return false
		}
	}
	return true
}

func main() {
	// 创建布隆过滤器，位数组大小为1000，使用3个哈希函数
	bf := NewBloomFilter(1000, 3)

	// 添加元素
	bf.Add("hello")
	bf.Add("world")

	// 检查元素
	fmt.Println(bf.Check("hello"))  // 可能为 true
	fmt.Println(bf.Check("world"))  // 可能为 true
	fmt.Println(bf.Check("golang")) // 可能为 false
}

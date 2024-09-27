请你设计并实现一个满足LRU(最近最少使用)缓存约束的数据结队列→最近使用放头部LRUCache 类:
- LRUCache (int capacity) 以正整数 作为容量 capacity 初始化LRU缓存
- int get(int key)如果关键字 key 存在于缓存中，则返回关键字的值, 否则返回-1 
- void put(int key,int value)如果关键字 key 已经存在，则变更民数据值value ;如果不存在，则向缓存中插入该组key-value 。如果插入操作导致关键字数量超过 capacity ，则应该 逐出 最久未使用的关键字 
函数get和put 必须以 0(1) 的平均时间复杂度运行。
```
package main

import "fmt"

// Node 代表双向链表的节点
type Node struct {
	key   int
	value int
	prev  *Node
	next  *Node
}

// LRUCache 代表 LRU 缓存
type LRUCache struct {
	capacity int
	cache    map[int]*Node
	head     *Node
	tail     *Node
}

// Constructor 初始化 LRUCache
func Constructor(capacity int) LRUCache {
	head := &Node{}
	tail := &Node{}
	head.next = tail
	tail.prev = head

	return LRUCache{
		capacity: capacity,
		cache:    make(map[int]*Node),
		head:     head,
		tail:     tail,
	}
}

// get 从缓存中获取值
func (lru *LRUCache) Get(key int) int {
	if node, found := lru.cache[key]; found {
		lru.moveToFront(node)
		return node.value
	}
	return -1
}

// put 向缓存中插入键值对
func (lru *LRUCache) Put(key int, value int) {
	if node, found := lru.cache[key]; found {
		// 更新值并移动到前面
		node.value = value
		lru.moveToFront(node)
	} else {
		if len(lru.cache) >= lru.capacity {
			lru.removeLeastUsed()
		}
		newNode := &Node{key: key, value: value}
		lru.cache[key] = newNode
		lru.addToFront(newNode)
	}
}

// moveToFront 将节点移动到链表前面
func (lru *LRUCache) moveToFront(node *Node) {
	lru.removeNode(node)
	lru.addToFront(node)
}

// removeNode 移除节点
func (lru *LRUCache) removeNode(node *Node) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

// addToFront 将节点添加到链表前面
func (lru *LRUCache) addToFront(node *Node) {
	node.prev = lru.head
	node.next = lru.head.next
	lru.head.next.prev = node
	lru.head.next = node
}

// removeLeastUsed 移除最久未使用的节点
func (lru *LRUCache) removeLeastUsed() {
	if lru.tail.prev == lru.head {
		return // 链表为空
	}
	node := lru.tail.prev
	lru.removeNode(node)
	delete(lru.cache, node.key)
}

func main() {
	lru := Constructor(2)
	lru.Put(1, 1) // 缓存是 {1=1}
	lru.Put(2, 2) // 缓存是 {1=1, 2=2}
	fmt.Println(lru.Get(1)) // 返回 1
	lru.Put(3, 3) // 该操作会使得关键字 2 被逐出缓存，缓存是 {1=1, 3=3}
	fmt.Println(lru.Get(2)) // 返回 -1 (未找到)
	lru.Put(4, 4) // 该操作会使得关键字 1 被逐出缓存，缓存是 {4=4, 3=3}
	fmt.Println(lru.Get(1)) // 返回 -1 (未找到)
	fmt.Println(lru.Get(3)) // 返回 3
	fmt.Println(lru.Get(4)) // 返回 4
}

```

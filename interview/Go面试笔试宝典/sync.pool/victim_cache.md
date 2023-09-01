## cache结构
1. 直接相连 (Direct mapping cache)

2. 全相联(Associative mapping cache) (全相关cache)

3. 组相联 (set-Associative mapping cache) (集相关cache)

## victim cache 概念
所谓受害者缓存（Victim Cache），是一个与直接匹配或低相联缓存并用的、容量很小的全相联缓存。
当一个数据块被逐出缓存时，并不直接丢弃，而是暂先进入受害者缓存。如果受害者缓存已满，就替换掉其中一项。
当进行缓存标签匹配时，在与索引指向标签匹配的同时，并行查看受害者缓存，如果在受害者缓存发现匹配，
就将其此数据块与缓存中的不匹配数据块做交换，同时返回给处理器。


受害者缓存的意图是弥补因为低相联度造成的频繁替换所损失的时间局部性。

## sync pool 的用法
sync.Pool拥有两个对象存储容器：local pool和victim cache。local pool与victim cache相似，
相当于primary cache。

当获取对象时，优先从local pool中查找，若未找到则再从victim cache中查找，
若也未获取到，则调用New方法创建一个对象返回。当对象放回sync.Pool时候，会放在local pool中。

当GC开始时候，首先将victim cache中所有对象清除，然后将local pool容器中所有对象都会移动到victim cache中，
所以说缓存池中的对象会在每2个GC循环中清除。

## redis和mongodb
### 单个文档大小
&ensp;&ensp;注：1KB=1024B 1MB=1024KB 1GB =1024MB
- mongobd
  + MongoDB存储在BSON(二进制JSON)中，最大BSON文档大小为16MB
- redis
  + 虽然redis的Key的大小上限为512M,但是一般建议key的大小不要超过1KB，这样既可以节约存储空间，又有利于Redis进行检索
  + redis的value的最大值也是512M。对于String类型的value值上限为512M，而集合、链表、哈希等key类型，单个元素的value上限也为512M

## redis可以当数据库吗?
### redis出现
&ensp;&ensp;大多数数据库，由于经常和磁盘打交道，在高并发场景下，响应会非常的慢。为了解决这种速度差异，大多数系统都习惯性的加入一个缓存层，来加速数据的读取。
同时redis也是一个基于内存的键值对数据库。如果数据量很大都存储到内存中会增加成本，而且一般redis都会开启持久化，如果数据量较大，那么持久化的就会变得很多，
增加了redis 的压力，同时会降低redis的性能，因为很大一部分资源都用于持久化数据了。加上数据量达到了上百G，那么就要耗费上百G的内存，成本上也划不过来。

### 能否做数据库用取决于如下几个条件
- 数据量，毕竟内存数据库，还是受限于内存的容量
- 数据的结构，是否能够将关系型数据结构都转换为key/value的形式
- 查询的效率，对范围查询等，是否能转换为高效的hash索引查询
- 如果打算存储一些临时数据，数据规模不大，不需要太复杂的查询，但是对性能的要求比较高，那可以拿redis当数据库使用


## Redis使用Lua脚本为什么能保证原子性

Redis使用同一个Lua解释器来执行所有命令，同时，Redis保证以一种原子性的方式来执行脚本：
当lua脚本在执行的时候，不会有其他脚本和命令同时执行，这种语义类似于 MULTI/EXEC。从别的客户端的视角来看，
一个lua脚本要么不可见，要么已经执行完。

url：https://redis.io/docs/interact/programmability/eval-intro/



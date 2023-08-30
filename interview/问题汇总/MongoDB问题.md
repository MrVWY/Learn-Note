## 过期索引
一种会过期的索引，在索引过期之后，索引对应的数据会被删除 
```
db.collection.createIndex({ expireField: 1 }, { expireAfterSeconds: 3600 })
```

# MySQL问题

## 0、为什么用MySQL？

- mysql性能卓越,服务稳定,很少出现异常宕机。
- mysql开放源代码且无版权制约,自主性及使用成本低。
- mysql历史悠久,社区及用户非常活跃,遇到问题,可以寻求帮助。
- mysql软件体积小,安装使用简单,并且易于维护,安装及维护成本低。
- mysql品牌口碑效应,使得企业无需考虑就直接用之,lamp,lnmp流行架构

## 1、讲一下数据库的表连接操作？什么是笛卡尔积？

左连接、右连接、内连接

## 3、MySQL走索引怎么查到数据？

`innodb`引擎以**页**的形式将数据储存到磁盘，查询时将页读入内存，在叶子节点中查取数据，叶节点内部通过二分法查找，找不到转到该页指向的下一个页继续查询。

## 4、介绍ACID？

**Atomicity：**事务本身被视为不可分割的最小单元，事务的操作要么全部成功要么全部失败回滚。

**Consistency：**数据库在事务的执行前后都保持一致，**所有事务对同一数据的读取结果都相同**。

**Isolation：**一个事务的操作在提交之前，对其他事务是不可见的

**Durability：**一旦事务提交之后对于数据库的更改就是永久不可回退的

## 5、事务的 ACID 特性在MySQL中的实现？
事务：MySQL事务是在引擎层实现的，MySQL原生myISAM存储引擎不支持事务。

原子性：利用undo log 实现的

持久性：利用redo log 实现的

一致性:是利用 原子性、持久性、隔离性来实现的。事务的四大特性中一致性是目的，其他都是保证一致性的手段。

**redo log :**记录了数据操作在物理层面的修改，事务进行中会不断的产生redo log 在事务进行提交时一次flush操作保存到磁盘中。

**undo log:** 记录事务的修改操作，可以实现事务的回滚。

**事务的隔离性由MVCC（多版本并发控制）与锁实现**：因而隔离性也可以叫做并发控制。
innodb存储引擎中实现了三种隔离级别，分别为读未提交、读已提交、可重复读。其中后两者的实现均基于MVCC，其原理为根据read view（当前未提交事务视图）在事务回滚连中往寻找，直到找到了合适的记录。在聚簇索引中存在两个隐藏列为trx_id：当前行最近的事务改变、roll_pointer：当前旧版本在undo log（事务回滚链路）中位置的指针。
RU: 读未提交
RC: 读已提交 每次读语句开始时新建视图
RR: 可重复读（解决不可重复读问题） 每次事务开始前创建视图
串行化: 锁实现

`InnoDB`储存引擎标准实现的锁只有两种：**行级锁、意向锁**。

`InnoDB`实现了如下两种标准的行级锁：

- 共享锁（读锁 S Lock），允许事务读一行数据
- 排它锁（写锁 X Lock），允许事务删除一行数据或者更新一行数据

`InnoDB`支持两种意向锁（即为表级别的锁）：

- 意向共享锁（读锁 IS Lock），事务想要获取一张表的几行数据的共享锁，事务在给一个数据行加共享锁前必须先取得该表的IS锁。
- 意向排他锁（写锁 IX Lock），事务想要获取一张表中几行数据的排它锁，事务在给一个数据行加排他锁前必须先取得该表的IX锁。

**加意向锁表明某个事务正在锁定一行或者将要锁定一行**。首先申请意向锁的动作是`InnoDB`完成的，怎么理解意向锁呢？例如：事务A要对一行记录r进行上X锁，那么`InnoDB`会先申请表的IX锁，再锁定记录r的X锁。在事务A完成之前，事务B想要来个全表操作，此时直接在表级别的IX就告诉事务B需要等待而不需要在表上判断每一行是否有锁。**意向排它锁存在的价值在于节约`InnoDB`对于锁的定位和处理性能。**

`InnoDB`有**3种行锁的算法**：

- Record Lock：单个行记录上的锁
- Gap Lock：间隙锁，锁定一个范围，而非记录本身
- Next-Key Lock：结合Gap Lock和Record Lock，锁定一个范围，并且锁定记录本身。主要解决的问题是RR隔离级别下的**幻读问题**。

## 6、聚簇索引与非聚簇索引？

**页**的概念：一块小的且连续的内存空间

**聚簇索引：**

​	`InnoDB`储存引擎中，**聚簇索引就是按照每张表的主键构造一颗B+树**，同时**叶子节点中存放的就是整张表的行记录数据**，**也将聚簇索引的叶子节点称为数据页**。

　　一般建表会用一个**自增主键**做**聚簇索引**，**没有的话MySQL会默认创建**，但是这个主键如果更改代价较高（页撕裂），故建表时要考虑自增ID不能频繁update这点。

　　我们日常工作中，根据实际情况自行添加的索引都是辅助索引，辅助索引就是一个为了寻找主键索引的**二级索引**，先找到主键索引再通过主键索引找数据。非聚簇索引的叶子节点存储的是数据行的主键信息。

**聚簇索引的优点：**

- 数据访问更快，因为聚簇索引将索引和数据保存在同一个B+树中，因此从聚簇索引中获取数据比非聚簇索引更快
- 聚簇索引对于主键的排序查找和**范围查找**速度非常快

**聚簇索引的缺点：**

- 插入速度严重依赖于插入顺序，按照主键的**顺序插入**是最快的方式，否则将会出现页分裂，严重影响性能。因此，对于`InnoDB`表，我们一般都会定义一个**自增的ID列为主键**
- 更新主键的代价很高，因为将会导致被更新的行移动。因此，对于`InnoDB`表，我们一般定义主键为不可更新。
- 二级索引访问需要两次索引查找，第一次找到主键值，第二次根据主键值找到行数据。这种二次查询的方式叫做**回表查询**。

![](https://github.com/lokles/Web-Development-Interview-With-Java/blob/main/images/聚簇索引.png)



## 7、B+树的特性？

1. 所有的非叶子节点只保存索引，不保存数据。因此树结构更加矮胖，相较于普通N叉树减少磁盘I/O次数。
2. **对于范围查找来说，b+树只需遍历叶子节点链表即可，b树却需要重复地中序遍历。**
3. 索引维护过程可能会导致页分裂（当前数据页满了，需申请新数据页）、页合并 影响空间利用率。

## 8、非聚簇索引的叶子节点储存什么数据？

**非聚簇索引又叫二级索引**，该索引的叶子节点保存的是**数据行的主键值**，想要得到结果还需要使用主键值去聚簇索引（主键索引）中进行二次检索。这一过程称之为回表。

## 9、MySQL多字段查询?如何设置索引？索引的顺序？

最长搜索的字段放最右侧。范围搜索后面的字段的索引会失效。尽量使用覆盖索引。

## 10、MVCC是什么？原理？

全称**多版本并发控制**，与之相对的是**基于锁的并发控制**。

**MVCC最大的优势：读不加锁，读写不冲突。在读多写少的OLTP应用中，读写不冲突是非常重要的，极大的增加了系统的并发性能**

**MVCC实现**

而 MVCC 利用了多版本快照的思想，写操作更新最新的版本快照，而读操作去读旧版本快照(read view)根据隔离级别不同读取的规则也不同，没有互斥关系。

## 11、联合索引的数据结构？

联合索引的数据结构依然是B+树。其非叶子节点储存的是第一个关键字的索引。叶子节点存储的是三个关键字的顺序。且按照字段从左到右排序。

如图，index(年龄, 姓氏,名字)，叶节点上data域存储的是三个关键字的数据。且是按照年龄、姓氏、名字的顺序排列的。

![](https://github.com/lokles/Web-Development-Interview-With-Java/blob/main/images/联合索引.png)

如果跳过年纪按照后面两个字段搜索，会导致全表扫描。

## 12、explain查询到的字段？

1. `select_type` : 查询类型，有简单查询、联合查询、子查询等
2. `key` : 实际使用到的索引，如果为null，表示没有使用到索引。
3. possiable_key：
4. `type`：显示查询使用了何种索引类型，all < index < range < ref
5. `table`：显示这一行的数据是关于哪张表的
6. `rows` : 根据表统计信息及索引选用情况，大致估算出找到所需数据所需要读取的行数。
7. `id`： select查询的序列号,包含一组数字，表示查询中执行select子句的顺序。
8. extra：其他信息，显示如 using index 、`using filesort` 等等。 

## 13、MySQL的自增ID用完了怎么办？

数据库表的自增 ID 达到上限之后，再申请时它的值就不会在改变了，**继续插入数据时会导致报主键冲突错误**。因此在设计数据表时，尽量根据业务需求来选择合适的字段类型。可以考虑使用`bigint` 类型。

## 14、数据库中保存商品价格使用什么数据类型？

在java的开发中，货币在数据库中MySQL常用`Decimal`和`Numric`类型表示，这两种类型被MySQL实现为同样的类型。

DECIMAL和NUMERIC值作为**字符串**存储，而不是作为二进制浮点数，以便保存那些值的小数精度。

不使用float或者double的原因：因为float和double是以二进制存储的，所以有一定的误差。

## 15、 如果数据库出现了死锁，怎么去发现死锁？

通过**获取死锁日志**来获取死锁信息。

`mysql`使用几个特殊的表名来作为监控的开关。比如在数据库中创建一个表名为`innodb_monitor`的表用于开启标准监控。创建一个表名为 `innodb_lock_monitor` 的表**开启锁监控**。MySQL 通过检测是否存在这个表名来决定是否开启监控，至于表的结构和表里的内容无所谓。相反的，如果要关闭监控，则将这两个表删除即可。

## 16、你能用sql语句模拟一下幻读的情况吗？

```sql
事务1
select age from table where id > 2
事务2
Insert into table(id , age) values (5, 10)
commit
事务1 
select age from table where id > 2
commit
事务1两个相同的select语句执行了两次，两次的查询结果不相同，这就是产生了幻读。
```

## 17、redo undo log 的作用？

redo log 常用作MySQL服务器异常宕机后的数据恢复工作，复杂保证事务的持久性

undo log 常用于记录被改动的数据，负责事务的一致性。

## 18、MySQL中除了undo log 以外还有什么操作是为了保证事务的一致性？

各种**隔离级别**保证事务的一致性。

## 19、数据库是怎么去做持久性的，做持久性的时候可能会遇到什么问题？

利用 redo log 做持久性，redo log主要记录了data在物理层面的修改。redo log 在事务进行提交时**一次flush操作保存到磁盘中**。

## 20、如何保证MySQL的主从强一致性？

1. 在**主库**事务提交的时候，同时发起两个操作，操作一是将日志写到本地磁盘，操作二是将日志同步到从库并确保落盘。
2. **主库**此时等待两个操作全部成功返回之后，才返回给应用程序，事务提交成功。

## 21、mysql主从一致要求强一致会导致什么问题？

事务的每次提交都需要等到从机的落盘完成后才可以提交。

## 22、如何保证MySQL主从的高可用性？

HA（High Availability）检测工具应运而生。HA工具一般部署在第三台服务器上，同时连接主从，检测主从是否存活，如果主库宕机则及时将仓库升级为主库，将原来的主库降级为从库。

## 23、MySQL的日志除了redo undo log别的有了解吗？

bin log 是MySQL数据库的**二进制日志**，用于记录用户对数据库操作的SQL语句（(除了数据查询语句）信息。

## 24、bin log 与 redo log 的区别?

1. bin log是MySQL级别的日志文件，无论使用哪种存储引擎都会生成。而redo log 是`innodb`引擎独有的日志，用于记录事务操作的变化，记录的是数据修改之后的值，不管事务是否提交都会记录下来。在实例和介质失败时，redo log文件就能派上用场，如数据库掉电，`InnoDB`存储引擎会使用redo log恢复到掉电前的时刻，以此来保证数据的完整性。
2. 两种日志记录的内容形式不同。MySQL的bin log是逻辑日志，其记录是对应的SQL语句。而`innodb`存储引擎层面的重做日志是物理日志。
3. 两种日志与记录写入磁盘的时间点不同，二进制日志只在事务提交完成后进行一次写入。而`innodb`存储引擎的重做日志在事务进行中不断地被写入，并日志不是随事务提交的顺序进行写入的。
4. bin log可以作为恢复数据使用，主从复制搭建，redo log作为异常宕机或者介质故障后的数据恢复使用。

## 25、四种隔离级别解决的问题？

1. 读未提交:会导致 脏读、不可重复读、幻读。解决了更新丢失问题（两个事务对一条数据修改导致的更新覆盖问题）可以直接使用排它写锁实现。
2. 读已提交:会导致不可重复读、幻读。解决了脏读问题（innoDB基于MVCC实现）
3. 可重复度:解决了不可重复读问题，会导致幻读（innoDB基于MVCC实现）
4. 序列化:全部解决（提供严格的事务隔离，事务没有并发性可言）

## 26、读已提交隔离级别为什么会有不可重复读的问题出现？

```java 
//开启事务并设置隔离级别为读已提交,表count两个字段 name, money
A事务 select * from count  结果name = Tom money = 1000 
B事务 update money = 2000 from count where name = Tom B事务提交
A事务 select * from count  结果name = Tom money = 2000 显然A事务对一个数据行两次读操作结果不一致,这就导致了不可重复读问题
```

## 27、介绍几种索引吧？

1. 单一索引：
2. 复合索引：根据创建联合索引的顺序，以**最左前缀匹配原则**进行where检索。
3. 覆盖索引：查询的字段与建立索引的字段一一对应就叫做覆盖索引。

## 28、B+树实现索引的原理说一下？

巴拉巴拉

## 29、非叶子节点它的一个数据结构描述一下？

非叶子结点中仅含有其子节点的索引，不包含实际数据。

## 30、举一个读未提交导致的脏读、不可重复读、幻读的例子？



## 31、MySQL主、从机之间传输数据的网路I/O模型？



## 32、innodb与myisam的适用场景?

1. 大量读不需要事务控制的情况下使用`myisam`，写操作多的情况下使用`innodb`存储引擎。
2. 需要用到行锁的场景下要使用`innodb`。

## 33、什么是索引下推？（有赞）

索引下推在**非主键索引**上的优化，可以**有效减少回表的次数**，大大提升了查询的效率。
举例：
表k中建立有联合索引（name, age），执行如下语句 select * from k where name like '张%' and age = 10 and sex = 1;
1. 在MySQL5.6之前，在利用完name索引后，只能从根据name找到的第一个主键id逐一回表查询数据行，对比字段值。
2. 在MySQL5.6之后，在回表查询之前，会先对索引中包含的字段进行判断，过滤不满足条件的记录，从而减少回表次数。

## 34、MySQL中事务控制语法？

```mysql

```

## 35、不可重复读与幻读的区别？

不可重复读的重点是修改：在同一事务中，同样的条件，第一次读的数据和第二次读的「数据不一样」。（因为中间有其他事务提交了修改）

幻读的重点在于新增或者删除：在同一事务中，同样的条件，第一次和第二次读出来的「记录数不一样」。（因为中间有其他事务提交了插入/删除）

## 36、当前读与快照读？

在一个支持MVCC的系统中，读操作被分为当前读与快照读

快照读：简单的select操作，不加锁。

```MySQL
select * from table where ?;
```

当前读：插入/更新/删除操作，需要加锁

```mysql
select * from table where ? lock in share mode;
select * from table where ? for update;
insert into table values (…);
update table set ? where ?;
delete from table where ?;
```

## 37、DDL与DML？

- DML（data manipulation language）数据操纵语言：

　　　　就是我们最经常用到的 SELECT、UPDATE、INSERT、DELETE。 主要用来对数据库的数据进行一些操作。

- DDL（data definition language）数据库定义语言：

　　　　其实就是我们在创建表的时候用到的一些sql，比如说：CREATE、ALTER、DROP等。**DDL主要是用在定义或改变表的结构**，数据类型，表之间的链接和约束等初始化工作上

## 38、JDBC说一下？

1. 通过驱动建立一个连接，这个连接代表着一个真实的数据库连接。
2. 由conn建立一个`Statement`或`PreparedStatement`对象。
3. `stmt.executeUpdate(sql)`执行语句，返回即查询解决。

## 39、一条sql语句(非update)是如何执行的？
1. 连接器，管理连接、权限验证
2. 分析器，进行词法分析以及语法分析
3. 优化器，进行语句优化，索引选择
4. 执行器，操作底层的数据存储引擎，返回结果
5. 存储引擎，存储数据，对外提供读写接口。
## 40、一条update语句是如何执行的？
在39的基础上，进行update操作还涉及到redo log、binlog。
redo log: （innodb引擎引入的日志，本质上类似于记账账本，不直接在mysql中进行存储，而视在空闲时利用redo log进行数据到磁盘的写入，innodb引擎利用redo log保证数据的不丢失）。其是物理日志，记录某个数据页上做的修改，大小固定，循环写入，一旦空间用完会清除。
binlog: MySQL server提供的功能，逻辑日志，击记录的是sql语句的原始逻辑，如“给 id = 2 的一行数据的c字段增加1”。追加写入，binlog文件到一定大小后会切换到下一个
过程：
1. 执行器根据id找到数据，然后判断内存中是否存在该数据，不存在的话从磁盘中读取并读入内存。
2. 执行器拿到引擎给的数据进行更新，引擎将这条修改后的数据更新入内存，同时更新操作记录到redo log中，此时redo log处于prepare状态，告知执行器处理完成随时可以提交事务
3. 执行器生成该操作的binlog并写入磁盘
4. 执行器调用引擎的提交事务接口，引擎将redo log状态由prepare修改为commit，更新完成。
redo log的两阶段提交是为了保证redo log 与binlog之间的逻辑一致。这样来即便发生MySQL的crash也会保证两个日志的数据一致性。
## 41、覆盖索引
假设table k表中定义了两个索引分别为主键索引id以及非主键索引name。那么如下sql语句 select * from k where name between 3 and 5则需要进行回表查询，而select id from k where name between 3 and 5由于name索引的叶子节点存储的即为其主键id值，这一过程是不需要回表查询的。索引name覆盖了我们的查询需求，称之为覆盖索引。
## 42、如何避免多事务的锁冲突导致的死锁问题？
1. 进入等待，直到超时，超时时间通过innodb_lock_wait_timeout设置。
2. 发起死锁检测，主动回滚死锁链条中的某个事务，让其他事务得以执行，参数innodb_deadlock_detect设置为on。
在2的基础上，如果并发度很高的情况下进行死锁检测也是一个很费时的操作，可能的话在MySQL的服务端进行限流，在请求进入innodb前进行排队，限制同一时刻修改db中某一行的线程数。不仅降低了发生死锁的概率，即使发生了死锁进行死锁检测的效率也会提升很大。

## 43、普通索引与唯一索引该怎么选
在查询性能上两者差距微乎其微，在更新性能上由于普通索引可以利用change buffer的优化机制性能更优。
## 44、count(*)为什么这么慢？

## 45、order by如何工作？
两种情况，内存大小允许的情况下仅使用快排在内存中排序，否则的话需要用到外部硬盘空间，在这块空间中MySQL将数据分为12块进行归并排序。

## 如何在数据库中复制表（只复制结构）
select * into a from b  where 1<>1 （<> 含义与 != 相似,皆为不等于的意思）

## delete、drop、truncate有什么区别？

detele: 可用于删除表的部分或所有数据

    delete from table_name [where...] [order by...] [limit...]
    
    删除学生表中数学成绩排名最高的前 3 位学生，可以使用以下 SQL：
    delete from student order by math desc limit 3;
    
    在 InnoDB 引擎中，delete 操作并不是真的把数据删除掉了，
    而是给数据打上删除标记，标记为删除状态

truncate: 执行效果和 delete 类似，也是用来删除表中的所有行数据的

    truncate [table] table_name

drop: 和前两个命令只删除表的行数据不同，drop 会把整张表的行数据和表结构一起删除掉

    DROP [TEMPORARY(临时表)] TABLE [IF EXISTS] tbl_name [,tbl_name]

runcate 在使用上和 delete 最大的区别是，delete 可以使用条件表达式删除部分数据，而 truncate 不能加条件表达式，所以它只能删除所有的行数据
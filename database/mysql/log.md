## binlog、redolog、undolog
1. binlog二进制日志是mysql-server层的，主要是做主从复制，时间点恢复使用。
2. redolog重做日志是InnoDB存储引擎层的，用来保证事务安全。
3. undolog回滚日志保存了事务发生之前的数据的一个版本，undolog有两个作用：提供回滚和多个行版本控制(MVCC)
4. MVCC:MVCC 的英文全称是 Multiversion Concurrency Control ，中文意思是多版本并发控制技术。原理是，通过数据行的多个版本管理来实现数据库的并发控制，简单来说就是保存数据的历史版本。可以通过比较版本号决定数据是否显示出来。读取数据的时候不需要加锁可以保证事务的隔离效果。
    - 读写之间阻塞的问题，通过 MVCC 可以让读写互相不阻塞，读不相互阻塞，写不阻塞读，这样可以提升数据并发处理能力。
    - 降低了死锁的概率，这个是因为 MVCC 采用了乐观锁的方式，读取数据时，不需要加锁，写操作，只需要锁定必要的行。
    - 解决了一致性读的问题，当我们朝向某个数据库在时间点的快照是，只能看到这个时间点之前事务提交更新的结果，不能看到时间点之后事务提交的更新结果。

## binlog
binlog 是归档日志，属于 `Server 层`的日志，是一个二进制格式的文件，用于记录用户对数据库更新的SQL语句信息  
主要作用：主从复制、数据恢复

事务执行过程中，先把日志写到 binlog cache，事务提交的时候，再把 binlog cache 写到 binlog 文件(文件系统缓存page cache)中。一个事务的 binlog 是不能被拆开的，因此不论这个事务多大，也要确保一次性写入。

### 写入策略
1. sync_binlog=0 的时候，表示每次提交事务都只 write(把日志写入到文件系统的 page cache)，不 fsync(数据持久化到磁盘的操作)；实际的业务场景中，考虑到丢失日志量的可控性，一般不建议将这个参数设成 0
2. sync_binlog=1 的时候，表示每次提交事务都会执行 fsync(数据持久化到磁盘的操作)；
3. sync_binlog=N(N>1) 的时候，表示每次提交事务都 write(把日志写入到文件系统的 page cache)，但累积 N 个事务后才 fsync(数据持久化到磁盘的操作)。比较常见的是将其设置为 100~1000 中的某个数值。但是，将 sync_binlog 设置为 N，对应的风险是：如果主机发生异常重启，会丢失最近 N 个事务的 binlog 日志，主机都crash了，文件系统缓存里的日志当然也丢了，没法落盘

## undolog
undolog 是 InnoDB 存储引擎的日志，用于保证数据的原子性，保存了事务发生之前的数据的一个版本，也就是说记录的是数据是修改之前的数据，可以用于回滚，同时可以提供多版本并发控制下的读（MVCC）  
主要作用：事务回滚、实现多版本控制(MVCC)  
可以追加写入的，并不会覆盖以前的日志  

## relaylog
relaylog 是中继日志，在主从同步的时候使用到，它是一个中介临时的日志文件，用于存储从master节点同步过来的binlog日志内容  
master 主节点的 binlog 传到 slave 从节点后，被写入 relay log 里，从节点的 slave sql 线程从 relaylog 里读取日志然后应用到 slave 从节点本地。从服务器 I/O 线程将主服务器的二进制日志读取过来记录到从服务器本地文件，然后 SQL 线程会读取 relay-log 日志的内容并应用到从服务器，从而使从服务器和主服务器的数据保持一致

## redolog
用于服务器发生故障后重启 MySQL后，恢复事务已提交但未写入数据表的数据。

MySQL在做数据更新操作时，如果每次都需要写进磁盘的话，那么需要到磁盘中找到对应的那条记录，然后更新，这样下来，IO成本和查找成本都很高。为了解决这个问题，MySQL用到了WAL(Write-Ahead Logging)技术，它的关键点在于先写日志，再写磁盘，这就用到了redo log。  

InnoDB的redo log是固定大小的，可以配置为一组4个文件，每个文件1GB大小，那么总共就可以记录4GB的操作，当写满的时候，会淘汰掉当前最老的记录以得到空闲空间。

`redo log是InnoDB引擎特有的日志`，可以保证即使数据库发生异常重启，之前提交过的记录也不会丢失，这个能力称为crash-safe。

有特有的checkpoin, checkpoint前进方向区域一定是数据页还没落盘的提交，这样我们就能确定哪些数据即使事务commit也就是日志刷盘已经成功的状态下其实数据还没有刷盘

### redo log 的写入策略
InnoDB 提供了 innodb_flush_log_at_trx_commit 参数
1. 设置为 0 的时候，表示每次事务提交时都只是把 redo log 留在 redo log buffer 中
2. 设置为 1 的时候，表示每次事务提交时都将 redo log 直接持久化到磁盘，非常安全，但慢
3. 设置为 2 的时候，表示每次事务提交时都只是把 redo log 写到 page cache，写入文件系统的page cache，主机掉电后会丢数据，但是MySQL异常重启不会丢数据，风险较低，写入比较快

## redolog 和 binlog 的区别
1. redolog 是 Innodb 独有的日志，而 binlog 是 server 层的，所有的存储引擎都有使用到
2. redolog记录了具体的数值，对某个页做了什么修改，binlog 记录的操作内容
3. binlog大小达到上限或者 flush log 会生成一个新的文件，而 redolog 有固定大小只能循环利用
4. binlog 日志没有 crash-safe 的能力，只能用于归档。而 redo log 有 crash-safe 能力

### crash-safe

redolog的存在使得数据库具有crash-safe能力，即如果Mysql 进程异常重启了，系统会自动去检查redolog，将未写入到Mysql的数据从redo log恢复到Mysql中去。

## 两阶段提交机制
MySQL通过两阶段提交的机制，保证了redo log和bin log的逻辑一致性，进而保证了数据的不丢失以及主从库的数据一致。  

如果不用两阶段提交，要么就是先写完 redo log 再写 binlog，或者采用反过来的顺序。
- 先写 redo log 直接提交，然后写 binlog，假设写完 redo log 后，机器挂了，binlog 日志没有被写入，那么机器重启后，这台机器会通过 redo log 恢复数据，但是这个时候 binlog 并没有记录该数据，后续进行机器备份的时候，就会丢失这一条数据，同时主从同步也会丢失这一条数据。binlog丢，无法备份
- 先写 binlog，然后写 redo log，假设写完了 binlog，机器异常重启了，由于 redo log 还没写，崩溃恢复以后这个事务无效。本机是无法恢复这一条记录的，但是 binlog 又有记录，那么和上面同样的道理，就会产生数据不一致的情况。binlog多一条事务执行记录，备份归档或主从同步时产生脏事务

### 原理
redo log prepare -> 写 bin log -> redo log commit，这个流程就叫做两阶段提交。

![](两阶段提交机制.png)

好处:  
情景一，redo log处于prepare状态时，如果写bin log失败了，那么更新失败，此时redo log没有commit，bin log也没有记录，两者的状态是一致的，没有问题。  

情景二，redo log处于prepare状态时，写bin log成功，但是宕机导致commit失败了。此时bin log产生了记录，redo log没有写入成功，数据暂时不一致。   

但是当MySQL重启时，会检查redo log中处于prepare状态的记录。在redo log中，记录了一个叫做XID的字段，这个字段在bin log中也有记录，MySQL会通过这个XID，如果在bin log中找到了，那么就commit这个redo log，如果没有找到，说明bin log其实没有写成功，就放弃提交。  

通过这样的机制，保证了redo log和bin log的一致性。
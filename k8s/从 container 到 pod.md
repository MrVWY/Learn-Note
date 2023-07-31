
### linux 的 cgroup和namespace

#### namesapce
  用户可以创建指定类型的namespace并将程序放入该namespace中运行，这表示从当前的系统运行环境中隔离一个进程的运行环境，在此namespace中运行的进程将认为自己享有该namespace中的独立
资源
  linux 内核有提供了很多不同类型的 namespace：
  1. Mount(mnt) namespace 可以在不影响宿主机文件系统的基础上挂载或者取消挂载文件系统了；
  2. PID(Process ID) 在一个pid namespace中，第一个进程的 pid 为1，所有在此 namespace 中的其他进程都是此进程的子进程，操作系统级别的其他进程不需要在此进程中；
  3. network(netns) namespace 会虚拟出新的内核协议栈，每个 network namespace 包括自己的网络设备、路由表等；
  4. IPC(Interprocess Communication) namespace 用来隔离处于不同 IPC namespace 的进程间通信资源，如共享内存等；
  5. UTS: UTS namespace 用于隔离 hostname 与 domainname
  6. cgroup namespace：该namespace可单独管理自己的cgroup
  7. time namespace：该namespace有自己的启动时间点信息和单调时间，比如可设置某个namespace的开机时间点为1年前启动，再比如不同的namespace创建后可能流逝的时间不一样
  8. user namespace：该namespace有自己的用户权限管理机制(比如独立的UID/GID)，使得namespace更安全

#### cgroup
  一个进程树里面的进程会占有宿主机的资源(CPU/Memory/NetworkIO/DiskIO 等)，这样有可能导致其他进程得不到足够的资源。幸亏我们可以使用 linux 内核的 cgroup 特性来限制进程能占用的
资源数量，这些资源包括 CPU、内存、网络带宽、磁盘IO等。
    ```
    cgexec、cgcreate
    ```
### docker容器 和 pod
   docker 通过为每一个容器创建 namespace 和 cgroup 的组合来实现隔离, 2个容器可以通过挂载同一个namesapce和cgroup来互相连通，但如果不同程序都放在一个容器container里面又会变
得十分复杂。如果不同程序拆开成不同的container放进一个pod里面，又会面临谁去创建最初的namespace，使得这些container处在同一个namesapce下呢？答案便是Pause容器。
   Pause容器被创建后会初始化Network Namespace，之后其他容器就可以加入到Pause容器的namesapce中共享Pause容器的网络
   共享Pause容器网络:
   1. 容器之间能够直接用localhost通信
   2. Pod只有一个IP地址，也就是该Pod的Network Namespace对应的IP地址（由Pause容器初始化并创建）

### k8s 跨namespace通信
   可以参考k8s-service.md这笔记

### Reference
1、https://morven.life/posts/from-container-to-pod/
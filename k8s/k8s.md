### 集群架构组件

<img src="C:\Users\zhou jielun\Desktop\k8s概念\k8s-架构.png"  />

1. master（主控节点）
   - apiserver：集群统一入口，通过restful方式交给etcd存储
   - scheduler：负责节点node调度
   - controller-manager：处理集群中常规后台任务，一个资源对应一个控制器
   - etcd：存储，用于保存集群相关的数据
2. node（工作节点）
   - kubelet：管理本机容器
   - kube-proxy：提供网络代理，负载均衡等操作

### 核心概念

#### Pod

- K8s的最小工作单元，包含一个或多个容器。

#### Controller-manager

​	Controller Manager是Kubernetes集群内部的管理控制中心， 负责Kubernetes集群内的Node、 Pod、服务端点、 服务、 资源配额、 命名空间 、服务账号等资源的管理 、 自动化部署、健康监测， 并对异常资源执行自动化修复， 确保集群各资源始终处于预期的工作状态 。 比如， 当某个Node意外若机时，Controller Manager会根据资源调度策略选择集群内其他节点自动部署原右机节点上的Pod副本 。

​	Controller Manager是 一 个控制器集合，包含多个控制器，Controller Manager是这些控制器的核心管理者。

##### ReplicationController

ReplicationController会持续监控正在运行的pod列表，确保pod的数量始终与其标签选择器匹配，ReplicationController由三部分组成：

- label selector(标签选择器),用于确定ReplicationController作用域中有哪些pod
- replica count(副本个数),指定应运行的pod 数量
- pod template(pod模板),用于创建新的pod 副本

##### ReplicaSet

ReplicaSet的行为与ReplicationController完全相同

##### Deployment

Deployment为Pod和Replica Set(下一代 Replication Controller)提供声明式更新。你只需要在Deployment中描述你想要的目标状态是什么,Deployment controller就会帮你将Pod和Replica Set的实际状态改变到你的目标状态。你可以定义一个全新的Deployment,也可以创建一个新的替换旧的Deployment。

Deployment的典型应用场景 包括：

- 定义Deployment来创建Pod和ReplicaSet
- 滚动升级和回滚应用
- 扩容和缩容
- 暂停和继续Deployment

##### Job

​	Job 会创建一个或者多个 Pod，并将继续重试 Pod 的执行，直到指定数量的 Pod 成功终止。 随着 Pod 成功结束，Job 跟踪记录成功完成的 Pod 个数。 当数量达到指定的成功个数阈值时，任务（即 Job）结束。 删除 Job 的操作会清除所创建的全部 Pod。 挂起 Job 的操作会删除 Job 的所有活跃 Pod，直到 Job 被再次恢复执行。

##### CronJob

CronJob 用于执行周期性的动作，例如备份、报告生成等。 这些任务中的每一个都应该配置为周期性重复的（例如：每天/每周/每月一次）； 你可以定义任务开始执行的时间间隔

##### StatefulSet

###### 有状态应用

- 稳定的、唯一的网络标识符。
- 稳定的、持久的存储。
- 有序的、优雅的部署和扩缩。
- 有序的、自动的滚动更新。

比如MySQL、MongoDB集群，pod之间可能包含主从、主备的相互依赖关系，甚至对启动顺序也有要求等

###### 无状态应用

​	多个应用实例对于同一个用户请求的响应结果是完全一致的。这种多服务实例之间是没有依赖关系，比如web应用,在k8s控制器 中动态启停无状态服务的pod并不会对其它的pod产生影响

Refence

- https://kubernetes.io/zh-cn/docs/concepts/workloads/controllers/job/

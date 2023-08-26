### 什么是StorageClass
&ensp;&ensp;简单来说StorageClass对象定义了下面两部分内容:
```yaml
 1.PV的属性。比如，存储类型，Volume的大小等
 2.创建这种PV需要用到的存储插件，即存储制备器。
```
有了这两个信息之后，Kubernetes就能够根据用户提交的PVC, 找到一个对应的StorageClass,
之后Kubernetes就会调用该StorageClass声明的存储插件，进而创建出需要的PV

### 为什么需要StorageClass
&ensp;&ensp;在一个大规模的Kubernetes集群里, 可能有成千上万个PVC, 同时还有新的PVC, 而且不同的应用程序对于存储性能的要求可能也不尽相同，比如读写速度、并发性能等
因此必须快速创建PV, 否则新的Pod就会因为PVC绑定不到PV而导致创建失败。而且通过 PVC 请求到一定的存储空间也很有可能不足以满足应用对于存储设备的各种需求。
所以需要一套可以自动创建PV的机制. 而StorageClass就是这个机制的核心

### 主要字段
StorageClass 中包含 provisioner、parameters 和 reclaimPolicy 字段
- provisioner:用来决定使用哪个卷插件分配 PV. 该字段必须指定
- reclaimPolicy: 回收策略, 可以是 Delete 或者 Retain. 如果 StorageClass 对象被创建时没有指定 reclaimPolicy, 它将默认为 Delete.
```yaml
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: standard
provisioner: kubernetes.io/aws-ebs
parameters:
  type: gp2
reclaimPolicy: Retain
mountOptions:
  - debug
```
### 相关配置解析
配置：https://jimmysong.io/kubernetes-handbook/concepts/storageclass.html
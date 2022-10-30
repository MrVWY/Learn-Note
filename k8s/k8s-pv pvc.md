### 基本概念

​	**持久卷（PersistentVolume，PV）** 是集群中的一块存储，可以由管理员事先制备（**静态制备**）， 或者使用[存储类（Storage Class）](https://kubernetes.io/zh-cn/docs/concepts/storage/storage-classes/)来**动态制备**。 持久卷是集群资源，就像节点也是集群资源一样。PV 持久卷和普通的 Volume 一样， 也是使用卷插件来实现的，只是它们拥有独立于任何使用 PV 的 Pod 的生命周期。 此 API 对象中记述了存储的实现细节，无论其背后是 NFS、iSCSI 还是特定于云平台的存储系统。

​	**持久卷申领（PersistentVolumeClaim，PVC）** 表达的是用户对存储的请求。概念上与 Pod 类似。 Pod 会耗用节点资源，而 PVC 申领会耗用 PV 资源。Pod 可以请求特定数量的资源（CPU 和内存）；同样 PVC 申领也可以请求特定的大小和访问模式 （例如，可以要求 PV 卷能够以 ReadWriteOnce、ReadOnlyMany 或 ReadWriteMany 模式之一来挂载，参见[访问模式](https://kubernetes.io/zh-cn/docs/concepts/storage/persistent-volumes/#access-modes)）。

​	PV 属于集群中的资源。PVC 是对这些资源的请求，也作为对资源的请求的检查。

#### 创建PV的2种方式

- 静态制备（statically Provision）：一种是集群管理员通过手动方式静态创建应用所需要的 PV。
- 动态制备（dynamically Provision）：

#### 访问模式

- ReadWriteOnce -- 卷可以被一个节点以读写方式挂载；
- ReadOnlyMany -- 卷可以被多个节点以只读方式挂载；
- ReadWriteMany -- 卷可以被多个节点以读写方式挂载。

在命令行接口（CLI）中，访问模式也使用以下缩写形式：

- RWO -- ReadWriteOnce
- ROX -- ReadOnlyMany
- RWX -- ReadWriteMany

#### 绑定

​	用户创建一个**带有特定存储容量**和**特定访问模式**需求的 PersistentVolumeClaim 对象； 在动态制备场景下，这个 PVC 对象可能已经创建完毕。 主控节点中的控制回路监测新的 PVC 对象，寻找与之匹配的 PV 卷（如果可能的话）， 并将二者绑定到一起。 如果为了新的 PVC 申领动态制备了 PV 卷，则控制回路总是将该 PV 卷绑定到这一 PVC 申领。 否则，用户总是能够获得他们所请求的资源，只是所获得的 PV 卷可能会超出所请求的配置。 一旦绑定关系建立，则 PersistentVolumeClaim 绑定就是排他性的， 无论该 PVC 申领是如何与 PV 卷建立的绑定关系。 PVC 申领与 PV 卷之间的绑定是一种一对一的映射，实现上使用 ClaimRef 来记述 PV 卷与 PVC 申领间的双向绑定关系。

​	如果找不到匹配的 PV 卷，PVC 申领会无限期地处于未绑定状态。 当与之匹配的 PV 卷可用时，PVC 申领会被绑定。 例如，即使某集群上制备了很多 50 Gi 大小的 PV 卷，也无法与请求 100 Gi 大小的存储的 PVC 匹配。当新的 100 Gi PV 卷被加入到集群时， 该 PVC 才有可能被绑定。

![](C:\Users\zhou jielun\Desktop\111\k8s概念\1.png)

### PV回收策略

​	当用户不再使用其存储卷时，他们可以从 API 中将 PVC 对象删除， 从而允许该资源被回收再利用。PersistentVolume 对象的回收策略告诉集群， 当其被从申领中释放时如何处理该数据卷。 目前，数据卷可以被 Retained（保留）、Recycled（回收，已废弃）或 Deleted（删除）。

​	在yaml文件配置参数`persistentVolumeReclaimPolicy`

- **保留（Retain）** 保留策略允许手动回收资源，当删除PVC的时候，PV仍然存在，变为Realease状态，需要用户手动通过以下步骤回收卷（只有hostPath和nfs支持Retain回收策略）：
  - 1.删除PV。
  - 2.手动清理存储的数据资源。
- **回收（Resycle）** 该策略已废弃，推荐使用dynamic provisioning，回收策略会在 volume上执行基本擦除（rm -rf thevolume/*），可被再次声明使用。
- **删除（Delete）**
  - 当发生删除操作的时候，会从 Kubernetes 集群中删除 PV 对象，并执行外部存储资源的删除操作（根据不同的provisioner 定义的删除逻辑不同，有的是重命名而不是删除）。
  - 动态配置的卷继承其 StorageClass 的回收策略，默认为Delete，即当用户删除 PVC 的时候，会自动执行 PV 的删除策略。

### 保护使用中的存储对象

​	保护使用中的存储对象（Storage Object in Use Protection） 这一功能特性的目的是确保仍被 Pod 使用的 PersistentVolumeClaim（PVC） 对象及其所绑定的 PersistentVolume（PV）对象在系统中不会被删除，因为这样做可能会引起数据丢失。

**说明：** 当使用某 PVC 的 Pod 对象仍然存在时，认为该 PVC 仍被此 Pod 使用。

​	如果用户删除被某 Pod 使用的 PVC 对象，该 PVC 申领不会被立即移除。 PVC 对象的移除会被推迟，直至其不再被任何 Pod 使用。 此外，如果管理员删除已绑定到某 PVC 申领的 PV 卷，该 PV 卷也不会被立即移除。 PV 对象的移除也要推迟到该 PV 不再绑定到 PVC。

​	如果看到PVC和PV的状态为`Terminating` 且其 `Finalizers` 列表中包含 `kubernetes.io/pvc-protection`或者kubernetes.io/pv-protection 时PVC或者PV对象是处于被保护状态的。

```
kubectl describe pvc 

//output
Name:          hostpath
Namespace:     default
StorageClass:  example-hostpath
Status:        Terminating
Volume:
Labels:        <none>
Annotations:   volume.beta.kubernetes.io/storage-class=example-hostpath
               volume.beta.kubernetes.io/storage-provisioner=example.com/hostpath
Finalizers:    [kubernetes.io/pvc-protection]
...
```

### yaml文件配置参数

PV Yaml ：https://kubernetes.io/zh-cn/docs/reference/kubernetes-api/config-and-storage-resources/persistent-volume-v1/

PVC Yaml : https://kubernetes.io/zh-cn/docs/reference/kubernetes-api/config-and-storage-resources/persistent-volume-claim-v1/

### Reference

- https://kubernetes.io/zh-cn/docs/concepts/storage/persistent-volumes/#introduction

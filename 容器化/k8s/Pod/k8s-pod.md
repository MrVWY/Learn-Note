

​	Pod 是 kubernetes 中你可以创建和部署的最小也是最简的单位。Pod 代表着集群中运行的进程。

​	Pod 中封装着应用的容器（有的情况下是好几个容器），存储、独立的网络 IP，管理容器如何运行的策略选项。Pod 代表着部署的一个单位：kubernetes 中应用的一个实例，可能由一个或者多个容器组合在一起共享资源。

### Init容器

​	Pod 能够具有多个容器，应用运行在容器里面，但是它也可能有一个或多个先于应用容器启动的 Init 容器（理解：可以用于判断当前环境是否满足我们应用运行的要求）。

​	Init 容器与普通的容器非常像，除了如下两点：

- Init 容器总是运行到成功完成为止。
- 每个 Init 容器都必须在下一个 Init 容器启动之前成功完成。

​	如果 Pod 的 Init 容器失败，Kubernetes 会不断地重启该 Pod，直到 Init 容器成功为止。然而，如果 Pod 对应的 `restartPolicy` 为 Never，它不会重新启动。

​	如果为一个 Pod 指定了多个 Init 容器，那些容器会按顺序一次运行一个。只有当前面的 Init 容器必须运行成功后，才可以运行下一个 Init 容器。当所有的 Init 容器运行完成后，k8s才初始化 Pod 和运行应用容器。

#### Init 容器能做什么？

因为 Init 容器具有与应用程序容器分离的单独镜像，所以它们的启动相关代码具有如下优势：

- 它们可以包含并运行实用工具，但是出于安全考虑，是不建议在应用程序容器镜像中包含这些实用工具的。
- 它们可以包含使用工具和定制化代码来安装，但是不能出现在应用程序镜像中。例如，创建镜像没必要 `FROM` 另一个镜像，只需要在安装过程中使用类似 `sed`、 `awk`、 `python` 或 `dig` 这样的工具。
- 应用程序镜像可以分离出创建和部署的角色，而没有必要联合它们构建一个单独的镜像。
- Init 容器使用 Linux Namespace，所以相对应用程序容器来说具有不同的文件系统视图。因此，它们能够具有访问 Secret 的权限，而应用程序容器则不能。
- 它们必须在应用程序容器启动之前运行完成，而应用程序容器是并行运行的，所以 Init 容器能够提供了一种简单的阻塞或延迟应用容器的启动的方法，直到满足了一组先决条件。

#### 模板

使用`initContainers`字段指定所要的init容器

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: myapp-pod
  labels:
    app: myapp
spec:
  containers:
  - name: myapp-container
    image: busybox
    command: ['sh', '-c', 'echo The app is running! && sleep 3600']
  initContainers: #指定容器为 Init 容器
  - name: init-myservice
    image: busybox
    command: ['sh', '-c', 'until nslookup myservice; do echo waiting for myservice; sleep 2; done;']
  - name: init-mydb
    image: busybox
    command: ['sh', '-c', 'until nslookup mydb; do echo waiting for mydb; sleep 2; done;']
```

### Pause 容器

Pause 容器，又叫 Infra 容器

#### 特点

- 镜像非常小，目前在 700KB 左右
- 永远处于 Pause (暂停) 状态

#### 背景

​	像 Pod 这样一个东西，本身是一个逻辑概念。那在机器上，它究竟是怎么实现的呢？这就是我们要解释的一个问题。

​	既然说 Pod 要解决这个问题，核心就在于如何让一个 Pod 里的多个容器之间最高效的共享某些资源和数据。

​	因为容器之间原本是被 Linux Namespace 和 cgroups 隔开的，所以现在实际要解决的是怎么去打破这个隔离，然后共享某些事情和某些信息。这就是 Pod 的设计要解决的核心问题所在。

​	所以说具体的解法分为两个部分：网络和存储。

​	Pause 容器就是为解决 Pod 中的网络问题而生的。

#### 实现

Pod 里的多个容器怎么去共享网络？下面是个例子：

比如说现在有一个 Pod，其中包含了一个容器 A 和一个容器 B，它们两个就要共享 Network Namespace。在 Kubernetes 里的解法是这样的：它会在每个 Pod 里，额外起一个 Infra container 小容器来共享整个 Pod 的 Network Namespace。

Infra container 是一个非常小的镜像，大概 700KB 左右，是一个 C 语言写的、永远处于 “暂停” 状态的容器。由于有了这样一个 Infra container 之后，其他所有容器都会通过 Join Namespace 的方式加入到 Infra container 的 Network Namespace 中。

所以说一个 Pod 里面的所有容器，它们看到的网络视图是完全一样的。即：它们看到的网络设备、IP 地址、Mac 地址等等，跟网络相关的信息，其实全是一份，这一份都来自于 Pod 第一次创建的这个 Infra container。这就是 Pod 解决网络共享的一个解法。

在 Pod 里面，一定有一个 IP 地址，是这个 Pod 的 Network Namespace 对应的地址，也是这个 Infra container 的 IP 地址。所以大家看到的都是一份，而其他所有网络资源，都是一个 Pod 一份，并且被 Pod 中的所有容器共享。这就是 Pod 的网络实现方式。

由于需要有一个相当于说中间的容器存在，所以整个 Pod 里面，必然是 Infra container **第一个启动**。并且整个 Pod 的生命周期是等同于 Infra container 的生命周期的，与容器 A 和 B 是无关的。这也是为什么在 Kubernetes 里面，它是允许去单独更新 Pod 里的某一个镜像的，即：做这个操作，整个 Pod 不会重建，也不会重启，这是非常重要的一个设计。

### hook

​	Pod hook（钩子）是由 Kubernetes 管理的 kubelet 发起的，当容器中的进程启动前或者容器中的进程终止之前运行，这是包含在容器的生命周期之中。可以同时为 Pod 中的所有容器都配置 hook。

Hook 的类型包括两种：

- exec：执行一段命令
- HTTP：发送 HTTP 请求。

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: lifecycle-demo
spec:
  containers:
  - name: lifecycle-demo-container
    image: nginx
    lifecycle:
      postStart:
        exec:
          command: ["/bin/sh", "-c", "echo Hello from the postStart handler> /usr/share/message"]
      preStop:
        exec:
          command: ["/usr/sbin/nginx","-s","quit"]
```

​	postStart 在容器创建之后（但并不能保证钩子会在容器 ENTRYPOINT 之前）执行，这时候 Pod 已经被调度到某台 node 上，被某个 kubelet 管理了，这时候 kubelet 会调用 postStart 操作，该操作跟容器的启动命令是在同步执行的，也就是说**在 postStart 操作执行完成之前，kubelet 会锁住容器，不让应用程序的进程启动，只有在 postStart 操作完成之后容器的状态才会被设置成为 RUNNING**。

PreStop 在容器终止之前被同步阻塞调用，常用于在容器结束前优雅的释放资源。

如果 postStart 或者 preStop hook 失败，将会终止容器。

### 生命周期

Pod 的 `status` 字段是一个 PodStatus 对象，PodStatus中有一个 `phase` 字段。

Pod 的相位（phase）是 Pod 在其生命周期中的简单宏观概述。该字段并不是对容器或 Pod 的综合汇总，也不是为了做为综合状态机。

Pod 相位的数量和含义是严格指定的。除了本文档中列举的状态外，不应该再假定 Pod 有其他的 `phase` 值。

下面是 `phase` 可能的值：

- 挂起（Pending）：Pod 已被 Kubernetes 系统接受，但有一个或者多个容器镜像尚未创建。等待时间包括调度 Pod 的时间和通过网络下载镜像的时间，这可能需要花点时间。
- 运行中（Running）：该 Pod 已经绑定到了一个节点上，Pod 中所有的容器都已被创建。至少有一个容器正在运行，或者正处于启动或重启状态。
- 成功（Succeeded）：Pod 中的所有容器都被成功终止，并且不会再重启。
- 失败（Failed）：Pod 中的所有容器都已终止了，并且至少有一个容器是因为失败终止。也就是说，容器以非0状态退出或者被系统终止。
- 未知（Unknown）：因为某些原因无法取得 Pod 的状态，通常是因为与 Pod 所在主机通信失败。

![](kubernetes-pod-life-cycle.jpg)

#### 容器探针

### 镜像拉取策略

- IfNotPresent : 默认值，镜像在宿主机上不存在时才拉取
- Always ：每次创建Pod都会重新拉取一次镜像
- Never ：Pod 永远不会主动拉取这个镜像

<img src="C:\Users\zhou jielun\Desktop\k8s概念\3.jpg"  />

### 资源限制

<img src="C:\Users\zhou jielun\Desktop\k8s概念\1.jpg"  />

### 重启机制

```
restartPolicy ： ××
```

- Always ：当容器终止退出后，总是重启容器，默认策略
- OnFailure : 当容器异常退出（退出状态码非0）时，才重启容器
- Never : 当容器终止退出，从不重启容器

### 健康检查

- livenessProbe（存活检查）：如果检查失败，将杀死容器，根据Pod的restartPolicy来操作
- readinessProbe（就绪检查）：如果检查失败，k8s会把pod从service endpoints中剔除

<img src="C:\Users\zhou jielun\Desktop\k8s概念\2.jpg"  />

### 亲和性

#### 为什么需要亲和性

​	有时候某一些特定的pod需要部署在特定的机器上

#### 干预pod的调度方式

- PodSpec里面的nodeName字段指定
- nodeSelector是PodSpec中的一个字段，nodeSelector是最简单实现将pod运行在特定node节点的实现方式，其通过指定key和value键值对的方式实现，需要node设置上匹配的Labels，节点调度的时候指定上特定的labels即可。如下以node-2添加一个app:web的labels，调度pod的时候通过nodeSelector选择该labels：
- 节点亲和性

#### 节点（node）亲和性

- requiredDuringSchedulingIgnoredDuringExecution ：条件`必须`满足
- preferredDuringSchedulingIgnoredDuringExecution：条件`尽量`满足(不保证总是满足)

​	`IgnoredDuringExecution`的意思是节点标签发送变化时，并不会驱逐不符合条件的pod。

​	节点的亲和性通过PodSpec的`affinity`字段下的`nodeAffinity`字段进行指定，操作符支持`In`,`NotIn`,`Exists`,`DoesNotExist`,`Gt`,`Lt`。

反亲和性通过`NotIn`和`DoesNotExist`实现。

1. 如果同时指定了`nodeAffinity`和`nodeSelector`，则必须同时满足两个条件的node才是可调度的(AND关系)，
2. 如果指定了多个`nodeSelectorTerms`，则满足任意一个`nodeSelectTerms`的node均是可调度的(OR关系)。
3. 如果指定了多个`matchExpressions`，则必须同时满足所以`matchExpressions`的条件的node才是可调度的(AND关系)。
4. `preferredDuringSchedulingIgnoredDuringExecution`中的`weight`字段则是计算权重，多个条件权重加起来分数最高的node被调度的优先级最高。

```
apiVersion: v1
kind: Pod
metadata:
  name: mysql
spec:
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
          - matchExpressions:
              - key: mysql
                operator: In
                values:
                  - true
  containers:
    - name: mysql
      image: mysql:5.6.50
      imagePullPolicy: IfNotPresent
```

#### pod 亲和与反亲和

​	如果希望2个pod在一个node上，那么可以用到node的亲和性，或者使用pod亲和性。如果希望2个io密集型或者2个CPU密集型的pod不在一个node上，就可以用到pod的反亲和性。

- requiredDuringSchedulingIgnoredDuringExecution
- requiredDuringSchedulingIgnoredDuringExecution

​	pod亲和与反亲和通过`PodSpec`的`affinity`字段下的`podAffinity`和`podAntiAffinity`字段指定，前者为亲和性；后者为反亲和性。操作符仅支持`In`,`NotIn`,`Exists`,`DoesNotExist`

### 污点（Taints）和污点容忍（Tolerations）

​	Taint（污点）和 Toleration（容忍）可以作用于 node 和 pod 上，其目的是优化 pod 在集群间的调度，这跟节点亲和性类似，只不过它们作用的方式相反，具有 taint 的 node 和 pod 是**互斥关系**，而具有节点亲和性关系的 node 和 pod 是相吸的。另外还有可以给 node 节点设置 label，通过给 pod 设置 `nodeSelector` 将 pod 调度到具有匹配标签的节点上。

​	Taint 和 toleration 相互配合，可以用来避免 pod 被分配到不合适的节点上。每个节点上都可以应用**一个或多个** taint ，这表示对于那些不能容忍这些 taint 的 pod，是不会被该节点接受的。如果将 toleration 应用于 pod 上，则表示这些 pod 可以（但不要求）被调度到具有相应 taint 的节点上

可以通过下面命令查看污点

```
kubectl describe node <name>
```

#### Taints污点的组成

​	使用kubectl taint命令可以给某个**Node节点**设置污点，Node被设置污点之后就和Pod之间存在一种**相斥**的关系，可以让Node拒绝Pod的调度执行，甚至将Node上已经存在的Pod驱逐出去。

每个污点的组成如下：

```
key=value:effect
```

每个污点有一个key和value作为污点的标签，effect描述污点的作用。当前taint effect支持如下选项：

- NoSchedule：表示K8S将不会把Pod调度到具有该污点的Node节点上
- PreferNoSchedule：表示K8S将尽量避免把Pod调度到具有该污点的Node节点上
- NoExecute：表示K8S将不会把Pod调度到具有该污点的Node节点上，同时会将Node上已经存在的Pod驱逐出去

若taint 的 effect 值 NoExecute，它会影响已经在节点上运行的 pod：

- 如果 pod 不能容忍 effect 值为 NoExecute 的 taint，那么 pod 将马上被驱逐
- 如果 pod 能够容忍 effect 值为 NoExecute 的 taint，且在 toleration 定义中**没有指定 tolerationSeconds**，则 **pod** 会一直在这个节点上运行。
- 如果 pod 能够容忍 effect 值为 NoExecute 的 taint，但是在toleration定义中**指定了 tolerationSeconds**，则表示 **pod **还能在这个节点上继续运行的时间长度。

#### Tolerations容忍

​	设置了污点的Node将根据taint的effect：NoSchedule、PreferNoSchedule、NoExecute和Pod之间产生互斥的关系，Pod将在一定程度上不会被调度到Node上。

​	**但我们可以在Pod上设置容忍（Tolerations），意思是设置了容忍的Pod将可以容忍污点的存在，可以被调度到存在污点的Node上。**

```
tolerations:
- key: "key"
  operator: "Equal"
  value: "value"
  effect: "NoSchedule"
---
tolerations:
- key: "key"
  operator: "Exists"
   effect: "NoSchedule"
---
tolerations:
- key: "key"
  operator: "Equal"
  value: "value"
  effect: "NoExecute"
  tolerationSeconds: 3600
```

上面的参数说明

- 其中key、value、effect要与Node上设置的taint保持一致
- operator的值为Exists时，将会忽略value；只要有key和effect就行
- tolerationSeconds：表示pod 能够容忍 effect 值为 NoExecute 的 taint；当指定了 tolerationSeconds（容忍时间），则表示 pod 还能在这个节点上继续运行的时间长度。

##### 当不指定key值时

当不指定key值和effect值时，且operator为Exists，表示容忍所有的污点（能匹配污点所有的keys，values和effects）

##### 当不指定effect值时

当不指定effect值时，则能匹配污点key对应的所有effects情况

### Pod 和 Controller

​	Controller 可以创建和管理多个 Pod，提供副本管理、滚动升级和集群级别的自愈能力。例如，如果一个 Node 故障，Controller 就能自动将该节点上的 Pod 调度到其他健康的 Node 上。

​	通常，Controller 会用你提供的 Pod Template 来创建相应的 Pod。
​	当Pod宕机后重新生成时，其IP等状态信息可能会变动、不固定，Service会根据Pod的Label对这些状态信息进行监控和变更，保证上游服务不受Pod的变动而影响。因此由于pod的IP是可能变动的，k8s提供了service资源，service会对提供同一个服务的多个pod进行聚合，并且提供一个统一的入口地址。通过访问service的入口地址就能访问到后面的pod服务。

### Service 

service通过label和slector标签与pod建立关联

#### 使用场景

- 当客户端想要访问K8S集群中的pod时，需要知道pod的ip以及端口，那K8S中如何在不知道pod的地址信息的情况下进行pod服务的快速连接？


- 若某一node上的pod发生故障，K8S最大的特点就是能够给感知和重启该pod，但是pod重启后ip会发生变化，那么客户端如何感知并保持对pod的访问？


- 如果多个pod组合在一起形成pod组，如何在被访问时达到负载均衡的效果？

#### 具体的作用和场景如下

- 通过Pod的Label Selector访问Pod组。
- Service的IP保持不变（Headless Servcie除外），保证了访问接口的稳定性，屏蔽了Pod的IP地址变化带来的影响，进而实现解耦合。虽然这样，还是建议使用ServiceName进行访问。

#### 工作机制



#### 类型

##### ClusterIP

​	通过集群的内部 IP 暴露服务，选择该值时服务只能够在集群内部访问。 这也是默认的 `ServiceType`

##### NodePort

​	通过每个节点上的 IP 和静态端口（`NodePort`）暴露服务。 `NodePort` 服务会路由到自动创建的 `ClusterIP` 服务。 通过请求 `<节点 IP>:<节点端口>`，你可以从集群的外部访问一个 `NodePort` 服务。

##### LoadBalancer

​	使用云提供商的负载均衡器向外部暴露服务。 外部负载均衡器可以将流量路由到自动创建的 `NodePort` 服务和 `ClusterIP` 服务上

##### ExternalName

​	Service的ExternalName方式实现，即设置Service的type为ExternalName。这样做的好处就是内部服务访问外部服务的时候是通过别名来访问的，屏蔽了外部服务的真实信息，外部服务对内部服务透明，外部服务的修改基本上不会影响到内部服务的访问，做到了内部服务和外部服务解耦合。

#### Headless Service

​	有时不需要或不想要负载均衡，以及单独的 Service IP。 遇到这种情况，可以通过指定 Cluster IP（`spec.clusterIP`）的值为 `"None"` 来创建 `Headless` Service。
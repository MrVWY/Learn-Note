&ensp;&ensp;Pod 是在 Kubernetes 集群中运行部署应用或服务的最小单元，它是可以支持多容器的。Pod 的设计理念是支持多个容器在一个 Pod 中共享网络地址和文件系统，可以通过进程间通信和文件共享这种简单高效的方式组合完成服务。Pod 对多容器的支持是 K8 最基础的设计理念。

&ensp;&ensp;简单来说就是，不同的团队各自开发构建自己的容器镜像，在部署的时候组合成一个微服务对外提供服务。

&ensp;&ensp;关于Pod的概念

- Init容器
- Pause容器
- hook钩子
- 生命周期
- 镜像拉取策略
- 资源限制
- 重启机制
- 健康检查
- 亲和性
- 污点(Taints)和污点容忍(Tolerations)
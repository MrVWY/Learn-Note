&ensp;&ensp;Secret 解决了密码、token、密钥等敏感数据的配置问题，而不需要把这些敏感数据暴露到镜像或者 Pod Spec 中。Secret 可以以 Volume 或者环境变量的方式使用。

Secret 有三种类型：

- Service Account ：用来访问 Kubernetes API，由 Kubernetes 自动创建，并且会自动挂载到 Pod 的 /run/secrets/kubernetes.io/serviceaccount 目录中；
- Opaque ：base64 编码格式的 Secret，用来存储密码、密钥等；
- kubernetes.io/dockerconfigjson ：用来存储私有 docker registry 的认证信息

&ensp;&ensp;具体内容：https://jimmysong.io/kubernetes-handbook/concepts/secret.html
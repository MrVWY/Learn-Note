### 什么是helm
&ensp;&ensp;helm可以把这些零零散散的应用资源yaml文件放在一起进行统一配置，也就是说helm能把你部署的应用所关联的yaml文件（
例如PVC、PV、Pod、Server、Ingress等） 给整合起来，下次直接用helm一键部署即可

### helm基础概念
- Charts：Helm使用的打包格式，一个Chart包含了一组K8s资源集合的描述文件。Chart有特定的文件目录结构，
  如果开发者想自定义一个新的 Chart，只需要使用Helm create命令生成一个目录结构即可进行开发
- Release：通过Helm将Chart部署到 K8s集群时创建的特定实例，包含了部署在容器集群内的各种应用资源
- Tiller：Helm 2.x版本中，Helm采用Client/Server的设计，Tiller就是Helm的Server部分，需要具备集群管理员权限才能安装到K8s集群中运行。Tiller与Helm client进行交互，接收client的请求，再与K8s API Server通信，
  根据传递的Charts来生成Release。而在最新的Helm 3.x中，据说是为了安全性考虑移除了Tiller
- Chart Repository：Helm Chart包仓库，提供了很多应用的Chart包供用户下载使用，官方仓库的地址是https://hub.helm.sh


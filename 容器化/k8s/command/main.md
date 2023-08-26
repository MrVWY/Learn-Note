k8s基础命令操作

## Pod
查看所有 pod 列表, -n 后跟 namespace, 查看指定的命名空间
kubectl get pod
kubectl get pod -n kube
kubectl get pod -o wide
kubectl describe pod <pod-name>  //显示 Pod 的详细信息, 特别是查看 pod 无法创建的时候的日志
kubectl delete pod --all //删除所有 Pod
kubectl scale deploy/nginx-1 --replicas=0 //伸缩 pod 副本
kubectl scale deploy/nginx-1 --replicas=1 //伸缩 pod 副本
kubectl get pod <POD名称> -n <NAMESPACE名称> -o yaml | kubectl replace --force -f - //重启 pod

## 查看 RC 和 service 列表， -o wide 查看详细信息
kubectl get rc,svc
kubectl get pod,svc -o wide
kubectl get pod -o yaml

## 显示 Node 的详细信息
kubectl describe node 192.168.0.212

## 根据 yaml 创建资源, apply 可以重复执行，create 不行
kubectl create -f pod.yaml
kubectl apply -f pod.yaml

## 基于 pod.yaml 定义的名称删除 pod
kubectl delete -f pod.yaml

## 删除所有包含某个 label 的pod 和 service
kubectl delete pod,svc -l name=<>

## 查看 endpoint 列表
kubectl get endpoints

## 执行 pod 的 date 命令
kubectl exec – date
kubectl exec – bash
kubectl exec – ping 10.24.51.9

## 通过bash获得 pod 中某个容器的TTY，相当于登录容器
kubectl exec -it -c – bash
eg:
kubectl exec -it redis-master-cln81 – bash

## 查看容器的日志
kubectl logs
kubectl logs -f # 实时查看日志
kubectl log -c <container_name> # 若 pod 只有一个容器，可以不加 -c

kubectl logs -l app=frontend # 返回所有标记为 app=frontend 的 pod 的合并日志。

## 查看注释
kubectl explain pod
kubectl explain pod.apiVersion

## 查看节点 labels
kubectl get node --show-labels

## 修改网络类型
kubectl patch service istio-ingressgateway -n istio-system -p ‘{“spec”:{“type”:“NodePort”}}’

## 查看前一个 pod 的日志，logs -p 选项
kubectl logs --tail 100 -p user-klvchen-v1.0-6f67dcc46b-5b4qb > pre.log

第一章主要讲了控制循环，事件驱动等几个基本概念。

## 控制循环

1. 读取资源状态（watch资源）
2. 更新集群状态
3. 更新资源状态
4. 回到第一步，继续控制循环

控制循环是由多个 `kubernetes` 组件负责协调的，在未满足期望值（spec）之前，各个控制器会不断调整，这个过程叫做调谐（reconcile）；当应用满足期望值时，控制循环终止，接下来会以一个 deployment pod 创建的例子，顺带领会控制循环这种设计思路。

## edge-driven （事件驱动）

deployment 的 pod  周期：

1. deployment controller （kube-controller-manager内部组件) 创建 deployment
1. replicaset controller （kube-controller-manager内部组件) 创建 replica （deployment的本质是 replica 的封装）
1. kube-scheduler 创建一个带空白 spec.nodeName 字段的pod，并把该 pod 加入调度队列
1. kubelet 空转（因为 nodeName 尚未指定）
1. scheduler 把调度队列里面的pod 取出来，并分配一个节点，写入 API server
1. 节点上的 kubelet 唤醒并拉起容器，并将状态实时更新到 status 字段，提交到 API server 
1. replicaset controller 空转
1. pod 终止，kubelet 设置 status 字段，标记pod 为 terminated ，提交到 API server 
1. replicaset controller 删除旧的pod，创建新的pod

从这个例子可以初步体会 kubernetes 这种 控制循环的设计思路。

## level driven（状态驱动）

level driven 有点拗口，我更愿意称之为状态驱动。

以replica set controller 为例，该控制器会比对应用当前状态与期望状态，根据差异进行调谐。

https://github.com/kubernetes/kubernetes/blob/16d33c49858579fbe13df52c065dbea6627e9aef/pkg/controller/replicaset/replica_set.go

## edge-driven or level driven ？

作者通过对比 edge-driven 和 level driven 的不足与特点，继而引出 kubernetes 这种 edge-driven 和 level driven 相结合的设计思路。



## 安装

https://github.com/kubernetes-sigs/kubebuilder/releases

## 初始化项目

kubebuilder init --domain bullshitprogram.com

```bash
.
├── Dockerfile
├── Makefile
├── PROJECT
├── bin
│   └── manager
├── config
│   ├── certmanager
│   │   ├── certificate.yaml
│   │   ├── kustomization.yaml
│   │   └── kustomizeconfig.yaml
│   ├── default
│   │   ├── kustomization.yaml
│   │   ├── manager_auth_proxy_patch.yaml
│   │   ├── manager_webhook_patch.yaml
│   │   └── webhookcainjection_patch.yaml
│   ├── manager
│   │   ├── kustomization.yaml
│   │   └── manager.yaml
│   ├── prometheus
│   │   ├── kustomization.yaml
│   │   └── monitor.yaml
│   ├── rbac
│   │   ├── auth_proxy_client_clusterrole.yaml
│   │   ├── auth_proxy_role.yaml
│   │   ├── auth_proxy_role_binding.yaml
│   │   ├── auth_proxy_service.yaml
│   │   ├── kustomization.yaml
│   │   ├── leader_election_role.yaml
│   │   ├── leader_election_role_binding.yaml
│   │   └── role_binding.yaml
│   └── webhook
│       ├── kustomization.yaml
│       ├── kustomizeconfig.yaml
│       └── service.yaml
├── go.mod
├── go.sum
├── hack
│   └── boilerplate.go.txt
└── main.go

9 directories, 30 files
```

kubebuilder create api --group bullshitprogram.com --version v1alpha1 --kind KubeKillerServer

之后会多出一些文件

```bash
.
├── Dockerfile
├── Makefile
├── PROJECT
├── api
│   └── v1alpha1
│       ├── groupversion_info.go
│       ├── kubekillerserver_types.go
│       └── zz_generated.deepcopy.go
├── bin
│   └── manager
├── config
│   ├── certmanager
│   │   ├── certificate.yaml
│   │   ├── kustomization.yaml
│   │   └── kustomizeconfig.yaml
│   ├── crd
│   │   ├── kustomization.yaml
│   │   ├── kustomizeconfig.yaml
│   │   └── patches
│   │       ├── cainjection_in_kubekillerservers.yaml
│   │       └── webhook_in_kubekillerservers.yaml
│   ├── default
│   │   ├── kustomization.yaml
│   │   ├── manager_auth_proxy_patch.yaml
│   │   ├── manager_webhook_patch.yaml
│   │   └── webhookcainjection_patch.yaml
│   ├── manager
│   │   ├── kustomization.yaml
│   │   └── manager.yaml
│   ├── prometheus
│   │   ├── kustomization.yaml
│   │   └── monitor.yaml
│   ├── rbac
│   │   ├── auth_proxy_client_clusterrole.yaml
│   │   ├── auth_proxy_role.yaml
│   │   ├── auth_proxy_role_binding.yaml
│   │   ├── auth_proxy_service.yaml
│   │   ├── kubekillerserver_editor_role.yaml
│   │   ├── kubekillerserver_viewer_role.yaml
│   │   ├── kustomization.yaml
│   │   ├── leader_election_role.yaml
│   │   ├── leader_election_role_binding.yaml
│   │   └── role_binding.yaml
│   ├── samples
│   │   └── bullshitprogram.com_v1alpha1_kubekillerserver.yaml
│   └── webhook
│       ├── kustomization.yaml
│       ├── kustomizeconfig.yaml
│       └── service.yaml
├── controllers
│   ├── kubekillerserver_controller.go
│   └── suite_test.go
├── go.mod
├── go.sum
├── hack
│   └── boilerplate.go.txt
└── main.go

15 directories, 42 files
```

## 实现

主要的方法在
    func (r *KubeKillerServerReconciler) Reconcile

这里我CRUD的逻辑是：

创建的逻辑：已有 CR 但无 deploy，创建资源（deploy，svc）
读取的逻辑：这部分可以忽略，代码体现是最后一行
更新的逻辑：CR里面的Spec（这里我只简单设定了 replica）跟 deploy 不一致时
删除的逻辑：找不到 CR 但找到 deploy 时做资源清理（删除deploy，删除svc）

对比了以往 sample-controller 生成 typed clinet 的操作，现在直接用 controller-runtime 实现CRUD 还是比较简便的。

## 运行

kubebuilder 是用 Makefile 来做构建操作的，基本的打包和编译都在 Makefile 里面。

比如 `make manifests` 是生成 CRD , `make run` 则是在本地运行 operator 。

```bash
➜  samples ka bullshitprogram.com_v1alpha1_kubekillerserver.yaml
kubekillerserver.bullshitprogram.com/sample-0 created
➜  samples kgsvc
NAME                           TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)   AGE
kube-killer-operatorsample-0   ClusterIP   172.16.253.146   <none>        80/TCP    3s
kubernetes                     ClusterIP   172.16.252.1     <none>        443/TCP   121d
➜  samples kgdep
NAME                           READY   UP-TO-DATE   AVAILABLE   AGE
kube-killer-operatorsample-0   1/1     1            1           6s
➜  samples krm kubekillerservers sample-0
kubekillerserver.bullshitprogram.com "sample-0" deleted
➜  samples kgdep
No resources found in default namespace.
➜  samples kgsvc
NAME         TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)   AGE
kubernetes   ClusterIP   172.16.252.1   <none>        443/TCP   121d
```



这一章节主要讲几种 kubernetes 程序的包管理方式和生产注意事项。

## 包管理方式

1. helm
2. kustomize
3. 

## 生产注意事项

1. 最小权限原则。对自己的 operator 创建一个 service account ，并在此基础上设置 RBAC；
2. 为 pod 的崩溃兜底
3. 日志和可观测性
4. 自动化构建，测试


这一章主要讲 CRD 这种东西。

apiextensions-apiserver 负责检查 CRD。

kubectl api-resources 可以查看 resources 归属的 KIND 和范围（比如是否namespace隔离）

## additionalPrinterColumns

可以通过添加 additionalPrinterColumns 这个字段，修改 kubectl get 返回的字段

```yaml
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
name: ats.cnat.programming-kubernetes.info
spec:
additionalPrinterColumns: (optional)
- name: kubectl column name
type: OpenAPI type for the column
format: OpenAPI format for the column (optional)
description: human-readable description of the column (optional)
priority: integer, always zero supported by kubectl
JSONPath: JSON path inside the CR for the displayed value
```

## Subresources

1. /api/v1/namespace/namespace/pods/name/logs
1. /api/v1/namespace/namespace/pods/name/portforward
1. /api/v1/namespace/namespace/pods/name/exec
1. /api/v1/namespace/namespace/pods/name/status

## 访问CRD的几种方式

1. client-go dynamic client (see “Dynamic Client”)
2. kubernetes-sigs/controller-runtime
3. client-gen
4. 

## Dynamic Client

k8s.io/client-go/dynamic

根据 gvr （group version resources） 动态获取资源

  client.Resource(gvr).Namespace(namespace).Get("foo", metav1.GetOptions{})
  
动态获取的资源类型是 *unstructured.Unstructured ，需要用内置方法获取对象属性

name, found, err := unstructured.NestedString(u.Object, "metadata", "name")



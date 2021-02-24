
# client-go

## protobuf 支持

```go
cfg, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
cfg.AcceptContentTypes = "application/vnd.kubernetes.protobuf,
application/json"
cfg.ContentType = "application/vnd.kubernetes.protobuf"
clientset, err := kubernetes.NewForConfig(cfg)
```

注意：crd不支持。

## 版本问题

由于 kubernetes 内的资源都是带版本的，所以在请求的时候也要注意版本的问题，如果服务器没有，而客户端没有，请求会出现问题。

像是版本号带 alpha 的，一般视为不稳定版本，beta 版本比 alpha 稍微好一点，而v1，v2则视为稳定版。

## runtime.Object & ObjectKind

```go
// Object interface must be supported by all API types registered with Scheme.
// Since objects in a scheme are expected to be serialized to the wire, the
// interface an Object must provide to the Scheme allows serializers to set
// the kind, version, and group the object is represented as. An Object may
// choose to return a no-op ObjectKindAccessor in cases where it is not
// expected to be serialized.
type Object interface {
GetObjectKind() schema.ObjectKind
DeepCopyObject() Object
}

// All objects that are serialized from a Scheme encode their type information.
// This interface is used by serialization to set type information from the
1// Scheme onto the serialized version of an object. For objects that cannot
// be serialized or have unique requirements, this interface may be a no-op.
type ObjectKind interface {
// SetGroupVersionKind sets or clears the intended serialized kind of an
// object. Passing kind nil should clear the current setting.
SetGroupVersionKind(kind GroupVersionKind)
// GroupVersionKind returns the stored group, version, and kind of an
// object, or nil if the object does not expose or provide these fields.
GroupVersionKind() GroupVersionKind
}
```

可以说，实现了 GVK 和深拷贝接口的，就能视为 kubernetes 对象（比如pod，deployment等）

## TypeMeta

类型元数据 TypeMeta 包含 Kind 和 APIVersion 2个字段，但在kubernetes 对象中，一般取内联的值。

```go
// Pod is a collection of containers that can run on a host. This resource is
// created by clients and scheduled onto hosts.
type Pod struct {
metav1.TypeMeta `json:",inline"`
// Standard object's metadata.
// +optional
metav1.ObjectMeta `json:"metadata,omitempty"`
// Specification of the desired behavior of the pod.
// +optional
Spec PodSpec `json:"spec,omitempty"`
// Most recently observed status of the pod.
// This data may not be up to date.
// Populated by the system.
// Read-only.
// +optional
Status PodStatus `json:"status,omitempty"`
}
```

以pod 为例，已yaml形式获取该资源对象时，结果一般是这样的：

```go
apiVersion: v1
kind: Pod
......
```
YAML 序列化器使用了JSON encoder 的标签，使得 TypeMeta 内的字段也放进来，这就叫做内联。

## ObjectMeta

```go
type ObjectMeta struct {
Name string `json:"name,omitempty"`
Namespace string `json:"namespace,omitempty"`
UID types.UID `json:"uid,omitempty"`
ResourceVersion string `json:"resourceVersion,omitempty"`
CreationTimestamp Time `json:"creationTimestamp,omitempty"`
DeletionTimestamp *Time `json:"deletionTimestamp,omitempty"`
Labels map[string]string `json:"labels,omitempty"`
Annotations map[string]string `json:"annotations,omitempty"`
...
}
```

这个对应 metadata 里面的字段

## RESTMapping

传入 GVR ，转化成 GVK

RESTMapping(gk schema.GroupKind, versions ...string) (*RESTMapping, error)

## shared informer

```go
// NewFilteredSharedInformerFactory constructs a new instance of
// sharedInformerFactory. Listers obtained via this sharedInformerFactory will be
// subject to the same filters as specified here.// j p
func NewFilteredSharedInformerFactory(
client versioned.Interface, defaultResync time.Duration,
namespace string,
tweakListOptions internalinterfaces.TweakListOptionsFunc
) SharedInformerFactor
type TweakListOptionsFunc func(*v1.ListOptions)
```

为了减少watch的负载，一般建议使用 shared informer 去监听资源变化。

通过 ObjectMeta.resourceVersion 这个字段可以判定资源是否被更新。


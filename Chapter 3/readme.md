
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

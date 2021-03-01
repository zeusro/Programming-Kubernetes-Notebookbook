
这一章主要讲与CR相关的4类代码生成器

## 生成器

1. deepcopy-gen
  生成 DeepCopy 和 DeepCopyInto 方法
1. client-gen
  生成强类型客户端。由于 custom resources 属于拓展类型，原生客户端没有，所以需要对 cr 生成一个强类型的客户端 
1. informer-gen
  提供一个基于事件的 informer，用于watch cr 变化（CRUD）
1. lister-gen
  只读缓存，用于监听 cr 的 get 和list 请求
  
## tag

由生成器生成的代码一般会通过注释显示其生成来源和特性

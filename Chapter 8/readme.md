
这一章节主要 自定义 API Server 。

## CRD 的问题

1. 只支持 etcd
2. 只支持JSON，不支持 protobuf （一种高性能的序列化语言）
3. 只支持2种子资源接口 （ /status 和 /scale）
4. 不支持优雅删除
5. 显著增加 api server 负担
6. 只支持 CRUD 原语
7. 不支持跨  API groups 共享存储

## 自定义 API Server 相比 CRD 的优势

1. 底层存储无关（像metrics server 存在内存里面）
2. 支持 protobuf
3. 支持自定义子资源
4. 可以实现优雅删除
5. 支持复杂验证
6. 支持自定义语义

等。

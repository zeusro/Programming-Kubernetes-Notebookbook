
这一章节主要讲了 Kubebuilder 和 Operator SDK 这2种生成 CRD 的方式。

但实测，现在的 operator-sdk （版本1.4.2) 为了管理 operator 引入了 OLM ，我个人不是很喜欢，所以还是推荐使用 kubebuilder 或者原生的代码生成器（https://github.com/kubernetes/code-generator） 作为生成 operator 的手段。
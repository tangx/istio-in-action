# 初始化第一个项目

1. 项目代码在 `https://github.com/tangx/istio-in-action`
2. 命令中有很多快捷键， 参考 [install and prepare](01-install.md)

## 1. 创建 namespace 并开启整体 istio 注入

### 1.1 创建 namespace myistio

```bash
kc ns myistio
    namespace/myistio created

kns myistio
    Context "default" modified.
    Active namespace is "myistio".
```

### 1.2 向 namespace 中开启 istio 注入

```bash
# 向 ns 加入标签 istio-injection=enabled ， 开启注入
kubectl label namespace myistio istio-injection=enabled
    namespace/myistio labeleds


# 查看具有 istio-injection 标签的 ns
kgall ns -L istio-injection
    NAME              STATUS   AGE   ISTIO-INJECTION
    kube-system       Active   42d
    kube-public       Active   42d
    istio-system      Active   10m
    myistio           Active   11s   enabled
    default           Active   42d
```


## 2. 创建第一个项目

```go
.
├── cmd
│   └── prod   // 项目名称
├── dockerfiles  // 编译镜像使用的 dockerfile
├── scripts
│   └── deployment  // k8s 发布时用的 yaml 文件。 通过渲染发布
├── .version  // 版本编号管理
└── version   // go 程序版本注入
```


### 2.1 程序说明

程序功能很简单， 就是在请求地址 **http://servername/prods/list** 是返回一个固定结果, 如下。

```json
{
  "data": {
    "Name": "istio in action",
    "Price": 300,
    "Reviews": null
  },
  "version": "v1.0.0"
}
```

1. data 的值是在 gin handler 中固定写死的。 `/cmd/prod/main.go`
2. version 是通过 `/version/version.go` 在编译时注入的， 其值来源于文件 `.version`。

使用如下命令进行编译发布

```bash
make apply.docker
```

## 3. 简单测试

在 myistio namespace 下创建一个容器， 作为客户端。

```bash
ksn myistio
k create deployment toolbox --image=nginx:alpine
```

进入创建的工具容器， 使用 curl 调用 prod 服务。 确认调用无异常。


```bash
keti toolbox-77889d56fd-dnfbz sh
    kubectl exec [POD] [COMMAND] is DEPRECATED and will be removed in a future version. Use kubectl exec [POD] -- [COMMAND] instead.

curl svc-prod/prods/list
    {"data":{"Name":"istio in action","Price":300,"Reviews":null},"version":"v1.0.0"}
```

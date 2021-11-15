# VirtualService 使用路径重写

有了 VirtualService 的路径重写功能后， 就更符合 Ingress 的标准定义了。 

但 VirtualService 不仅仅如此， 路径重写包含了三种方式

1. `prefix`: 前缀匹配。 只要 uri 路径的 **前段** 匹配则转发。 **后端** 自动补齐。
2. `exact`: 精确匹配。 只有 uri **全部** 匹配才转发， 并且只能转发某一个固定地址。
    + **精确匹配**
3. `regex`: 正则匹配。 只有 uri 全部路径 **满足正则规则** 才转发。
    + 正则规则: https://github.com/google/re2/wiki/Syntax
    + **精确匹配**， 正则模式也是精确匹配目标路径。

> 补充: 关于正则匹配模式官网资料也很少。

```yaml
# https://istio.io/latest/docs/reference/config/networking/virtual-service/#HTTPRewrite
---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: vs-prod
  namespace: myistio
spec:
  gateways:
    - istio-tangx-in
  hosts:
    - svc-prod
    - istio.tangx.in
  http:
  - name: "prefix-match"  # 规则名称
    match:
    - uri:
        prefix: "/p1"  # 新路径, prefix 前缀匹配， 满足 /p1 的都要被重写
    rewrite:
      uri: "/prod"    # 老路径
    route:
    - destination:
        host: svc-prod  # 后端服务

  - name: "exact-match"
    match:
    - uri:
        exact: "/p2-list" # 新路径， exact 精确匹配， 只能满足 /p2-list
    rewrite:
      uri: "/prod/list"   # 精确匹配
    route:
    - destination:
        host: svc-prod

  - name: "regex-match"
    match:
    - uri:
        regex: "/pr[1-3]/.*" # 新路径, regex 正则匹配。 正则匹配的整个 uri，因此允许所有要 使用 `.*`。 正则规则使用: https://github.com/google/re2/wiki/Syntax
    rewrite:
      uri: "/prod/list"  ## 精确匹配路径
    route:
    - destination:
        host: svc-prod

```

> https://istio.io/latest/docs/reference/config/networking/virtual-service/#HTTPRewrite


## 测试

执行命令， 部署环境。 （快捷键见第一章）

```bash
ka -f istio-samples/05
```

使用 [05.http](/istio-samples/05/05.http) 中的测试用例， 进行测试。

```bash

### GET，原访问地址
#     现在已经 404
GET http://istio.tangx.in/prod/list


### GET 使用路径重写: prefix 前缀匹配
GET http://istio.tangx.in/p1/list


### GET 使用路径重写: exact 精确匹配
GET http://istio.tangx.in/p2-list


### GET 使用路径重写: regex 正则匹配(有效)
GET http://istio.tangx.in/pr3/list

### GET 使用路径重写: regex 正则匹配(无效)
GET http://istio.tangx.in/pr4/list
```


## 不同的 404 not found

客户端请求后得到的 **404 not found** 有两种

1. istio 没有匹配到路由规则而返回的 404.

```bash

### GET，原访问地址
#     现在已经 404, istio 返回
GET http://istio.tangx.in/prod/list

    # HTTP/1.1 404 Not Found
    # date: Mon, 15 Nov 2021 04:19:43 GMT
    # server: istio-envoy
    # connection: close
    # content-length: 0

```

2. istio 成功将请求转发到后端server， 后端 server 找不到路由而返回的 404。

```bash
### GET 使用路径重写: prefix 前缀匹配
#     404 not found, server 返回。
GET http://istio.tangx.in/p1/list2

    # HTTP/1.1 404 Not Found
    # content-type: text/plain
    # date: Mon, 15 Nov 2021 04:20:09 GMT
    # content-length: 18
    # x-envoy-upstream-service-time: 0
    # server: istio-envoy
    # connection: close

    # 404 page not found
```
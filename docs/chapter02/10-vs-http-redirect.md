# VirtualService 路由重定向

在 VirtualService 配置中， 除了 http rewrite 路由重写之外， 还有 http redirect 路由重定向。 即常说的 30x。

> https://istio.io/latest/docs/reference/config/networking/virtual-service/#HTTPRedirect


## http redirect

VirtualService 重定向配置如下。 有三个重要参数

1. uri: 重定向后的 **uri**
2. redirectCode: 重定向时的 http response code。 ex: 301, 302。 默认值为 **301** 。
3. authority: 重定向后的 http host。 即 http response header 中的 location 字段。

```yaml
---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: review-http-redirect
  namespace: myistio
spec:
  gateways:
    - istio-tangx-in
  hosts:
    - svc-review
    - istio.tangx.in
  http:
    - match:
        - uri:
            exact: /review
      redirect:
        uri: /review/all
        redirectCode: 302
        authority: svc-review  # 重定向后的地址。
```

使用 `curl` 命令请求测试， 结果如下。

```bash
curl -I  http://istio.tangx.in/review

    HTTP/1.1 302 Found
    location: http://svc-review/review/all
    date: Mon, 15 Nov 2021 10:32:59 GMT
    server: istio-envoy
    transfer-encoding: chunked
```

可以看到已经正常实现重定向。 


## 兼顾内群内外的重定向

但是 `location: http://svc-review/review/all` 结果是集群内部地址， 而我们的请求时从集群外部发起的访问。

虽然可以将 `authority` 字段的值修改为 **集群外部地址**。

```yaml
  http:
    - match:
        - uri:
            exact: /review
      redirect:
        uri: /review/all
        redirectCode: 302
        # authority: svc-review
        authority: istio.tangx.in
```

但这是一个 **蠢到爆** 的方式。

1. 每次请求都必须要走 **外部网关**
2. 外部地址与 VirtualService 强耦合， 无法适配多地址的情况。


### 相同路由规则下 redirect 和 route 互斥

下面这个规则是不合法的， 在 **同一条** 路由规则下， redirect 和 route 互斥。

```yaml
---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: review-http-redirect
  namespace: myistio
spec:
  gateways:
    - istio-tangx-in
  hosts:
    - svc-review
    - istio.tangx.in
  http:
    - match:
        - uri:
            exact: /review
      redirect:
        uri: /review/all
        redirectCode: 302
      route:    # redirect 和 route 在同一条规则下互斥
          - destination:
              host: svc-review

```

报错如下

```
for: "istio-samples/10-http-redirect/vs.yml": admission webhook "validation.istio.io" denied the request: configuration is invalid: HTTP route cannot contain both route and redirect
```

### 使用多路由规则无法兼顾鱼和熊掌

> 遗留问题: 虽然 `redirect` 和 `route` 不能在 **同一个** 规则下。 但是他们可以在 **不同** 规则下。 因此使用 **多条** 路由规则即可兼得鱼和熊掌 ??? 

经测试发现， 如下包含 gateway 字段的 VirtualService 定义， 无法完成内网的 http-redirect。 

```yaml
---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: review-http-redirect
  namespace: myistio
spec:
  gateways:
    - istio-tangx-in
  hosts:
    - svc-review
    - istio.tangx.in
  http:
    # 规则重定向
    - match:
        - uri:
            exact: /review
      redirect:
        uri: /review/all
        redirectCode: 302

    # 路由转发
    - match:
        - uri:
            prefix: /
      route:
          - destination:
              host: svc-review
```

在集群内部的 toolbox 容器中的执行命令， 出现  not found 错误。

```bash
curl -I http://svc-review/review

    HTTP/1.1 404 Not Found
    date: Mon, 15 Nov 2021 11:08:00 GMT
    server: istio-envoy
    transfer-encoding: chunked
```

### 使用多配置兼得鱼和熊掌（不优雅）

没办法， 只能创建两个配置实现内外网的重定向

1. **不包含 gateway** 的 [vs.yml](/istio-samples/10-http-redirect/vs.yml) 
2. **包含 gateway** 的 [vs-gateway.yml](/istio-samples/10-http-redirect/vs-gateway.yml) 

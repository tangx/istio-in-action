# VirtualService 重试机制

在 Istio VirtualService 中， 有一个很关键的机制： **重试**。 

发起重试不需要业务本身实现， 而是 istio 通过 envoy 发起的。

其中有几个关键参数

1. `attempts`: 重试次数（不含初始请求）， 即最大请求次数为 n+1。

2. `perTryTimeout`: 发起重试的间隔时间。
    + 必须大于 1ms。 
    + 默认于 http route 中的 timeout 一致， 即无 timeout 时间

3. `retryOn`: 执行重试的触发条件。
    + 条件值有 envoy 提供: https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/router_filter#x-envoy-retry-on

## http retry

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
            prefix: /
      route:
        - destination:
            host: svc-review
      retries:  # 重试
        attempts: 3 # 重试次数（不含本身一次）， 共计 4 次。
        perTryTimeout: 2s # 间隔时间， 默认 25ms。必须大于 1ms
        retryOn: gateway-error,connect-failure,refused-stream # 触发条件
```

## 测试

部署用例 11 进行测试。

```bash
ka -f istio-samples/11-http-retry
```

执行 curl 请求命令， 通过结果可以看到， 总共耗时 8 秒。

```bash
time curl http://istio.tangx.in/review/delay?delay=3

    upstream request timeout

    real    0m8.118s
    user    0m0.000s
    sys     0m0.010s
```

通过 review 的日志可以看到， 总共请求了 **4次 (1+3)**, 每次间隔 **2秒** 。 刚好 8 秒超时

```log
[GIN] 2021/11/15 - 15:56:08 | 200 |  3.000822016s |       10.42.0.1 | GET      "/review/delay?delay=3"
[GIN] 2021/11/15 - 15:56:10 | 200 |  3.000916703s |       10.42.0.1 | GET      "/review/delay?delay=3"
[GIN] 2021/11/15 - 15:56:12 | 200 |  3.000723194s |       10.42.0.1 | GET      "/review/delay?delay=3"
[GIN] 2021/11/15 - 15:56:14 | 200 |  3.000565097s |       10.42.0.1 | GET      "/review/delay?delay=3"
```


## 设置 timeout

如下, 增加 http route 的全局 timeout 参数。

```yaml
# ... 略
      timeout: 5s # 总请求时间不会操作 timeout 时常

      retries:  # 重试
        attempts: 3 
        perTryTimeout: 2s 
        retryOn: gateway-error,connect-failure,refused-stream 
```

虽然按照 **重试** 逻辑依旧需要 4次 8秒。 但受限于 timeout 的阈值， 请求在 5秒 后超时退出。

```bash
time curl  http://istio.tangx.in/review/delay?delay=3

    upstream request timeout

    real    0m5.012s
    user    0m0.009s
    sys     0m0.000s
```

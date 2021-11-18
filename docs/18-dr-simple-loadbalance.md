# 使用 DestionationRule 流量控制策略 - 简单负载均衡

**简单负载均衡** 策略， 官方指定名称。

1. `ROUND_ROBIN`: 轮训策略， 默认。
2. `LEAST_CONN`: 最小连接数。 **随机** 选择 **两个健康** 后端， 通过 O(1) 算法选择连接数最少的后端。
3. `RANDOM`: 随机选择了一个 **健康** 后端。 如果 **没有配置健康检查策略**， 随机策略比轮训更好。
4. `PASSTHROUGH`: 此选项会将连接转发到调用者请求的原始 IP 地址，而不进行任何形式的负载平衡。必须谨慎使用此选项。它适用于高级用例。有关更多详细信息，请参阅 Envoy 中的原始目标负载均衡器。


```yaml
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
  - name: "v1-subset"
    route:
    - destination:
        host: svc-prod


---

apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: dr-prod
spec:
  host: svc-prod
  trafficPolicy:
    loadBalancer:
      # simple: RANDOM
      simple: ROUND_ROBIN
      # simple: PASSTHROUGH
      # simple: LEAST_CONN
```

## 部署测试

```bash
ka -f istio-samples/18-dr-simple-loadbalance
```

10000 次请求， 2个后端， 差别不是很大

```bash

root@toolbox-54f88c8c95-82f4p:/tmp# ./18-dr-simple-loadbalance PASSTHROUGH
    v2.0.0 => 5227 次
    v1.0.0 => 4772 次

root@toolbox-54f88c8c95-82f4p:/tmp# ./18-dr-simple-loadbalance RANDOM
    v2.0.0 => 4943 次
    v1.0.0 => 5056 次

root@toolbox-54f88c8c95-82f4p:/tmp# ./18-dr-simple-loadbalance ROUND_ROBIN
    v2.0.0 => 5019 次
    v1.0.0 => 4981 次

root@toolbox-54f88c8c95-82f4p:/tmp# ./18-dr-simple-loadbalance LEAST_CONN
    v2.0.0 => 4990 次
    v1.0.0 => 5009 次
```
# 使用 istio Gateway 运维外部访问 




## 2.1. VirutalService

正常情况下， 
```yaml
---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: vs-prod
  namespace: myistio
spec:
  gateways: # 选择 gateway
    - istio-tangx-in
  hosts:
    - svc-prod
    - istio.tangx.in # 使用的外部地址 FSDN
  http:
  - route:
    - destination:
        host: svc-prod

```

## 2.2. Gateway

```yaml

---

# https://istio.io/latest/docs/reference/config/networking/gateway/

apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: istio-tangx-in
  namespace: myistio
spec:
  selector:
    istio: ingressgateway # 选择 ingressgateway, 省略则兼容所有
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
      # - myistio/istio.tangx.in # 只针对 ns myistio 有效
      - istio.tangx.in # 所有 ns 都有效
```
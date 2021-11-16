# Gateway 支持 https 访问 - 标准模式


https 证书在创建时， 需要 **保持** 与 **istio-ingressgateway** 服务在 **相同** 的 namespace。




## 配置 k8s secret

> 1. https://istio.io/latest/docs/reference/config/networking/gateway/#ServerTLSSettings
> 2. https://istio.io/latest/docs/tasks/traffic-management/ingress/ingress-control/



```bash

kubectl create secret generic wild-tangx-in \
    --from-file=key=./certificates/_.tangx.in.key  \
    --from-file=cert=./certificates/_.tangx.in.crt  \
    -n istio-system

```

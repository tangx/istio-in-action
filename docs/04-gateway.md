# 使用 istio Gateway 允许外部访问 

仅仅是简单的创建了 VirtualService 是不能实现集群外部的访问的。

在 Istio 中， 还有一个 Gateway 的概念。 顾名思义， Gateway 就是大门保安， 只允许具有特定特征的流量通过。



## 1.1. 创建 Gateway

先来创建一个 Gateway

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
      # - myistio/istio.tangx.in # 只针对特定的 namespace myistio 有效
      - istio.tangx.in # 所有 ns 都有效
```

上述 gateway 注意以下几点。

1. 使用 `.spec.selector` 选择了绑定的 ingressgateway。 如果 **省略** 则绑定到所有的 ingressgateway。

```bash
kgall deployment -l istio=ingressgateway

    NAMESPACE      NAME                   READY   UP-TO-DATE   AVAILABLE   AGE
    istio-system   istio-ingressgateway   1/1     1            1           3d23h
```

2. `.spec.servers.port` 指定了 gateway 允许的 **端口** 和 **协议**。 
    + 截止 `istio v1.11.4` 只支持 `HTTP|HTTPS|GRPC|HTTP2|MONGO|TCP|TLS` 7中。

3. `.sepc.servers.hosts` 指定了允许通过的 **域名**。
    + 如果使用 `ns_name/istio.tangx.in` **namespace** 字段， 则表示只有 **特定** 的namespace 中生效。 
    + `istio.tangx.in` 如果没有 ns 字段， 则表示所有 ns 中生效。


## 1.2. VirutalService 定义

随后， 更新 VirtualService 配置

```yaml
---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: vs-prod
  namespace: myistio
spec:
  gateways: # 选择 gateway
    - istio-tangx-in  # 这里的名字要与 gateway 的名字匹配
  hosts:
    - svc-prod
    - istio.tangx.in  # 使用的外部地址 FQDN。 这里的域名是 gateway hosts 中定义的
  http:
  - route:
    - destination:
        host: svc-prod
```

需要注意

1. `.spec.gateways` 的列表值必须是存在的 gateway 名称
2. `.spec.hosts` 的值， 必须是上述选中的 gateway 中定义的。


## 2. 测试

运行如下命令创建相关环境

```bash
kubectl apply -f istio-samples/04/
```

使用 [04.http](/istio-samples/04/04.http) 的 GET 请求进行测试

![gateway result](`./imgs/04/04-gateway.png)

> 注意:  使用访问的外部域名 `istio.tangx.in` 一定要进行 dns 解析。 或使用 `/etc/resolv.conf` 进行绑定。
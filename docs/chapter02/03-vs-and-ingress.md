# isti VirtualService 和 k8s Ingress

可以简单的认为 Ingress 是 k8s 中提出的流量入口转发的一个 **标准定义规范**（只是认为）。 怎么实现， 需要根据不同的 IngressController 的逻辑。

VirtualService 的部分功能就承担了 Ingress 的这一功能。

## 1. Ingress 与 VirtualService 的定义

k8s Ingress 配置

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  creationTimestamp: null
  name: simple
spec:
  rules:
  - host: foo.com  # 访问的域名
    http:
      paths:
      - backend:
          service:
            name: svc1  # 后端服务名称
            port:
              number: 80
        path: /bar  # 匹配的路由
        pathType: Exact  # 匹配模式
```

对应的 istio virtual service 配置如下

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: vs-prods
  namespace: myistio
spec:
  hosts:
  - svc-prod    # 访问地址
  http:
  - route:
    - destination:
        host: svc-prod  # 后端服务名称
```

可以看到，  VirtualService 的配置大致与 Ingress 相似。 这里 VS 中没有指定 path 规则， 将在以后路由转发中提到。

> https://istio.io/latest/zh/docs/concepts/traffic-management/#the-hosts-field

虚拟服务主机名可以是 IP 地址、DNS 名称，或者依赖于平台的一个简称（例如 Kubernetes 服务的短名称）， **隐式或显式地指向一个完全限定域名（FQDN）** 。您也可以使用通配符（“*”）前缀，让您创建一组匹配所有服务的路由规则。虚拟服务的 hosts 字段实际上不必是 Istio 服务注册的一部分，它只是虚拟的目标地址。这让您可以为没有路由到网格内部的虚拟主机建模。

> 隐式或显式地指向一个完全限定域名（FQDN）: 即， **客户端** 需要能解析该域名。  因此在内网时通常需要使用 **service name**， k8s coredns 完成了域名解析。 在使用非集群地址时 （ex. 公网域名时） 可以通过 dns 解析， 也可以是修改 /etc/hosts 文件。


## 2. 测试 VirtualService

使用如如下命令, 创建 VirtualService

```bash
kubectl apply -f istio-samples/03/
```

部署完成后， 进入到 toolbox 请求 svc-prod

```bash
curl svc-prod/prod/list

{"data":{"Name":"istio in action","Price":300,"Reviews":null},"version":"v1.0.0"}
```

这个时候能请求， 但是结果看不出什么有什么特别s， 打开之前部署的 kiali

进入到 `Graph -> namespace (myistio) -> traffic -> service grpha` 就可以看到流量请求了

![soolbox-svc-prod](../imgs/03/toolbox-svc-prod.png)


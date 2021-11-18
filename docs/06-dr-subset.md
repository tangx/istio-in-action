# 使用 DestinationRule Subset 进行路由分组(版本控制)

当一个程序并行发布多个版本的时候， 如 `prod-v1 / prod-v2`

```bash
kgd
    NAME      READY   UP-TO-DATE   AVAILABLE   AGE
    toolbox   1/1     1            1           3d22h
    prod-v1   1/1     1            1           16m
    prod-v2   1/1     1            1           16m
```

```json5
// 两个版本的测试结果， 仅定义为 version 不一致
{
  "data": {
    "Name": "istio in action",
    "Price": 300,
    "Reviews": null
  },
  "version": "v2.0.0"  //   "version": "v1.0.0" 
}
```


**k8s Service** 依旧实现最根本的 **服务级别的 Selector**。

```yaml
---
# Service
apiVersion: v1
kind: Service
metadata:
  labels:
    app: prod
  name: svc-prod
  namespace: myistio
spec:
  ports:
  - name: 80-8080
    port: 80
    protocol: TCP
    targetPort: 8080
  selector:
    app: prod
  type: ClusterIP
```

默认情况下会根据 VirtualService 的默认规则 **轮询** 到后端的所有服务。


## 使用 subset 实现路由控制

但是在一些特定的环境下，需要对路由或者流量进行精确的认为控制。 这个时候就需要对后端服务进行 **分组** 处理。 

这个时候就可以使用 istio 的 subset 功能。 subset 的定义为 `Service Version (服务版本)`， 产生的目的就是为了在持续集成场景中， 可以通过 **路由、 请求头(Header)、权重等** 等方式进行路由或流量控制，以便进行 A/B 测试、金丝雀测试等。

## `DestinationRule` 服务分组

在 vs 使用 subset 的时候， **必须依赖** `DestinationRule` 控制器进行 **后端服务的分组**。

DR 通过 label 规则对后端进行服务分组。

这样当流量达到 envoy 的以后， 进一步根据 `label-> version:v1` 选择真是的后端服务。

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: dr-prod
spec:
  host: svc-prod
  subsets:
  - name: groupv1
    labels:
      version: v1
  - name: groupv2
    labels:
      version: v2
```


### 1. 使用 **流量权重** 实现分组

在 **同一个** 路由规则下， 可以使用 **权重模式** ， 将流量分发到不同的后端 subset 组中。

> 注意: 权重值的总和必须是 100 。

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
    route:                  # 同一个 route 下面的两个 destination
    - destination:
        host: svc-prod
        subset: groupv1     # subset 的值与 DestinationRule 中定义一致
      weight: 25
    - destination:
        host: svc-prod
        subset: groupv2
      weight: 75
```


使用如下命令进行测试

```bash
ka -f istio-samples/06-dr-subset/
ka -f istio-samples/06-dr-subset/vs/03-subset-weight.yml
```


### 2. 使用 **路由重写** 实现分组

**路由重写** 只是路由分组其中一个小的分支。 同样还可以使用 header， queryParams 参数。  逻辑都是类似的。

在 VirtualService 配置中， 使用多个 route 规则， 将流量转发到不同的后端组。


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
  - name: "v2-routes"      # 路由重写分组， 是针对不同的路由匹配规则
    match:
      - uri:
          prefix: "/v2/prod"    # 新增一个路由匹配规则， 只有 uri 满足 /v2/prod 才会访问 v2 版本的 pod
    rewrite:
      uri: "/prod"
    route:
    - destination:
        host: svc-prod
        subset: groupv2
  - name: "default-routes"   # 可以说是默认分组
    route:
    - destination:
        host: svc-prod
        subset: groupv1
```

使用如下命令进行测试

```bash
ka -f istio-samples/06/
ka -f istio-samples/06/vs/02-subset-rewrite-path.yml
```


## 流量的目的地址

这里总结一下， 无论是在 `VirtualService` 还是在 `DestinationRule` 中， 流量的目的地址都是 `k8s service`。

> 注意: 这里的 `k8s service` 指的是在 istio 以外能满足 FQDN


```yaml
---
# DestinationRule
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
spec:
  host: svc-prod            # 目的地址是 svc
  # ....


---
# VirtualService
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
spec:
  http:
  - name: "default-routes"
    route:
    - destination:
        host: svc-prod      # 目的地址是 svc
        subset: groupv1
  # .....
```
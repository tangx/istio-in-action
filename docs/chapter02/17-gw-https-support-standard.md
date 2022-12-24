# Gateway 支持 https 访问 - 标准模式

> https://istio.io/latest/docs/reference/config/networking/gateway/#ServerTLSSettings
>> `credentialName`: The secret (of type `generic`) should contain the following keys and values: `key: <privateKey>` and `cert: <serverCert>`

## 创建证书 k8s secret

1. 在 **标准模式** 下， **必须使用** `key` 作为私钥文件名， `cert` 作为证书文件名。
2. 证书文件需要 **保持** 与 **istio-ingressgateway** 服务在 **相同** 的命名空间。


因此证书文件的创建命令如下

```bash
kubectl create secret generic wild-tangx-in \
    --from-file=key=./certificates/_.tangx.in.key  \
    --from-file=cert=./certificates/_.tangx.in.crt  \
    -n istio-system
```

其中

1. `wild-tangx-in`: 是 secret name。 之后 istio gateway 需要使用
2. `./certificates/_.tangx.in.key(crt)` 是证书私钥/文件所在的路径。

```bash
kg secret -n istio-system

    NAME                       TYPE                 DATA   AGE
    wild-tangx-in              Opaque               2      175m
```

## 创建支持 https 的 istio Gateway


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
    istio: ingressgateway
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
      - istio.tangx.in
    tls:
      httpsRedirect: true   # 开启 http -> https 301 重定向

  - port:
      number: 443
      name: https
      protocol: HTTPS   # 匹配协议
    hosts:
      - "*.tangx.in"    # 匹配域名， 这部分和 http 一样
    tls:
      mode: SIMPLE      # tls 模式
      credentialName: wild-tangx-in # 创建在 istio-system 下的证书 secret 
```

1. `.tls.httpRedirect`: 是否开启 http -> https 的 301 重定向。 
2. `.tls.mode`: tls 模式。 https 使用 `SIMPLE`  模式。  支持所有模式为 `PASSTHROUGH / SIMPLE / MUTUAL / AUTO_PASSTHROUGH / ISTIO_MUTUAL`。 
3. `.tls.credentialName`: 在 k8s 环境下， 证书使用的 secret name。 不用特意挂载到 istio-ingressgateway 服务中。


## 测试

通过请求可以看到

![17-standard-istio](/docs/imgs/17/17-standard-https.png)


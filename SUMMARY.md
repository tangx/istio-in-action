# 目录

## 环境准备

1. [安装 docker, k3s, istio 环境](./docs/01-install.md)
2. [初始化第一个项目 - prod](./docs/02-initial-project.md)
7. [升级项目 - prod and review](./docs/07-upgrade-project.md)
16. [使用 lego 创建 https 证书](./docs/16-lego-create-server-certificate.md)

## VirtualService

3. [istio VirtualService 和 k8s Ingress](./docs/03-vs-and-ingress.md)
4. [创建 Gateway 允许外部访问](./docs/04-gateway.md)
5. [VirtualService 给予 uri 重写路由](./docs/05-vs-http-rewrite-by-uri.md)
6. [使用 Subset 进行路由分组(版本控制)](./docs/06-subset.md)
8. [VirtualService 基于 Header 重写路由](./docs/08-vs-http-rewrite-by-header.md)
9. VirtualService 支持重写路由的所有方式
10. [VirtualService 路由重定向](./docs/10-vs-http-redirect.md)
11. [VirtualService 的重试机制](./docs/11-vs-http-retry.md)
12. [VirtualService 注入错误实现混沌测试](./docs/12-vs-http-fault-injection.md)
13. VirtualService 委托: 测试失败
14. [VirtualService Header 管理](./docs/14-vs-http-header-operation.md)
15. VirutalService 根据客户端 Label 转发路由(sourceLabels): 待测试

## Gateway

17. [Gateway 支持 https 访问 - 标准模式](./docs/17-gw-https-support-standard.md)


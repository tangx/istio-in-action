# VirtualService 混沌测试/错误注入



## 测试

```bash
time curl http://istio.tangx.in/prod/list

    {"data":{"Name":"istio in action","Price":300,"Reviews":{"1":{"id":"1","name":"zhangsan","commment":"istio 功能很强大， 就是配置太麻烦"},"2":{"id":"1","name":"wangwu","commment":"《istio in action》 真是一本了不起的书"}}},"version":"v1.1.0"}

    real    0m0.011s
    user    0m0.004s
    sys     0m0.005s


time curl http://istio.tangx.in/prod/list

    upstream request timeout

    real    0m3.014s
    user    0m0.004s
    sys     0m0.005s
```
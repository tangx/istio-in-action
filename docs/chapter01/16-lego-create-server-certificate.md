# 使用 lego 创建 https 证书

> https://go-acme.github.io/lego/dns/

```bash
#!/bin/bash
#

cd $(dirname $0)

source .env

lego  --email="${EMAIL}" \
      --key-type rsa2048 \
      --domains="${DOMAIN1}" \
      --path=$(pwd) --dns $DNS_PROVIDER --accept-tos run
```

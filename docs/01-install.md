# 安装 docker-k3s-istio 开发环境

## 1. 安装 docker


配置 docker 加速仓库

```json
{
  "registry-mirrors": [
    "https://mirror.ccs.tencentyun.com",
    "https://wlzfs4t4.mirror.aliyuncs.com"
  ],
  "bip": "169.253.32.1/24"
}
```

上述是腾讯云和阿里云的加速仓库， 根据需求自行调整。


## 2. 安装 k3s

### 2.1 安装 k3s

1. k3s 使用 `--docker` 模式是为了方便 docker build 产生的镜像可以直接用在 k3s 中。 否则在 k3s 和 docker 各自使用自己的 containerd runtime， 在程序发布的时候还需要再实现一个镜像 push 和 pull。 麻烦

2. 仅用 `--disable=traefik` 其一是为了保证集群的干净， 只有一个 ingress 控制器。 其二是 traefik 和 istio 默认都使用 LB 控制器， 会抢占 80/443 端口。 直接禁用，懒得再卸载。

```bash
curl -sfL http://rancher-mirror.cnrancher.com/k3s/k3s-install.sh | INSTALL_K3S_MIRROR=cn INSTALL_K3S_EXEC="server --docker --disable=traefik" sh -
```

如果之前已经安装过 k3s 的， 直接更在 systemd 启动文件后 

```bash
# /etc/systemd/system/k3s.service

ExecStart=/usr/local/bin/k3s \
    server \
	'--docker' \
	'--disable=traefik' \
```

使用如果下命令重载配置并重启 k3s

```bash
systemctl daemon-reload
systemctl restart k3s
```

### 2.2 复制 k3s config 文件

方便 kubectl 和之后的 istioctl 调用

```bash
mkdir -p ~/.kube
cp -a /etc/rancher/k3s/k3s.yaml ~/.kube/config
```


## 3. 安装 istio

### 3.1 安装 istioctl 

```bash
export ISTIO_VERSION=1.11.4
mkdir -p /data/istio-install && cd $_
wget -c https://github.com/istio/istio/releases/download/${ISTIO_VERSION}/istio-${ISTIO_VERSION}-linux-amd64.tar.gz && tar xf istio-${ISTIO_VERSION}-linux-amd64.tar.gz && mv istio-${ISTIO_VERSION} /usr/local/istio
export PATH=/usr/local/istio/bin/:$PATH
```

### 3.2 安装 istio 控制器

这里安装预设的 demo 文件

```bash
istioctl install --set profile=demo -y

    ✔ Istio core installed
    ✔ Istiod installed
    ✔ Egress gateways installed
    ✔ Ingress gateways installed
    ✔ Installation complete
    Thank you for installing Istio 1.11.  Please take a few minutes to tell us about your install/upgrade experience!  https://forms.gle/kWULBRjUv7hHci7T6
```


### 3.3 安装 kiali 和 prometheus

kiali 是 istio 的一个可视化 dashboard， 必须配合 prometheus 一起使用才能达到最佳效果。

幸运的是 istio 已经为我们准备好了所有东西。


```bash
# 之前已经将 istio 安装包移动到了 /usr/local/istio
export ISTIO_HOME=/usr/local/istio
ka -f ${ISTIO_HOME}/samples/addons/kiali.yaml
ka -f ${ISTIO_HOME}/samples/addons/prometheus.yaml
```

## 4. 安装快捷键

```bash
curl https://raw.githubusercontent.com/ahmetb/kubectx/master/kubens -o /usr/local/bin/kubens && chmod +x /usr/local/bin/kubens
```

方便使用命令行查看。

```bash
alias k=kubectl
alias ka='k apply'
alias kaf='kubectl apply -f'
alias kc='k create'
alias kca='_kca(){ kubectl "$@" --all-namespaces;  unset -f _kca; }; _kca'
alias kccc='kubectl config current-context'
alias kcdc='kubectl config delete-context'
alias kcgc='kubectl config get-contexts'
alias kcn='kubectl config set-context --current --namespace'
alias kcp='kubectl cp'
alias kcsc='kubectl config set-context'
alias kctx=kubectx
alias kcuc='kubectl config use-context'
alias kd='k describe'
alias kdcj='kubectl describe cronjob'
alias kdcm='kubectl describe configmap'
alias kdd='kubectl describe deployment'
alias kdds='kubectl describe daemonset'
alias kdel='kubectl delete'
alias kdelcj='kubectl delete cronjob'
alias kdelcm='kubectl delete configmap'
alias kdeld='kubectl delete deployment'
alias kdelds='kubectl delete daemonset'
alias kdelf='kubectl delete -f'
alias kdeli='kubectl delete ingress'
alias kdelno='kubectl delete node'
alias kdelns='kubectl delete namespace'
alias kdelp='kubectl delete pods'
alias kdelpvc='kubectl delete pvc'
alias kdels='kubectl delete svc'
alias kdelsa='kubectl delete sa'
alias kdelsec='kubectl delete secret'
alias kdelss='kubectl delete statefulset'
alias kdi='kubectl describe ingress'
alias kdno='kubectl describe node'
alias kdns='kubectl describe namespace'
alias kdp='kubectl describe pods'
alias kdpvc='kubectl describe pvc'
alias kds='kubectl describe svc'
alias kdsa='kubectl describe sa'
alias kdsec='kubectl describe secret'
alias kdss='kubectl describe statefulset'
alias kecj='kubectl edit cronjob'
alias kecm='kubectl edit configmap'
alias ked='kubectl edit deployment'
alias keds='kubectl edit daemonset'
alias kei='kubectl edit ingress'
alias keno='kubectl edit node'
alias kens='kubectl edit namespace'
alias kep='kubectl edit pods'
alias kepvc='kubectl edit pvc'
alias kes='kubectl edit svc'
alias kess='kubectl edit statefulset'
alias keti='kubectl exec -ti'
alias kg='k get'
alias kga='kubectl get all'
alias kgaa='kubectl get all --all-namespaces'
alias kgall='kg --all-namespaces'
alias kgcj='kubectl get cronjob'
alias kgcm='kubectl get configmaps'
alias kgcma='kubectl get configmaps --all-namespaces'
alias kgd='kubectl get deployment'
alias kgda='kubectl get deployment --all-namespaces'
alias kgds='kubectl get daemonset'
alias kgdsw='kgds --watch'
alias kgdw='kgd --watch'
alias kgdwide='kgd -o wide'
alias kgi='kubectl get ingress'
alias kgia='kubectl get ingress --all-namespaces'
alias kgno='kubectl get nodes'
alias kgns='kubectl get namespaces'
alias kgp='kubectl get pods'
alias kgpa='kubectl get pods --all-namespaces'
alias kgpall='kubectl get pods --all-namespaces -o wide'
alias kgpl='kgp -l'
alias kgpn='kgp -n'
alias kgpvc='kubectl get pvc'
alias kgpvca='kubectl get pvc --all-namespaces'
alias kgpvcw='kgpvc --watch'
alias kgpw='kgp -o wide'
alias kgpwide='kgp -o wide'
alias kgpy='kgp -o yaml'
alias kgrs='kubectl get rs'
alias kgs='kubectl get svc'
alias kgsa='kubectl get svc --all-namespaces'
alias kgsec='kubectl get secret'
alias kgseca='kubectl get secret --all-namespaces'
alias kgss='kubectl get statefulset'
alias kgssa='kubectl get statefulset --all-namespaces'
alias kgssw='kgss --watch'
alias kgsswide='kgss -o wide'
alias kgsw='kgs --watch'
alias kgswide='kgs -o wide'
alias khelp='cat /Users/tangxin/.zshrc.d/k8s.profile'
alias kk='k kustomize'
alias kl='kubectl logs'
alias kl1h='kubectl logs --since 1h'
alias kl1m='kubectl logs --since 1m'
alias kl1s='kubectl logs --since 1s'
alias klf='kubectl logs -f'
alias klf1h='kubectl logs --since 1h -f'
alias klf1m='kubectl logs --since 1m -f'
alias klf1s='kubectl logs --since 1s -f'
alias kns=kubens
alias kpf='kubectl port-forward'
alias krh='kubectl rollout history'
alias krm='k delete'
alias krsd='kubectl rollout status deployment'
alias krsss='kubectl rollout status statefulset'
alias kru='kubectl rollout undo'
alias ksd='kubectl scale deployment'
alias ksss='kubectl scale statefulset'
```
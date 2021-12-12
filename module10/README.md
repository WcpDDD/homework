### 项目改造
1. 在index.go的handle中添加一个随机的延时
2. 在handle的wrapper函数中添加指标采集代码

### 项目发布
1. docker打包，使用新版本tag，并发布到私有仓库
2. 修改module8中的deploy文件，将其修改为新的版本
3. 更新k8s中的deploy

### prometheus
1. helm安装
```
default      loki                                 ClusterIP   10.101.145.87    <none>        3100/TCP                     55m
default      loki-grafana                         NodePort    10.107.148.167   <none>        80:32160/TCP                 55m
default      loki-headless                        ClusterIP   None             <none>        3100/TCP                     55m
default      loki-kube-state-metrics              ClusterIP   10.104.205.196   <none>        8080/TCP                     55m
default      loki-prometheus-alertmanager         ClusterIP   10.109.174.237   <none>        80/TCP                       55m
default      loki-prometheus-node-exporter        ClusterIP   None             <none>        9100/TCP                     55m
default      loki-prometheus-pushgateway          ClusterIP   10.100.108.243   <none>        9091/TCP                     55m
default      loki-prometheus-server               ClusterIP   10.108.24.185    <none>        80/TCP                       55m
```
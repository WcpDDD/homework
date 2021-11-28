### 事前考量

![]()![作业](./作业.png)

####  http_server 项目改造

1. 读取外部配置文件，采取yaml格式
2. 监听外部配置文件变更
3. 配置文件中修改比如端口号的配置时，可能需要重启server，这里为了保证应用的访问，启动了一个新服务器，然后关闭旧的
4. 上条修改端口号本身没有什么意义，修改端口号本身就要修改service，容器应当是需要重新发布的，单纯的应用重启没有意义
5. glog模块进行日志分级 需要在启动时指定-v
4. 优雅终止
4. 优雅终止因为sh会忽略信号，这里修改dockerfile，引入tini
5. 将项目推送到github，并打上版本tag，用于项目本身的版本管理
5. 因为配置文件是通过configmap挂载为一个volume实现的，为了配置文件路径更灵活，应当支持通过命令行参数覆盖默认的配置文件路径

#### namespace创建

1. 创建一个namespace，做好resource quota限制

#### 日志处理

1. 日志打到stdout
1. 日志收集这里先不做 
1. 为了避免日志爆炸，去kubelet中进行containerLogMaxSize及containerLogMaxFiles的配置

####  配置文件分离
1. 创建一个用于保存配置文件的项目 并打上版本tag，主要用于配置文件的版本管理

2. 编写configmap的yaml

#### 应用发布失败的处理
1. 在应用容器打好tag，利用k8s自身的deployment的回滚功能

####  编写deployment
1. 当前项目为一个线上项目 所以有滚动发布的需求
2. 由于个人不是很擅长于评估项目的资源消耗，这里采用Burstable的QoS等级，同时在不同的worker上部署多个副本 保证服务的高可用
3. 同时定义一个优先级，保证线上业务的优先调度
3. 配置maxSurge及maxUnvailable保证应用升级时，失败及时停止滚动升级
4. 将configmap挂载为一个volumn
5. 这里发现一个问题 configmap在我的想法中是不一定挂载到应用所在目录的，所以这里对http_server做出改造，从环境变量中读取config文件路径
5. 定义readiness探针，保证应用启动成功后，容器才ready
#### service

1. 配合deployment中的label和service中的selector实现金丝雀

#### ingress



### 作业过程
```
alias k=kubectl
```
1. 完成k8s部署
1. 在kubelet中进行日志文件的配置
2. 完成代码改写
3. 改写dockerfile，引入tini
4. aliyun镜像加速创建镜像
5. 推送至阿里云镜像
6. 创建namespace 并配置resourcequota
```
k create namespace cncamp
k apply -f resource-quota.yaml
```
7. 创建一个github项目用于模拟配置文件的版本管理
8. 提交配置文件，并打上版本tag
9. 克隆指定tag的配置文件项目到本地 并创建configmap
```
git clone https://gitee.com/chenxinpei/homework-config.git
k create configmap http-server-config --from-file ./homework-config/http-server.yaml -n cncamp
```

10. 这种方式感觉违背了基准代码原则，遂对其进行修改，修改为k8s yaml的形式，并打上新的tag
```
git pull
k apply -f ./homework-config/http-server.yaml
```

11. 根据事前考量的内容进行deployment的编写
```
k apply -f priority.yaml
k apply -f deployment.yaml
```
12. 发现pod一直为notready，因为配置文件配置的端口号为8081，与探针不符，更新配置文件，并重新发布configmap
13. 发现pod中应用没有监听到配置文件变化，检查代码，发现没有监听指定的配置文件
14. 修改代码，打包为镜像，并修改tag，修改deployment，重新发布，顺带调整replicaset
15. 成功ready
16. 编写service spec配置
```
subsets:
  - addresses:
    - ip: 10.0.1.18
      nodeName: chenxinpei-worker-1
      targetRef:
        kind: Pod
        name: http-server-68ff4764db-sc22d
        namespace: cncamp
        resourceVersion: "1741461"
        uid: 6da38a0f-616f-4c4f-8238-d111be013a66
    - ip: 10.0.1.204
      nodeName: chenxinpei-worker-1
      targetRef:
        kind: Pod
        name: http-server-68ff4764db-dqbb2
        namespace: cncamp
        resourceVersion: "1741494"
        uid: c3802628-b346-48c1-a698-52ae1e2ae744
    - ip: 10.0.1.61
      nodeName: chenxinpei-worker-1
      targetRef:
        kind: Pod
        name: http-server-68ff4764db-dsbl7
        namespace: cncamp
        resourceVersion: "1741377"
        uid: a9cea676-61a9-441d-9fc3-22b216680e9f
    - ip: 10.0.1.88
      nodeName: chenxinpei-worker-1
      targetRef:
        kind: Pod
        name: http-server-68ff4764db-5q22h
        namespace: cncamp
        resourceVersion: "1741436"
        uid: b9e662b1-36c1-4205-8941-a6f8c44ba3fc
    ports:
    - port: 80
      protocol: TCP
```

17. 生成https证书
18. ingress安装，其中镜像拉不出来，本地科学上网拉下来后提交到阿里云的镜像加速，修改ingress部署脚本的镜像地址
19. 配置secret
20. 配置ingress访问service
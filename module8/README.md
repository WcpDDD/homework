### 事前考量

![]()![作业](./作业.png)

####  http_server 项目改造

1. 读取外部配置文件，采取yaml格式

2. 监听外部配置文件变更

3. glog模块进行日志分级

4. 优雅终止
5. 将项目推送到github，并打上版本tag，用于项目本身的版本管理

#### namespace创建

1. 创建一个namespace，做好resource quota限制

#### 日志处理
1. 为了避免日志爆炸，去kubelet中进行containerLogMaxSize及containerLogMaxFiles的配置

####  配置文件分离
1. 创建一个用于保存配置文件的项目 并打上版本tag，主要用于配置文件的版本管理
2. 编写configmap的yaml

####  编写deployment
1. 当前项目为一个线上项目 所以有滚动发布的需求
2. 由于个人不是很擅长于评估项目的资源消耗，这里采用Burstable的QoS等级，同时在不同的worker上部署多个副本 保证服务的高可用
3. 同时定义一个优先级，保证线上业务的优先调度
3. 配置maxSurge及maxUnvailable
4. 将configmap挂载为一个volumn
5. 这里发现一个问题 configmap在我的想法中是不一定挂载到应用所在目录的，所以这里对http_server做出改造，从环境变量中读取config文件路径
6. 日志直接输出到stdout，方便后续的集群日志管理
### service
### ingress
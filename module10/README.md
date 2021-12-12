### 项目改造
1. 在index.go的handle中添加一个随机的延时
2. 在handle的wrapper函数中添加指标采集代码

### 项目发布
1. docker打包，使用新版本tag，并发布到私有仓库
2. 修改module8中的deploy文件，将其修改为新的版本
3. 更新k8s中的deploy
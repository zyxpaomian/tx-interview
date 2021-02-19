# docker 管理 demo
___
 
 对docker资源进行管理，腾讯蓝鲸的某面试题：

> 设计并实现一个简单的存储资源分配系统实现一个服务，该服务可以接收用户端请求，为用户申请 MySQL 与 Redis 两类资源。分配给用户的资源实例必须是真实、可以连接使用的。用户可以通过接口查看分配给自己的资源配置信息。- 服务以 HTTP REST API 的方式提供接口，部分示例接口：- 申请一个新的 MySQL/Redis 资源实例- 查看某个实例的配置信息- MySQL、Redis 服务可以在服务端用 Docker 容器启动，也可以使用其他方式- 分配出的不同实例之间需要避免端口等资源冲突- 资源的连接、鉴权等信息应该随机生成，部分必须的信息- MySQL 连接地址、数据库名称、用户号、密码- Redis 连接地址、密码加分项：- 完整的项目架构图、项目安装、使用以及 README 文档- MySQL 与 Redis 实例支持不同的个性化配置，比如：- Redis 可以由用户设置数据最大占用空间- MySQL 可以由用户设置数据库字符集

### 安装
``` shell
go env -w GOPROXY=https://goproxy.cn,https://goproxy.io,direct
go env -w GO111MODULE=on

# 运行
make run 

# 编译
make compile
```
 
### 主要目录结构
* common 基本库，如日志，mysql 驱动, 字符转换，配置加载等
* controller 控制层，控制agent基本的逻辑代码位置，如查看所有image 创建容器等等
* dao dao层，数据层
* http mux+http封装的http 层，接口封装
* structs 统一的结构体位置

### 相关接口:
*查看镜像*
```
curl -X GET http://127.0.0.1:8080/v1/api/user/getallimage -H "Auth-Token: f6b76e8e-313f-4b9f-937f-35b7e28e1c71"
```

*创建docker容器*
containertype 指的是申请资源的类型, 目前仅支持mysql/redis
```shell
curl -X POST http://127.0.0.1:8080/v1/api/user/containercreate -H  "Content-Type: application/json" -d '{"containertype":"mysql"}'
```

###  待完成:
* 端口占用目前是临时使用map + lock的方式进行避免重复分配的，且每次重启服务会清零，后续需要找个地方存一下，或者直接通过etcd的分布式锁的方式避免重复分配/分配冲突
* redis 相关代码
* 用户自定义配置，可以将docker env 再单独拎出来传参


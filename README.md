# go-server
gonet 游戏服务器架构，mmo架构，包含数学库(box,matrix,point2d,point3d),[Recast Navigation寻路模块](https://blog.csdn.net/mango9126/article/details/79390543)，
a星寻路模块。

分布式雪花uuid,ai行为树，ai状态机，[excel导出配置](https://github.com/bobohume/gonet/tree/master/tool/data),raft同步模块，分片raft同步模块，hashring分布式一致性算法。

gonet核心思想是actor模式,消息驱动

服务器之间通过rpc，rpc无需注册，支持通用数据(int,[]int,[3]int),map数据,以及struct数据，[rpc性能](https://github.com/bobohume/gonet/blob/master/src/gonet/rpc/rpc_test.go)

sql封装简单的orm(orm支持pb结构体做mysql blob,orm支持结构体做mysql json类型)具体看[demo](https://github.com/bobohume/gonet/blob/master/src/gonet/db/db_test.go)

统一websocket和socket消息格式

客户端和网关之间通过protobuf + rpc，客户端tcp遵从如下消息包头

    前四位包体大小,再四位protobuf name 的 crc，中间protobuf字节流
    //另外支持特殊结束标志,前四位 protobuf name 的 crc，中间protobuf字节流， 尾部+结束标志💞♡ (结束标志也可以自己定义在base.TCP_END控制)（搜索tcp粘包特殊结束标志）


1.下载etcd做服发现

2.bin目录下的gonet_server.cfg配置数据库以及端口

3.数据库在sql文件目录下生产

4.win下执行build.bat,start.bat

5.linux下执行build.sh,start.sh

# pb协议生成

1.proto下载教程 https://blog.csdn.net/weixin_42117918/article/details/88920221

2.网关加入消息防火墙:在 ipacket.go 中 添加RegisterPacket(&message)

3.win下拷贝protoc.exe,protoc-gen-go.exe到项目bin目录,再执行proto.bat

4.linux下拷贝protoc.exe到项目bin目录,再执行proto.sh

5.生成后的pb文件在message目录对应的*.go


# 目前游戏库分类：

1.actor核心库，actor模式的雏形。

2.base基础库，分装rpc以及其他基础库。

3.db库，mysql，支持简单orm，没有重度gorm，更加轻便，还在受gorm 0 nil “” 数据库更新就失败的痛苦吗。还在忍受重度gorm带来sql语句都不知道怎么写，没错这个是轻度的。

4.rpc，服务器之间rpc通信。

5.nework库，网络库，tcp，websocket网络管理。rd库，redis库，做一些集群唯一缓存用。

6.raft 分布式同步

7.common 集群相关库



# 目前游戏模块：

1.account账号服务，提供注册账号，登录校验，集群服务。

2.natgate网关服务，对外连接，消息防火墙，对内消息转发，集群服务。

3.world世界服务，所有逻辑，集群服务。

4.login 登陆服，网关负载以及a，b切换

5.第三方中间件：etcd分布式服发现，redis分布式缓存。

# 交流

QQ群:950288306

# 服务器架构如下：
![image](https://github.com/bobohume/go-server/blob/master/框架.jpg)

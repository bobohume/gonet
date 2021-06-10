# go-server
gonet 游戏服务器架构，mmo架构，包含数学库(box,matrix,point2d,point3d),[Recast Navigation寻路模块](https://blog.csdn.net/mango9126/article/details/79390543)，
a星寻路模块。

分布式雪花uuid,ai行为树，ai状态机，[excel导出配置](https://github.com/bobohume/gonet/tree/master/tool/data),raft同步模块，分片raft同步模块，hashring分布式一致性算法。

gonet核心思想是actor模式,消息驱动,采用无锁队列mailbox替换channel阻塞模型

微服务，微服务之间使用分布式消息队列

[WIKI](https://github.com/bobohume/gonet/wiki)

# 交流

QQ群:950288306

# 服务器架构如下：
![image](框架.jpg)

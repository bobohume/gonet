# go-server
gonet 游戏服务器架构，mmo架构，包含数学库(box,matrix,point2d,point3d),[Recast Navigation寻路模块](https://blog.csdn.net/mango9126/article/details/79390543)，
a星寻路模块。

分布式雪花uuid,ai行为树，ai状态机，[excel导出配置](https://github.com/bobohume/gonet/tree/master/tool/data),raft同步模块，分片raft同步模块，hashring分布式一致性算法。

gonet核心思想是actor模式,消息驱动,采用mpsc替换channel.

channel在队列满了,会阻塞produce,mpsc类似mailbox.

分布式一致性采用lease一致性,每个player actor维护自己的lease.
~~~~
基准测试如下: i7 10700 2.9GHZ 16核 执行10万次生产和消费
BenchmarkChanPushPop/100000_1-16             252           4715405 ns/op
BenchmarkChanPushPop/100000_2-16             174           6881893 ns/op
BenchmarkChanPushPop/100000_4-16             180           6635532 ns/op
BenchmarkChanPushPop/100000_8-16             142           8440178 ns/op
BenchmarkChanPushPop/100000_16-16            129           9312717 ns/op
BenchmarkPushPopActor/100000_1-16            252           4748732 ns/op
BenchmarkPushPopActor/100000_2-16            205           5818782 ns/op
BenchmarkPushPopActor/100000_4-16            278           4348560 ns/op
BenchmarkPushPopActor/100000_8-16            297           3986121 ns/op
BenchmarkPushPopActor/100000_16-16           304           3839012 ns/op
~~~~

另外go的定时器是大小堆,对高精度10毫秒定时器会吃掉大部分cpu,这里采用5级时间轮定时器优化到O(1)

微服务，微服务之间使用分布式消息队列

[WIKI](https://github.com/bobohume/gonet/wiki)

# 交流

QQ群:950288306

# 服务器架构如下：
![image](框架.jpg)

# gm stub
![image](gm_stub.jpg)

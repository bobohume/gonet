# go-server
gonet 游戏服务器架构，mmo架构，分布式snowflake64为整形uuid,ai行为树，配置data，游戏大部分都在内存运算,分布式缓存redis,增加db模块读取blob数据。

设计之初，建立在actor模式下的；rpc，以及消息驱动，rpc无需注册，支持int，数据，struct（struct必须要注册结构题即可）；sql封装简单的orm，具体看demo

websocket模式下，要在net，websocket注视掉如下代码：https://studygolang.com/articles/14842

代码除了mysql，protobuf，redis这几个第三方库以外，其他都是自己写的，方便性能和修改，主动权在自己手里

服务器之间rpc，客户端服务器之间protobuf + rpc，客户端tcp遵从如下消息包头(支持json，考虑到性能，两种传输协议不兼容，请用项目中的json3替换到message，写同名的json替换pb的结构体，json继承messagebase类即可,再搜索//解析json)

    前四位 protobuf name 的 crc，中间protobuf字节流， 尾部+结束标志💞♡ (结束标志也可以自己定义在base.TCP_END控制)
    //另外支持包头大小- 前四位包体大小,再四位protobuf name 的 crc，中间protobuf字节流,代码注视掉,（搜索tcp粘包固定包头）

1.支持go mod, gopath可以不需要设置。（也支持go vendor（删除项目下的go.mod文件），下载三个基础库，mysql，protobuf，redis）

// go get github.com/golang/net

// go get github.com/go-sql-driver/mysql

// go get github.com/gomodule/redigo/redis

2.bin目录下的sxz_server.cfg配置数据库以及端口

3.数据库在sql文件目录下生产

4.win下执行build.bat,start.bat

5.linux下执行build.sh,start.sh

有问题可以加qq群：950288306


服务器架构如下：
![image](https://github.com/bobohume/go-server/blob/master/框架.jpg)

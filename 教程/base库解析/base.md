base 基础类库

    bitstream字节流库 rpc基于这个字节流库实现
        支持int,float,string,bool,float64,int64字节流。默认小端字节序,你可能看到很多<<3不要惊讶，做字节处理
        
    config配置文件库
        读取服务中的配置文件cfg文件,#为注释
        
    datafile文件
        data文件,游戏用的策划配置文件,策划用excel生成文件,参考地下tool目录下的源码,将excel生成二进制文件流
               
    gob文件
        gob作为go内部解析协议,在大数据的情况下比pb性能更加,但是gob生成加密耗费挺大,可以用缓存,gob解析数据很快

    log日志文件
        服务里面在bin目录下log有各个服务器的log文件,根据时间,每天生成不同的log和err文件
        
     memorycheck
        pprof用户内存检测或者grotinue,cpu性能监测
        
     random
        随机数,一个更好的随机数,RandI(lower,upper)
        
     rpc
        引擎里面的rpc,不需要知道类型的魔力全靠它,底层字节流
     
     UUID
        著名的snowflake,分布式64位整型uuid,妈妈再也不用担心,mysql的自增量了,里面的workerid基于etcd实现分布式workerid分配
        
      Vector
        一个高性能的vector,双端队列
        
        
[data工具地址](https://github.com/bobohume/gonet/tree/master/src/gonet/tool/data)
        
        

       

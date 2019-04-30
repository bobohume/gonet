base 基础类库

    bitstream字节流库 rpc基于这个字节流库实现
        支持int,float,string,bool,float64,int64字节流。默认小端字节序,你可能看到很多<<3不要惊讶，做字节处理
        
    config配置文件库
        读取服务中的配置文件cfg文件,#为注释
        
    datafile文件
        data文件,游戏用的策划配置文件,策划用excel生成文件,参考地下tool目录下的源码,将excel生成二进制文件流
        
[data工具地址](https://github.com/bobohume/gonet/tree/master/src/gonet/tool/data)

    log日志文件
        服务里面在bin目录下log有各个服务器的log文件,根据时间,每天生成不同的log和err文件
       
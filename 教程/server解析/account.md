account server
     顾名思义,账号服务器,负责账号登录,创建,验证,不同于一般服务器,账号可能对外。
     在gonet里面对外只有网关,网关转发给具体服务
     
     accountserver
        main
        
      evnetprocess
        包处理函数,tcp到包处理,具体看m_PacketFuncList,以及actor的 PacketFunc回调函数
        
      accountmgr
        账号管理类
        
       serversocketmgr
        记录区分分布式各个服务器的链接管理
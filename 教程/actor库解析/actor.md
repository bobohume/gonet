actor,基类对象,所有包处理或者逻辑处理都基于这个actor

    actor已事件驱动模式,消息队列,定时器,生命结束信号,都在一个goroutine
    
    actor并没有这么复杂,看下代码
    
    func (this *Actor) loop() bool{
    
    	defer func() {
    	
    		if err := recover(); err != nil{
    		
    			base.TraceCode(err)
    			
    		}
    		
    	}()
    	
    	select {//io阻塞
    	
    	case io := <-this.m_CallChan://消息队列，rpc回调
    	
    		this.call(io)
    		
    	case msg := <-this.m_AcotrChan ://生命周期控制信息
    	
    		if msg == DESDORY_EVENT{
    		
    			return true
    			
    		}
    		
    	case <- this.m_pTimer.C://定时器
    	
    		if this.m_TimerCall != nil{
    		
    			this.m_TimerCall()
    			
    		}
    		
    	}
    	
    	return false
    	
    }
    
    同在一个goroutine,不存在多个goroutine访问共享变量的问题，即你可以简单的看成actor模式下，变量只对actor可见，其他actor
    想要访问他,rpc消息发过来。
    
    
actor消息队列

    actor.SendMsg(funcName string, params ...interface{})//actor靠它来传输消息
    
    funcName为回调函数名字
    
    params可以是基础类型:int,uint32,*int,*uint32,[]int,[]uint32,[6]int,[6]uint32,支持结构体不需要注册
     
        
actor消息队列回调

    actor.RegisterCall(funcName string, call interface{})//actor靠他来进行消息回调
    
    funcName为回调函数名字
    
    call为SendMsg传过来的params
    
    
actor列子

    SendMsg("COMMON_RegisterRequest",ServerType, Ip, Port)//发送给actor一个消息队列
    
    回调actor,没错,正如你看到的,只要在初始化调用回调函数，并且 nType int, Ip string, Port int 参数和sendmsg
    
    发送的参数一致就可以了。和一般消息回调不同,你要注册消息回调,在写回调函数,而且一般的rpc都是参数类型要确定,参数长度也要一致。
    
    在这里,没有这些拘泥
    
    func (this *ServerSocketManager) Init(num int){
    
    	this.Actor.Init(num)
    	
    	this.m_GateMap 		= make(HashSocketMap)
    	
    	this.m_SocketMap 	= make(HashSocketMap)
    	
    	this.m_Locker		= &sync.RWMutex{}
    	
    	this.RegisterCall("COMMON_RegisterRequest", func(nType int, Ip string, Port int) {
    	
    		pServerInfo := new(common.ServerInfo)
    		
    		pServerInfo.SocketId = this.GetSocketId()
    		
    		pServerInfo.Type = nType
    		
    		pServerInfo.Ip = Ip
    		
    		pServerInfo.Port = Port
    		
    		this.AddServerMap(pServerInfo)
    		
    		switch pServerInfo.Type {
    		
    		case int(message.SERVICE_GATESERVER):
    		
    			SERVER.GetServer().SendMsgByID(this.GetSocketId(), "COMMON_RegisterResponse")
    		
    		}
    	
    	})
    	
    	this.Actor.Start()
   
    }

    
actor消息队列为什么能知道,我传过来什么,解析什么呢？

    这个得力于rpc模块，rpc是什么？其实就是字节流。
    
    rpc模块请看/gonet/base/rpc.go
    
    func GetPacket(funcName string, params ...interface{})[]byte {
    	
    	defer func() {
    		
    		if err := recover(); err != nil {
    			
    			fmt.Println("GetPacket", err)
    		
    		}
    	
    	}()
    	
    	msg := make([]byte, 1024)
    	
    	bitstream := NewBitStream(msg, 1024)
    	
    	bitstream.WriteString(funcName)
    	
    	bitstream.WriteInt(len(params), 8)
    	
    	for _, param := range params {
    		
    		sType := GetTypeString(param)
    		
    		switch sType {
    		
    		case "bool":
    			
    			bitstream.WriteInt(1, 8)
    			
    			bitstream.WriteFlag(param.(bool))
    		
    		case "float64":
    			
    			bitstream.WriteInt(2, 8)
    			
    			bitstream.WriteFloat64(param.(float64))
    		
    那么问题来了,rpc字节流性能怎么样呢？
    
        请单元测试下 [rpc的性能代码](https://github.com/bobohume/gonet/blob/master/src/gonet/test/client_test.go)
       
        测试性能如下:测试100万次压一个长度为25的数组
        
        go test -v client_test.go
        
        === RUN   TestJson          //json加密
       
        --- PASS: TestJson (1.41s)
       
        === RUN   TestUJson         //json解密
       
        --- PASS: TestUJson (8.11s)
        
        === RUN   TestPB            //pb加密
        
        --- PASS: TestPB (0.41s)
        
        === RUN   TestUPB           //pb解密
        
        --- PASS: TestUPB (0.58s)
        
        === RUN   TestRpc           //rpc加密
        
        --- PASS: TestRpc (0.62s)
        
        === RUN   TestURpc          //rpc解密
        
        --- PASS: TestURpc (0.53s)


actor我怎么知道我要往那个actor发送呢

    对于actor绑定network,在network里面分析。
    
    那actor对actor呢 ActorMgr.go里面有,通过AddActor类型名字注册actor,这个适用于全局的actor,考虑性能
    
    没用锁模式,直接map多线程访问,主要还是全局的
    
    那么不是全局的怎么办,这种需求少,目前只有玩家才有。比如玩家管理统筹所有玩家actor,登录和离开过playermgr,
    
    其他消息通过玩家id分发到玩家actor

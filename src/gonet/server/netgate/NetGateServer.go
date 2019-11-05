 package netgate

 import (
	 "github.com/golang/protobuf/proto"
	 "gonet/base"
	 "gonet/message"
	 "gonet/network"
	 "gonet/server/common/cluster"
	 "time"
 )

type(
	ServerMgr struct{
		m_pService	*network.ServerSocket
		m_Inited bool
		m_config base.Config
		m_Log	base.CLog
		m_TimeTraceTimer *time.Ticker
		m_PlayerMgr *PlayerManager
		m_WorldCluster *cluster.Cluster
		m_AccountCluster *cluster.Cluster
		m_Cluster *cluster.Service
	}

	IServerMgr interface{
		Init() bool
		GetLog() *base.CLog
		GetServer() *network.ServerSocket
		GetWorldCluster() *cluster.Cluster
		GetAccountCluster() *cluster.Cluster
		GetPlayerMgr() *PlayerManager
		InitWorldSocket()
		AddWorldSocket(string, int, int)
		DelWorldSocket(int)
		OnServerStart()
	}

	BitStream base.BitStream
)

var(
	UserNetIP string
	UserNetPort string
	AccountServerIp string
	AccountServerPort string
	EtcdEndpoints []string

	SERVER ServerMgr
)

func (this *ServerMgr) GetLog() *base.CLog{
 	return &this.m_Log
}

func (this *ServerMgr) GetServer() *network.ServerSocket{
 	return this.m_pService
}

func (this *ServerMgr) GetWorldCluster() *cluster.Cluster{
	return this.m_WorldCluster
}

func (this *ServerMgr) GetAccountCluster() *cluster.Cluster{
 	return this.m_AccountCluster
}

func (this *ServerMgr) GetPlayerMgr() *PlayerManager{
	return this.m_PlayerMgr
}

func (this *ServerMgr)Init() bool{
	if(this.m_Inited){
		return true
	}

	//初始化log文件
	this.m_Log.Init("netgate")
	//初始ini配置文件
	this.m_config.Read("GONET_SERVER.CFG")

	EtcdEndpoints = this.m_config.Get5("Etcd_Cluster", ",")
	UserNetIP, UserNetPort 	= this.m_config.Get2("NetGate_WANAddress", ":")
	AccountServerIp, AccountServerPort 	= this.m_config.Get2("Account_LANAddress", ":")
	ShowMessage := func(){
		this.m_Log.Println("**********************************************************")
		this.m_Log.Printf("\tNetGateServer Version:\t%s",base.BUILD_NO)
		this.m_Log.Printf("\tNetGateServerIP(LAN):\t%s:%s", UserNetIP, UserNetPort)
		this.m_Log.Printf("\tAccountServerIP(LAN):\t%s:%s", AccountServerIp, AccountServerPort)
		this.m_Log.Println("**********************************************************");
	}
	ShowMessage()

	//初始化socket
	this.m_pService = new(network.ServerSocket)
	port := base.Int(UserNetPort)
	this.m_pService.Init(UserNetIP, port)
	this.m_pService.SetConnectType(network.CLIENT_CONNECT)
	//this.m_pService.Start()
	packet := new(UserPrcoess)
	packet.Init(1000)
	this.m_pService.BindPacketFunc(packet.PacketFunc)
	this.m_pService.Start()

	//websocket
	/*this.m_pService = new(network.WebSocket)
	port,_:=strconv.Atoi(UserNetPort)
	this.m_pService.Init(UserNetIP, port)
	this.m_pService.SetConnectType(network.CLIENT_CONNECT)
	//this.m_pService.Start()
	packet := new(UserPrcoess)
	packet.Init(1000)
	this.m_pService.BindPacketFunc(packet.PacketFunc)
	this.m_pService.Start()*/
	//注册到集群服务器
	this.m_Cluster = cluster.NewService(int(message.SERVICE_GATESERVER), UserNetIP, base.Int(UserNetPort), EtcdEndpoints)

	//世界服务器集群
	this.m_WorldCluster = new(cluster.Cluster)
	this.m_WorldCluster.Init(1000, int(message.SERVICE_WORLDSERVER), UserNetIP, base.Int(UserNetPort), EtcdEndpoints)
	this.m_WorldCluster.BindPacket(&WorldProcess{})
	this.m_WorldCluster.BindPacketFunc(DispatchPacketToClient)

	//账号服务器集群
	this.m_AccountCluster = new(cluster.Cluster)
	this.m_AccountCluster.Init(1000, int(message.SERVICE_ACCOUNTSERVER), UserNetIP, base.Int(UserNetPort), EtcdEndpoints)
	this.m_AccountCluster.BindPacket(&AccountProcess{})

	//初始玩家管理
	this.m_PlayerMgr = new(PlayerManager)
	this.m_PlayerMgr.Init(1000)
	return  false
}

func (this *ServerMgr) OnServerStart(){
	this.m_pService.Start()
}

func SendToClient(socketId int, packet proto.Message){
	buff, err := proto.Marshal(packet)
	if err == nil {
		SERVER.GetServer().SendById(socketId, buff)
	}
}
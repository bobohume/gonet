 package netgate

 import (
	 "gonet/base"
	 "gonet/common"
	 "gonet/common/cluster"
	 "gonet/network"
	 "gonet/rpc"
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
		m_pCluster *cluster.Cluster
	}

	IServerMgr interface{
		Init() bool
		GetLog() *base.CLog
		GetServer() *network.ServerSocket
		GetCluster () *cluster.Service
		GetPlayerMgr() *PlayerManager
		OnServerStart()
	}
)

var(
	UserNetIP string
	UserNetPort string
	EtcdEndpoints []string
	Nats_Cluster string
	SERVER ServerMgr
)

func (this *ServerMgr) GetLog() *base.CLog{
 	return &this.m_Log
}

func (this *ServerMgr) GetServer() *network.ServerSocket{
 	return this.m_pService
}

func (this *ServerMgr) GetCluster () *cluster.Cluster {
 	return this.m_pCluster
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
	Nats_Cluster = this.m_config.Get("Nats_Cluster")
	ShowMessage := func(){
		this.m_Log.Println("**********************************************************")
		this.m_Log.Printf("\tNetGateServer Version:\t%s",base.BUILD_NO)
		this.m_Log.Printf("\tNetGateServerIP(LAN):\t%s:%s", UserNetIP, UserNetPort)
		this.m_Log.Println("**********************************************************");
	}
	ShowMessage()

	//初始化socket
	this.m_pService = new(network.ServerSocket)
	port := base.Int(UserNetPort)
	this.m_pService.Init(UserNetIP, port)
	this.m_pService.SetMaxReceiveBufferSize(base.MAX_CLIENT_PACKET)
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

	var packet1 EventProcess
	packet1.Init(1000)
	this.m_pCluster = new (cluster.Cluster)
	this.m_pCluster.Init(1000, &common.ClusterInfo{Type: rpc.SERVICE_GATESERVER, Ip:UserNetIP, Port:int32(base.Int(UserNetPort))}, EtcdEndpoints, Nats_Cluster)
	this.m_pCluster.BindPacketFunc(packet1.PacketFunc)
	this.m_pCluster.BindPacketFunc(DispatchPacket)

	//初始玩家管理
	this.m_PlayerMgr = new(PlayerManager)
	this.m_PlayerMgr.Init(1000)
	return  false
}

func (this *ServerMgr) OnServerStart(){
	this.m_pService.Start()
}

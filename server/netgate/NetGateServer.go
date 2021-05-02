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
		m_WorldCluster *cluster.Cluster
		m_AccountCluster *cluster.Cluster
		m_ZoneCluster *cluster.Cluster
		m_Cluster *cluster.Service
	}

	IServerMgr interface{
		Init() bool
		GetLog() *base.CLog
		GetServer() *network.ServerSocket
		GetCluster () *cluster.Service
		GetWorldCluster() *cluster.Cluster
		GetAccountCluster() *cluster.Cluster
		GetZoneCluster() *cluster.Cluster
		GetPlayerMgr() *PlayerManager
		OnServerStart()
	}
)

var(
	UserNetIP string
	UserNetPort string
	EtcdEndpoints []string

	SERVER ServerMgr
)

func (this *ServerMgr) GetLog() *base.CLog{
 	return &this.m_Log
}

func (this *ServerMgr) GetServer() *network.ServerSocket{
 	return this.m_pService
}

func (this *ServerMgr) GetCluster () *cluster.Service {
 	return this.m_Cluster
}

func (this *ServerMgr) GetWorldCluster() *cluster.Cluster {
	return this.m_WorldCluster
}

func (this *ServerMgr) GetAccountCluster() *cluster.Cluster {
 	return this.m_AccountCluster
}

func (this *ServerMgr) GetZoneCluster() *cluster.Cluster {
	return this.m_ZoneCluster
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
	this.m_Cluster = cluster.NewService( &common.ClusterInfo{Type: rpc.SERVICE_GATESERVER, Ip:UserNetIP, Port:int32(base.Int(UserNetPort))}, EtcdEndpoints)

	//世界服务器集群
	this.m_WorldCluster = new(cluster.Cluster)
	this.m_WorldCluster.Init(1000, &common.ClusterInfo{Type: rpc.SERVICE_WORLDSERVER}, EtcdEndpoints)
	this.m_WorldCluster.BindPacket(&WorldProcess{})
	this.m_WorldCluster.BindPacketFunc(DispatchPacket)

	//账号服务器集群
	this.m_AccountCluster = new(cluster.Cluster)
	this.m_AccountCluster.Init(1000,  &common.ClusterInfo{Type: rpc.SERVICE_ACCOUNTSERVER}, EtcdEndpoints)
	this.m_AccountCluster.BindPacket(&AccountProcess{})
	this.m_AccountCluster.BindPacketFunc(DispatchPacket)

	//战斗服务器集群
	this.m_ZoneCluster = new(cluster.Cluster)
	this.m_ZoneCluster.Init(1000,  &common.ClusterInfo{Type: rpc.SERVICE_ZONESERVER}, EtcdEndpoints)
	this.m_ZoneCluster.BindPacket(&ZoneProcess{})
	this.m_ZoneCluster.BindPacketFunc(DispatchPacket)

	//初始玩家管理
	this.m_PlayerMgr = new(PlayerManager)
	this.m_PlayerMgr.Init(1000)
	return  false
}

func (this *ServerMgr) OnServerStart(){
	this.m_pService.Start()
}

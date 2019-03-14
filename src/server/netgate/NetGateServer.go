 package netgate

 import (
	 "gonet/base"
	 "github.com/golang/protobuf/proto"
	 "gonet/message"
	 "gonet/network"
	 "gonet/server/common"
	 "time"
 )

type(
	ServerMgr struct{
		m_pService	*network.ServerSocket
		m_pAccountClient *network.ClientSocket
		m_pMonitorClient *common.MonitorClient
		m_Inited bool
		m_config base.Config
		m_Log	base.CLog
		m_TimeTraceTimer *time.Ticker
		m_PlayerMgr *PlayerManager
		m_WorldSocketMgr *common.DispatchMgr
	}

	IServerMgr interface{
		Init() bool
		GetLog() *base.CLog
		GetServer() *network.ServerSocket
		GetWorldSocketMgr() *common.DispatchMgr
		GetAccountSocket() *network.ClientSocket
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

	SERVER ServerMgr
)

 func (this *ServerMgr) GetLog() *base.CLog{
	 return &this.m_Log
 }

 func (this *ServerMgr) GetServer() *network.ServerSocket{
	 return this.m_pService
 }

 func (this *ServerMgr) GetWorldSocketMgr() *common.DispatchMgr{
 	return this.m_WorldSocketMgr
 }

 func (this *ServerMgr) GetPlayerMgr() *PlayerManager{
 	return this.m_PlayerMgr
 }

 func (this *ServerMgr) GetAccountSocket() *network.ClientSocket{
 	return this.m_pAccountClient
 }

func (this *ServerMgr)Init() bool{
	if(this.m_Inited){
		return true
	}

	//初始化log文件
	this.m_Log.Init("netgate")
	//初始ini配置文件
	this.m_config.Read("SXZ_SERVER.CFG")

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

	//链接monitor
	this.m_pMonitorClient = new(common.MonitorClient)
	monitorIp, monitroPort := this.m_config.Get2("Monitor_LANAddress", ":")
	this.m_pMonitorClient.Connect(int(message.SERVICE_GATESERVER), monitorIp, monitroPort, UserNetIP, UserNetPort)

	//初始化socket
	this.m_pService = new(network.ServerSocket)
	port := base.Int(UserNetPort)
	this.m_pService.Init(UserNetIP, port)
	this.m_pService.SetConnectType(network.CLIENT_CONNECT)
	//this.m_pService.Start()
	packet := new(UserPrcoess)
	packet.Init(1000)
	packet1 := new(UserServerProcess)
	packet1.Init(1000)
	this.m_pService.BindPacketFunc(packet.PacketFunc)
	this.m_pService.BindPacketFunc(packet1.PacketFunc)
	this.m_pService.Start()

	//websocket
	/*this.m_pService = new(network.WebSocket)
	port,_:=strconv.Atoi(UserNetPort)
	this.m_pService.Init(UserNetIP, port)
	this.m_pService.SetConnectType(network.CLIENT_CONNECT)
	//this.m_pService.Start()
	packet := new(UserPrcoess)
	packet.Init(1000)
	packet1 := new(UserServerProcess)
	packet1.Init(1000)
	this.m_pService.BindPacketFunc(packet.PacketFunc)
	this.m_pService.BindPacketFunc(packet1.PacketFunc)
	this.m_pService.Start()*/

	this.m_WorldSocketMgr = new(common.DispatchMgr)
	this.m_WorldSocketMgr.Init(1000)
	this.m_WorldSocketMgr.BindPacket(&WorldProcess{})
	this.m_WorldSocketMgr.BindPacketFunc(DispatchPacketToClient)

	//连接account
	this.m_pAccountClient = new(network.ClientSocket)
	port = base.Int(AccountServerPort)
	this.m_pAccountClient.Init(AccountServerIp,port)
	packet3 := new(AccountProcess)
	packet3.Init(1000)
	this.m_pAccountClient.BindPacketFunc(packet3.PacketFunc)
	this.m_pAccountClient.BindPacketFunc(this.m_WorldSocketMgr.PacketFunc)
	this.m_pAccountClient.Start()


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
		SERVER.GetServer().SendByID(socketId, buff)
	}
}
 package netgate

 import (
	 "base"
	 "github.com/golang/protobuf/proto"
	 "network"
	 "strconv"
	 "time"
 )

type(
	ServerMgr struct{
		m_pService	*network.ServerSocket
		m_pWorldClient *network.ClientSocket
		m_pAccountClient *network.ClientSocket
		m_pMonitorClient *network.ClientSocket
		m_Inited bool
		m_config base.Config
		m_Log	base.CLog
		m_TimeTraceTimer *time.Ticker
		m_PlayerMgr *PlayerManager
	}

	IServerMgr interface{
		Init() bool
		GetLog() *base.CLog
		GetServer() *network.ServerSocket
		GetWorldSocket() *network.ClientSocket
		GetAccountSocket() *network.ClientSocket
		GetMonitorSocket() *network.ClientSocket
		GetPlayerMgr() *PlayerManager
		OnServerStart()
	}

	BitStream base.BitStream
)

var(
	UserNetIP string
	UserNetPort string
	WorldServerIP string
	WorldServerPort string
	AccountServerIp string
	AccountServerPort string
	NetGateId string

	SERVER ServerMgr
)

 func (this *ServerMgr) GetLog() *base.CLog{
	 return &this.m_Log
 }

 func (this *ServerMgr) GetServer() *network.ServerSocket{
	 return this.m_pService
 }

 func (this *ServerMgr) GetPlayerMgr() *PlayerManager{
 	return this.m_PlayerMgr
 }

 func (this *ServerMgr) GetWorldSocket() *network.ClientSocket{
 	return this.m_pWorldClient
 }

 func (this *ServerMgr) GetAccountSocket() *network.ClientSocket {
 	return this.m_pAccountClient
 }

 func (this *ServerMgr) GetMonitorSocket() *network.ClientSocket{
 	return this.m_pMonitorClient
 }

func (this *ServerMgr)Init() bool{
	if(this.m_Inited){
		return true
	}

	//初始化log文件
	this.m_Log.Init("netgate")
	//初始ini配置文件
	this.m_config.Read("SXZ_SERVER.CFG")

	NetGateId  = this.m_config.Get("GateID")
	UserNetIP, UserNetPort 	= this.m_config.Get2("NetGate_WANAddress", ":")
	AccountServerIp, AccountServerPort 	= this.m_config.Get2("Account_WANAddress", ":")
	WorldServerIP, WorldServerPort 	= this.m_config.Get2("World_LANAddress", ":")
	ShowMessage := func(){
		this.m_Log.Println("**********************************************************")
		this.m_Log.Printf("\tNetGateServer Version:\t%s",base.BUILD_NO)
		this.m_Log.Printf("\tNetGateServerIP(LAN):\t%s:%s", UserNetIP, UserNetPort)
		this.m_Log.Printf("\tWorldServerIP(LAN):\t%s:%s", WorldServerIP, WorldServerPort)
		this.m_Log.Printf("\tAccountServerIP(LAN):\t%s:%s", AccountServerIp, AccountServerPort)
		this.m_Log.Println("**********************************************************");
	}
	ShowMessage()

	//初始化socket
	this.m_pService = new(network.ServerSocket)
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
	this.m_pService.BindPacketFunc(packet1.PacketFunc)
	this.m_pService.BindPacketFunc(packet.PacketFunc)
	this.m_pService.Start()*/

	//连接world
	this.m_pWorldClient = new(network.ClientSocket)
	port,_ = strconv.Atoi(WorldServerPort)
	this.m_pWorldClient.Init(WorldServerIP,port)
	packet2 := new(WorldProcess)
	packet2.Init(1000)
	this.m_pWorldClient.BindPacketFunc(packet2.PacketFunc)
	this.m_pWorldClient.BindPacketFunc(DispatchPacketToClient)
	//this.m_pWorldClient.Start()

	//连接account
	this.m_pAccountClient = new(network.ClientSocket)
	port,_ = strconv.Atoi(AccountServerPort)
	this.m_pAccountClient.Init(AccountServerIp,port)
	packet3 := new(AccountProcess)
	packet3.Init(1000)
	this.m_pAccountClient.BindPacketFunc(packet3.PacketFunc)
	//this.m_pAccountClient.Start()


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

func SendToWorld(msg string, params ...interface{}){
	SERVER.GetWorldSocket().SendMsg(msg, params...)
}

 func SendToAccount(msg string, params ...interface{}){
	 SERVER.GetAccountSocket().SendMsg(msg, params...)
 }
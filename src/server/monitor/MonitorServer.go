package monitor

import (
	"gonet/base"
	"github.com/golang/protobuf/proto"
	"gonet/message"
	"gonet/network"
)

type(
	ServerMgr struct{
		m_pService	*network.ServerSocket
		m_pServerMgr *ServerSocketManager
		m_Inited bool
		m_config base.Config
		m_Log	base.CLog
	}

	IServerMgr interface{
		Init() bool
		GetLog() *base.CLog
		GetServer() *network.ServerSocket
		GetServerMgr() *ServerSocketManager
	}

	BitStream base.BitStream
)

var(
	UserNetIP string
	UserNetPort string

	SERVER ServerMgr
)

func (this *ServerMgr)Init() bool{
	if(this.m_Inited){
		return true
	}

	//初始化log文件
	this.m_Log.Init("monitor")
	//初始ini配置文件
	this.m_config.Read("SXZ_SERVER.CFG")
	UserNetIP, UserNetPort = this.m_config.Get2("Monitor_LANAddress", ":")

	ShowMessage := func(){
		this.m_Log.Println("**********************************************************")
		this.m_Log.Printf("\tMonitorServer Version:\t%s",base.BUILD_NO)
		this.m_Log.Printf("\tMonitorServerIP(LAN):\t%s:%s", UserNetIP, UserNetPort)
		this.m_Log.Println("**********************************************************");
	}
	ShowMessage()

	//初始化socket
	this.m_pService = new(network.ServerSocket)
	port := base.Int(UserNetPort)
	this.m_pService.Init(UserNetIP, port)
	this.m_pService.SetMaxReceiveBufferSize(1024)
	this.m_pService.SetMaxSendBufferSize(1024)
	this.m_pService.Start()
	var packet EventProcess
	packet.Init(1000)
	this.m_pService.BindPacketFunc(packet.PacketFunc)

	this.m_pServerMgr = new(ServerSocketManager)
	this.m_pServerMgr.Init(1000)

	return  false
}

func (this *ServerMgr) GetLog() *base.CLog{
	return &this.m_Log
}

func (this *ServerMgr) GetServer() *network.ServerSocket{
	return this.m_pService
}

func (this *ServerMgr) GetServerMgr() *ServerSocketManager{
	return this.m_pServerMgr
}

func SendToClient(socketId int, packet proto.Message){
	bitstream := base.NewBitStream(make([]byte, 1024), 1024)
	if !message.GetProtoBufPacket(packet, bitstream) {
		return
	}
	SERVER.GetServer().SendByID(socketId, bitstream.GetBuffer())
}


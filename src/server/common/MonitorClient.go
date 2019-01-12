package common

import (
	"actor"
	"base"
	"network"
)

type (
	MonitorClient struct {
		actor.Actor
		m_LostTimer *SimpleTimer
		m_pClient *network.ClientSocket
		m_Ip string
		m_Port int
		m_ServerType int
	}

	IMonitorClient interface {
		actor.IActor

		RegisterServer(int, string, int)
		Connect(int, string, string, string, string)
	}
)

func (this *MonitorClient) RegisterServer(ServerType int, Ip string, Port int)  {
	this.m_pClient.SendMsg("COMMON_RegisterRequest",ServerType, Ip, Port)
}

func (this *MonitorClient) Connect(ServerType int, Ip string, Port string, oIp string, oPort string){
	this.m_ServerType = ServerType
	this.m_pClient = new(network.ClientSocket)
	this.m_pClient.Init(Ip, base.Int(Port))
	this.m_Ip = oIp
	this.m_Port = base.Int(oPort)
	this.Init(1000)
	this.m_pClient.BindPacketFunc(this.PacketFunc)
	this.m_pClient.Start()
}

func (this *MonitorClient) Init(num int) {
	this.Actor.Init(num)
	this.m_LostTimer = NewSimpleTimer(10)
	this.m_LostTimer.Start()
	this.RegisterTimer(1 * 1000 * 1000 * 1000, this.Update)
	this.RegisterCall("COMMON_RegisterRequest", func() {
		this.RegisterServer(this.m_ServerType, this.m_Ip, this.m_Port)
	})

	this.RegisterCall("COMMON_RegisterResponse", func() {
		this.m_LostTimer.Stop()
	})

	this.RegisterCall("SnowFlake_WorkId", func(workId int) {
		base.UUID.Init(int64(workId))
		base.GLOG.Println("收到monitor初始化 snowflakes [", workId, "]", base.UUID.UUID())
	})

	this.RegisterCall("DISCONNECT", func(socketId int) {
		this.m_LostTimer.Start()
	})

	this.Actor.Start()
}

func (this* MonitorClient) Update(){
	if this.m_LostTimer.CheckTimer(){
		this.m_pClient.Start()
	}
}


package netgate

import (
	"gonet/actor"
	"gonet/base"
	"gonet/message"
	"gonet/rpc"
	"gonet/server/common"
	"strconv"
)

type (
	WorldProcess struct {
		actor.Actor
		m_LostTimer *common.SimpleTimer

		m_ClusterId int
	}

	IWorldlProcess interface {
		actor.IActor

		RegisterServer(int, string, int)
		SetClusterId(int)
	}
)

func (this * WorldProcess) RegisterServer(ServerType int, Ip string, Port int)  {
	SERVER.GetWorldCluster().SendMsg(this.m_ClusterId, "COMMON_RegisterRequest",ServerType, Ip, Port)
}

func (this * WorldProcess) SetClusterId(clusterId int){
	this.m_ClusterId = clusterId
}

func (this *WorldProcess) Init(num int) {
	this.Actor.Init(num)
	this.m_LostTimer = common.NewSimpleTimer(3)
	this.m_LostTimer.Start()
	this.m_ClusterId = 0
	this.RegisterTimer(1 * 1000 * 1000 * 1000, this.Update)
	this.RegisterCall("COMMON_RegisterRequest", func() {
		port,_:=strconv.Atoi(UserNetPort)
		this.RegisterServer(int(message.SERVICE_GATESERVER), UserNetIP, port)
	})

	this.RegisterCall("COMMON_RegisterResponse", func() {
		//收到worldserver对自己注册的反馈
		this.m_LostTimer.Stop()
		SERVER.GetLog().Println("收到world对自己注册的反馈")
	})

	this.RegisterCall("STOP_ACTOR", func() {
		this.Stop()
	})

	this.RegisterCall("DISCONNECT", func(socketId int) {
		this.m_LostTimer.Start()
	})

	this.Actor.Start()
}

func (this* WorldProcess) Update(){
	if this.m_LostTimer.CheckTimer(){
		SERVER.GetWorldCluster().GetCluster(this.m_ClusterId).Start()
	}
}

func DispatchWorldPacketToClient(id int, buff []byte) bool{
	defer func(){
		if err := recover(); err != nil{
			base.TraceCode(err)
		}
	}()

	rpcPacket := rpc.UnmarshalHead(buff)
	socketId := SERVER.GetPlayerMgr().GetSocket(rpcPacket.RpcHead.Id)
	SERVER.GetServer().SendById(socketId, rpcPacket.RpcBody)
	return true
}

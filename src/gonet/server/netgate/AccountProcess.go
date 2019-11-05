package netgate

import (
	"gonet/actor"
	"gonet/message"
	"strconv"
	"gonet/server/common"
)

type (
	AccountProcess struct {
		actor.Actor
		m_LostTimer *common.SimpleTimer

		m_Id uint32
	}

	IAccountProcess interface {
		actor.IActor

		RegisterServer(int, int, string, int)
		SetSocketId(uint32)
	}
)

func (this * AccountProcess) SetSocketId(socketId uint32){
	this.m_Id = socketId
}

func (this * AccountProcess)RegisterServer(ServerType int, Ip string, Port int)  {
	SERVER.GetAccountCluster().GetCluster(this.m_Id).SendMsg("COMMON_RegisterRequest",ServerType, Ip, Port)
}

func (this *AccountProcess) Init(num int) {
	this.Actor.Init(num)
	this.m_LostTimer = common.NewSimpleTimer(3)
	this.m_LostTimer.Start()
	this.RegisterTimer(1 * 1000 * 1000 * 1000, this.Update)
	this.RegisterCall("COMMON_RegisterRequest", func() {
		port,_:=strconv.Atoi(UserNetPort)
		this.RegisterServer(int(message.SERVICE_GATESERVER), UserNetIP, port)
	})

	this.RegisterCall("COMMON_RegisterResponse", func() {
		this.m_LostTimer.Stop()
	})

	this.RegisterCall("STOP_ACTOR", func() {
		this.Stop()
	})

	this.RegisterCall("DISCONNECT", func(socketId int) {
		this.m_LostTimer.Start()
	})

	this.RegisterCall("A_G_Account_Login", func(accountId int64, socketId int) {
		SERVER.GetPlayerMgr().SendMsg("ADD_ACCOUNT", accountId, socketId)
	})

	this.RegisterCall("A_C_RegisterResponse", func(packet *message.A_C_RegisterResponse) {
		buff := message.Encode(packet)
		SERVER.GetServer().SendById(int(packet.GetSocketId()), buff)
	})

	this.RegisterCall("A_C_LoginRequest", func(packet *message.A_C_LoginRequest) {
		buff := message.Encode(packet)
		SERVER.GetServer().SendById(int(packet.GetSocketId()), buff)
	})

	this.Actor.Start()
}

func (this* AccountProcess) Update(){
	if this.m_LostTimer.CheckTimer(){
		SERVER.GetAccountCluster().GetCluster(this.m_Id).Start()
	}
}
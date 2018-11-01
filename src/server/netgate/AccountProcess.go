package netgate

import (
	"actor"
	"message"
	"strconv"
	"server/common"
)

type (
	AccountProcess struct {
		actor.Actor
		m_LostTimer *common.SimpleTimer
	}

	IAccountProcess interface {
		actor.IActor

		RegisterServer(int, int, string, int)
	}
)

func (this * AccountProcess)RegisterServer(ServerType int, Ip string, Port int)  {
	SERVER.GetAccountSocket().SendMsg("COMMON_RegisterRequest",ServerType, Ip, Port)
}

func (this *AccountProcess) Init(num int) {
	this.Actor.Init(num)
	this.m_LostTimer = common.NewSimpleTimer(10)
	this.m_LostTimer.Start()
	this.RegisterTimer(1 * 1000 * 1000 * 1000, this.Update)
	this.RegisterCall("COMMON_RegisterRequest", func() {
		port,_:=strconv.Atoi(UserNetPort)
		this.RegisterServer(int(message.SERVICE_GATESERVER), UserNetIP, port)
	})

	this.RegisterCall("COMMON_RegisterResponse", func() {
		this.m_LostTimer.Stop()
	})

	this.RegisterCall("DISCONNECT", func(socketId int) {
		this.m_LostTimer.Start()
	})

	this.RegisterCall("G_ClientLost", func(accountId int) {
		SERVER.GetWorldSocket().SendMsg("G_ClientLost", accountId)
	})

	this.RegisterCall("A_G_Account_Login", func(accountId int, socketId int) {
		SERVER.GetPlayerMgr().SendMsg("ADD_ACCOUNT", socketId, accountId)
	})

	this.RegisterCall("A_C_RegisterResponse", func(packet *message.A_C_RegisterResponse) {
		buff := message.Encode(packet)
		SERVER.GetServer().SendByID(int(*packet.SocketId), buff)
	})

	this.RegisterCall("A_C_LoginRequest", func(packet *message.A_C_LoginRequest) {
		buff := message.Encode(packet)
		SERVER.GetServer().SendByID(int(*packet.SocketId), buff)
	})

	this.RegisterCall("A_W_CreatePlayer", func(accountId int, playerId int, playername string, sex int32) {
		SERVER.GetWorldSocket().SendMsg("A_W_CreatePlayer", accountId, playerId, playername, sex)
	})

	this.Actor.Start()
}

func (this* AccountProcess) Update(){
	if this.m_LostTimer.CheckTimer(){
		SERVER.GetAccountSocket().Start()
	}
}
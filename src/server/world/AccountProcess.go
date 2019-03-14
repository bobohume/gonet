package world

import (
	"gonet/actor"
	"gonet/base"
	"gonet/message"
	"gonet/server/common"
)

type (
	AccountProcess struct {
		actor.Actor
		m_LostTimer *common.SimpleTimer
	}

	IAccountProcess interface {
		actor.IActor

		RegisterServer(int, string, int)
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
		this.RegisterServer(int(message.SERVICE_WORLDSERVER), UserNetIP, base.Int(UserNetPort))
	})

	this.RegisterCall("COMMON_RegisterResponse", func() {
		this.m_LostTimer.Stop()
	})

	this.RegisterCall("DISCONNECT", func(socketId int) {
		this.m_LostTimer.Start()
	})

	this.RegisterCall("G_ClientLost", func(accountId int64) {
		SERVER.GetServer().CallMsg("G_ClientLost", accountId)
	})

	this.RegisterCall("A_W_CreatePlayer", func(accountId int64, playerId int64, playername string, sex int32, socketId int) {
		SERVER.GetServer().CallMsg("A_W_CreatePlayer", accountId, playerId, playername, sex, socketId)
	})

	this.Actor.Start()
}

func (this* AccountProcess) Update(){
	if this.m_LostTimer.CheckTimer(){
		SERVER.GetAccountSocket().Start()
	}
}

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

		m_Id uint32
	}

	IAccountProcess interface {
		actor.IActor

		RegisterServer(int, string, int)
		SetSocketId(uint32)
	}
)

func (this * AccountProcess) SetSocketId(socketId uint32){
	this.m_Id = socketId
}

func (this * AccountProcess) RegisterServer(ServerType int, Ip string, Port int)  {
	SERVER.GetAccountCluster().SendMsg(this.m_Id, "COMMON_RegisterRequest",ServerType, Ip, Port)
}

func (this *AccountProcess) Init(num int) {
	this.Actor.Init(num)
	this.m_LostTimer = common.NewSimpleTimer(3)
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

	this.RegisterCall("STOP_ACTOR", func() {
		this.Stop()
	})

	this.RegisterCall("G_ClientLost", func(accountId int64) {
		actor.MGR.SendMsg("playermgr", "G_ClientLost", accountId)
	})

	this.RegisterCall("A_W_CreatePlayer", func(accountId int64, playerId int64, playername string, sex int32, socketId int) {
		actor.MGR.SendMsg("playermgr", "A_W_CreatePlayer", accountId, playerId, playername, sex, socketId)
	})

	this.Actor.Start()
}

func (this* AccountProcess) Update(){
	if this.m_LostTimer.CheckTimer(){
		SERVER.GetAccountCluster().GetCluster(this.m_Id).Start()
	}
}
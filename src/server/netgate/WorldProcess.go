package netgate

import (
	"actor"
	"message"
	"strconv"
	"server/common"
)

type (
	WorldProcess struct {
		actor.Actor
		m_LostTimer *common.SimpleTimer
	}

	IWorldlProcess interface {
		actor.IActor

		RegisterServer(int, int, string, int)
	}
)

func (this * WorldProcess)RegisterServer(ServerType int,  ServerId int, Ip string, Port int)  {
	SERVER.GetWorldScoket().SendMsg("COMMON_RegisterRequest",ServerType, ServerId, Ip, Port)
}

func (this *WorldProcess) Init(num int) {
	this.Actor.Init(num)
	this.m_LostTimer = common.NewSimpleTimer(3)
	this.RegisterTimer(1 * 1000 * 1000 * 1000, this.Update)
	this.RegisterCall("COMMON_RegisterRequest", func(caller *actor.Caller) {
		port,_:=strconv.Atoi(UserNetPort)
		this.RegisterServer(int(message.SERVICE_GATESERVER), SERVER.m_GateId, UserNetIP, port)
	})

	this.RegisterCall("COMMON_RegisterResponse", func(caller *actor.Caller) {
			//收到worldserver对自己注册的反馈
			this.m_LostTimer.Stop()
			SERVER.GetLog().Println("收到world对自己注册的反馈")
			SERVER.GetPlayerMgr().SendMsg(caller.SocketId, "Account_Relink")
	})

	this.RegisterCall("G_ClientLost", func(caller *actor.Caller, accountId int) {
		SERVER.GetAccountScoket().SendMsg("G_ClientLost", accountId)
	})

	this.RegisterCall("DISCONNECT", func(caller *actor.Caller, socketId int) {
		this.m_LostTimer.Start()
	})

	this.RegisterCall("W_A_CreatePlayer", func(caller *actor.Caller, accountId int, playername string, sex int32) {
		SERVER.GetAccountScoket().SendMsg("W_A_CreatePlayer", accountId, playername, sex)
	})

	this.RegisterCall("W_A_DeletePlayer", func(caller *actor.Caller, accountId int, playerId int) {
		SERVER.GetAccountScoket().SendMsg("W_A_DeletePlayer", accountId, playerId)
	})

	this.Actor.Start()
}

func (this* WorldProcess) Update(){
	if this.m_LostTimer.CheckTimer(){
		SERVER.GetWorldScoket().Start()
	}
}



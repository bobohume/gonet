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

func (this * AccountProcess)RegisterServer(ServerType int,  ServerId int, Ip string, Port int)  {
	SERVER.GetAccountScoket().SendMsg("COMMON_RegisterRequest",ServerType, ServerId, Ip, Port)
}

func (this *AccountProcess) Init(num int) {
	this.Actor.Init(num)
	this.m_LostTimer = common.NewSimpleTimer(3)
	this.RegisterTimer(1 * 1000 * 1000 * 1000, this.Update)
	this.RegisterCall("COMMON_RegisterRequest", func(caller *actor.Caller) {
		port,_:=strconv.Atoi(UserNetPort)
		this.RegisterServer(int(message.SERVICE_GATESERVER), SERVER.m_GateId, UserNetIP, port)
	})

	this.RegisterCall("COMMON_RegisterResponse", func(caller *actor.Caller) {
		this.m_LostTimer.Stop()
	})

	this.RegisterCall("G_ClientLost", func(caller *actor.Caller, accountId int) {
		SERVER.GetWorldScoket().SendMsg("G_ClientLost", accountId)
	})

	this.RegisterCall("A_G_Account_Login", func(caller *actor.Caller, accountId int, socketId int) {
		SERVER.GetPlayerMgr().SendMsg(caller.SocketId, "ADD_ACCOUNT", socketId, accountId)
	})

	this.RegisterCall("DISCONNECT", func(caller *actor.Caller, socketId int) {
		this.m_LostTimer.Start()
	})

	this.RegisterCall("A_C_RegisterResponse", func(caller *actor.Caller, packet *message.A_C_RegisterResponse) {
		buff := message.Encode(packet)
		SERVER.GetServer().SendByID(int(*packet.SocketId), buff)
	})

	this.RegisterCall("A_C_LoginRequest", func(caller *actor.Caller, packet *message.A_C_LoginRequest) {
		buff := message.Encode(packet)
		SERVER.GetServer().SendByID(int(*packet.SocketId), buff)
	})

	this.RegisterCall("A_W_CreatePlayer", func(caller *actor.Caller, accountId int, playerId int, playername string, sex int32) {
		SERVER.GetWorldScoket().SendMsg("A_W_CreatePlayer", accountId, playerId, playername, sex)
	})

	this.Actor.Start()
}

func (this* AccountProcess) Update(){
	if this.m_LostTimer.CheckTimer(){
		SERVER.GetAccountScoket().Start()
	}
}




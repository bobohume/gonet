package worlddb

import (
	"gonet/actor"
	"gonet/base"
	"gonet/message"
	"gonet/server/common"
)

type (
	WorldProcess struct {
		actor.Actor
		m_LostTimer *common.SimpleTimer
	}

	IWorldProcess interface {
		actor.IActor

		RegisterServer(int, string, int)
	}
)

func (this * WorldProcess)RegisterServer(ServerType int, Ip string, Port int)  {
	SERVER.GetWorldSocket().SendMsg("COMMON_RegisterRequest",ServerType, Ip, Port)
}

func (this *WorldProcess) Init(num int) {
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

	this.Actor.Start()
}

func (this* WorldProcess) Update(){
	if this.m_LostTimer.CheckTimer(){
		SERVER.GetWorldSocket().Start()
	}
}

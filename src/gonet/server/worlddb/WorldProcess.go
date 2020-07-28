package worlddb

import (
	"context"
	"gonet/actor"
	"gonet/base"
	"gonet/message"
	"gonet/rpc"
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
func (this *WorldProcess) Init(num int) {
	this.Actor.Init(num)
	this.m_LostTimer = common.NewSimpleTimer(3)
	this.m_LostTimer.Start()
	this.RegisterTimer(1 * 1000 * 1000 * 1000, this.Update)
	this.RegisterCall("COMMON_RegisterRequest", func(ctx context.Context) {
		SERVER.GetWorldSocket().SendMsg(rpc.RpcHead{},"COMMON_RegisterRequest", &message.ClusterInfo{Type:message.SERVICE_WORLDDBSERVER, Ip:UserNetIP, Port:int32(base.Int(UserNetPort))})
	})

	this.RegisterCall("COMMON_RegisterResponse", func(ctx context.Context) {
		this.m_LostTimer.Stop()
	})

	this.RegisterCall("DISCONNECT", func(ctx context.Context, socketId uint32) {
		this.m_LostTimer.Start()
	})

	this.Actor.Start()
}

func (this* WorldProcess) Update(){
	if this.m_LostTimer.CheckTimer(){
		SERVER.GetWorldSocket().Start()
	}
}

package gate

import (
	"context"
	"gonet/actor"
	"gonet/rpc"
)

type (
	EventProcess struct {
		actor.Actor
	}

	IEventProcess interface {
		actor.IActor
	}
)

func (this *EventProcess) Init() {
	this.Actor.Init()
	actor.MGR.RegisterActor(this)
	this.Actor.Start()
}

func (this *EventProcess) A_G_Account_Login(ctx context.Context, socketId uint32, clusterInfo rpc.PlayerClusterInfo) {
	actor.MGR.SendMsg(rpc.RpcHead{}, "ADD_ACCOUNT", socketId, clusterInfo)
}
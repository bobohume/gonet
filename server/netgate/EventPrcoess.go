package netgate

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

func (this *EventProcess) Init(num int) {
	this.Actor.Init(num)

	this.RegisterCall("A_G_Account_Login", func(ctx context.Context, accountId int64, socketId uint32) {
		SERVER.GetPlayerMgr().SendMsg(rpc.RpcHead{},"ADD_ACCOUNT", accountId, socketId)
	})


	this.Actor.Start()
}

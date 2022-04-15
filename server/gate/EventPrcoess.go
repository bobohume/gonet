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

func (this *EventProcess) G_Player_Login(ctx context.Context, socketId uint32, mailbox rpc.MailBox) {
	actor.MGR.SendMsg(rpc.RpcHead{}, "PlayerMgr.ADD_ACCOUNT", socketId, mailbox)
}

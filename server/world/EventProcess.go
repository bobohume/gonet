package world

import (
	"context"
	"database/sql"
	"gonet/actor"
)

type (
	EventProcess struct {
		actor.Actor
		m_db *sql.DB
	}

	IEventProcess interface {
		actor.IActor
	}
)

func (this *EventProcess) Init() {
	this.Actor.Init()
	this.m_db = SERVER.GetDB()
	
	this.RegisterCall("G_ClientLost", func(ctx context.Context, accountId int64) {
		head := this.GetRpcHead(ctx)
		head.ActorName = "playermgr"
		actor.MGR.SendMsg(head, "G_ClientLost", accountId)
	})

	this.Actor.Start()
}

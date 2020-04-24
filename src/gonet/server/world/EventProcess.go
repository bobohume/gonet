package world

import (
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

func (this *EventProcess) Init(num int) {
	this.Actor.Init(num)
	this.m_db = SERVER.GetDB()
	
	this.RegisterCall("G_ClientLost", func(accountId int64) {
		head := this.GetRpcHead()
		head.ActorName = "playermgr"
		actor.MGR.SendMsg(head, "G_ClientLost", accountId)
	})

	this.Actor.Start()
}

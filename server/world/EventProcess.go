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

func (this *EventProcess) Init() {
	this.Actor.Init()
	this.m_db = SERVER.GetDB()
	actor.MGR.RegisterActor(this)
	this.Actor.Start()
}

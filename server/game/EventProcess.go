package game

import (
	"gonet/actor"
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

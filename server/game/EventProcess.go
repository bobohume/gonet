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

func (e *EventProcess) Init() {
	e.Actor.Init()
	actor.MGR.RegisterActor(e)
	e.Actor.Start()
}

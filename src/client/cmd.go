package main

import (
	"actor"
	"server/common"
)

type (
	CmdProcess struct {
		actor.Actor
	}

	ICmdProcess interface {
		actor.IActor
	}
)

func (this *CmdProcess) Init(num int) {
	this.Actor.Init(num)
	this.RegisterCall("move", func(args string) {
		//PACKET.Move(cm.M_2PI / 4, 100.0)
	})

	this.Actor.Start()
}

var(
	g_Cmd *CmdProcess
)

func InitCmd(){
	g_Cmd = &CmdProcess{}
	g_Cmd.Init(1000)
	common.StartConsole(g_Cmd)
}


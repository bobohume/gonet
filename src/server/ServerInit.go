package main

import (
	"actor"
	"server/world"
	"server/world/chat"
	"server/world/cmd"
	"server/world/data"
	"server/world/mail"
	"server/world/player"
	"server/world/social"
	"server/world/toprank"
)

func InitMgr(serverName string){
	if serverName == "account"{
	}else if serverName == "monitor"{
	}else if serverName == "netgate"{
	}else if serverName == "world"{
		cmd.Init()
		data.InitRepository()
		player.PLAYERMGR.Init(1000)
		chat.CHATMGR.Init(1000)
		mail.MAILMGR.Init(1000)
		toprank.MGR().Init(1000)
		player.PLAYERSIMPLEMGR.Init(1000)
		social.MGR().Init(1000)
		actor.MGR().InitActorHandle(world.SERVER.GetServer())
	}
}
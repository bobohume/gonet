package main

import (
	"server/world/player"
	"server/world"
	"server/world/mail"
	"server/world/toprank"
	"server/world/cmd"
	"actor"
	"server/world/chat"
	"server/world/social"
	"server/world/data"
)

func InitMgr(serverName string){
	if serverName == "account"{
	}else if serverName == "netgate"{
	} else if serverName == "world"{
		cmd.Init()
		data.InitRepository()
		player.PLAYERMGR.Init(1000)
		chat.CHATMGR.Init(1000)
		mail.MAILMGR.Init(1000)
		toprank.TOPMGR.Init(1000)
		player.PLAYERSIMPLEMGR.Init(1000)
		social.SOCIALMGR.Init(1000)
		actor.GetGActorList().InitGActorListHandle(world.SERVER.GetServer())
	}
}

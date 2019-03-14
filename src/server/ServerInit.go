package main

import (
	"gonet/actor"
	"gonet/server/world"
	"gonet/server/world/chat"
	"gonet/server/world/cmd"
	"gonet/server/world/data"
	"gonet/server/world/mail"
	"gonet/server/world/player"
	"gonet/server/world/social"
	"gonet/server/world/toprank"
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
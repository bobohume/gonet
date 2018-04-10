package main

import (
	"server/world/player"
	"server/world"
	"server/world/mail"
	"server/world/data"
	"server/world/toprank"
	"server/world/cmd"
	"actor"
	"server/world/chat"
	"server/world/social"
)

func InitMgr(serverName string){
	if serverName == "account"{
	}else if serverName == "netgate"{
	} else if serverName == "world"{
		cmd.Init()
		player.PLAYERMGR.Init(1000)
		//world.SERVER.GetServer().BindPacketFunc(player.PLAYERMGR.PacketFunc)
		chat.CHATMGR.Init(1)
		mail.MAILMGR.Init(1)
		toprank.TOPMGR.Init(1)
		player.PLAYERSIMPLEMGR.Init(1)
		social.SOCIALMGR.Init(1)
		actor.GetGActorList().InitGActorListHandle(world.SERVER.GetServer())
		data.InitRepository()
	}
}

package main

import (
	"gonet/actor"
	"gonet/common"
	"gonet/server/world"
	"gonet/server/world/chat"
	"gonet/server/world/cmd"
	"gonet/server/world/data"
	"gonet/server/world/mail"
	"gonet/server/world/player"
	"gonet/server/world/social"
	"gonet/server/world/toprank"
	"gonet/server/zone"
	"gonet/server/zone/game"
	data2 "gonet/server/zone/game/data"
)


func InitMgr(serverName string){
	//一些共有数据量初始化
	common.Init()
	if serverName == "account"{
	}else if serverName == "netgate"{
	}else if serverName == "world"{
		cmd.Init()
		data.InitRepository()
		player.MGR.Init(1000)
		chat.MGR.Init(1000)
		mail.MGR.Init(1000)
		toprank.MGR().Init(1000)
		player.SIMPLEMGR.Init(1000)
		social.MGR().Init(1000)
		actor.MGR.InitActorHandle(world.SERVER.GetServer())
	}else if serverName == "zone"{
		data2.InitRepository()
		game.MAPMGR.Init(1000)
		actor.MGR.InitActorHandle(zone.SERVER.GetServer())
	}
}

//程序退出后执行
func ExitMgr(serverName string){
	if serverName == "account"{
	}else if serverName == "netgate"{
	}else if serverName == "world"{
	}
}
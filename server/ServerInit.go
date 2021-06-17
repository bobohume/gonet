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
)


func InitMgr(serverName string){
	//一些共有数据量初始化
	common.Init()
	if serverName == "account"{
	}else if serverName == "netgate"{
	}else if serverName == "world"{
		cmd.Init()
		data.InitRepository()
		player.MGR.Init()
		chat.MGR.Init()
		mail.MGR.Init()
		toprank.MGR().Init()
		player.SIMPLEMGR.Init()
		social.MGR().Init()
		actor.MGR.InitActorHandle(world.SERVER.GetCluster())
	}
}

//程序退出后执行
func ExitMgr(serverName string){
	if serverName == "account"{
	}else if serverName == "netgate"{
	}else if serverName == "world"{
	}
}
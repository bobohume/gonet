package main

import (
	"gonet/common"
	"gonet/server/game/cmd"
	"gonet/server/game/data"
	"gonet/server/game/mail"
	"gonet/server/game/player"
	"gonet/server/gm/chat"
	"gonet/server/gm/login"
)


func InitMgr(serverName string){
	//一些共有数据量初始化
	common.Init()
	if serverName == "gm"{
		login.Init()
		chat.MGR.Init()
	}else if serverName == "gate"{
	}else if serverName == "game"{
		cmd.Init()
		data.InitRepository()
		player.MGR.Init()
		mail.MGR.Init()
	}else if serverName == "db"{
	}
}

//程序退出后执行
func ExitMgr(serverName string){
	if serverName == "gm"{
	}else if serverName == "gate"{
	}else if serverName == "game"{
	}else if serverName == "db"{
	}
}